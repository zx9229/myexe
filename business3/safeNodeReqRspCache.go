package main

import (
	"sync"

	"github.com/zx9229/myexe/txdata"
)

type nodeReqRsp struct {
	sync.Mutex
	condVar *myCondVariable
	key     UniSym
	reqData ProtoMessage
	rspData []ProtoMessage
}

func newNodeReqRsp() *nodeReqRsp {
	//return &nodeReqRsp{condVar: newMyCondVariable(), rspData: make([]ProtoMessage, 0)}
	return &nodeReqRsp{condVar: newMyCondVariable()}
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

func (thls *nodeReqRsp) xyz() (slcOut []*txdata.CommonRsp) {
	thls.Lock()
	if thls.rspData != nil {
		slcOut = make([]*txdata.CommonRsp, 0)
		for _, node := range thls.rspData {
			slcOut = append(slcOut, node.(*txdata.CommonRsp))
		}
	}
	thls.Unlock()
	return
}

type safeNodeReqRspCache struct {
	sync.Mutex
	M map[UniSym]*nodeReqRsp
}

func newSafeNodeReqRspCache() *safeNodeReqRspCache {
	return &safeNodeReqRspCache{M: make(map[UniSym]*nodeReqRsp)}
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

func (thls *safeNodeReqRspCache) deleteNode(sym *UniSym) (node *nodeReqRsp, isExist bool) {
	thls.Lock()
	if node, isExist = thls.M[*sym]; isExist {
		delete(thls.M, *sym)
	}
	thls.Unlock()
	return
}

func (thls *safeNodeReqRspCache) queryNode(sym *UniSym) (node *nodeReqRsp, isExist bool) {
	thls.Lock()
	node, isExist = thls.M[*sym]
	thls.Unlock()
	return
}

func (thls *safeNodeReqRspCache) operateNode(sym *UniSym, rspData ProtoMessage, doNotify bool) (isSuccess bool) {
	assert4true(rspData != nil)
	var node *nodeReqRsp
	if doNotify {
		node, isSuccess = thls.deleteNode(sym)
	} else {
		node, isSuccess = thls.queryNode(sym)
	}
	if isSuccess {
		node.appendRspData(rspData)
		if doNotify {
			node.condVar.notifyAll()
		}
	}
	return
}
