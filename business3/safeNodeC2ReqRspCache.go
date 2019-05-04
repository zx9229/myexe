package main

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/zx9229/myexe/txdata"
)

type nodeC2ReqRsp struct {
	sync.Mutex
	condVar *myCondVariable
	key     UniSym
	reqData ProtoMessage
	rspData []ProtoMessage
}

func newNodeC2ReqRsp() *nodeC2ReqRsp {
	//return &nodeC2ReqRsp{condVar: newMyCondVariable(), rspData: make([]ProtoMessage, 0)}
	return &nodeC2ReqRsp{condVar: newMyCondVariable()}
}

func (thls *nodeC2ReqRsp) setReqData(msg ProtoMessage) {
	thls.reqData = msg
}

func (thls *nodeC2ReqRsp) appendRspData(msg ProtoMessage) {
	assert4true(msg != nil)
	thls.Lock()
	if thls.rspData == nil {
		thls.rspData = make([]ProtoMessage, 0)
	}
	thls.rspData = append(thls.rspData, msg)
	thls.Unlock()
}

func (thls *nodeC2ReqRsp) xyz() (slcOut []*txdata.Common2Rsp) {
	thls.Lock()
	if thls.rspData != nil {
		slcOut = make([]*txdata.Common2Rsp, 0)
		for _, node := range thls.rspData {
			slcOut = append(slcOut, node.(*txdata.Common2Rsp))
		}
	}
	thls.Unlock()
	return
}

type safeNodeC2ReqRspCache struct {
	sync.Mutex
	M map[UniSym]*nodeC2ReqRsp
}

func newSafeNodeC2ReqRspCache() *safeNodeC2ReqRspCache {
	return &safeNodeC2ReqRspCache{M: make(map[UniSym]*nodeC2ReqRsp)}
}

func (thls *safeNodeC2ReqRspCache) insertNode(node *nodeC2ReqRsp) (isSuccess bool) {
	thls.Lock()
	if _, isSuccess = thls.M[node.key]; !isSuccess {
		thls.M[node.key] = node
	}
	thls.Unlock()
	isSuccess = !isSuccess
	return
}

func (thls *safeNodeC2ReqRspCache) deleteNode(sym *UniSym) (node *nodeC2ReqRsp, isExist bool) {
	thls.Lock()
	if node, isExist = thls.M[*sym]; isExist {
		delete(thls.M, *sym)
	}
	thls.Unlock()
	return
}

func (thls *safeNodeC2ReqRspCache) queryNode(sym *UniSym) (node *nodeC2ReqRsp, isExist bool) {
	thls.Lock()
	node, isExist = thls.M[*sym]
	thls.Unlock()
	return
}

func (thls *safeNodeC2ReqRspCache) operateNode(uniKey *txdata.UniKey, rspData ProtoMessage, doNotify bool) (isSuccess bool) {
	sym := &UniSym{UserID: uniKey.UserID, MsgNo: uniKey.MsgNo, SeqNo: 0} //从Rsp转成Req要置SeqNo为0才行.
	assert4true(rspData != nil)
	var node *nodeC2ReqRsp
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
func (thls *safeNodeC2ReqRspCache) MarshalJSON() ([]byte, error) {
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
