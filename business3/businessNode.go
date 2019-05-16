package main

import (
	"encoding/json"
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

//MarshalJSON 为了能通过[json.Marshal(obj)]而编写的函数.
func (thls *businessNode) MarshalJSON() (byteSlice []byte, err error) {
	tmpObj := struct {
		LetUpCache bool
		OwnInfo    *txdata.ConnectionInfo
		IamRoot    bool
		ParentInfo *safeFatherData
		RootOnline bool
		CacheUser  *safeConnInfoMap
		CacheSock  *safeWsSocketMap
		CacheSync  *safeSynchCache
		CacheC2RR  *safeNodeC2ReqRspCache
		CacheC1RR  *safeNodeC1ReqRspCache
		ZZZXML     *safeUniSymCache
		BAQYTDHC   *safeMemoryTmpCache
		AQZXJG     *safeSynchCache
		OwnMsgNo   string //用(int64)时会出现BUG(明明int64的值在慢慢地增加,但是Marshal后的字符串中,我们可以看到,其值一直没有变化,迫于无奈,我使用了string)
		//chanSync   chan string
	}{LetUpCache: thls.letUpCache, OwnInfo: &thls.ownInfo, IamRoot: thls.iAmRoot, ParentInfo: &thls.parentInfo, RootOnline: thls.rootOnline, CacheUser: thls.cacheUser, CacheSock: thls.cacheSock, CacheSync: thls.cacheSync, CacheC2RR: thls.cacheC2RR, CacheC1RR: thls.cacheC1RR, ZZZXML: thls.cacheZZZXML, BAQYTDHC: thls.cBAQYTDHC, AQZXJG: thls.cAQZXJG, OwnMsgNo: strconv.FormatInt(atomic.LoadInt64(&thls.ownMsgNo), 10)}
	byteSlice, err = json.Marshal(&tmpObj)
	return
}

type businessNode struct {
	letUpCache  bool                  //(一经设置,不再修改)让上游缓存数据;TODO:做检查(此时它必须是叶子节点).
	ownInfo     txdata.ConnectionInfo //(一经设置,不再修改)
	iAmRoot     bool                  //(一经设置,不再修改)(I am root node)
	parentInfo  safeFatherData
	rootOnline  bool
	cacheUser   *safeConnInfoMap
	cacheSock   *safeWsSocketMap
	cacheSync   *safeSynchCache        //(db)要绝对的投递过去而缓存+因为UpCache而缓存.
	cacheC2RR   *safeNodeC2ReqRspCache //异步转同步,让请求和响应对应起来.
	cacheC1RR   *safeNodeC1ReqRspCache
	ownMsgNo    int64
	chanSync    chan string         //哪个UserID连通了,就投递过来一个信号.
	cacheZZZXML *safeUniSymCache    //(正在执行命令)的Req消息的UniSym.给续传模式的命令使用,因为只有它们可能续传.
	cBAQYTDHC   *safeMemoryTmpCache //(ROOT模式有用)(不安全用途的缓存)(Common2Req+Common2Rsp+!IsSafe)的数据缓存在这里.
	cAQZXJG     *safeSynchCache     //(db)(ROOT模式有用)安全执行结果.将续传模式的结果写入此表.
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
	curData.letUpCache = false //TODO:
	//
	curData.ownInfo.UserID = cfg.UserID
	curData.ownInfo.BelongID = cfg.BelongID
	curData.ownInfo.Version = "Version20190411"
	curData.ownInfo.LinkMode = txdata.ConnectionInfo_Zero2
	curData.ownInfo.ExePid = int32(os.Getpid())
	curData.ownInfo.ExePath, _ = filepath.Abs(os.Args[0])
	curData.ownInfo.Remark = ""
	//
	curData.iAmRoot = (curData.ownInfo.UserID == EMPTYSTR)
	//
	curData.parentInfo.setData(nil, nil, true)
	//
	if curData.iAmRoot {
		curData.setRootOnline(true)
	}
	//
	curData.cacheUser = newSafeConnInfoMap()
	curData.cacheSock = newSafeWsSocketMap()
	curData.cacheSync = newSafeSynchCache()
	curData.cacheC2RR = newSafeNodeC2ReqRspCache()
	curData.cacheC1RR = newSafeNodeC1ReqRspCache()
	curData.cacheZZZXML = newSafeUniSymCache()
	curData.cBAQYTDHC = newSafeMemoryTmpCache()
	curData.cAQZXJG = newSafeSynchCache()
	//
	curData.refreshSeqNo()
	//
	curData.chanSync = make(chan string, 256)
	go curData.backgroundWork()
	//
	return curData
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

func (thls *businessNode) onConnected(msgConn *wsnet.WsSocket, isAccepted bool) {
	glog.Warningf("[   onConnected] msgConn=%p, isAccepted=%v, LocalAddr=%v, RemoteAddr=%v", msgConn, isAccepted, msgConn.LocalAddr(), msgConn.RemoteAddr())
	if !thls.cacheSock.insertData(msgConn, isAccepted) {
		glog.Fatalf("onConnected, already cached msgConn=%p", msgConn)
	}
	if !isAccepted { //(协议规则)主动connect的socket要主动发送连接请求.
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
		if sonNum != 1 { //一个socket只允许有1个儿子登录.
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
	//thls.deleteConnectionFromAll(msgConn, false)
	thls.cacheSock.deleteData(msgConn)
}

//func (thls *businessNode) deleteConnectionFromAll(conn *wsnet.WsSocket, closeIt bool) {
//	if closeIt {
//		conn.Close()
//	}
//	if thls.parentInfo.conn == conn {
//		//不需要在这里处理它,因为主动断开连接,也会触发onDisconnected回调,回调里面已经有这个逻辑了.
//		//thls.parentInfo.setData(nil, nil, true)
//		//if thls.rootOnline {
//		//	thls.setRootOnline(false)
//		//}
//	}
//	thls.cacheSock.deleteData(conn)
//	thls.cacheUser.deleteDataByConn(conn)
//}

func (thls *businessNode) onMessage(msgConn *wsnet.WsSocket, msgData []byte, msgType int) {
	txMsgType, txMsgData, err := package2msg(msgData)
	if err != nil {
		glog.Errorln("onMessage", txMsgType, txMsgData, err)
		return
	}

	//glog.Infof("onMessage, msgConn=%p, txMsgType=%v, txMsgData=%v", msgConn, txMsgType, txMsgData)

	switch txMsgType {
	case txdata.MsgType_ID_Common2Ack:
		thls.handle_MsgType_ID_Common2Ack(txMsgData.(*txdata.Common2Ack), msgConn)
	case txdata.MsgType_ID_Common2Req:
		thls.handle_MsgType_ID_Common2Req(txMsgData.(*txdata.Common2Req), msgConn)
	case txdata.MsgType_ID_Common2Rsp:
		thls.handle_MsgType_ID_Common2Rsp(txMsgData.(*txdata.Common2Rsp), msgConn)
	case txdata.MsgType_ID_Common1Req:
		thls.handle_MsgType_ID_Common1Req(txMsgData.(*txdata.Common1Req), msgConn)
	case txdata.MsgType_ID_Common1Rsp:
		thls.handle_MsgType_ID_Common1Rsp(txMsgData.(*txdata.Common1Rsp), msgConn)
	case txdata.MsgType_ID_DisconnectedData:
		thls.handle_MsgType_ID_DisconnectedData(txMsgData.(*txdata.DisconnectedData), msgConn)
	case txdata.MsgType_ID_ConnectReq:
		thls.handle_MsgType_ID_ConnectReq(txMsgData.(*txdata.ConnectReq), msgConn)
	case txdata.MsgType_ID_ConnectRsp:
		thls.handle_MsgType_ID_ConnectRsp(txMsgData.(*txdata.ConnectRsp), msgConn)
	case txdata.MsgType_ID_OnlineNotice:
		thls.handle_MsgType_ID_OnlineNotice(txMsgData.(*txdata.OnlineNotice), msgConn)
	case txdata.MsgType_ID_SystemReport:
		thls.handle_MsgType_ID_SystemReport(txMsgData.(*txdata.SystemReport), msgConn)
	default:
		glog.Errorf("onMessage, unknown txdata.MsgType, msgConn=%p, txMsgType=%v, txMsgData=%v", msgConn, txMsgType, txMsgData)
	}
}

func (thls *businessNode) handle_MsgType_ID_Common2Ack(msgData *txdata.Common2Ack, msgConn *wsnet.WsSocket) {
	if msgData.IsLog {
		glog.Infof("handle_MsgType_ID_Common2Ack, msgConn=%p, msgData=%v", msgConn, msgData)
	}

	assert4true(msgData.Key != nil)
	if pconn := thls.parentInfo.conn; pconn != nil {
		assert4true((msgConn != pconn) == msgData.ToRoot) //如果是(儿子)发过来的数据,那么(ToRoot)必为真.
	}

	if msgData.RecverID == thls.ownInfo.UserID {
		thls.cacheSync.deleteData(msgData.Key)
		return
	}
	if thls.iAmRoot {
		msgData.ToRoot = !msgData.ToRoot
		assert4false(msgData.ToRoot)
	}
	thls.sendDataEx2(msgData, nil, msgData.ToRoot, msgData.RecverID)
}

func (thls *businessNode) handle_MsgType_ID_Common2Req(msgData *txdata.Common2Req, msgConn *wsnet.WsSocket) {
	if msgData.IsLog {
		glog.Infof("handle_MsgType_ID_Common2Req, msgConn=%p, msgData=%v", msgConn, msgData)
	}

	assert4true(msgData.Key != nil)
	assert4true(msgData.Key.SeqNo == 0)
	assert4false(msgData.UpCache && !msgData.ToRoot) //从(根节点)发往(叶子节点)只允许一次性到位,不允许中间再有托管环节了.
	if pconn := thls.parentInfo.conn; pconn != nil {
		assert4true((msgConn != pconn) == msgData.ToRoot) //如果是(儿子)发过来的数据,那么(ToRoot)必为真.
	}

	if thls.iAmRoot {
		//TODO:ROOT节点应当有去重的功能,已经存在的,应当直接回应ACK,然后直接丢弃.
		if msgData.IsSafe { //TODO:留痕.
			if isExist, isInsert := thls.cAQZXJG.insertData(msgData.Key, msgData.ToRoot, msgData.RecverID, msgData); !isExist && !isInsert {
				//TODO:报警.
			}
		} else {
			thls.cBAQYTDHC.insertReqData(msgData.Key, msgData)
		}
	}

	//请求消息一定要经过ROOT节点之后,才可以被处理,因为该消息要在ROOT留痕.
	if (!msgData.ToRoot || thls.iAmRoot) && (msgData.RecverID == thls.ownInfo.UserID) {
		//能插入成功,表示尚未执行过此命令;已经存在了,表示已经执行过该命令了.
		if msgData.IsSafe {
			//TODO:要执行命令了,结果崩溃了,然后消息丢失了,我也没办法,我不准备"程序重启之后继续执行该命令",崩了就算了.
			dataAck := thls.genAck4Common2Req(msgData)
			thls.sendDataEx2(dataAck, msgConn, dataAck.ToRoot, dataAck.RecverID)
			//ROOT发出去请求消息,NODE接收后,在即将发送ACK的时候,ROOT与之断开,ACK丢失,
			//NODE开始执行请求命令,请求命令非常耗时(约耗时1分钟),期间没有响应消息发出.
			//数秒之后,ROOT与NODE重连成功,ROOT再次发送请求消息,NODE就会再次受到该消息.
			//此时,应当:NODE开始执行请求的时候,需要在内存中有一个"正在执行命令的map",命令执行结束后,将其删除.
			//RSP发送到ROOT后,应当根据RSP的主键清理掉ROOT的待同步表.
			//TODO:即使如此,还是有可能重复执行,如果真的出现了这种理论上的情况,那么就:
			//对于(发送叶子节点的Req消息)查询沿途的所有节点的(待同步表)如果有结果,那么就丢弃,也不发送ACK,等待RSP清理掉ROOT的缓存.
			if !thls.cacheZZZXML.insertData(msgData.Key) {
				return //命令正在执行中,则认为重复收到该请求,则丢弃该请求.
			}
			if thls.cacheSync.queryCount(msgData.Key.UserID, msgData.Key.MsgNo) != 0 {
				thls.cacheZZZXML.deleteData(msgData.Key)
				return //从待同步表中能查到对应的响应,则认为该请求已经执行过了,这次为重复获取,则丢弃该请求.
			}
		}
		thls.handle_MsgType_ID_Common2Req_exec(msgData, msgConn)
		return
	}

	if msgData.IsSafe {
		if msgData.UpCache || thls.iAmRoot {
			dataAck := thls.genAck4Common2Req(msgData)
			msgData.SenderID = thls.ownInfo.UserID //缓存,可能是在内存中缓存起来,也可能插入数据库,所以这里需要先修改数据,再进行缓存.
			if thls.iAmRoot {
				msgData.ToRoot = !msgData.ToRoot
				assert4false(msgData.ToRoot) //此时要从ROOT往叶子节点发送.
			}
			msgData.UpCache = false

			if isExist, isInsert := thls.cacheSync.insertData(msgData.Key, msgData.ToRoot, msgData.RecverID, msgData); isExist {
				thls.sendDataEx2(dataAck, msgConn, dataAck.ToRoot, dataAck.RecverID) //已经存在了,就发送ACK让对方别再续传了,已经在待同步表里了,它自会同步,也不用再发送了.
				return
			} else if isInsert { //不存在,又插入失败,估计硬盘满了,赶紧报警吧;(如果插入成功了,肯定要正常往下走,然后发送出去).
				thls.sendDataEx2(dataAck, msgConn, dataAck.ToRoot, dataAck.RecverID)
			} else {
				//TODO:报警.
				return
			}
		}
	} else {
		assert4false(msgData.UpCache)
		if thls.iAmRoot {
			msgData.ToRoot = !msgData.ToRoot
			assert4false(msgData.ToRoot) //此时要从ROOT往叶子节点发送.
		}
	}

	err := thls.sendDataEx2(msgData, nil, msgData.ToRoot, msgData.RecverID)

	if (err != nil) && !msgData.IsSafe && !msgData.IsPush {
		tmpTxRspData := thls.genRsp4Common2Req(msgData, 1, &txdata.CommonErr{ErrNo: 1, ErrMsg: err.Error()}, true)
		if thls.iAmRoot {
			//如果我是ROOT,那么原始的Req消息肯定是ToRoot,在ROOT转发失败,要回复的Rsp消息肯定要发往叶子节点.
			tmpTxRspData.ToRoot = !tmpTxRspData.ToRoot
			assert4false(tmpTxRspData.ToRoot)
		}
		thls.sendData(msgConn, tmpTxRspData)
	}
}

func (thls *businessNode) handle_MsgType_ID_Common2Rsp(msgData *txdata.Common2Rsp, msgConn *wsnet.WsSocket) {
	if msgData.IsLog {
		glog.Infof("handle_MsgType_ID_Common2Rsp, msgConn=%p, msgData=%v", msgConn, msgData)
	}

	assert4true(msgData.Key != nil)
	assert4false(msgData.UpCache && !msgData.ToRoot) //从(根节点)发往(叶子节点)只允许一次性到位,不允许中间再有托管环节了.
	assert4true(msgData.Key.SeqNo != 0)
	if pconn := thls.parentInfo.conn; pconn != nil {
		assert4true((msgConn != pconn) == msgData.ToRoot) //如果是(儿子)发过来的数据,那么(ToRoot)必为真.
	}

	if thls.iAmRoot {
		if true {
			//如果经过了ROOT节点,那么就删除ROOT节点的(待同步表)防止ROOT重复发送Req数据.
			kkk := cloneUniKey(msgData.Key)
			kkk.SeqNo = 0
			thls.cacheSync.deleteData(kkk)
		}
		if msgData.IsSafe { //TODO:留痕.
			//留痕表,虽然重要,但是它是留痕作用的,不应用它忽略消息???.
			if isExist, isInsert := thls.cAQZXJG.insertData(msgData.Key, msgData.ToRoot, msgData.RecverID, msgData); !isExist && !isInsert {
				//TODO:报警.
			}
		} else {
			thls.cBAQYTDHC.appendRspData(msgData.Key, msgData)
		}
	}

	//因为东西都需要在ROOT那里留痕,所以,从ROOT发过来的消息,是走完整个流程的,此时才应当被处理.
	if (!msgData.ToRoot || thls.iAmRoot) && (msgData.RecverID == thls.ownInfo.UserID) {
		if true { //如果有续传,就删除请求的续传.删除(待同步表)防止重复发送Req数据.(刚才是删除ROOT节点的,这次是删除原始节点的)
			kkk := cloneUniKey(msgData.Key)
			kkk.SeqNo = 0
			thls.cacheSync.deleteData(kkk)
		}
		thls.cacheC2RR.operateNode(msgData.Key, msgData, msgData.IsLast)
		//TODO:应当添加响应的回调函数,供外部使用.
		if msgData.IsSafe { //应当所有的操作都处理完了,再回应ACK.
			dataAck := thls.genAck4Common2Rsp(msgData)
			thls.sendDataEx2(dataAck, msgConn, dataAck.ToRoot, dataAck.RecverID)
		}
		return
	}

	if msgData.IsSafe {
		if msgData.UpCache || thls.iAmRoot {
			dataAck := thls.genAck4Common2Rsp(msgData)
			msgData.SenderID = thls.ownInfo.UserID //缓存,可能是在内存中缓存起来,也可能插入数据库,所以这里需要先修改数据,再进行缓存.
			if thls.iAmRoot {
				msgData.ToRoot = !msgData.ToRoot
				assert4false(msgData.ToRoot) //此时要从ROOT往叶子节点发送.
			}
			msgData.UpCache = false
			if isExist, isInsert := thls.cacheSync.insertData(msgData.Key, msgData.ToRoot, msgData.RecverID, msgData); isExist {
				thls.sendDataEx2(dataAck, msgConn, dataAck.ToRoot, dataAck.RecverID)
				return
			} else if isInsert {
				thls.sendDataEx2(dataAck, msgConn, dataAck.ToRoot, dataAck.RecverID)
			} else {
				//TODO:报警
				return
			}
		}
	} else {
		assert4false(msgData.UpCache)
		if thls.iAmRoot {
			msgData.ToRoot = !msgData.ToRoot
			assert4false(msgData.ToRoot) //此时要从ROOT往叶子节点发送.
		}
	}

	thls.sendDataEx2(msgData, nil, msgData.ToRoot, msgData.RecverID)
}

func (thls *businessNode) handle_MsgType_ID_Common1Req(msgData *txdata.Common1Req, msgConn *wsnet.WsSocket) {
	if msgData.IsLog {
		glog.Infof("handle_MsgType_ID_Common1Req, msgConn=%p, msgData=%v", msgConn, msgData)
	}

	if pconn := thls.parentInfo.conn; pconn != nil {
		assert4true((msgConn != pconn) == msgData.ToRoot) //如果是(儿子)发过来的数据,那么(ToRoot)必为真.
	}

	if msgData.RecverID == thls.ownInfo.UserID {
		go thls.handle_MsgType_ID_Common1Req_exec(msgData, msgConn)
		return
	}

	originalToRoot := msgData.ToRoot

	var err error
	if connEx, isExist := thls.cacheUser.queryData(msgData.RecverID); isExist {
		msgData.ToRoot = false //我即将发送给自己的子节点,所以不是发往根节点的方向.
		err = thls.sendDataEx(connEx.conn, msgData, false)
	} else {
		if !thls.iAmRoot && msgData.ToRoot {
			err = thls.sendDataEx(thls.parentInfo.conn, msgData, true)
		} else {
			//一旦(!ToRoot)则说明该消息已经在某一个节点找到了RecverID,然后扭头(!ToRoot)发往目标节点,在奔往目标节点的过程中,目标节点离线了.
			//此时目标节点不可达了,因此,直接发送响应即可.
			err = errors.New("node is offline")
		}
	}
	if (err != nil) && (!msgData.IsPush) {
		msgData.ToRoot = originalToRoot
		tmpTxRspData := thls.genRsp4Common1Req(msgData, 1, &txdata.CommonErr{ErrNo: 1, ErrMsg: err.Error()}, true)
		thls.sendData(msgConn, tmpTxRspData) //TODO:打印日志.
	}
}

func (thls *businessNode) handle_MsgType_ID_Common1Rsp(msgData *txdata.Common1Rsp, msgConn *wsnet.WsSocket) {
	if msgData.IsLog {
		glog.Infof("handle_MsgType_ID_Common1Rsp, msgConn=%p, msgData=%v", msgConn, msgData)
	}

	if pconn := thls.parentInfo.conn; pconn != nil {
		assert4true((msgConn != pconn) == msgData.ToRoot)
	}
	if msgData.RecverID == thls.ownInfo.UserID {
		thls.cacheC1RR.operateNode(msgData.MsgNo, msgData, msgData.IsLast)
		//TODO:回调函数.
		return
	}

	var err error
	if connEx, isExist := thls.cacheUser.queryData(msgData.RecverID); isExist {
		msgData.ToRoot = false
		err = thls.sendDataEx(connEx.conn, msgData, false)
	} else {
		if msgData.ToRoot {
			err = thls.sendDataEx(thls.parentInfo.conn, msgData, true)
		} else {
			err = errors.New("node is offline")
		}
	}
	if err != nil {
		glog.Infof("handle_MsgType_ID_Common1Rsp,err=%v,msgConn=%p,msgData=%v", err, msgConn, msgData)
	}
}

func (thls *businessNode) handle_MsgType_ID_DisconnectedData(msgData *txdata.DisconnectedData, msgConn *wsnet.WsSocket) {
	if msgConn == thls.parentInfo.conn { //协议规定,它必须是从儿子的方向发过来的.
		glog.Errorf("recv DisconnectedData from father, msgConn=%p, msgData=%v", msgConn, msgData)
		return
	}
	if !thls.cacheUser.deleteData(msgData.Info.UserID) {
		//极端情况下,可能会出现:儿子和孙子连接好了,儿子和父亲连接起来了,儿子即将和父亲同步连接信息的时候,儿子和孙子连接断开了,
		//儿子向父亲发送DisconnectedData,父亲接收了DisconnectedData,父亲无法清理缓存.
		glog.Errorf("cache cleanup failed, msgConn=%p, msgData=%v", msgConn, msgData)
	}
	thls.sendData(thls.parentInfo.conn, msgData)
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
			thls.reportErrorMsg(errMsg)
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

	thls.feedToChan(curData.Info.UserID)

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
			thls.reportErrorMsg(errMsg)
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

	if true { //和父亲建立连接了,要把自己的缓存发送给父亲,更新父亲的缓存.
		thls.cacheUser.Lock()
		for _, cInfoEx := range thls.cacheUser.M {
			tmpTxData := txdata.ConnectReq{InfoReq: &cInfoEx.Info, Pathway: append(cInfoEx.Pathway, thls.ownInfo.UserID)}
			thls.sendData(msgConn, &tmpTxData)
		}
		thls.cacheUser.Unlock()
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

		//thls.deleteConnectionFromAll(msgConn, true)
		msgConn.Close()

		if true { //UserID冲突,应当立即上报该情况.
			errMsg := fmt.Sprintf("UserID conflict, msgData=%v", msgData)
			thls.reportErrorMsg(errMsg)
		}
		sendToParent = false
	} else {
		sendToParent = true
		thls.feedToChan(curData.Info.UserID)
	}

	return
}

func (thls *businessNode) handle_MsgType_ID_ConnectRsp(msgData *txdata.ConnectRsp, msgConn *wsnet.WsSocket) {
	if msgData.ErrNo != 0 {
		glog.Errorln("handle_MsgType_ID_ConnectRsp", msgData, msgConn)
		//thls.deleteConnectionFromAll(msgConn, true)
		msgConn.Close()
	}
}

func (thls *businessNode) handle_MsgType_ID_OnlineNotice(msgData *txdata.OnlineNotice, msgConn *wsnet.WsSocket) {
	pConn := thls.parentInfo.conn
	if msgConn != pConn {
		glog.Errorf("handle_MsgType_ID_OnlineNotice, OnlineNotice not from parent, msgConn=%p, pConn=%p", msgConn, pConn)
		return
	}
	thls.setRootOnline(msgData.RootIsOnline)

	if thls.rootOnline {
		thls.feedToChan("")
	}
}

func (thls *businessNode) handle_MsgType_ID_SystemReport(msgData *txdata.SystemReport, msgConn *wsnet.WsSocket) {
	if msgConn == thls.parentInfo.conn { //协议规定,它必须是从儿子的方向发过来的.
		glog.Errorf("recv SystemReport from father, msgConn=%p, msgData=%v", msgConn, msgData)
		return
	}
	if thls.iAmRoot {
		glog.Infoln(msgData)
	} else {
		thls.sendData(thls.parentInfo.conn, msgData)
	}
}

func (thls *businessNode) genAck4Common2Req(dataReq *txdata.Common2Req) (dataAck *txdata.Common2Ack) {
	//一定要"刚从socket里面接收过来,未经任何修改,然后立即调用该函数"
	//(Common2Req.Key)不会被修改,所以不用clone一个副本.
	assert4true(dataReq.IsSafe)
	return &txdata.Common2Ack{Key: dataReq.Key, SenderID: thls.ownInfo.UserID, RecverID: dataReq.SenderID, ToRoot: !dataReq.ToRoot, IsLog: dataReq.IsLog}
}

func (thls *businessNode) genAck4Common2Rsp(dataRsp *txdata.Common2Rsp) (dataAck *txdata.Common2Ack) {
	//一定要"刚从socket里面接收过来,未经任何修改,然后立即调用该函数"
	//(Common2Rsp.Key)不会被修改,所以不用clone一个副本.
	assert4true(dataRsp.IsSafe)
	assert4false(dataRsp.IsPush)
	return &txdata.Common2Ack{Key: dataRsp.Key, SenderID: thls.ownInfo.UserID, RecverID: dataRsp.SenderID, ToRoot: !dataRsp.ToRoot, IsLog: dataRsp.IsLog}
}

func (thls *businessNode) genRsp4Common2Req(dataReq *txdata.Common2Req, seqno int32, pm ProtoMessage, isLast bool) (dataRsp *txdata.Common2Rsp) {
	dataRsp = &txdata.Common2Rsp{}
	dataRsp.Key = cloneUniKey(dataReq.Key)
	dataRsp.Key.SeqNo = seqno
	dataRsp.SenderID = thls.ownInfo.UserID
	dataRsp.RecverID = dataRsp.Key.UserID
	dataRsp.ToRoot = !dataReq.ToRoot //TODO:好像在ROOT的时候有问题.
	dataRsp.IsLog = dataReq.IsLog
	dataRsp.IsSafe = dataReq.IsSafe
	dataRsp.IsPush = dataReq.IsPush
	dataRsp.UpCache = thls.letUpCache && dataReq.IsSafe //只有在续传模式下,才允许设置UpCache字段.
	dataRsp.RspType = CalcMessageType(pm)
	dataRsp.RspData = msg2slice(pm)
	dataRsp.RspTime, _ = ptypes.TimestampProto(time.Now())
	dataRsp.IsLast = isLast
	//
	return
}

func (thls *businessNode) genRsp4Common1Req(dataReq *txdata.Common1Req, seqno int32, pm ProtoMessage, isLast bool) (dataRsp *txdata.Common1Rsp) {
	dataRsp = &txdata.Common1Rsp{}
	dataRsp.MsgNo = dataReq.MsgNo
	dataRsp.SeqNo = seqno
	dataRsp.SenderID = thls.ownInfo.UserID
	dataRsp.RecverID = dataReq.SenderID
	dataRsp.ToRoot = !dataReq.ToRoot
	dataRsp.IsLog = dataReq.IsLog
	dataRsp.IsPush = dataReq.IsPush
	dataRsp.RspType = CalcMessageType(pm)
	dataRsp.RspData = msg2slice(pm)
	dataRsp.RspTime, _ = ptypes.TimestampProto(time.Now())
	dataRsp.IsLast = isLast
	//
	return
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
	//                              20060102150405          86400
	str4int64 := time.Now().Format("20060102150405")[3:] + "00000000"
	val4int64, err := strconv.ParseInt(str4int64, 10, 64)
	assert4true(err == nil)
	atomic.SwapInt64(&thls.ownMsgNo, val4int64)
}

func (thls *businessNode) increaseSeqNo() int64 {
	return atomic.AddInt64(&thls.ownMsgNo, 1)
}

func (thls *businessNode) setRootOnline(newValue bool) {
	if oldValue := thls.rootOnline; oldValue == newValue { //我的目标是:消息无冗余无重复,很显然这里消息重复了.
		glog.Errorf("setRootOnline, oldValue=%v, newValue=%v", oldValue, newValue)
	}
	thls.rootOnline = newValue
	if thls.cacheUser != nil {
		thls.cacheUser.sendDataToSon(&txdata.OnlineNotice{RootIsOnline: newValue})
	}
}

func (thls *businessNode) reportErrorMsg(message string) {
	tmpTxData := txdata.SystemReport{UserID: thls.ownInfo.UserID, Pathway: []string{thls.ownInfo.UserID}, Message: message}
	thls.sendData(thls.parentInfo.conn, &tmpTxData)
}

func (thls *businessNode) syncExecuteCommon2ReqRsp(reqInOut *txdata.Common2Req, d time.Duration) (rspSlice []*txdata.Common2Rsp) {
	if true { //修复请求结构体的相关字段.
		reqInOut.Key = &txdata.UniKey{UserID: thls.ownInfo.UserID, MsgNo: thls.increaseSeqNo(), SeqNo: 0}
		reqInOut.SenderID = thls.ownInfo.UserID
		//reqInOut.RecverID
		reqInOut.ToRoot = true
		//reqInOut.IsLog
		//reqInOut.IsSafe
		//reqInOut.IsPush
		reqInOut.UpCache = thls.letUpCache
		//reqInOut.ReqType
		//reqInOut.ReqData
		reqInOut.ReqTime, _ = ptypes.TimestampProto(time.Now())
	}

	if reqInOut.IsSafe {
		if isExist, isInsert := thls.cacheSync.insertData(reqInOut.Key, reqInOut.ToRoot, reqInOut.RecverID, reqInOut); isExist || !isInsert {
			eMsg := fmt.Sprintf("isExist=%v,isInsert=%v", isExist, isInsert)
			rspSlice = []*txdata.Common2Rsp{thls.genRsp4Common2Req(reqInOut, 1, &txdata.CommonErr{ErrNo: 1, ErrMsg: eMsg}, true)} //TODO:
			return
		}
	}

	node := newNodeC2ReqRsp()
	node.key.fromUniKey(reqInOut.Key)
	node.reqData = reqInOut
	if !thls.cacheC2RR.insertNode(node) {
		panic(node) //TODO:
	}

	var rspData *txdata.Common2Rsp
	var err error
	for range FORONCE {
		err = thls.sendDataEx(thls.parentInfo.conn, reqInOut, true)
		//如果推送,等待字段无效,直接返回(等待字段没有意义).
		//如果推送,等待字段有效,直接返回(等待字段没有意义).
		//如果应答,等待字段无效,直接返回.
		//如果应答,等待字段有效,视发送结果而定.
		if reqInOut.IsPush || (d <= 0) {
			if err != nil {
				rspData = thls.genRsp4Common2Req(reqInOut, 1, &txdata.CommonErr{ErrNo: 1, ErrMsg: "(simulate)" + err.Error()}, true)
			} else {
				rspData = thls.genRsp4Common2Req(reqInOut, 1, &txdata.CommonErr{ErrNo: 0, ErrMsg: "(simulate)SUCCESS"}, true)
			}
			break
		}
		//如果应答,等待字段有效(0<d),如果安全执行( IsSafe),本次发送成功了,需要等待.
		//如果应答,等待字段有效(0<d),如果安全执行( IsSafe),本次发送失败了,需要等待(等待期间可能续传成功).
		//如果应答,等待字段有效(0<d),若非安全执行(!IsSafe),本次发送成功了,需要等待.
		//如果应答,等待字段有效(0<d),若非安全执行(!IsSafe),本次发送失败了,直接返回.
		if !reqInOut.IsSafe && (err != nil) {
			rspData = thls.genRsp4Common2Req(reqInOut, 1, &txdata.CommonErr{ErrNo: 1, ErrMsg: "(simulate)" + err.Error()}, true)
			break
		}
		if isTimeout := node.condVar.waitFor(d); isTimeout {
			rspData = thls.genRsp4Common2Req(reqInOut, 1, &txdata.CommonErr{ErrNo: 1, ErrMsg: "(simulate)timeout"}, true)
			break
		}
	}
	if rspData != nil {
		thls.cacheC2RR.operateNode(rspData.Key, rspData, rspData.IsLast)
	}
	rspSlice = node.xyz()
	//thls.cacheC2RR.deleteNode(&node.key)
	return
}

func (thls *businessNode) syncExecuteCommon1ReqRsp(reqInOut *txdata.Common1Req, d time.Duration) (rspSlice []*txdata.Common1Rsp) {
	if true { //修复请求结构体的相关字段.
		reqInOut.MsgNo = thls.increaseSeqNo()
		reqInOut.SeqNo = 0
		reqInOut.SenderID = thls.ownInfo.UserID
		//reqInOut.RecverID
		reqInOut.ToRoot = true
		//reqInOut.IsLog
		//reqInOut.IsPush
		//reqInOut.ReqType
		//reqInOut.ReqData
		reqInOut.ReqTime, _ = ptypes.TimestampProto(time.Now())
	}

	node := newNodeC1ReqRsp()
	node.MsgNo = reqInOut.MsgNo
	node.reqData = reqInOut
	if !thls.cacheC1RR.insertNode(node) {
		panic(node) //TODO:
	}

	var rspData *txdata.Common1Rsp
	var err error
	for range FORONCE {
		err = thls.sendDataEx(thls.parentInfo.conn, reqInOut, true)
		//如果推送,等待字段无效,直接返回(等待字段没有意义).
		//如果推送,等待字段有效,直接返回(等待字段没有意义).
		//如果应答,等待字段无效,直接返回.
		//如果应答,等待字段有效,视发送结果而定.
		if reqInOut.IsPush || (d <= 0) {
			if err != nil {
				rspData = thls.genRsp4Common1Req(reqInOut, 1, &txdata.CommonErr{ErrNo: 1, ErrMsg: "(simulate)" + err.Error()}, true)
			} else {
				rspData = thls.genRsp4Common1Req(reqInOut, 1, &txdata.CommonErr{ErrNo: 0, ErrMsg: "(simulate)SUCCESS"}, true)
			}
			break
		}
		//如果应答,等待字段有效(0<d),若非安全执行(!IsSafe),本次发送成功了,需要等待.
		//如果应答,等待字段有效(0<d),若非安全执行(!IsSafe),本次发送失败了,直接返回.
		if err != nil {
			rspData = thls.genRsp4Common1Req(reqInOut, 1, &txdata.CommonErr{ErrNo: 1, ErrMsg: "(simulate)" + err.Error()}, true)
			break
		}
		if isTimeout := node.condVar.waitFor(d); isTimeout {
			rspData = thls.genRsp4Common1Req(reqInOut, 1, &txdata.CommonErr{ErrNo: 1, ErrMsg: "(simulate)timeout"}, true)
			break
		}
	}
	if rspData != nil {
		thls.cacheC1RR.operateNode(rspData.MsgNo, rspData, rspData.IsLast)
	}
	rspSlice = node.xyz()
	//thls.cacheC2RR.deleteNode(&node.key)
	return
}

func (thls *businessNode) handle_MsgType_ID_Common2Req_exec(reqData *txdata.Common2Req, msgConn *wsnet.WsSocket) {
	stream := newCommon2RspWrapper(reqData, thls.cacheSync, thls.cacheZZZXML, thls.letUpCache, msgConn)

	objData, err := slice2msg(reqData.ReqType, reqData.ReqData)
	if err != nil {
		if !stream.sendData(&txdata.CommonErr{ErrNo: 1, ErrMsg: err.Error()}, true) {
			//TODO:报警.
		}
		return
	}

	switch reqData.ReqType {
	case txdata.MsgType_ID_EchoItem:
		thls.execute_MsgType_ID_EchoItem(objData.(*txdata.EchoItem), stream)
	case txdata.MsgType_ID_QueryRecordReq:
		thls.execute_MsgType_ID_QueryRecordReq(objData.(*txdata.QueryRecordReq), stream)
	case txdata.MsgType_ID_ExecCmdReq:
	case txdata.MsgType_ID_QryConnInfoReq:
		thls.execute_MsgType_ID_QryConnInfoReq(objData.(*txdata.QryConnInfoReq), stream)
	default:
		if !stream.sendData(&txdata.CommonErr{ErrNo: 1, ErrMsg: "unknown_txdata.MsgType"}, true) {
			//TODO:报警.
		}
	}

	stream.doRemainder()
}

func (thls *businessNode) handle_MsgType_ID_Common1Req_exec(reqData *txdata.Common1Req, msgConn *wsnet.WsSocket) {
	stream := newCommon1RspWrapper(reqData, msgConn)

	objData, err := slice2msg(reqData.ReqType, reqData.ReqData)
	if err != nil {
		if !stream.sendData(&txdata.CommonErr{ErrNo: 1, ErrMsg: err.Error()}, true) {
			//TODO:报警.
		}
		return
	}

	switch reqData.ReqType {
	case txdata.MsgType_ID_EchoItem:
		thls.execute_MsgType_ID_EchoItem(objData.(*txdata.EchoItem), stream)
	case txdata.MsgType_ID_QueryRecordReq:
		thls.execute_MsgType_ID_QueryRecordReq(objData.(*txdata.QueryRecordReq), stream)
	case txdata.MsgType_ID_ExecCmdReq:
	case txdata.MsgType_ID_QryConnInfoReq:
		thls.execute_MsgType_ID_QryConnInfoReq(objData.(*txdata.QryConnInfoReq), stream)
	default:
		if !stream.sendData(&txdata.CommonErr{ErrNo: 1, ErrMsg: "unknown_txdata.MsgType"}, true) {
			//TODO:报警.
		}
	}

	stream.doRemainder()
}

func (thls *businessNode) execute_MsgType_ID_EchoItem(reqData *txdata.EchoItem, stream CommonRspWrapper) {
	if reqData.RspCnt <= 0 || reqData.SecGap < 0 {
		stream.sendData(&txdata.CommonErr{ErrNo: 1, ErrMsg: "field (RspCnt and/or SecGap) error"}, true)
	} else {
		for i := int32(1); i <= reqData.RspCnt; i++ {
			isLast := (i == reqData.RspCnt)
			stream.sendData(&txdata.EchoItem{Data: fmt.Sprintf("%v.%v", reqData.Data, i), RspCnt: reqData.RspCnt, SecGap: reqData.SecGap}, isLast)
			if !isLast {
				time.Sleep(time.Second * time.Duration(reqData.SecGap))
			}
		}
	}
}

func (thls *businessNode) execute_MsgType_ID_QueryRecordReq(reqData *txdata.QueryRecordReq, stream CommonRspWrapper) {
}

func (thls *businessNode) execute_MsgType_ID_QryConnInfoReq(reqData *txdata.QryConnInfoReq, stream CommonRspWrapper) {
	data := &txdata.QryConnInfoRsp{UserID: thls.ownInfo.UserID, Cache: thls.cacheUser.tmpF1()}
	data.Cache[thls.ownInfo.UserID] = &txdata.ConnectReq{InfoReq: &thls.ownInfo, Pathway: []string{}}
	for _, v := range data.Cache {
		v.Pathway = append(v.Pathway, thls.ownInfo.UserID)
	}
	stream.sendData(data, true)
}

func (thls *businessNode) backgroundWork() {
	var userID string
	var isOk bool
	var cie *connInfoEx
	var nodeSlice []*node4sync
	var nodeItem *node4sync
	for {
		if userID, isOk = <-thls.chanSync; isOk {
			if userID == EMPTYSTR {
				if nodeSlice = thls.cacheSync.queryDataByToRoot(true); nodeSlice != nil {
					for _, nodeItem = range nodeSlice {
						thls.sendData(thls.parentInfo.conn, nodeItem.data)
					}
				}
			} else {
				if cie, isOk = thls.cacheUser.queryData(userID); isOk {
					if nodeSlice = thls.cacheSync.queryData(false, userID); nodeSlice != nil {
						for _, nodeItem = range nodeSlice {
							thls.sendData(cie.conn, nodeItem.data)
						}
					}
				}
			}
		} else {
			glog.Warningf("the channel may be closed")
			break
		}
	}
}

func (thls *businessNode) feedToChan(userID string) {
	thls.chanSync <- userID
}
