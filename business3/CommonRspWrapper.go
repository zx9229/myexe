package main

import (
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/zx9229/myexe/txdata"
	"github.com/zx9229/myexe/wsnet"
)

//CommonRspWrapper omit
type CommonRspWrapper struct {
	sync.Mutex
	upCache bool
	conn    *wsnet.WsSocket
	isLast  bool
	rspIdx  int32
	reqData *txdata.CommonReq
	cache   *safeSynchCache
}

func newCommonRspWrapper(req *txdata.CommonReq, cache *safeSynchCache, upCache bool, conn *wsnet.WsSocket) *CommonRspWrapper {
	return &CommonRspWrapper{upCache: upCache, conn: conn, cache: cache, reqData: req}
}

//doRemainder 把剩余的事情做完. 执行(善后/清理)工作.
func (thls *CommonRspWrapper) doRemainder() {
	//执行(善后/清理)工作
	thls.Lock()
	defer thls.Unlock()
	if thls.isLast {
		return
	}
	thls.sendDataWithoutLock(&txdata.CommonErr{ErrNo: 1, ErrMsg: "handler not implemented"}, true)
}

func (thls *CommonRspWrapper) sendDataWithoutLock(data ProtoMessage, isLast bool) bool {
	thls.rspIdx++
	thls.isLast = isLast

	curRspData := txdata.CommonRsp{}
	curRspData.Key = cloneUniKey(thls.reqData.Key)
	curRspData.SenderID = thls.reqData.RecverID
	if curRspData.Key != nil {
		curRspData.Key.SeqNo = thls.rspIdx
		curRspData.RecverID = curRspData.Key.UserID //中间可能因为缓存而修改了(req.SenderID)
	} else {
		curRspData.RecverID = thls.reqData.SenderID
	}
	curRspData.TxToRoot = !thls.reqData.TxToRoot
	curRspData.UpCache = thls.upCache
	if data != nil {
		curRspData.RspType = CalcMessageType(data)
		curRspData.RspData = msg2slice(data)
	}
	curRspData.RspTime, _ = ptypes.TimestampProto(time.Now())
	curRspData.IsLast = thls.isLast
	curRspData.IsLog = thls.reqData.IsLog
	curRspData.IsSafe = thls.reqData.IsSafe
	curRspData.IsPush = thls.reqData.IsPush

	if !thls.reqData.IsPush {
		if curRspData.IsSafe {
			isOk := thls.cache.insertData(curRspData.Key, curRspData.TxToRoot, curRspData.RecverID, &curRspData)
			assert4true(isOk)
		}
		thls.conn.Send(msg2package(&curRspData))
	}

	return true
}

func (thls *CommonRspWrapper) sendData(data ProtoMessage, isLast bool) bool {
	thls.Lock()
	defer thls.Unlock()
	if thls.isLast {
		assert4true(thls.isLast == false)
		return false
	}
	return thls.sendDataWithoutLock(data, isLast)
}
