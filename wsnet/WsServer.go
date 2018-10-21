package wsnet

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type socketSet struct {
	sync.Mutex
	M map[*WsSocket]bool
}

//WsServer omit
type WsServer struct {
	upgrader       websocket.Upgrader
	connSet        socketSet
	CbConnected    NetConnected
	CbDisconnected NetDisconnected
	CbReceive      NetMessage
}

//NewWsServer omit
func NewWsServer() *WsServer {
	curData := new(WsServer)
	curData.connSet.M = make(map[*WsSocket]bool)
	return curData
}

//WsHandler omit
func (thls *WsServer) WsHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var wsConnection *websocket.Conn
	if wsConnection, err = thls.upgrader.Upgrade(w, r, nil); err != nil {
		log.Print("upgrade:", err)
		return
	}
	wsSock := new(WsSocket)
	wsSock.setConnection(wsConnection)

	defer func() {
		if thls.CbDisconnected != nil {
			thls.CbDisconnected(wsSock, err)
		}

		thls.connSet.Lock()
		delete(thls.connSet.M, wsSock)
		thls.connSet.Unlock()

		wsSock.wsConn.Close()
	}()

	if true {
		thls.connSet.Lock()
		thls.connSet.M[wsSock] = true
		thls.connSet.Unlock()
	}

	if thls.CbConnected != nil {
		thls.CbConnected(wsSock, true)
	}

	go thls.doHeartbeat(wsSock) //在服务端心跳,不应在客户端心跳,因为客户端有可能telnet过来的之类的.

	var msgType int
	var msgData []byte
	for {
		if msgType, msgData, err = wsSock.wsConn.ReadMessage(); err != nil {
			break
		}
		thls.CbReceive(wsSock, msgData, msgType)
	}
}

func (thls *WsServer) doHeartbeat(wsSock *WsSocket) {
	//这个变量是一个指针,里面的socket指针可能因为断线重连而被外部修改成新的指针,
	//如果修改成新的指针,就意味着这个连接又启动了一个心跳协程,所以老的协程一定要退出才行.
	realConnection := wsSock.wsConn
	//TODO:这个心跳有问题,超时断开
	msgData := []byte("heartbeat")
	for (wsSock.wsConn == realConnection) && (wsSock.WriteMessageSafe(websocket.PingMessage, msgData) == nil) {
		time.Sleep(time.Second * 60 * 2)
	}
}
