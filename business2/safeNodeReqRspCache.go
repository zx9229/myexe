package main

import (
	"sync"
)

type nodeReqRsp struct {
	sync.Mutex
	key struct {
		UID string
		NUM int64
	}
	condVar *myCondVariable
	reqData ProtoMessage
	rspData []ProtoMessage
}

func newNodeReqRsp() *nodeReqRsp {
	//return &nodeReqRsp{condVar: newMyCondVariable(), rspData: make([]ProtoMessage, 0)}
	return &nodeReqRsp{condVar: newMyCondVariable()}
}

func (thls *nodeReqRsp) setKey(uid string, num int64) {
	thls.key.UID = uid
	thls.key.NUM = num
}

func (thls *nodeReqRsp) setReqData(msg ProtoMessage) {
	thls.reqData = msg
}

func (thls *nodeReqRsp) appendRspData(msg ProtoMessage) {
	assert4true(msg != nil)
	thls.Lock()
	if thls.rspData == nil {
		thls.rspData = make([]ProtoMessage, 0)
	}
	thls.rspData = append(thls.rspData, msg)
	thls.Unlock()
}

type safeNodeReqRspCache struct {
	sync.Mutex
	M map[struct {
		UID string
		NUM int64
	}]*nodeReqRsp
}

func newSafeNodeReqRspCache() *safeNodeReqRspCache {
	return &safeNodeReqRspCache{M: make(map[struct {
		UID string
		NUM int64
	}]*nodeReqRsp)}
}

func (thls *safeNodeReqRspCache) insertNode(node *nodeReqRsp) (isSuccess bool) {
	thls.Lock()
	if _, isSuccess = thls.M[node.key]; !isSuccess {
		thls.M[node.key] = node
	}
	thls.Unlock()
	isSuccess = !isSuccess
	return
}

func (thls *safeNodeReqRspCache) deleteNode(uid string, num int64) (node *nodeReqRsp, isExist bool) {
	key := struct {
		UID string
		NUM int64
	}{uid, num}
	thls.Lock()
	if node, isExist = thls.M[key]; isExist {
		delete(thls.M, key)
	}
	thls.Unlock()
	return
}
