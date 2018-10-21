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

func handleReportData(w http.ResponseWriter, r *http.Request) {
	rspData := new(struct {
		UniqueID string
		SeqNo    int64
		ErrNo    int32
		ErrMsg   string
	})

	var err error
	var byteSlice []byte

	for range "1" {

		if byteSlice, err = ioutil.ReadAll(r.Body); err != nil {
			rspData.ErrMsg = fmt.Sprintf("read Request with err = %v", err)
			break
		}
		r.Body.Close()

		reqData := new(struct {
			Topic   string
			Data    string
			Timeout int
		})
		if err = json.Unmarshal(byteSlice, reqData); err != nil {
			rspData.ErrMsg = fmt.Sprintf("Unmarshal Request with err = %v", err)
			break
		}

		reqInOut := txdata.ReportDataReq{Topic: reqData.Topic, Data: reqData.Data}
		rspOut := globalA.reportData(&reqInOut, time.Duration(reqData.Timeout)*time.Second)
		rspData.UniqueID = reqInOut.UniqueID
		rspData.SeqNo = rspOut.SeqNo
		rspData.ErrNo = rspOut.ErrNo
		rspData.ErrMsg = rspOut.ErrMsg
	}

	if byteSlice, err = json.Marshal(rspData); err != nil {
		glog.Fatalln(err)
	}

	fmt.Fprintf(w, string(byteSlice))
}

func handlePushData(w http.ResponseWriter, r *http.Request) {
	rspData := new(struct {
		ErrNo  int32
		ErrMsg string
	})

	var err error
	var byteSlice []byte

	for range "1" {

		if byteSlice, err = ioutil.ReadAll(r.Body); err != nil {
			rspData.ErrMsg = fmt.Sprintf("read Request with err = %v", err)
			break
		}
		r.Body.Close()

		reqData := new(struct {
			Topic string
			Data  string
		})
		if err = json.Unmarshal(byteSlice, reqData); err != nil {
			rspData.ErrMsg = fmt.Sprintf("Unmarshal Request with err = %v", err)
			break
		}

		reqInOut := txdata.PushData{Topic: reqData.Topic, Data: reqData.Data}
		if err = globalA.pushData(&reqInOut); err != nil {
			rspData.ErrNo = -1
			rspData.ErrMsg = err.Error()
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
	cs.GetSimpleHttpServer().GetHttpServeMux().HandleFunc("/reportData", handleReportData)
	cs.GetSimpleHttpServer().GetHttpServeMux().HandleFunc("/pushData", handlePushData)
	cs.GetSimpleHttpServer().GetHttpServeMux().HandleFunc("/cacheAgent4a", cacheAgent4a)
	cs.Run()
}
