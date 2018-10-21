/*
搜索"Windows 应用程序 注册 系统服务"可知,微软有一个名为"instsrv.exe",可知,微软比较喜欢吧server缩写成srv
*/

package wsnet

import (
	"net/url"
)

//WsCliSrv websocket的client和server的联合体
type WsCliSrv struct {
	cliURL         map[string]url.URL
	srvURL         url.URL
	shs            *SimpleHttpServer
	wss            *WsServer
	wcli           map[string]*WsClient
	CbConnected    NetConnected
	CbDisconnected NetDisconnected
	CbReceive      NetMessage
}

func NewWsCliSrv() *WsCliSrv {
	curData := new(WsCliSrv)
	curData.cliURL = make(map[string]url.URL)
	curData.wcli = make(map[string]*WsClient)
	return curData
}

//GetSimpleHttpServer omit
func (thls *WsCliSrv) GetSimpleHttpServer() *SimpleHttpServer {
	return thls.shs
}

//Init 只允许调用一次
func (thls *WsCliSrv) Init(cliData []url.URL, srvData url.URL) (err error) {
	//检查入参
	tmpMap := make(map[string]url.URL)
	if cliData != nil {
		for _, u := range cliData {
			if len(u.String()) == 0 {
				err = errPlaceholder
				return
			}
			tmpMap[u.String()] = u
		}
		if len(tmpMap) != len(cliData) {
			err = errPlaceholder
			return
		}
	}
	if len(srvData.String()) == 0 {
		err = errPlaceholder
		return
	}
	//检查缓存
	if (len(thls.srvURL.String()) != 0) || (len(thls.cliURL) != 0) {
		err = errPlaceholder
		return
	}
	//初始化
	thls.cliURL = tmpMap
	thls.srvURL = srvData
	//
	thls.shs = New_SimpleHttpServer(thls.srvURL.Host)
	//
	return
}

//Run omit
func (thls WsCliSrv) Run() {
	for _, u := range thls.cliURL {
		cli := NewWsClient()
		cli.CbConnected = thls.CbConnected
		cli.CbDisconnected = thls.CbDisconnected
		cli.CbReceive = thls.CbReceive
		thls.wcli[u.String()] = cli
		cli.Connect(u, true)
	}

	thls.wss = NewWsServer()
	thls.shs.GetHttpServeMux().HandleFunc(thls.srvURL.Path, thls.wss.WsHandler)
	thls.wss.CbConnected = thls.CbConnected
	thls.wss.CbDisconnected = thls.CbDisconnected
	thls.wss.CbReceive = thls.CbReceive
	thls.shs.Run()
}
