package main

import (
	"sync"

	"github.com/golang/glog"
	"github.com/zx9229/myexe/txdata"
	"github.com/zx9229/myexe/wsnet"
)

type SubInfo struct {
	sync.Mutex
	userID string //订阅者.
	nodeID string //从哪个节点订阅.
	toRoot bool   //节点发送数据到订阅者的方向.
	isLog  bool
	isPush bool
	conn   *wsnet.WsSocket
	cache  []*txdata.PushWrap //nil就直接发送,非nil需要缓存.
}

func tmpMerge(srcOld, srcNew []*txdata.PushWrap) []*txdata.PushWrap {
	if srcOld == nil || len(srcOld) == 0 {
		return srcNew
	}
	if srcNew == nil || len(srcNew) == 0 {
		return srcOld
	}
	srcOldLast := srcOld[len(srcOld)-1]
	srcNewIdx := -1
	for idx, srcNewItem := range srcNew {
		if srcOldLast.MsgNo < srcNewItem.MsgNo {
			srcNewIdx = idx
			break
		}
	}
	if 0 <= srcNewIdx {
		srcOld = append(srcOld, srcNew[srcNewIdx:]...)
	}
	return srcOld
}

func (thls *SubInfo) xxx(qwert *safePushCache, msgNo int64) {
	glog.Infof("sub, thls=%v", thls)
	glog.Infof("sub, msgNo=%v", msgNo)
	tmpResults := qwert.Select(msgNo, -1)
	glog.Infof("sub, tmpResults.len=%v", len(tmpResults))
	thls.Lock()
	tmpResults = tmpMerge(tmpResults, thls.cache)
	glog.Infof("sub, tmpResults.2.len=%v", len(tmpResults))
	for _, item := range tmpResults {
		c1req := &txdata.Common1Req{SenderID: thls.nodeID, RecverID: thls.userID, ToRoot: thls.toRoot, IsLog: thls.isLog, IsPush: thls.isPush, ReqType: CalcMessageType(item), ReqData: msg2slice(item)}
		thls.conn.Send(msg2package(c1req))
	}
	thls.cache = nil
	thls.Unlock()
}

func (thls *SubInfo) Send(data *txdata.PushWrap) {
	thls.Lock()
	if thls.cache == nil {
		c1req := &txdata.Common1Req{SenderID: thls.nodeID, RecverID: thls.userID, ToRoot: thls.toRoot, ReqType: CalcMessageType(data), ReqData: msg2slice(data)}
		thls.conn.Send(msg2package(c1req))
	} else {
		thls.cache = append(thls.cache, data)
	}
	thls.Unlock()
}

type safeSubscriberMap struct {
	sync.Mutex
	qw *safePushCache
	M  map[string]*SubInfo
}

func newSafeSubscriberMap(qw *safePushCache) *safeSubscriberMap {
	return &safeSubscriberMap{qw: qw, M: make(map[string]*SubInfo)}
}

func (thls *safeSubscriberMap) insertData(uID string, nID string, toR bool, isLog bool, conn *wsnet.WsSocket, msgNo int64) (isSuccess bool) {
	sInfo := &SubInfo{userID: uID, nodeID: nID, toRoot: toR, isLog: isLog, conn: conn, cache: make([]*txdata.PushWrap, 0)}
	thls.Lock()
	if _, isSuccess = thls.M[sInfo.userID]; !isSuccess {
		thls.M[sInfo.userID] = sInfo
		go sInfo.xxx(thls.qw, msgNo)
	}
	thls.Unlock()
	isSuccess = !isSuccess
	return
}

func (thls *safeSubscriberMap) deleteData(userID string) {
	thls.Lock()
	delete(thls.M, userID)
	thls.Unlock()
}

func (thls *safeSubscriberMap) deleteByConn(conn *wsnet.WsSocket) {
	thls.Lock()
	for key, val := range thls.M {
		if val.conn == conn {
			delete(thls.M, key)
		}
	}
	thls.Unlock()
}

func (thls *safeSubscriberMap) Send(data *txdata.PushWrap) {
	thls.Lock()
	for _, xx := range thls.M {
		xx.Send(data)
	}
	thls.Unlock()
}
