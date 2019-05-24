package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	_ "github.com/mattn/go-sqlite3"
	"github.com/zx9229/myexe/txdata"
	"github.com/zx9229/myexe/wsnet"
)

func main() {
	init4glog()
	var cfgNode *configNode
	//////////////////////////////////////////////////////////////////////////
	if true {
		var err error
		var argHelp bool
		var argJSON string
		flag.BoolVar(&argHelp, "help", false, "[M] display this help and exit.")
		flag.StringVar(&argJSON, "json", "", "[M] json configuration file.")
		flag.Parse()
		if argHelp {
			flag.Usage()
			return
		}
		if cfgNode, err = toConfigNode(argJSON); err != nil {
			glog.Errorf("filename=[%v], err=[%v]", argJSON, err)
			return
		}
	}
	defer glog.Flush()
	//////////////////////////////////////////////////////////////////////////
	glog.Infoln(os.Args)
	//cfgNode := toConfigNode(os.Args[1])
	globalNode := newBusinessNode(cfgNode)
	cs := wsnet.NewWsCliSrv()
	cs.CbConnected = globalNode.onConnected
	cs.CbDisconnected = globalNode.onDisconnected
	cs.CbReceive = globalNode.onMessage
	cs.Init(cfgNode.ClientURL, cfgNode.ServerURL)

	cs.GetSimpleHttpServer().GetHttpServeMux().HandleFunc("/cache", func(w http.ResponseWriter, r *http.Request) { handleNodeCache(globalNode, w, r) })
	cs.GetSimpleHttpServer().GetHttpServeMux().HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) { handleEchoItem(globalNode, w, r) })
	cs.Run()
}

func init4glog() {
	//-alsologtostderr
	//	log to standard error as well as files
	//	日志写入标准错误和文件.
	//-log_backtrace_at value
	//	when logging hits line file:N, emit a stack trace
	//-log_dir string
	//	If non-empty, write log files in this directory
	//	如果非空,写日志文件到此目录,而不是默认的临时目录.
	//-logtostderr
	//	log to standard error instead of files
	//	日志写入标准错误而不是文件.
	//-stderrthreshold value    (INFO/WARNING/ERROR)
	//	logs at or above this threshold go to stderr
	//	达到或高于此等级的日志将写入标准错误(和文件).
	//-v value
	//	log level for V logs
	//-vmodule value
	//	comma-separated list of pattern=N settings for file-filtered logging
	//备注:
	//Log line format: [IWEF]mmdd hh:mm:ss.uuuuuu threadid file:line] msg
	//[IWEF]警告级别的首字母(INFO/WARNING/ERROR/FATAL)
	log2dir := new(string)
	stderrthreshold := new(string)
	cmdLine := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	cmdLine.Usage = func() {}
	cmdLine.StringVar(log2dir, "log_dir", "", "")
	cmdLine.StringVar(stderrthreshold, "stderrthreshold", "", "")
	cmdLine.Parse(os.Args[1:])
	if *log2dir != EMPTYSTR {
		fmt.Println("init4glog", "set", "MkdirAll", *log2dir, os.MkdirAll(*log2dir, os.ModePerm))
		fmt.Println("init4glog", "set", "alsologtostderr")
		flag.Set("alsologtostderr", "true")
	} else {
		fmt.Println("init4glog", "set", "logtostderr")
		flag.Set("logtostderr", "true")
	}
	if *stderrthreshold == EMPTYSTR {
		fmt.Println("init4glog", "set", "stderrthreshold")
		flag.Set("stderrthreshold", "INFO")
	}
	fmt.Println("========================================")
}

func handleNodeCache(node *businessNode, w http.ResponseWriter, r *http.Request) {
	var jsonContent string
	//////////////////////////////////////////////////////////////////////////
	if byteSlice, err := json.Marshal(node); err != nil {
		glog.Fatalln(err, node)
	} else {
		jsonContent = string(byteSlice)
	}
	//////////////////////////////////////////////////////////////////////////
	const templateContent = `
<html>
<head>
<title>{{.MyTitle}}</title>
</head>
<body>
<p id="content">{{.MyContent}}</p>
</body>
<script type="text/javascript">
  var text = document.getElementById('content').innerText;//获取字符串.
  var obje = JSON.stringify(JSON.parse(text), null, 2);//将json字符串转换成json对象.
  document.getElementById('content').innerText = obje;
</script>
</html>
`
	pageVariables := struct {
		MyTitle   string
		MyContent string
	}{MyTitle: "nodeCache", MyContent: jsonContent}

	t := template.Must(template.New("name").Parse(templateContent)) //不理解[template.New("name")]的意义.
	if err := t.Execute(w, pageVariables); err != nil {
		glog.Fatalln(err)
	}
}

func handleCommon2Fun(node *businessNode, w http.ResponseWriter, r *http.Request, obj interface{}, Obj2Msg func(obj interface{}) (*txdata.Common1Req, *txdata.Common2Req, int)) {
	var err error
	var byteSlice []byte

	resultSlice := make([]struct {
		Name string
		Data ProtoMessage
	}, 0)
	resultNode := struct {
		Name string
		Data ProtoMessage
	}{}

	ceData := txdata.CommonErr{}

	for range FORONCE {
		if byteSlice, err = ioutil.ReadAll(r.Body); err != nil {
			ceData.ErrNo = 1
			ceData.ErrMsg = fmt.Sprintf("read Request with err = %v", err)
			break
		}
		r.Body.Close()
		if err = json.Unmarshal(byteSlice, obj); err != nil {
			ceData.ErrNo = 1
			ceData.ErrMsg = fmt.Sprintf("Unmarshal Request with err = %v", err)
			break
		}

		var reqData *txdata.Common2Req
		var secTimeout int
		if true {
			_, reqData, secTimeout = Obj2Msg(obj)
			assert4true(reqData != nil)
		}
		rspSlice := node.syncExecuteCommon2ReqRsp(reqData, time.Duration(secTimeout)*time.Second)

		assert4true(ceData.ErrNo == 0)
		assert4true(len(resultSlice) == 0)
		if true {
			resultNode.Data = reqData.Key
			resultNode.Name = reflect.TypeOf(resultNode.Data).Elem().Name()
			resultSlice = append(resultSlice, resultNode)
		}
		for _, rspItem := range rspSlice {
			if resultNode.Data, err = slice2msg(rspItem.RspType, rspItem.RspData); err != nil {
				assert4true(ceData.ErrNo == 0)
				ceData.ErrMsg = fmt.Sprintf("can_not_unmarshal_data(%v)", rspItem.RspType)
				resultNode.Data = &ceData
			}
			resultNode.Name = reflect.TypeOf(resultNode.Data).Elem().Name()
			resultSlice = append(resultSlice, resultNode)
		}
	}
	if ceData.ErrNo != 0 {
		resultNode.Data = &ceData
		resultNode.Name = reflect.TypeOf(resultNode.Data).Elem().Name()
		assert4true(len(resultSlice) == 0)
		resultSlice = append(resultSlice, resultNode)
	}

	if byteSlice, err = json.Marshal(resultSlice); err != nil {
		glog.Fatalln(err)
	}

	fmt.Fprintf(w, string(byteSlice))
}

func handleCommonFun(node *businessNode, w http.ResponseWriter, r *http.Request, obj interface{}, Obj2Msg func(obj interface{}) (*txdata.Common1Req, *txdata.Common2Req, int)) {
	var err error
	var byteSlice []byte

	resultSlice := make([]struct {
		Name string
		Data ProtoMessage
	}, 0)
	resultNode := struct {
		Name string
		Data ProtoMessage
	}{}

	ceData := txdata.CommonErr{}

	for range FORONCE {
		if byteSlice, err = ioutil.ReadAll(r.Body); err != nil {
			ceData.ErrNo = 1
			ceData.ErrMsg = fmt.Sprintf("read Request with err = %v", err)
			break
		}
		r.Body.Close()
		if err = json.Unmarshal(byteSlice, obj); err != nil {
			ceData.ErrNo = 1
			ceData.ErrMsg = fmt.Sprintf("Unmarshal Request with err = %v", err)
			break
		}

		var c1reqData *txdata.Common1Req
		var c1rspSlice []*txdata.Common1Rsp
		var c2reqData *txdata.Common2Req
		var c2rspSlice []*txdata.Common2Rsp
		var secTimeout int

		c1reqData, c2reqData, secTimeout = Obj2Msg(obj)
		if c1reqData != nil {
			c1rspSlice = node.syncExecuteCommon1ReqRsp(c1reqData, time.Duration(secTimeout)*time.Second)
		} else if c2reqData != nil {
			c2rspSlice = node.syncExecuteCommon2ReqRsp(c2reqData, time.Duration(secTimeout)*time.Second)
		} else {
			panic("logic_error")
		}

		assert4true(ceData.ErrNo == 0)
		assert4true(len(resultSlice) == 0)
		if true {
			if c1reqData != nil {
				resultNode.Data = &txdata.UniKey{UserID: "SIM", MsgNo: c1reqData.MsgNo, SeqNo: 0}
			} else {
				resultNode.Data = &txdata.UniKey{UserID: "SIM", MsgNo: c2reqData.Key.MsgNo, SeqNo: 0}
			}
			resultNode.Name = reflect.TypeOf(resultNode.Data).Elem().Name()
			resultSlice = append(resultSlice, resultNode)
		}
		if c1rspSlice != nil {
			for _, rspItem := range c1rspSlice {
				if resultNode.Data, err = slice2msg(rspItem.RspType, rspItem.RspData); err != nil {
					assert4true(ceData.ErrNo == 0)
					ceData.ErrMsg = fmt.Sprintf("can_not_unmarshal_data(%v)", rspItem.RspType)
					resultNode.Data = &ceData
				}
				resultNode.Name = reflect.TypeOf(resultNode.Data).Elem().Name()
				resultSlice = append(resultSlice, resultNode)
			}
		} else {
			for _, rspItem := range c2rspSlice {
				if resultNode.Data, err = slice2msg(rspItem.RspType, rspItem.RspData); err != nil {
					assert4true(ceData.ErrNo == 0)
					ceData.ErrMsg = fmt.Sprintf("can_not_unmarshal_data(%v)", rspItem.RspType)
					resultNode.Data = &ceData
				}
				resultNode.Name = reflect.TypeOf(resultNode.Data).Elem().Name()
				resultSlice = append(resultSlice, resultNode)
			}
		}
	}
	if ceData.ErrNo != 0 {
		resultNode.Data = &ceData
		resultNode.Name = reflect.TypeOf(resultNode.Data).Elem().Name()
		assert4true(len(resultSlice) == 0)
		resultSlice = append(resultSlice, resultNode)
	}

	if byteSlice, err = json.Marshal(resultSlice); err != nil {
		glog.Fatalln(err)
	}

	fmt.Fprintf(w, string(byteSlice))
}

func toC1C2(rID string, pm ProtoMessage, isLog bool, mode int, isC1NotC2 bool) (c1req *txdata.Common1Req, c2req *txdata.Common2Req) {
	var err error
	if isC1NotC2 {
		c1req = &txdata.Common1Req{}
		//c1req.RequestID
		//c1req.SenderID
		c1req.RecverID = rID
		//c1req.ToRoot
		c1req.IsLog = isLog
		c1req.IsPush, _ = int2mode(mode)
		c1req.ReqType = CalcMessageType(pm)
		if c1req.ReqData, err = proto.Marshal(pm); err != nil {
			glog.Fatalln(err, pm)
		}
		//c1req.ReqTime
	} else {
		c2req = &txdata.Common2Req{}
		//c2req.Key
		//c2req.SenderID
		c2req.RecverID = rID
		//c2req.ToRoot
		c2req.IsLog = isLog
		c2req.IsPush, c2req.IsSafe = int2mode(mode)
		//c2req.UpCache
		c2req.ReqType = CalcMessageType(pm)
		if c2req.ReqData, err = proto.Marshal(pm); err != nil {
			glog.Fatalln(err, pm)
		}
		//c2req.ReqTime
	}
	return
}

func handleEchoItem(node *businessNode, w http.ResponseWriter, r *http.Request) {
	curObj := new(struct {
		txdata.EchoItem
		Recver  string
		Timeout int
		Mode    int
		IsLog   bool
		IsC2    bool
	})
	obj2msg := func(obj interface{}) (c1req *txdata.Common1Req, c2req *txdata.Common2Req, sec int) {
		theObj := obj.(*struct {
			txdata.EchoItem
			Recver  string
			Timeout int
			Mode    int
			IsLog   bool
			IsC2    bool
		})
		c1req, c2req = toC1C2(theObj.Recver, &theObj.EchoItem, theObj.IsLog, theObj.Mode, !theObj.IsC2)
		sec = theObj.Timeout
		return
	}
	handleCommonFun(node, w, r, curObj, obj2msg)
}
