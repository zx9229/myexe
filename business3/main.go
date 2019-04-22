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
			glog.Errorf("filename=%v, err=%v", argJSON, err)
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

func handleCommonFun(node *businessNode, w http.ResponseWriter, r *http.Request, obj interface{}, Obj2Msg func(obj interface{}) (*txdata.CommonReq, int)) {
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

		var reqData *txdata.CommonReq
		var secTimeout int
		if true {
			reqData, secTimeout = Obj2Msg(obj)
		}
		rspSlice := node.syncExecuteCommonReqRsp(reqData, time.Duration(secTimeout)*time.Second)

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

func handleEchoItem(node *businessNode, w http.ResponseWriter, r *http.Request) {
	curObj := new(struct {
		txdata.EchoItem
		Recver  string
		Timeout int
		Mode    int
		IsLog   bool
	})
	obj2msg := func(obj interface{}) (req *txdata.CommonReq, sec int) {
		theObj := obj.(*struct {
			txdata.EchoItem
			Recver  string
			Timeout int
			Mode    int
			IsLog   bool
		})
		var err error
		req = &txdata.CommonReq{}
		//req.Key
		//req.SenderID
		req.RecverID = theObj.Recver
		//req.TxToRoot
		//req.UpCache
		req.ReqType = CalcMessageType(&theObj.EchoItem)
		if req.ReqData, err = proto.Marshal(&theObj.EchoItem); err != nil {
			glog.Fatalln(err, obj)
		}
		//req.ReqTime
		req.IsLog = theObj.IsLog
		req.IsPush, req.IsSafe = int2mode(theObj.Mode)
		sec = theObj.Timeout
		return
	}
	handleCommonFun(node, w, r, curObj, obj2msg)
}
