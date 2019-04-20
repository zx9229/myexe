package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/golang/glog"
	"github.com/golang/protobuf/ptypes"
	"github.com/zx9229/myexe/txdata"
	"github.com/zx9229/myexe/wsnet"
)

type configNode struct {
	IsRoot         bool   //这个Node是RootNode.
	UserID         string //为空表示是ROOT节点.
	BelongID       string
	ServerURL      url.URL
	ClientURL      []url.URL
	DataSourceName string //数据源的名字.
	LocationName   string //数据源的时区的名字.
}

type businessNode struct {
	letUpCache bool //让上游缓存数据;TODO:做检查(此时它必须是叶子节点).
	ownInfo    txdata.ConnectionInfo
	iAmRoot    bool //(I am root node)一经设置,不允许修改.
	parentInfo safeFatherData
	rootOnline bool
	cacheUser  *safeConnInfoMap
	cacheSock  *safeWsSocketMap
	cacheSync  *safeSynchCache //要绝对的投递过去而缓存+因为UpCache而缓存,所以它绝对会在ROOT的发送侧,不会处于ROOT的对端;即,从sync里取出数据后,肯定要无脑往parent那里发.而不会往孩子那里发送.
	cacheRR    *safeNodeReqRspCache
	ownSeqNo   int64
}

func newBusinessNode(cfg *configNode) *businessNode {
	if false ||
		(cfg.IsRoot && (cfg.UserID != EMPTYSTR || cfg.BelongID != EMPTYSTR)) || //(为防误操作).
		(!cfg.IsRoot && cfg.UserID == EMPTYSTR) ||
		(!cfg.IsRoot && cfg.UserID == cfg.BelongID) {
		glog.Fatalf("newBusinessNode fail with cfg=%v", cfg)
	}

	curData := new(businessNode)
	//
	curData.ownInfo.UserID = cfg.UserID
	curData.ownInfo.BelongID = cfg.BelongID
	curData.ownInfo.Version = "Version20190411"
	curData.ownInfo.LinkMode = txdata.ConnectionInfo_Zero3
	curData.ownInfo.ExePid = int32(os.Getpid())
	curData.ownInfo.ExePath, _ = filepath.Abs(os.Args[0])
	curData.ownInfo.Remark = ""
	//
	curData.iAmRoot = (curData.ownInfo.UserID == EMPTYSTR)
	//
	curData.parentInfo.setData(nil, nil, true)
	//
	curData.cacheSock = newSafeWsSocketMap()
	curData.cacheUser = newSafeConnInfoMap()
	curData.cacheSync = newSafeSynchCache()
	curData.cacheRR = newSafeNodeReqRspCache()
	//
	curData.refreshSeqNo()
	//
	if curData.iAmRoot {
		curData.setRootOnline(true)
	}
	//
	return curData
}

func (thls *businessNode) onConnected(msgConn *wsnet.WsSocket, isAccepted bool) {
	glog.Warningf("[   onConnected] msgConn=%p, isAccepted=%v, LocalAddr=%v, RemoteAddr=%v", msgConn, isAccepted, msgConn.LocalAddr(), msgConn.RemoteAddr())
	if !thls.cacheSock.insertData(msgConn, isAccepted) {
		glog.Fatalf("onConnected, already cached msgConn=%p", msgConn)
	}
	if !isAccepted {
		tmpTxData := txdata.ConnectReq{InfoReq: &thls.ownInfo, Pathway: []string{thls.ownInfo.UserID}}
		thls.sendData(msgConn, &tmpTxData)
	}
}

func (thls *businessNode) onDisconnected(msgConn *wsnet.WsSocket, err error) {
	checkSunWhenDisconnected := func(dataSlice []*connInfoEx) {
		sonNum := 0
		for _, data := range dataSlice {
			if len(data.Pathway) == 1 { //步长为1的是儿子.
				sonNum++
			}
			if len(data.Pathway) == 0 {
				glog.Fatalf("onDisconnected, empty Pathway with data=%v", data)
			}
		}
		if sonNum != 1 {
			glog.Fatalf("onDisconnected, there should be only one son and sonNum=%v", sonNum)
		}
	}
	glog.Warningf("[onDisconnected] msgConn=%p, err=%v", msgConn, err)
	if thls.parentInfo.conn == msgConn {
		//如果与父亲断开连接,就清理父亲的数据,这样就不用sendDataToParent了.
		glog.Infof("onDisconnected, disconnected with father, msgConn=%p", msgConn)
		thls.parentInfo.setData(nil, nil, true)
		if thls.rootOnline {
			thls.setRootOnline(false)
		}
	}
	if dataSlice := thls.cacheUser.deleteDataByConn(msgConn); dataSlice != nil { //儿子和我断开连接,我要清理掉儿子和孙子的缓存.
		checkSunWhenDisconnected(dataSlice)
		for _, data := range dataSlice { //发给父亲,让父亲也清理掉对应的缓存.
			tmpTxData := txdata.DisconnectedData{Info: &data.Info}
			thls.sendData(thls.parentInfo.conn, &tmpTxData)
		}
	}
	thls.deleteConnectionFromAll(msgConn, false)
}

func (thls *businessNode) deleteConnectionFromAll(conn *wsnet.WsSocket, closeIt bool) {
	if closeIt {
		conn.Close()
	}
	if thls.parentInfo.conn == conn {
		//不需要在这里处理它,因为主动断开连接,也会触发onDisconnected回调,回调里面已经有这个逻辑了.
		//thls.parentInfo.setData(nil, nil, true)
		//if thls.rootOnline {
		//	thls.setRootOnline(false)
		//}
	}
	thls.cacheSock.deleteData(conn)
	thls.cacheUser.deleteDataByConn(conn)
}

func (thls *businessNode) sendData(sock *wsnet.WsSocket, data ProtoMessage) {
	if sock != nil {
		sock.Send(msg2package(data))
	}
}

func (thls *businessNode) sendDataEx(sock *wsnet.WsSocket, data ProtoMessage, isParentSock bool) error {
	//如果不是父代的socket那么不会出现nil的情况,此时就让它崩溃.
	if sock == nil && isParentSock {
		return errors.New("parent is offline")
	}
	return sock.Send(msg2package(data))
}

func (thls *businessNode) sendDataEx2(data ProtoMessage, sock *wsnet.WsSocket, txToRoot bool, rID string) error {
	if sock != nil {
		return sock.Send(msg2package(data))
	}
	if txToRoot {
		assert4false(thls.iAmRoot) //此时我一定不是ROOT,否则入参就已经填写错误了.
		return thls.sendDataEx(thls.parentInfo.conn, data, true)
	}
	return thls.cacheUser.sendDataToUser(data, rID)
}

func (thls *businessNode) sendAck(key *txdata.UniKey, rID string, txToRoot bool, conn *wsnet.WsSocket) {
	//UpCache了Req或Rsp,而发送对应的ACK,此时应发往叶子节点.
	//ROOT缓存了数据,而发送对应的ACK,此时应发往叶子节点.
	//Recver收到了ROOT的数据,而发送对应的ACK,此时应发往ROOT节点.
	dataACK := &txdata.MessageAck{Key: key, SenderID: thls.ownInfo.UserID, RecverID: rID, TxToRoot: txToRoot}
	if conn != nil {
		thls.sendData(conn, dataACK)
	} else {
		if dataACK.TxToRoot {
			assert4true(thls.ownInfo.UserID != EMPTYSTR) //我一定不能是ROOT,否则肯定传参错误了.
			thls.sendData(thls.parentInfo.conn, dataACK)
		} else {
			thls.cacheUser.sendDataToUser(dataACK, dataACK.RecverID)
		}
	}
}

func (thls *businessNode) setRootOnline(newValue bool) {
	oldValue := thls.rootOnline
	if oldValue == newValue { //我的目标是:消息无冗余无重复,很显然这里消息重复了.
		glog.Errorf("setRootOnline, oldValue=%v, newValue=%v", oldValue, newValue)
	}
	thls.rootOnline = newValue
	thls.cacheUser.sendDataToSon(&txdata.OnlineNotice{RootIsOnline: newValue})
}

func (thls *businessNode) reportCommonErrMsg(message string) {
	tmpTxData := txdata.CommonErrMsg{UserID: thls.ownInfo.UserID, Pathway: []string{thls.ownInfo.UserID}}
	tmpTxData.Message = message
	thls.sendData(thls.parentInfo.conn, &tmpTxData)
}

func (thls *businessNode) onMessage(msgConn *wsnet.WsSocket, msgData []byte, msgType int) {
	txMsgType, txMsgData, err := package2msg(msgData)
	if err != nil {
		glog.Errorln(txMsgType, txMsgData, err)
		return
	}

	//glog.Infof("onMessage, msgConn=%p, txMsgType=%v, txMsgData=%v", msgConn, txMsgType, txMsgData)

	switch txMsgType {
	case txdata.MsgType_ID_ConnectReq:
		thls.handle_MsgType_ID_ConnectReq(txMsgData.(*txdata.ConnectReq), msgConn)
	case txdata.MsgType_ID_ConnectRsp:
		thls.handle_MsgType_ID_ConnectRsp(txMsgData.(*txdata.ConnectRsp), msgConn)
	case txdata.MsgType_ID_OnlineNotice:
		thls.handle_MsgType_ID_OnlineNotice(txMsgData.(*txdata.OnlineNotice), msgConn)
	case txdata.MsgType_ID_CommonReq:
		thls.handle_MsgType_ID_CommonReq(txMsgData.(*txdata.CommonReq), msgConn)
	case txdata.MsgType_ID_CommonRsp:
		thls.handle_MsgType_ID_CommonRsp(txMsgData.(*txdata.CommonRsp), msgConn)
	default:
		glog.Errorf("onMessage, unknown txdata.MsgType, msgConn=%p, txMsgType=%v, txMsgData=%v", msgConn, txMsgType, txMsgData)
	}
}

func (thls *businessNode) handle_MsgType_ID_ConnectReq(msgData *txdata.ConnectReq, msgConn *wsnet.WsSocket) {
	sendToParent := false

	tmpTxdata := txdata.ConnectRsp{InfoReq: msgData.InfoReq, InfoRsp: &thls.ownInfo, ErrNo: 0}
	if msgData.Pathway == nil || len(msgData.Pathway) == 0 {
		tmpTxdata.ErrNo = 1
		tmpTxdata.ErrMsg = "req.Pathway is empty"
	} else if len(msgData.Pathway) == 1 {
		sendToParent = thls.handle_MsgType_ID_ConnectReq_stepOne(msgData, msgConn, &tmpTxdata)
	} else {
		sendToParent = thls.handle_MsgType_ID_ConnectReq_stepMulti(msgData, msgConn)
	}

	if msgData.Pathway == nil || len(msgData.Pathway) <= 1 { //非ConnectReq_stepMulti时要发送回应.
		thls.sendData(msgConn, &tmpTxdata)
	}

	if sendToParent {
		msgData.Pathway = append(msgData.Pathway, thls.ownInfo.UserID)
		thls.sendData(thls.parentInfo.conn, msgData)
	}
}

func (thls *businessNode) handle_MsgType_ID_ConnectReq_stepOne(msgData *txdata.ConnectReq, msgConn *wsnet.WsSocket, rspData *txdata.ConnectRsp) (sendToParent bool) {
	assert4true(len(msgData.Pathway) == 1)

	for range FORONCE {
		rspData.ErrNo = 1
		if msgData.InfoReq.UserID != msgData.Pathway[0] {
			rspData.ErrMsg = "req.UserID != req.Pathway[0]"
			break
		}
		if (msgData.InfoReq.UserID == EMPTYSTR) && (msgData.InfoReq.BelongID != EMPTYSTR) { //ROOT节点的UserID和BelongID皆为空.
			rspData.ErrMsg = "(req.UserID == EMPTYSTR) && (req.BelongID != EMPTYSTR)"
			break
		}
		if (msgData.InfoReq.UserID != EMPTYSTR) && (msgData.InfoReq.UserID == msgData.InfoReq.BelongID) {
			rspData.ErrMsg = "(req.UserID != EMPTYSTR) && (req.UserID == req.BelongID)"
			break
		}
		if msgData.InfoReq.UserID == thls.ownInfo.UserID {
			rspData.ErrMsg = "req.UserID == rsp.UserID"
			break
		}
		if (msgData.InfoReq.UserID != thls.ownInfo.BelongID) && (msgData.InfoReq.BelongID != thls.ownInfo.UserID) {
			rspData.ErrMsg = "(req.UserID != rsp.BelongID) && (req.BelongID != rsp.UserID)"
			break
		}
		rspData.ErrNo = 0
	}
	if rspData.ErrNo != 0 {
		sendToParent = false
		return
	}

	if msgData.InfoReq.BelongID == thls.ownInfo.UserID {
		sendToParent = thls.handle_MsgType_ID_ConnectReq_stepOne_forSon(msgData, msgConn, rspData)
	} else if msgData.InfoReq.UserID == thls.ownInfo.BelongID {
		sendToParent = thls.handle_MsgType_ID_ConnectReq_stepOne_forParent(msgData, msgConn, rspData)
	} else {
		glog.Errorf("ConnectReq_stepOne, run into unreachable code, msgConn=%p, msgData=%v", msgConn, msgData)
		rspData.ErrNo = 1
		rspData.ErrMsg = "rsp internal error occurred"
		sendToParent = false
	}

	return
}

func (thls *businessNode) handle_MsgType_ID_ConnectReq_stepOne_forSon(msgData *txdata.ConnectReq, msgConn *wsnet.WsSocket, rspData *txdata.ConnectRsp) (sendToParent bool) {
	assert4true(len(msgData.Pathway) == 1)
	assert4true(msgData.InfoReq.BelongID == thls.ownInfo.UserID)

	//我们先假定cacheSock缓存了msgConn.
	curData := new(connInfoEx)
	curData.conn = msgConn
	curData.Info = *msgData.InfoReq
	curData.Pathway = msgData.Pathway

	if isSuccess := thls.cacheUser.insertData(curData); !isSuccess {
		glog.Errorf("ConnectReq_stepOne_forSon, UserID conflict, msgConn=%p, msgData=%v", msgConn, msgData)
		if true { //UserID冲突,应当立即上报该情况.
			errMsg := fmt.Sprintf("UserID conflict, msgData=%v", msgData)
			thls.reportCommonErrMsg(errMsg)
		}
		rspData.ErrNo = 1
		rspData.ErrMsg = "req.UserID is already online"
		sendToParent = false
		return
	}

	var isAccepted bool
	var isExist bool
	if isAccepted, isExist = thls.cacheSock.deleteData(msgConn); !isExist {
		if true { //先cacheUser.insertData,然后发现有异常,需cacheUser.deleteData,以恢复成原来的状态.
			thls.cacheUser.deleteData(curData.Info.UserID) //确保它一定能回退(cacheUser.insertData)操作.
		}
		rspData.ErrNo = 1
		rspData.ErrMsg = "rsp internal error occurred"
		glog.Errorf("ConnectReq_stepOne_forSon, msgConn not found in cache, msgConn=%p, msgData=%v", msgConn, msgData)
		sendToParent = false
		return
	}

	if isAccepted {
		//有如下通信规则:
		//连接建立后,_connect方要主动发送ConnectReq给accepted方.
		//校验通过后,accepted方要主动发送ConnectReq给_connect方.
		tmpTxData := txdata.ConnectReq{InfoReq: &thls.ownInfo, Pathway: []string{thls.ownInfo.UserID}}
		thls.sendData(msgConn, &tmpTxData)
	}

	if thls.rootOnline { //如果我能连通ROOT那么我就把这个消息通知(新建立连接的这个)儿子.
		thls.sendData(msgConn, &txdata.OnlineNotice{RootIsOnline: true})
	}

	sendToParent = true
	return
}

func (thls *businessNode) handle_MsgType_ID_ConnectReq_stepOne_forParent(msgData *txdata.ConnectReq, msgConn *wsnet.WsSocket, rspData *txdata.ConnectRsp) (sendToParent bool) {
	assert4true(len(msgData.Pathway) == 1)
	assert4true(msgData.InfoReq.UserID == thls.ownInfo.BelongID)
	sendToParent = false //来者就是父亲,我不能把父代的请求再发给父代了.

	if isSuccess := thls.parentInfo.setData(msgConn, msgData.InfoReq, false); !isSuccess {
		glog.Errorf("ConnectReq_stepOne_forParent, UserID conflict, msgConn=%p, msgData=%v", msgConn, msgData)
		if true { //UserID冲突,应当立即上报该情况.
			errMsg := fmt.Sprintf("UserID conflict, msgData=%v", msgData)
			thls.reportCommonErrMsg(errMsg)
		}
		rspData.ErrNo = 1
		rspData.ErrMsg = "req.UserID is already online"
		return
	}

	var isAccepted bool
	var isExist bool
	if isAccepted, isExist = thls.cacheSock.deleteData(msgConn); !isExist {
		if true { //先parentInfo.setData,然后发现有异常,需回退操作,以恢复成原来的状态.
			thls.parentInfo.setData(nil, nil, true) //确保它一定能回退(parentInfo.setData)操作.
		}
		rspData.ErrNo = 1
		rspData.ErrMsg = "rsp internal error occurred"
		glog.Errorf("ConnectReq_stepOne_forParent, msgConn not found in cache, msgConn=%p, msgData=%v", msgConn, msgData)
		return
	}

	if isAccepted {
		//有如下通信规则:
		//连接建立后,_connect方要主动发送ConnectReq给accepted方.
		//校验通过后,accepted方要主动发送ConnectReq给_connect方.
		tmpTxData := txdata.ConnectReq{InfoReq: &thls.ownInfo, Pathway: []string{thls.ownInfo.UserID}}
		thls.sendData(msgConn, &tmpTxData)
	}

	return
}

func (thls *businessNode) handle_MsgType_ID_ConnectReq_stepMulti(msgData *txdata.ConnectReq, msgConn *wsnet.WsSocket) (sendToParent bool) {
	assert4true(len(msgData.Pathway) > 1)

	curData := new(connInfoEx)
	curData.conn = msgConn
	curData.Info = *msgData.InfoReq
	curData.Pathway = msgData.Pathway

	//孙子级别的UserID都ConnectReq过来了,那么儿子的ConnectReq肯定已经处理了,或者马上就会发过来,所以此时不用处理cacheSock.

	if isSuccess := thls.cacheUser.insertData(curData); !isSuccess {
		glog.Errorf("ConnectReq_stepMulti, UserID conflict, msgConn=%p, msgData=%v", msgConn, msgData)
		//孙子级别的UserID冲突了,因为不用对ConnectReq做出ConnectRsp的回应,所以迫于无奈,只能和儿子断开连接.
		thls.deleteConnectionFromAll(msgConn, true)
		if true { //UserID冲突,应当立即上报该情况.
			errMsg := fmt.Sprintf("UserID conflict, msgData=%v", msgData)
			thls.reportCommonErrMsg(errMsg)
		}
		sendToParent = false
	} else {
		sendToParent = true
	}

	return
}

func (thls *businessNode) handle_MsgType_ID_ConnectRsp(msgData *txdata.ConnectRsp, msgConn *wsnet.WsSocket) {
	if msgData.ErrNo != 0 {
		glog.Errorln("handle_MsgType_ID_ConnectRsp", msgData, msgConn)
		thls.deleteConnectionFromAll(msgConn, true)
	}
}

func (thls *businessNode) handle_MsgType_ID_OnlineNotice(msgData *txdata.OnlineNotice, msgConn *wsnet.WsSocket) {
	pConn := thls.parentInfo.conn
	if msgConn != pConn {
		glog.Errorf("handle_MsgType_ID_OnlineNotice, OnlineNotice not from parent, msgConn=%p, pConn=%p", msgConn, pConn)
		return
	}
	thls.setRootOnline(msgData.RootIsOnline)
}

func (thls *businessNode) handle_MsgType_ID_CommonReq(msgData *txdata.CommonReq, msgConn *wsnet.WsSocket) {
	if pconn := thls.parentInfo.conn; pconn != nil {
		assert4true((msgConn != pconn) == msgData.TxToRoot) //如果是(儿子)发过来的数据,那么(TxToRoot)必为真.
	}

	if thls.iAmRoot {
		//TODO:留痕.
	}

	if (!msgData.TxToRoot || thls.iAmRoot) && (msgData.RecverID == thls.ownInfo.UserID) {
		thls.handle_MsgType_ID_CommonReq_exec(msgData, msgConn)
		return
	}

	if msgData.UpCache || thls.iAmRoot {
		dataAck := thls.genAck4CommonReq(msgData)
		msgData.SenderID = thls.ownInfo.UserID
		if thls.iAmRoot {
			msgData.TxToRoot = !msgData.TxToRoot
			assert4false(msgData.TxToRoot) //此时要从ROOT往叶子节点发送.
		}
		msgData.UpCache = false
		//缓存,可能是在内存中缓存起来,也可能插入数据库,所以这里需要先修改数据,再进行缓存.
		thls.cacheSync.insertData(msgData.Key, msgData) //缓存.
		//插入成功了,自然成功,插入失败了,说明已经存在了,其实也是接收成功了.
		thls.sendDataEx2(dataAck, msgConn, dataAck.TxToRoot, dataAck.RecverID)
	}
	thls.sendDataEx2(msgData, nil, msgData.TxToRoot, msgData.RecverID)
}

func (thls *businessNode) handle_MsgType_ID_CommonRsp(msgData *txdata.CommonRsp, msgConn *wsnet.WsSocket) {
	if pconn := thls.parentInfo.conn; pconn != nil {
		assert4true((msgConn != pconn) == msgData.TxToRoot) //如果是(儿子)发过来的数据,那么(TxToRoot)必为真.
	}

	if thls.iAmRoot {
		//TODO:留痕.
	}

	//因为东西都需要在ROOT那里留痕,所以,从ROOT发过来的消息,是走完整个流程的,此时才应当被处理.
	if (!msgData.TxToRoot || thls.iAmRoot) && (msgData.RecverID == thls.ownInfo.UserID) {
		thls.cacheRR.operateNode(toUniSym(msgData.Key), msgData, msgData.IsLast)
		//TODO:是否需要发送Rsp的Ack?
		//TODO:如果有续传,就删除请求的续传.
		return
	}

	if msgData.UpCache || thls.iAmRoot {
		dataAck := thls.genAck4CommonRsp(msgData)
		msgData.SenderID = thls.ownInfo.UserID
		if thls.iAmRoot {
			msgData.TxToRoot = !msgData.TxToRoot
			assert4false(msgData.TxToRoot) //此时要从ROOT往叶子节点发送.
		}
		msgData.UpCache = false
		//缓存,可能是在内存中缓存起来,也可能插入数据库,所以这里需要先修改数据,再进行缓存.
		thls.cacheSync.insertData(msgData.Key, msgData) //缓存.
		//插入成功了,自然成功,插入失败了,说明已经存在了,其实也是接收成功了.
		thls.sendDataEx2(dataAck, msgConn, dataAck.TxToRoot, dataAck.RecverID)
	}
	thls.sendDataEx2(msgData, nil, msgData.TxToRoot, msgData.RecverID)
}

func (thls *businessNode) genAck4CommonReq(dataReq *txdata.CommonReq) (dataAck *txdata.MessageAck) {
	//一定要"刚从socket里面接收过来,未经任何修改,然后立即调用该函数"
	//(CommonReq.Key)不会被修改,所以不用clone一个副本.
	return &txdata.MessageAck{Key: dataReq.Key, SenderID: thls.ownInfo.UserID, RecverID: dataReq.SenderID, TxToRoot: !dataReq.TxToRoot}
}

func (thls *businessNode) genAck4CommonRsp(dataRsp *txdata.CommonRsp) (dataAck *txdata.MessageAck) {
	//一定要"刚从socket里面接收过来,未经任何修改,然后立即调用该函数"
	//(CommonRsp.Key)不会被修改,所以不用clone一个副本.
	return &txdata.MessageAck{Key: dataRsp.Key, SenderID: thls.ownInfo.UserID, RecverID: dataRsp.SenderID, TxToRoot: !dataRsp.TxToRoot}
}

func (thls *businessNode) genRsp4CommonReq(dataReq *txdata.CommonReq, seqno int32, pm ProtoMessage, isLast bool) (dataRsp *txdata.CommonRsp) {
	dataRsp = &txdata.CommonRsp{}
	dataRsp.Key = cloneUniKey(dataReq.Key)
	dataRsp.Key.SeqNo = seqno
	dataRsp.SenderID = thls.ownInfo.UserID
	dataRsp.RecverID = dataRsp.Key.UserID
	dataRsp.TxToRoot = !dataReq.TxToRoot //TODO:好像在ROOT的时候有问题.
	dataRsp.UpCache = false
	dataRsp.RspType = CalcMessageType(pm)
	dataRsp.RspData = msg2slice(pm)
	dataRsp.RspTime, _ = ptypes.TimestampProto(time.Now())
	dataRsp.IsLast = isLast
	//
	return
}

func (thls *businessNode) handle_MsgType_ID_CommonReq_exec(reqData *txdata.CommonReq, msgConn *wsnet.WsSocket) {
	stream := newCommonRspWrapper(reqData, thls.cacheSync, thls.letUpCache, msgConn)

	objData, err := slice2msg(reqData.ReqType, reqData.ReqData)
	if err != nil {
		stream.sendData(&txdata.CommonErr{ErrNo: 1, ErrMsg: err.Error()}, true)
		return
	}

	switch reqData.ReqType {
	case txdata.MsgType_ID_QueryRecordReq:
		thls.execute_MsgType_ID_QueryRecordReq(objData.(*txdata.QueryRecordReq), stream)
	case txdata.MsgType_ID_ExecCmdReq:
	case txdata.MsgType_ID_EchoItem:
	default:
		stream.sendData(&txdata.CommonErr{ErrNo: 1, ErrMsg: "unknown_txdata.MsgType"}, true)
	}

	stream.doRemainder()
}

func (thls *businessNode) execute_MsgType_ID_QueryRecordReq(reqData *txdata.QueryRecordReq, stream *CommonRspWrapper) {
}

func (thls *businessNode) refreshSeqNo() {
	//9223372036854775807(int64.max)
	//91231      00000000|
	//yMMddHHmmSS        |
	//60102150405     |  |
	//           86400000
	//可以每隔(1天)重新获取该值.
	//服务端如果遇到冲突的情况,应当立即报警(发邮件等)
	//10年之内将当前表的数据迁移到历史表.
	//                           20060102150405      86400
	str4int64 := time.Now().Format("60102150405") + "00000000"
	val4int64, err := strconv.ParseInt(str4int64, 10, 64)
	assert4true(err == nil)
	atomic.SwapInt64(&thls.ownSeqNo, val4int64)
}

func (thls *businessNode) increaseSeqNo() int64 {
	return atomic.AddInt64(&thls.ownSeqNo, 1)
}

func (thls *businessNode) syncExecuteCommonReqRsp(reqInOut *txdata.CommonReq, d time.Duration) (slcOut []*txdata.CommonRsp) {
	if true { //修复请求结构体的相关字段.
		reqInOut.Key = &txdata.UniKey{UserID: thls.ownInfo.UserID, MsgNo: thls.increaseSeqNo(), SeqNo: 0}
		reqInOut.SenderID = thls.ownInfo.UserID
		//reqInOut.RecverID
		reqInOut.TxToRoot = true
		reqInOut.UpCache = false
		//reqInOut.ReqType
		//reqInOut.ReqData
		reqInOut.ReqTime, _ = ptypes.TimestampProto(time.Now())
	}

	for range FORONCE {
		node := newNodeReqRsp()
		node.key = *toUniSym(reqInOut.Key)
		node.reqData = reqInOut
		if !thls.cacheRR.insertNode(node) {
			panic(node)
		}
		if d <= 0 {
			//TODO:
		} else {
			if err := thls.sendDataEx(thls.parentInfo.conn, reqInOut, true); err != nil {
				rspData := thls.genRsp4CommonReq(reqInOut, 1, &txdata.CommonErr{ErrNo: 1, ErrMsg: err.Error()}, true)
				thls.cacheRR.operateNode(&node.key, rspData, rspData.IsLast)
				slcOut = node.xyz()
				break
			}
			if isTimeout := node.condVar.waitFor(d); isTimeout {
				rspData := thls.genRsp4CommonReq(reqInOut, 1, &txdata.CommonErr{ErrNo: 1, ErrMsg: "timeout"}, true)
				thls.cacheRR.operateNode(&node.key, rspData, rspData.IsLast)
				slcOut = node.xyz()
				break
			}
			slcOut = node.xyz()
		}
	}
	thls.cacheRR.deleteNode(toUniSym(reqInOut.Key))
	return
}
