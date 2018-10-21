package main

import (
	"github.com/zx9229/myexe/wsnet"
)

func main() {
	shs := wsnet.New_SimpleHttpServer("localhost:8080")
	wss := wsnet.NewWsServer()
	shs.GetHttpServeMux().HandleFunc("/echo", wss.WsHandler)
	wss.CbConnected = wsnet.Example_CbConnected
	wss.CbDisconnected = wsnet.Example_CbDisconnected
	wss.CbReceive = wsnet.Example_CbReceive
	shs.Run()
}
