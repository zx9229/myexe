package main

import (
	"my_code/_jiankong/wsnet"

	simplehttpserver "github.com/zx9229/zxgo_push/SimpleHttpServer"
)

func main() {
	shs := simplehttpserver.New_SimpleHttpServer("localhost:8080")
	wss := wsnet.NewWsServer()
	shs.GetHttpServeMux().HandleFunc("/echo", wss.WsHandler)
	wss.CbConnected = wsnet.Example_CbConnected
	wss.CbDisconnected = wsnet.Example_CbDisconnected
	wss.CbReceive = wsnet.Example_CbReceive
	shs.Run()
}
