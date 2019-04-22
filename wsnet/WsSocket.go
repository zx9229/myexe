/*
谷歌关于websocket的网站:
https://godoc.org/golang.org/x/net/websocket

部分内容如下所示:

package websocket
import "golang.org/x/net/websocket"
Package websocket implements a client and server for the WebSocket protocol as specified in RFC 6455.
This package currently lacks some features found in an alternative and more actively maintained WebSocket package:
https://godoc.org/github.com/gorilla/websocket

可知,谷歌官方比较推荐这个第三方库.
*/

package wsnet

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	errPlaceholder = errors.New("placeholder")
	errOffline     = errors.New("offline")
)

//NetConnected 连接成功的回调函数
type NetConnected func(wSock *WsSocket, isAccepted bool)

//NetDisconnected 连接断开的回调函数
type NetDisconnected func(wSock *WsSocket, err error)

//NetMessage 收到消息的回调函数
type NetMessage func(wSock *WsSocket, msgData []byte, msgType int)

//WsSocket omit
type WsSocket struct {
	sync.Mutex
	wsConn *websocket.Conn
}

func (thls *WsSocket) setConnection(conn *websocket.Conn) {
	//只允许用这个函数修改指针(wsConn)的值,(因为golang不像C++那样可以用权限(私有成员)约束,所以只能人为约束).
	//这样的话,一旦指针被赋值(非nil),那么它将永远有值(非nil)
	if conn == nil {
		return
	}
	thls.Lock()
	thls.wsConn = conn
	thls.Unlock()
}

//LocalAddr omit
func (thls *WsSocket) LocalAddr() net.Addr {
	if thls.wsConn == nil {
		return nil
	}
	return thls.wsConn.LocalAddr()
}

//RemoteAddr omit
func (thls *WsSocket) RemoteAddr() net.Addr {
	if thls.wsConn == nil {
		return nil
	}
	return thls.wsConn.RemoteAddr()
}

//Close omit
func (thls *WsSocket) Close() {
	thls.Lock()
	if thls.wsConn != nil {
		thls.wsConn.Close()
	}
	thls.Unlock()
}

//Send 可能因为断线/入参有误/等原因,返回错误.
func (thls *WsSocket) Send(msgData []byte) error {
	time.Sleep(time.Microsecond * 100) //TODO:待删除(为了能方便地按照时间戳查看日志,在此睡眠一会,临时代码)
	return thls.WriteMessageSafe(websocket.BinaryMessage, msgData)
}

//WriteMessageSafe omit
func (thls *WsSocket) WriteMessageSafe(msgType int, msgData []byte) (err error) {
	thls.Lock()
	if thls.wsConn != nil {
		err = thls.wsConn.WriteMessage(msgType, msgData)
	} else {
		err = errOffline
	}
	thls.Unlock()
	return
}

//Example_CbConnected omit
func Example_CbConnected(wSock *WsSocket, isAccepted bool) {
	log.Println(fmt.Sprintf("[   Connected][%p]LocalAddr=%v,RemoteAddr=%v,isAccepted=%v", wSock, wSock.wsConn.LocalAddr(), wSock.wsConn.RemoteAddr(), isAccepted))
}

//Example_CbDisconnected omit
func Example_CbDisconnected(wSock *WsSocket, err error) {
	log.Println(fmt.Sprintf("[Disconnected][%p]err=%v", wSock, err))
}

//Example_CbReceive omit
func Example_CbReceive(wSock *WsSocket, msgData []byte, msgType int) {
	log.Println(fmt.Sprintf("[     Receive][%p]data=%v,type=%v", wSock, string(msgData), msgType))
}
