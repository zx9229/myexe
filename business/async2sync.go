package main

import (
	"sync"
)

type nodeReqRsp struct {
	requestID int64
	condVar   *myCondVariable
	reqData   ProtoMessage
	rspData   ProtoMessage
}

func newNodeReqRsp() *nodeReqRsp {
	return &nodeReqRsp{condVar: newMyCondVariable()}
}

type safeNodeReqRspCache struct {
	sync.Mutex
	reqID int64
	M     map[int64]*nodeReqRsp
}

func newSafeNodeReqRspCache() *safeNodeReqRspCache {
	return &safeNodeReqRspCache{M: make(map[int64]*nodeReqRsp)}
}

func (thls *safeNodeReqRspCache) generateElement() (node *nodeReqRsp) {
	thls.Lock()
	for isOk := true; isOk; {
		thls.reqID++
		if _, isOk = thls.M[thls.reqID]; isOk {
			continue
		}
		node = newNodeReqRsp()
		node.requestID = thls.reqID
		thls.M[node.requestID] = node
	}
	thls.Unlock()
	return
}

func (thls *safeNodeReqRspCache) deleteElement(key int64) (node *nodeReqRsp, isExist bool) {
	thls.Lock()
	if node, isExist = thls.M[key]; isExist {
		delete(thls.M, key)
	}
	thls.Unlock()
	return
}
