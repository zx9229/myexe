package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/golang/glog"
	"github.com/zx9229/myexe/txdata"
	"github.com/zx9229/myexe/wsnet"
)

var globalA *businessAgent

func handleCommonAtos(w http.ResponseWriter, r *http.Request) {
	var err error
	var byteSlice []byte
	//
	rspData := new(struct {
		UniqueID string
		SeqNo    int64
		ErrNo    int32
		ErrMsg   string
	})
	for range "1" {
		if byteSlice, err = ioutil.ReadAll(r.Body); err != nil {
			rspData.ErrMsg = fmt.Sprintf("read Request with err = %v", err)
			break
		}
		r.Body.Close()
		reqData := new(struct {
			Endeavour bool
			DataType  string
			Data      string
			Timeout   int
		})
		if err = json.Unmarshal(byteSlice, reqData); err != nil {
			rspData.ErrMsg = fmt.Sprintf("Unmarshal Request with err = %v", err)
			break
		}
		reqInOut := txdata.CommonAtosReq{Endeavour: reqData.Endeavour, DataType: reqData.DataType, Data: []byte(reqData.Data)}
		rspOut := globalA.commonAtos(&reqInOut, time.Duration(reqData.Timeout)*time.Second)
		if true {
			rspData.UniqueID = reqInOut.UniqueID
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

func cacheAgent4a(w http.ResponseWriter, r *http.Request) {
	jsonContent := globalA.cacheAgent.humanReadable()
	fmt.Fprintf(w, jsonContent)
}

func runAgent(cfg *configAgent) {
	globalA = newBusinessAgent(cfg)
	cs := wsnet.NewWsCliSrv()
	cs.CbConnected = globalA.onConnected
	cs.CbDisconnected = globalA.onDisconnected
	cs.CbReceive = globalA.onMessage
	cs.Init(cfg.ClientURL, cfg.ServerURL)
	cs.GetSimpleHttpServer().GetHttpServeMux().HandleFunc("/commonAtos", handleCommonAtos)
	cs.GetSimpleHttpServer().GetHttpServeMux().HandleFunc("/cacheAgent4a", cacheAgent4a)
	cs.Run()
}
