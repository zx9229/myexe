package main

import (
	"encoding/json"
	"fmt"
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

func (thls *safeNodeReqRspCache) operateNode(uniKey *txdata.UniKey, rspData ProtoMessage, doNotify bool) (isSuccess bool) {
	sym := &UniSym{UserID: uniKey.UserID, MsgNo: uniKey.MsgNo, SeqNo: 0} //从Rsp转成Req要置SeqNo为0才行.
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

//MarshalJSON 为了能通过[json.Marshal(obj)]而编写的函数.
func (thls *safeNodeReqRspCache) MarshalJSON() ([]byte, error) {
	tmpMap := make(map[string]string)
	thls.Lock()
	var rspLen int
	for k, v := range thls.M {
		if v.rspData != nil {
			rspLen = len(v.rspData)
		} else {
			rspLen = 0
		}
		tmpK := fmt.Sprintf("(%v|%v|%v)", k.UserID, k.MsgNo, k.SeqNo)
		tmpV := fmt.Sprintf("req=%v,rspLen=%v", CalcMessageType(v.reqData), rspLen)
		tmpMap[tmpK] = tmpV
	}
	thls.Unlock()
	return json.Marshal(tmpMap)
}
