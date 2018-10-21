package wsnet

import (
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

//WsClient omit
type WsClient struct {
	WsSocket
	doReconnect    bool
	u              url.URL
	CbConnected    NetConnected
	CbDisconnected NetDisconnected
	CbReceive      NetMessage
}

//NewWsClient omit
func NewWsClient() *WsClient {
	curData := new(WsClient)
	return curData
}

//Connect omit
func (thls *WsClient) Connect(u url.URL, doReconnect bool) (err error) {
	if len(thls.u.String()) != 0 {
		err = errPlaceholder
		return
	}
	thls.u = u
	thls.doReconnect = doReconnect
	if thls.doReconnect {
		go thls.reConnect()
	} else {
		err = thls.reConnect()
	}
	return
}

//Stop omit
func (thls *WsClient) Stop() {
	thls.WsSocket.Close()
	thls.doReconnect = false
	thls.u = url.URL{}
}

func (thls *WsClient) reConnect() (err error) {
	err = errPlaceholder
	var wsConnection *websocket.Conn
	var URLstr string
	URLstr = thls.u.String()
	for err != nil {
		if wsConnection, _, err = websocket.DefaultDialer.Dial(URLstr, nil); err != nil {
			time.Sleep(time.Second * 5)
		} else {
			go thls.doRecv(wsConnection)
		}
		if !thls.doReconnect {
			break
		}
	}
	return
}

func (thls *WsClient) doRecv(wsConnection *websocket.Conn) {
	thls.WsSocket.setConnection(wsConnection)

	if thls.CbConnected != nil {
		thls.CbConnected(&thls.WsSocket, false)
	}

	var err error
	defer func() {
		wsConnection.Close()
		if thls.CbDisconnected != nil {
			thls.CbDisconnected(&thls.WsSocket, err)
		}
		if thls.doReconnect {
			go thls.reConnect()
		}
	}()

	var msgType int
	var msgData []byte
	for {
		if msgType, msgData, err = thls.WsSocket.wsConn.ReadMessage(); err != nil {
			break
		}
		thls.CbReceive(&thls.WsSocket, msgData, msgType)
	}
}
