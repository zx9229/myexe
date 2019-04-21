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

type node4sync struct {
	Key      UniSym
	TxToRoot bool
	RecverID string
	data     ProtoMessage
}

type safeSynchCache struct {
	sync.Mutex
	M map[UniSym]*node4sync
}

func newSafeSynchCache() *safeSynchCache {
	return &safeSynchCache{M: make(map[UniSym]*node4sync)}
}

//insertData 入参uniKey是入参pm的一个数据成员.
func (thls *safeSynchCache) insertData(uniKey *txdata.UniKey, toR bool, rID string, pm ProtoMessage) (isSuccess bool) {
	node := new(node4sync)
	node.Key.fromUniKey(uniKey)
	node.TxToRoot = toR
	node.RecverID = rID
	node.data = pm
	thls.Lock()
	if _, isSuccess = thls.M[node.Key]; !isSuccess {
		thls.M[node.Key] = node
	}
	thls.Unlock()
	isSuccess = !isSuccess
	return
}

func (thls *safeSynchCache) deleteData(uniKey *txdata.UniKey) (node *node4sync, isSuccess bool) {
	var sym UniSym
	sym.fromUniKey(uniKey)
	thls.Lock()
	if node, isSuccess = thls.M[sym]; isSuccess {
		delete(thls.M, sym)
	}
	thls.Unlock()
	return
}

func (thls *safeSynchCache) queryData(toR bool, rID string) (slcOut []*node4sync) {
	thls.Lock()
	for _, node := range thls.M {
		if node.TxToRoot == toR && node.RecverID == rID {
			if slcOut == nil {
				slcOut = make([]*node4sync, 0)
			}
			slcOut = append(slcOut, node)
		}
	}
	thls.Unlock()
	return
}

func (thls *safeSynchCache) queryDataByTxToRoot(toR bool) (slcOut []*node4sync) {
	thls.Lock()
	for _, node := range thls.M {
		if node.TxToRoot == toR {
			if slcOut == nil {
				slcOut = make([]*node4sync, 0)
			}
			slcOut = append(slcOut, node)
		}
	}
	thls.Unlock()
	return
}
