package main

import (
	"sync"

	"github.com/zx9229/myexe/txdata"
)

type tmpVal struct {
	Key    UniSym
	Req    ProtoMessage
	RspSlc []ProtoMessage
}

type safeMemoryTmpCache struct {
	sync.Mutex
	capacity int
	M        map[UniSym]*tmpVal
	Slc      []*tmpVal
}

func newSafeMemoryTmpCache() *safeMemoryTmpCache {
	return &safeMemoryTmpCache{capacity: 10000, M: make(map[UniSym]*tmpVal), Slc: make([]*tmpVal, 0)}
}

func (thls *safeMemoryTmpCache) insertReqData(k *txdata.UniKey, reqData ProtoMessage) (isSuccess bool) {
	assert4true(k.SeqNo == 0)
	var sym UniSym
	sym.fromUniKey(k)
	thls.Lock()
	if _, isSuccess = thls.M[sym]; !isSuccess {
		for (0 < thls.capacity) && (thls.capacity <= len(thls.Slc)) {
			delete(thls.M, thls.Slc[0].Key)
			thls.Slc = thls.Slc[1:]
		}
		tVal := &tmpVal{Key: sym, Req: reqData, RspSlc: make([]ProtoMessage, 0)}
		thls.M[tVal.Key] = tVal
		thls.Slc = append(thls.Slc, tVal)
	}
	thls.Unlock()
	isSuccess = !isSuccess
	return
}

func (thls *safeMemoryTmpCache) appendRspData(k *txdata.UniKey, rspData ProtoMessage) (isSuccess bool) {
	assert4true(k.SeqNo != 0)
	var sym UniSym
	sym.fromUniKey(k)
	var tVal *tmpVal
	thls.Lock()
	if tVal, isSuccess = thls.M[sym]; isSuccess {
		if tVal.RspSlc == nil {
			tVal.RspSlc = make([]ProtoMessage, 0)
		}
		tVal.RspSlc = append(tVal.RspSlc, rspData)
	}
	thls.Unlock()
	return isSuccess
}

func (thls *safeMemoryTmpCache) queryData(k *txdata.UniKey) (isSuccess bool) {
	assert4true(k.SeqNo == 0)
	var sym UniSym
	sym.fromUniKey(k)
	var tVal *tmpVal
	thls.Lock()
	if tVal, isSuccess = thls.M[sym]; isSuccess {
		tVal.Req = tVal.Req
		//TODO:
	}
	thls.Unlock()
	return isSuccess
}
