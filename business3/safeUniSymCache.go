package main

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/zx9229/myexe/txdata"
)

type safeUniSymCache struct {
	sync.Mutex
	M map[UniSym]bool
}

func newSafeUniSymCache() *safeUniSymCache {
	return &safeUniSymCache{M: make(map[UniSym]bool)}
}

func (thls *safeUniSymCache) insertData(k *txdata.UniKey) (isSuccess bool) {
	assert4true(k.SeqNo == 0)
	var sym UniSym
	sym.fromUniKey(k)
	thls.Lock()
	if _, isSuccess = thls.M[sym]; !isSuccess {
		thls.M[sym] = true
	}
	thls.Unlock()
	isSuccess = !isSuccess
	return
}

func (thls *safeUniSymCache) deleteData(k *txdata.UniKey) (isSuccess bool) {
	assert4true(k.SeqNo == 0)
	var sym UniSym
	sym.fromUniKey(k)
	thls.Lock()
	if _, isSuccess = thls.M[sym]; isSuccess {
		delete(thls.M, sym)
	}
	thls.Unlock()
	return
}
func (thls *safeUniSymCache) deleteDataByField(uID string, msgNo int64) (isSuccess bool) {
	sym := UniSym{UserID: uID, MsgNo: msgNo, SeqNo: 0}
	thls.Lock()
	if _, isSuccess = thls.M[sym]; isSuccess {
		delete(thls.M, sym)
	}
	thls.Unlock()
	return
}

//MarshalJSON 为了能通过[json.Marshal(obj)]而编写的函数.
func (thls *safeUniSymCache) MarshalJSON() ([]byte, error) {
	tmpSlc := make([]string, 0)
	var sym UniSym
	thls.Lock()
	for sym = range thls.M {
		tmpStr := fmt.Sprintf("(%v|%v|%v)", sym.UserID, sym.MsgNo, sym.SeqNo)
		tmpSlc = append(tmpSlc, tmpStr)
	}
	thls.Unlock()
	return json.Marshal(tmpSlc)
}
