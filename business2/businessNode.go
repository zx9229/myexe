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
	ownInfo    txdata.ConnectionInfo
	parentInfo safeFatherData
	rootOnline bool
	cacheUser  *safeConnInfoMap
	cacheSock  *safeWsSocketMap
	cachePsh   *safeDataPshCache
	cacheExec  *safeDataPshCache
	a2sPsh     *safeNodeReqRspCache //async2sync
	a2sReqRsp  *safeNodeReqRspCache //async2sync
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
	curData.parentInfo.setData(nil, nil, true)
	//
	curData.cacheSock = newSafeWsSocketMap()
	curData.cacheUser = newSafeConnInfoMap()
	curData.cachePsh = newSafeDataPshCache(curData.ownInfo.UserID == EMPTYSTR)
	curData.cacheExec = newSafeDataPshCache(true)
	curData.a2sPsh = newSafeNodeReqRspCache()
	curData.a2sReqRsp = newSafeNodeReqRspCache()
	//
	if curData.ownInfo.UserID == EMPTYSTR {
		curData.setRootOnline(true)
	}
	//
	curData.refreshSeqNo()
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
	case txdata.MsgType_ID_CommonErrMsg:
		thls.handle_MsgType_ID_CommonErrMsg(txMsgData.(*txdata.CommonErrMsg), msgConn)
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

func (thls *businessNode) handle_MsgType_ID_CommonErrMsg(msgData *txdata.CommonErrMsg, msgConn *wsnet.WsSocket) {
	if msgConn == thls.parentInfo.conn {
		glog.Errorf("handle_MsgType_ID_CommonErrMsg, CommonErrMsg from parent, msgConn=%p", msgConn)
		return
	}
	if thls.ownInfo.UserID != EMPTYSTR {
		thls.sendData(thls.parentInfo.conn, msgData)
	} else {
		glog.Infoln("handle_MsgType_ID_CommonErrMsg", msgData)
	}
}

func (thls *businessNode) refreshSeqNo() {
	//9223372036854775807(int64.max)
	//91231              |
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

func (thls *businessNode) handle_MsgType_ID_DataPsh(msgData *txdata.DataPsh, msgConn *wsnet.WsSocket) {
	//该函数,只允许,收到数据后被回调,不允许某处主动调用它.
	//即:只允许在onMessage里面调用该函数.
	//TODO:检查数据合法性.
	//if !isValidDataPsh(msgData) {
	//	if true {
	//		errMsg := fmt.Sprintf("invalid DataPsh=%v", msgData)
	//		thls.reportCommonErrMsg(errMsg)
	//	}
	//	dataA := DataPsh2DataAck(msgData)
	//	dataA.ErrNo = 1
	//	dataA.ErrMsg = "invalid DataPsh"
	//	thls.sendData(msgConn, dataA)
	//	return
	//}

	if msgData.RecverID == thls.ownInfo.UserID { //本次传输到达(接收者)
		if msgData.RecvUID == thls.ownInfo.UserID { //整个传输到达(最终者)
			//我是(接收者)和(最终者)
			if thls.ownInfo.UserID == EMPTYSTR {
				//TODO:
			} else {
				assert4true(msgData.SenderID == EMPTYSTR) //本次传输的发送者一定是ROOT.
			}
			thls.handle_MsgType_ID_DataPsh_Exec(msgData, msgConn) //TODO:
		} else {
			//我是(接收者)但不是(最终者),那么我是ROOT,此时我应当缓存和转发数据.
			assert4true(thls.ownInfo.UserID == EMPTYSTR) //我一定是ROOT.
			thls.handle_MsgType_ID_DataPsh_CacheAckPsh(msgData, msgConn, true)
		}
	} else {
		//有socket发过来了一个DataPsh,我不是RecverID,说明我在中途,不在终点,也不在起点,所以我一定不是ROOT.
		assert4true(thls.ownInfo.UserID != EMPTYSTR)         //我不是ROOT.
		assert4true(thls.ownInfo.UserID != msgData.SenderID) //我不在起点.
		thls.handle_MsgType_ID_DataPsh_CacheAckPsh(msgData, msgConn, msgData.UpCache)
	}
}

func (thls *businessNode) handle_MsgType_ID_DataPsh_CacheAckPsh(dataP *txdata.DataPsh, msgConn *wsnet.WsSocket, doCache bool) {
	//缓存,发送DataAck,发送DataPsh,
	dataA := DataPsh2DataAck(dataP)
	if true {
		if dataP.RecverID == thls.ownInfo.UserID {
			dataP.RecverID = dataP.RecvUID
		} else {
			assert4true(dataP.RecverID == EMPTYSTR)
			assert4true(dataP.UpCache && doCache)
			dataP.UpCache = false
		}
		dataP.SenderID = thls.ownInfo.UserID
	}
	assert4false(dataP.UpCache) //执行到这里时,若UpCache曾经为true,则此时已经被置为false了.
	if doCache {
		//feedDataPsh成功,缓存里面就有这一条数据了.
		//feedDataPsh失败,DataAck是一条失败的消息,对端还会不断的进行同步.
		thls.cachePsh.feedDataPsh(dataP, dataA)
		thls.sendData(msgConn, dataA)
		if dataA.ErrNo == 0 {
			thls.tempSendDataPsh(dataP)
		} else {
			errMsg := fmt.Sprintf("DataPsh=%v, DataAck=%v", dataP, dataA)
			thls.reportCommonErrMsg(errMsg)
		}
	} else {
		//这里是纯转发,可以认为这里就是一个socket代理,发送成功则不应有任何回应,发送失败则相当于socket断开了,此时要有回应.
		if err := thls.tempSendDataPsh(dataP); err != nil {
			dataA.ErrNo = 1
			dataA.ErrMsg = err.Error()
			thls.sendData(msgConn, dataA)
		}
	}
}

func (thls *businessNode) tempSendDataPsh(dataP *txdata.DataPsh) (err error) {
	//TODO:这是一个临时函数,正式使用的时候,需要做一些assert的检查.
	if dataP.RecverID == EMPTYSTR { //要发往ROOT,所以要发往父亲的方向.
		err = thls.sendDataEx(thls.parentInfo.conn, dataP, true)
	} else { //从ROOT发过来的,所以要发往儿子的方向.
		assert4true(dataP.SenderID == EMPTYSTR)
		if connEx, isExist := thls.cacheUser.queryData(dataP.RecverID); isExist {
			err = thls.sendDataEx(connEx.conn, dataP, false)
		} else {
			err = errors.New("children is offline")
		}
	}
	return
}

func (thls *businessNode) handle_MsgType_ID_DataPsh_Exec(msgData *txdata.DataPsh, msgConn *wsnet.WsSocket) {
	dataA := DataPsh2DataAck(msgData)
	thls.cacheExec.feedDataPsh(msgData, dataA)
	thls.sendData(msgConn, dataA) //回应.
	if dataA.ErrNo != 0 {         //只要没缓存成功,就认为没有接收到,就不处理.
		return
	}
	msgObj, err := slice2msg(msgData.PshType, msgData.PshData)
	if err != nil {
		dstP := thls.createDataPsh4Rsp(msgData, &txdata.CommonErrMsg{Message: err.Error()})
		tmpA := DataPsh2DataAck(dstP)
		thls.cachePsh.feedDataPsh(dstP, tmpA)
		assert4true(tmpA.ErrNo == 0)
		return
	}
	switch msgData.PshType {
	case txdata.MsgType_ID_EchoItem:
		thls.tempExecute_EchoItem(msgData, msgObj.(*txdata.EchoItem))
	default:
	}
}

func (thls *businessNode) tempExecute_EchoItem(dataP *txdata.DataPsh, msgData *txdata.EchoItem) {
	dataA := &txdata.DataAck{}
	msgData.Data = msgData.Data + "(byPeer)"
	curDst := thls.createDataPsh4Rsp(dataP, msgData)
	thls.cachePsh.feedDataPsh(curDst, dataA)
	assert4true(dataA.ErrNo == 0)
	thls.tempSendDataPsh(curDst)
}

func (thls *businessNode) createDataPsh4Rsp(src *txdata.DataPsh, protoMessage ProtoMessage) (dst *txdata.DataPsh) {
	assert4true(thls.ownInfo.UserID == src.RecverID)
	assert4true(thls.ownInfo.UserID == src.RecvUID)
	dst = new(txdata.DataPsh)
	dst.SenderID = thls.ownInfo.UserID
	if thls.ownInfo.UserID == EMPTYSTR {
		dst.RecverID = src.SendUID
	} else {
		dst.RecverID = EMPTYSTR
	}
	dst.SendUID = src.RecvUID
	dst.SendNo = thls.increaseSeqNo()
	dst.RecvUID = src.SendUID
	dst.RecvNo = src.SendNo
	dst.PshType = CalcMessageType(protoMessage)
	dst.PshData = msg2slice(protoMessage)
	//dst.UpCache
	return
}

func (thls *businessNode) xyzPshData(reqInOut *txdata.DataPsh, d time.Duration) (rspOut []*txdata.DataPsh) {
	rspOut = make([]*txdata.DataPsh, 0)
	if true {
		reqInOut.SenderID = thls.ownInfo.UserID
		reqInOut.RecverID = EMPTYSTR
		reqInOut.SendUID = thls.ownInfo.UserID
		reqInOut.SendNo = thls.increaseSeqNo()
		//reqInOut.RecvUID //外部赋值.
		reqInOut.RecvNo = 0 //reset
		//reqInOut.PshType //外部赋值.
		//reqInOut.PshData //外部赋值.
		reqInOut.UpCache = false
	}
	tmpA := DataPsh2DataAck(reqInOut)
	thls.cachePsh.feedDataPsh(reqInOut, tmpA)
	if tmpA.ErrNo != 0 {
		rspOut=append(rspOut,)
	}
}
