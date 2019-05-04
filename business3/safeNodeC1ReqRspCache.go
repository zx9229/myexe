package main

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/zx9229/myexe/txdata"
)

type nodeC1ReqRsp struct {
	sync.Mutex
	condVar   *myCondVariable
	RequestID int64
	reqData   ProtoMessage
	rspData   []ProtoMessage
}

func newNodeC1ReqRsp() *nodeC1ReqRsp {
	return &nodeC1ReqRsp{condVar: newMyCondVariable()}
}

func (thls *nodeC1ReqRsp) setReqData(msg ProtoMessage) {
	thls.reqData = msg
}

func (thls *nodeC1ReqRsp) appendRspData(msg ProtoMessage) {
	assert4true(msg != nil)
	thls.Lock()
	if thls.rspData == nil {
		thls.rspData = make([]ProtoMessage, 0)
	}
	thls.rspData = append(thls.rspData, msg)
	thls.Unlock()
}

func (thls *nodeC1ReqRsp) xyz() (slcOut []*txdata.Common1Rsp) {
	thls.Lock()
	if thls.rspData != nil {
		slcOut = make([]*txdata.Common1Rsp, 0)
		for _, node := range thls.rspData {
			slcOut = append(slcOut, node.(*txdata.Common1Rsp))
		}
	}
	thls.Unlock()
	return
}

type safeNodeC1ReqRspCache struct {
	sync.Mutex
	M map[int64]*nodeC1ReqRsp
}

func newSafeNodeC1ReqRspCache() *safeNodeC1ReqRspCache {
	return &safeNodeC1ReqRspCache{M: make(map[int64]*nodeC1ReqRsp)}
}

func (thls *safeNodeC1ReqRspCache) insertNode(node *nodeC1ReqRsp) (isSuccess bool) {
	thls.Lock()
	if _, isSuccess = thls.M[node.RequestID]; !isSuccess {
		thls.M[node.RequestID] = node
	}
	thls.Unlock()
	isSuccess = !isSuccess
	return
}

func (thls *safeNodeC1ReqRspCache) deleteNode(reqID int64) (node *nodeC1ReqRsp, isExist bool) {
	thls.Lock()
	if node, isExist = thls.M[reqID]; isExist {
		delete(thls.M, reqID)
	}
	thls.Unlock()
	return
}

func (thls *safeNodeC1ReqRspCache) queryNode(reqID int64) (node *nodeC1ReqRsp, isExist bool) {
	thls.Lock()
	node, isExist = thls.M[reqID]
	thls.Unlock()
	return
}

func (thls *safeNodeC1ReqRspCache) operateNode(reqID int64, rspData ProtoMessage, doNotify bool) (isSuccess bool) {
	assert4true(rspData != nil)
	var node *nodeC1ReqRsp
	if doNotify {
		node, isSuccess = thls.deleteNode(reqID)
	} else {
		node, isSuccess = thls.queryNode(reqID)
	}
	if isSuccess {
		node.appendRspData(rspData)
		if doNotify {
			node.condVar.notifyAll()
		}
	}
	return
}

//MarshalJSON 为了能通过[json.Marshal(obj)]而编写的函数.
func (thls *safeNodeC1ReqRspCache) MarshalJSON() ([]byte, error) {
	tmpMap := make(map[int64]string)
	thls.Lock()
	var rspLen int
	for k, v := range thls.M {
		if v.rspData != nil {
			rspLen = len(v.rspData)
		} else {
			rspLen = 0
		}
		tmpV := fmt.Sprintf("req=%v,rspLen=%v", CalcMessageType(v.reqData), rspLen)
		tmpMap[k] = tmpV
	}
	thls.Unlock()
	return json.Marshal(tmpMap)
}
