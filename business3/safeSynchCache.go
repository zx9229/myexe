package main

import (
	"sync"

	"github.com/zx9229/myexe/txdata"
)

//UniSym omit
type UniSym struct {
	UserID string
	MsgNo  int64
	SeqNo  int32
}

func (thls *UniSym) fromUniKey(src *txdata.UniKey) {
	thls.UserID = src.UserID
	thls.MsgNo = src.MsgNo
	thls.SeqNo = src.SeqNo
	assert4true(0 <= thls.MsgNo)
	assert4true(0 <= thls.SeqNo)
}

type safeSynchCache struct {
	sync.Mutex
	M map[UniSym]ProtoMessage
}

func newSafeSynchCache() *safeSynchCache {
	return &safeSynchCache{M: make(map[UniSym]ProtoMessage)}
}

//insertData 入参uniKey是入参pm的一个数据成员.
func (thls *safeSynchCache) insertData(uniKey *txdata.UniKey, pm ProtoMessage) (isSuccess bool) {

	var sym UniSym
	sym.fromUniKey(uniKey)
	thls.Lock()
	if _, isSuccess = thls.M[sym]; !isSuccess {
		thls.M[sym] = pm
	}
	thls.Unlock()
	isSuccess = !isSuccess
	return
}

func (thls *safeSynchCache) deleteData(uniKey *txdata.UniKey) (pm ProtoMessage, isSuccess bool) {
	var sym UniSym
	sym.fromUniKey(uniKey)
	thls.Lock()
	if pm, isSuccess = thls.M[sym]; isSuccess {
		delete(thls.M, sym)
	}
	thls.Unlock()
	return
}
