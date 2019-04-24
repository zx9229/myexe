package main

import (
	"encoding/json"
	"fmt"
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
	sym.SeqNo = 0
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

//MarshalJSON 为了能通过[json.Marshal(obj)]而编写的函数.
func (thls *safeMemoryTmpCache) MarshalJSON() ([]byte, error) {
	tmpMap := make(map[string]string)
	thls.Lock()
	var rspLen int
	for k, v := range thls.M {
		if v.RspSlc != nil {
			rspLen = len(v.RspSlc)
		} else {
			rspLen = 0
		}
		tmpK := fmt.Sprintf("(%v|%v|%v)", k.UserID, k.MsgNo, k.SeqNo)
		tmpV := fmt.Sprintf("req=%v,rspLen=%v", CalcMessageType(v.Req), rspLen)
		tmpMap[tmpK] = tmpV
	}
	thls.Unlock()
	return json.Marshal(tmpMap)
}
