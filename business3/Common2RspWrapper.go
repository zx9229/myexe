package main

import (
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/zx9229/myexe/txdata"
	"github.com/zx9229/myexe/wsnet"
)

//Common2RspWrapper omit
type Common2RspWrapper struct {
	sync.Mutex
	upCache bool
	conn    *wsnet.WsSocket
	isLast  bool
	rspIdx  int32
	reqData *txdata.Common2Req
	cache   *safeSynchCache
	zzzxml  *safeUniSymCache
}

func newCommon2RspWrapper(req *txdata.Common2Req, cache *safeSynchCache, zzzxml *safeUniSymCache, upCache bool, conn *wsnet.WsSocket) *Common2RspWrapper {
	return &Common2RspWrapper{upCache: upCache, conn: conn, cache: cache, zzzxml: zzzxml, reqData: req}
}

//doRemainder 把剩余的事情做完. 执行(善后/清理)工作.
func (thls *Common2RspWrapper) doRemainder() {
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

func (thls *Common2RspWrapper) sendDataWithoutLock(data ProtoMessage, isLast bool) bool {
	curRspData := txdata.Common2Rsp{}
	curRspData.Key = cloneUniKey(thls.reqData.Key)
	curRspData.Key.SeqNo = thls.rspIdx + 1
	curRspData.SenderID = thls.reqData.RecverID
	curRspData.RecverID = curRspData.Key.UserID //中间可能因为缓存而修改了(req.SenderID)
	curRspData.ToRoot = !thls.reqData.ToRoot
	curRspData.IsLog = thls.reqData.IsLog
	curRspData.IsSafe = thls.reqData.IsSafe
	curRspData.IsPush = thls.reqData.IsPush
	curRspData.UpCache = thls.upCache && thls.reqData.IsSafe //只有在续传模式下,才允许设置UpCache字段.
	if data != nil {
		curRspData.RspType = CalcMessageType(data)
		curRspData.RspData = msg2slice(data)
	}
	curRspData.RspTime, _ = ptypes.TimestampProto(time.Now())
	curRspData.IsLast = isLast

	if !thls.reqData.IsPush {
		if curRspData.IsSafe {
			var isExist, isInsert bool
			isExist, isInsert = thls.cache.insertData(curRspData.Key, curRspData.ToRoot, curRspData.RecverID, &curRspData)
			assert4false(isExist) //一定不会存在.
			if !isInsert {
				return false
			}
		}
		thls.conn.Send(msg2package(&curRspData))
	}
	thls.rspIdx = curRspData.Key.SeqNo
	thls.isLast = curRspData.IsLast
	if thls.isLast && thls.reqData.IsSafe {
		isOk := thls.zzzxml.deleteData(thls.reqData.Key)
		assert4true(isOk)
	}

	return true
}

func (thls *Common2RspWrapper) sendData(data ProtoMessage, isLast bool) bool {
	thls.Lock()
	defer thls.Unlock()
	if thls.isLast {
		assert4true(thls.isLast == false)
		return false
	}
	return thls.sendDataWithoutLock(data, isLast)
}
