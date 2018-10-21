package main

import (
	"fmt"
	"my_code/_jiankong/wsnet"
	"net/url"
	"time"
)

func main() {
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/echo"}
	cli := wsnet.NewWsClient()
	cli.CbConnected = wsnet.Example_CbConnected
	cli.CbDisconnected = wsnet.Example_CbDisconnected
	cli.CbReceive = wsnet.Example_CbReceive
	cli.Connect(u, true)
	fmt.Println(cli)
	for {
		time.Sleep(time.Second * 3)
		cli.Send([]byte("hello"))
	}
}
