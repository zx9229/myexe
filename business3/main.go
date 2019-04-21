package main

import (
	"encoding/json"
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
	glog.Infoln(os.Args)
	cfgNode := toConfigNode(os.Args[1])
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

func handleNodeCache(node *businessNode, w http.ResponseWriter, r *http.Request) {
	var jsonContent string
	//////////////////////////////////////////////////////////////////////////
	tmpObj := new(struct {
		OwnInfo    *txdata.ConnectionInfo
		ParentInfo *safeFatherData
		RootOnline bool
		CacheUser  *safeConnInfoMap
		CacheSock  *safeWsSocketMap
	})
	tmpObj.OwnInfo = &node.ownInfo
	tmpObj.ParentInfo = &node.parentInfo
	tmpObj.RootOnline = node.rootOnline
	tmpObj.CacheUser = node.cacheUser
	tmpObj.CacheSock = node.cacheSock

	tmpObj.ParentInfo.Lock()
	defer tmpObj.ParentInfo.Unlock()
	tmpObj.CacheUser.Lock()
	defer tmpObj.CacheUser.Unlock()
	tmpObj.CacheSock.Lock()
	defer tmpObj.CacheSock.Unlock()

	if byteSlice, err := json.Marshal(tmpObj); err != nil {
		glog.Fatalln(err, tmpObj)
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
		sliceRsp := node.syncExecuteCommonReqRsp(reqData, time.Duration(secTimeout)*time.Second)

		assert4true(ceData.ErrNo == 0)
		assert4true(len(resultSlice) == 0)
		if true {
			resultNode.Data = reqData.Key
			resultNode.Name = reflect.TypeOf(resultNode.Data).Elem().Name()
			resultSlice = append(resultSlice, resultNode)
		}
		for _, itemRsp := range sliceRsp {
			if resultNode.Data, err = slice2msg(itemRsp.RspType, itemRsp.RspData); err != nil {
				assert4true(ceData.ErrNo == 0)
				ceData.ErrMsg = fmt.Sprintf("can_not_unmarshal_data(%v)", itemRsp.RspType)
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
	})
	obj2msg := func(obj interface{}) (req *txdata.CommonReq, sec int) {
		theObj := obj.(*struct {
			txdata.EchoItem
			Recver  string
			Timeout int
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
		sec = theObj.Timeout
		return
	}
	handleCommonFun(node, w, r, curObj, obj2msg)
}
