package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"my_code/jiankong/txdata"
	"my_code/jiankong/wsnet"
	"net/http"
	"time"

	"github.com/golang/glog"
)

var globalS *businessServer

func handleExecuteCommand(w http.ResponseWriter, r *http.Request) {
	rspData := new(struct {
		RequestID int64
		Result    string
		ErrNo     int32
		ErrMsg    string
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
			UID     string
			Cmd     string
			Timeout int
		})
		if err = json.Unmarshal(byteSlice, reqData); err != nil {
			rspData.ErrMsg = fmt.Sprintf("Unmarshal Request with err = %v", err)
			break
		}

		reqInOut := txdata.ExecuteCommandReq{Pathway: []string{reqData.UID}, Command: reqData.Cmd}
		rspOut := globalS.executeCommand(&reqInOut, time.Duration(reqData.Timeout)*time.Second)
		rspData.RequestID = rspOut.RequestID
		rspData.Result = rspOut.Result
		rspData.ErrNo = rspOut.ErrNo
		rspData.ErrMsg = rspOut.ErrMsg
	}

	if byteSlice, err = json.Marshal(rspData); err != nil {
		glog.Fatalln(err)
	}

	fmt.Fprintf(w, string(byteSlice))
}

func cacheAgent4s(w http.ResponseWriter, r *http.Request) {
	jsonContent := globalS.cacheAgent.humanReadable()
	fmt.Fprintf(w, jsonContent)
}

func runServer(cfg *configServer) {
	globalS = newBusinessServer(cfg)
	cs := wsnet.NewWsCliSrv()
	cs.CbConnected = globalS.onConnected
	cs.CbDisconnected = globalS.onDisconnected
	cs.CbReceive = globalS.onMessage
	cs.Init(cfg.ClientURL, cfg.ServerURL)
	cs.GetSimpleHttpServer().GetHttpServeMux().HandleFunc("/executeCommand", handleExecuteCommand)
	cs.GetSimpleHttpServer().GetHttpServeMux().HandleFunc("/cacheAgent4s", cacheAgent4s)
	cs.Run()
}
