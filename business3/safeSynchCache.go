package main

import (
	"encoding/json"
	"fmt"
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
	Key      UniSym       //(一经设置,不再修改)
	ToRoot   bool         //(一经设置,不再修改)
	RecverID string       //(一经设置,不再修改)
	data     ProtoMessage //(一经设置,不再修改)
}

type safeSynchCache struct {
	sync.Mutex
	M map[UniSym]*node4sync
}

func newSafeSynchCache() *safeSynchCache {
	return &safeSynchCache{M: make(map[UniSym]*node4sync)}
}

//insertData 入参uniKey是入参pm的一个数据成员.
func (thls *safeSynchCache) insertData(uniKey *txdata.UniKey, toR bool, rID string, pm ProtoMessage) (isExist, isInsert bool) {
	//isExist, isInsert
	//true   , false    (已经存在了,肯定插不进去)
	//false  , true
	//false  , false (异常情况)
	node := new(node4sync)
	node.Key.fromUniKey(uniKey)
	node.ToRoot = toR
	node.RecverID = rID
	node.data = pm
	thls.Lock()
	if _, isExist = thls.M[node.Key]; !isExist {
		thls.M[node.Key] = node
		isInsert = true
	}
	thls.Unlock()
	return
}

func (thls *safeSynchCache) deleteData(uniKey *txdata.UniKey) (node *node4sync, isExist bool) {
	var sym UniSym
	sym.fromUniKey(uniKey)
	thls.Lock()
	if node, isExist = thls.M[sym]; isExist {
		delete(thls.M, sym)
	}
	thls.Unlock()
	return
}

func (thls *safeSynchCache) queryData(toR bool, rID string) (slcOut []*node4sync) {
	thls.Lock()
	for _, node := range thls.M {
		if node.ToRoot == toR && node.RecverID == rID {
			if slcOut == nil {
				slcOut = make([]*node4sync, 0)
			}
			slcOut = append(slcOut, node)
		}
	}
	thls.Unlock()
	return
}

func (thls *safeSynchCache) queryCount(uID string, msgNo int64) (cnt int) {
	thls.Lock()
	for _, node := range thls.M {
		if node.Key.UserID == uID && node.Key.MsgNo == msgNo {
			cnt++
		}
	}
	thls.Unlock()
	return
}

func (thls *safeSynchCache) queryDataByToRoot(toR bool) (slcOut []*node4sync) {
	thls.Lock()
	for _, node := range thls.M {
		if node.ToRoot == toR {
			if slcOut == nil {
				slcOut = make([]*node4sync, 0)
			}
			slcOut = append(slcOut, node)
		}
	}
	thls.Unlock()
	return
}

//MarshalJSON 为了能通过[json.Marshal(obj)]而编写的函数.
func (thls *safeSynchCache) MarshalJSON() ([]byte, error) {
	tmpMap := make(map[string]string)
	thls.Lock()
	for k, v := range thls.M {
		tmpK := fmt.Sprintf("(%v|%v|%v)", k.UserID, k.MsgNo, k.SeqNo)
		tmpV := fmt.Sprintf("(%v|%v|%v)", v.RecverID, v.ToRoot, CalcMessageType(v.data))
		tmpMap[tmpK] = tmpV
	}
	thls.Unlock()
	return json.Marshal(tmpMap)
}
