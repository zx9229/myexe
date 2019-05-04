package main

import (
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/zx9229/myexe/txdata"
	"github.com/zx9229/myexe/wsnet"
)

//CommonRspWrapper omit.
type CommonRspWrapper interface {
	sendData(data ProtoMessage, isLast bool) bool
}

//Common1RspWrapper omit
type Common1RspWrapper struct {
	sync.Mutex
	conn    *wsnet.WsSocket
	isLast  bool
	reqData *txdata.Common1Req
}

func newCommon1RspWrapper(req *txdata.Common1Req, conn *wsnet.WsSocket) *Common1RspWrapper {
	return &Common1RspWrapper{conn: conn, reqData: req}
}

//doRemainder 把剩余的事情做完. 执行(善后/清理)工作.
func (thls *Common1RspWrapper) doRemainder() {
	//执行(善后/清理)工作
	thls.Lock()
	defer thls.Unlock()
	if thls.isLast {
		return
	}
	if !thls.sendDataWithoutLock(&txdata.CommonErr{ErrNo: 1, ErrMsg: "handler not implemented"}, true) {
		//TODO:报警.
	}
}

func (thls *Common1RspWrapper) sendDataWithoutLock(data ProtoMessage, isLast bool) bool {
	curRspData := txdata.Common1Rsp{}
	curRspData.RequestID = thls.reqData.RequestID
	curRspData.SenderID = thls.reqData.RecverID
	curRspData.RecverID = thls.reqData.SenderID
	curRspData.TxToRoot = !thls.reqData.TxToRoot
	curRspData.IsLog = thls.reqData.IsLog
	curRspData.IsPush = thls.reqData.IsPush
	if data != nil {
		curRspData.RspType = CalcMessageType(data)
		curRspData.RspData = msg2slice(data)
	}
	curRspData.RspTime, _ = ptypes.TimestampProto(time.Now())
	curRspData.IsLast = isLast

	if !thls.reqData.IsPush {
		thls.conn.Send(msg2package(&curRspData))
	}
	thls.isLast = curRspData.IsLast

	return true
}

func (thls *Common1RspWrapper) sendData(data ProtoMessage, isLast bool) bool {
	thls.Lock()
	defer thls.Unlock()
	if thls.isLast {
		assert4true(thls.isLast == false)
		return false
	}
	return thls.sendDataWithoutLock(data, isLast)
}
