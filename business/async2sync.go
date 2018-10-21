package main

import (
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/zx9229/myexe/txdata"
)

type nodeReqRsp struct {
	requestID int64
	condVar   *myCondVariable
	reqType   txdata.MsgType
	reqData   proto.Message
	rspType   txdata.MsgType
	rspData   proto.Message
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
