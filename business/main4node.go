package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/zx9229/myexe/txdata"
	"github.com/zx9229/myexe/wsnet"
)

var globalA *businessNode

func handleReportData(w http.ResponseWriter, r *http.Request) {
	curObj := new(struct {
		txdata.ReportDataItem
		Cache   bool
		Timeout int
	})
	obj2msg := func(obj interface{}) (req *txdata.CommonNtosReq, saveDB bool, sec int) {
		theObj := obj.(*struct {
			txdata.ReportDataItem
			Cache   bool
			Timeout int
		})
		var err error
		req = &txdata.CommonNtosReq{ReqType: CalcMessageType(&curObj.ReportDataItem)}
		if req.ReqData, err = proto.Marshal(&theObj.ReportDataItem); err != nil {
			glog.Fatalln(err, obj)
		}
		saveDB = theObj.Cache
		sec = theObj.Timeout
		return
	}
	handleCommonFun(w, r, curObj, obj2msg)
}

func handleSendMail(w http.ResponseWriter, r *http.Request) {
	curObj := new(struct {
		txdata.SendMailItem
		Cache   bool
		Timeout int
	})
	obj2msg := func(obj interface{}) (req *txdata.CommonNtosReq, saveDB bool, sec int) {
		theObj := obj.(*struct {
			txdata.SendMailItem
			Cache   bool
			Timeout int
		})
		var err error
		req = &txdata.CommonNtosReq{ReqType: CalcMessageType(&curObj.SendMailItem)}
		if req.ReqData, err = proto.Marshal(&theObj.SendMailItem); err != nil {
			glog.Fatalln(err, obj)
		}
		saveDB = theObj.Cache
		sec = theObj.Timeout
		return
	}
	handleCommonFun(w, r, curObj, obj2msg)
}

func handleCommonFun(w http.ResponseWriter, r *http.Request, obj interface{}, Obj2Msg func(obj interface{}) (*txdata.CommonNtosReq, bool, int)) {
	var err error
	var byteSlice []byte
	//
	rspData := &CommRspData{}
	for range "1" {
		if byteSlice, err = ioutil.ReadAll(r.Body); err != nil {
			rspData.ErrMsg = fmt.Sprintf("read Request with err = %v", err)
			break
		}
		r.Body.Close()
		if err = json.Unmarshal(byteSlice, obj); err != nil {
			rspData.ErrMsg = fmt.Sprintf("Unmarshal Request with err = %v", err)
			break
		}
		var reqInOut *txdata.CommonNtosReq
		var saveDB bool
		var secTimeout int
		if true {
			reqInOut, saveDB, secTimeout = Obj2Msg(obj)
		}
		var rspOut *txdata.CommonNtosRsp
		if (globalA != nil) && (globalS == nil) {
			rspOut = globalA.commonAtos(reqInOut, saveDB, time.Duration(secTimeout)*time.Second)
		} else if (globalA == nil) && (globalS != nil) {
			rspOut = globalS.commonAtos(reqInOut, saveDB, time.Duration(secTimeout)*time.Second)
		} else {
			glog.Fatalln(globalA, globalS)
		}
		if true {
			rspData.UserID = reqInOut.UserID
			rspData.SeqNo = rspOut.SeqNo
			rspData.ErrNo = rspOut.ErrNo
			rspData.ErrMsg = rspOut.ErrMsg
		}
	}
	if byteSlice, err = json.Marshal(rspData); err != nil {
		glog.Fatalln(err)
	}
	fmt.Fprintf(w, string(byteSlice))
}

func cacheNode4a(w http.ResponseWriter, r *http.Request) {
	jsonContent := globalA.cacheUser.humanReadable()
	fmt.Fprintf(w, jsonContent)
}

func runNode(cfg *configNode) {
	globalA = newBusinessNode(cfg)
	cs := wsnet.NewWsCliSrv()
	cs.CbConnected = globalA.onConnected
	cs.CbDisconnected = globalA.onDisconnected
	cs.CbReceive = globalA.onMessage
	cs.Init(cfg.ClientURL, cfg.ServerURL)
	cs.GetSimpleHttpServer().GetHttpServeMux().HandleFunc("/cacheNode4a", cacheNode4a)
	cs.GetSimpleHttpServer().GetHttpServeMux().HandleFunc("/reportData", handleReportData)
	cs.GetSimpleHttpServer().GetHttpServeMux().HandleFunc("/sendMail", handleSendMail)
	cs.Run()
}
