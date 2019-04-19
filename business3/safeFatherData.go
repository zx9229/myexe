package main

import (
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
