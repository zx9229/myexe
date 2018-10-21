package main

import (
	"net/url"

	"github.com/zx9229/myexe/wsnet"
)

func main() {
	serverURL := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/echo"}
	clientURL := make([]url.URL, 0)
	clientURL = append(clientURL, url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/echo"})
	cs := wsnet.NewWsCliSrv()
	cs.Init(clientURL, serverURL)
	cs.CbConnected = wsnet.Example_CbConnected
	cs.CbDisconnected = wsnet.Example_CbDisconnected
	cs.CbReceive = wsnet.Example_CbReceive
	cs.Run()
}
