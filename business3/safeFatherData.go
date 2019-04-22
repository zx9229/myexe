package main

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/zx9229/myexe/txdata"
	"github.com/zx9229/myexe/wsnet"
)

type safeFatherData struct {
	sync.Mutex
	conn *wsnet.WsSocket
	Info txdata.ConnectionInfo
}

func (thls *safeFatherData) setData(newConn *wsnet.WsSocket, newInfo *txdata.ConnectionInfo, isForce bool) (isSuccess bool) {
	thls.Lock()
	if thls.conn == nil || isForce {
		thls.conn = newConn
		if newInfo == nil {
			thls.Info = txdata.ConnectionInfo{}
		} else {
			thls.Info = *newInfo
		}
		//
		isSuccess = true
	} else {
		isSuccess = false
	}
	thls.Unlock()
	return
}

//MarshalJSON 为了能通过[json.Marshal(obj)]而编写的函数.
func (thls *safeFatherData) MarshalJSON() (byteSlice []byte, err error) {
	thls.Lock()
	tmpObj := struct {
		Conn string
		Info *txdata.ConnectionInfo
	}{Conn: fmt.Sprintf("%p", thls.conn), Info: &thls.Info}
	byteSlice, err = json.Marshal(tmpObj)
	thls.Unlock()
	return
}
