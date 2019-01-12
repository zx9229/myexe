package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
	"unsafe"

	"github.com/go-xorm/builder"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/zx9229/myexe/txdata"
	"github.com/zx9229/myexe/wsnet"
	"github.com/zx9229/zxgo/zxxorm"
)

type businessNode struct {
	cacheSock   *safeWsSocketMap
	cacheNode   *safeConnInfoMap
	cacheReqRsp *safeNodeReqRspCache
	ownInfo     txdata.ConnectionInfo
	parentData  connInfoEx
	xEngine     *xorm.Engine
	workChan    chan int64
}

func newBusinessNode(cfg *configNode) *businessNode {
	if false ||
		atomicKeyIsValid(cfg.UserKey.toTxAtomicKey()) == false ||
		len(cfg.UserKey.NodeName) == 0 ||
		str2ProgramType(cfg.UserKey.ExecType) != txdata.ProgramType_NODE ||
		len(cfg.UserKey.ExecName) != 0 ||
		atomicKeyIsValid(cfg.BelongKey.toTxAtomicKey()) == false ||
		len(cfg.BelongKey.NodeName) == 0 ||
		(str2ProgramType(cfg.BelongKey.ExecType) != txdata.ProgramType_NODE &&
			str2ProgramType(cfg.BelongKey.ExecType) != txdata.ProgramType_SERVER) ||
		len(cfg.BelongKey.ExecName) != 0 ||
		atomicKey2Str(cfg.UserKey.toTxAtomicKey()) == atomicKey2Str(cfg.BelongKey.toTxAtomicKey()) {
		glog.Fatalf("newBusinessNode fail")
	}

	curData := new(businessNode)
	//
	curData.cacheSock = newSafeWsSocketMap()
	curData.cacheNode = newSafeConnInfoMap()
	curData.cacheReqRsp = newSafeNodeReqRspCache()
	//
	curData.ownInfo.UserKey = cfg.UserKey.toTxAtomicKey()
	curData.ownInfo.UserID = atomicKey2Str(curData.ownInfo.UserKey)
	curData.ownInfo.BelongKey = cfg.BelongKey.toTxAtomicKey()
	curData.ownInfo.BelongID = atomicKey2Str(curData.ownInfo.BelongKey)
	curData.ownInfo.Version = "Version20190107"
	curData.ownInfo.LinkMode = txdata.ConnectionInfo_Zero3
	curData.ownInfo.ExePid = int32(os.Getpid())
	curData.ownInfo.ExePath, _ = filepath.Abs(os.Args[0])
	curData.ownInfo.Remark = ""
	//
	curData.parentData = connInfoEx{}
	//
	curData.initEngine(cfg.DataSourceName, cfg.LocationName)
	curData.checkCachedDatabase()
	//
	curData.workChan = make(chan int64, 16)
	go curData.backgroundWork()
	//
	return curData
}

func (thls *businessNode) initEngine(dataSourceName string, locationName string) {
	var err error
	if thls.xEngine, err = xorm.NewEngine("sqlite3", dataSourceName); err != nil {
		glog.Fatalln(err)
	}
	//
	thls.xEngine.SetMapper(core.GonicMapper{}) //支持struct为驼峰式命名,表结构为下划线命名之间的转换,同时对于特定词支持更好.
	//
	if 0 < len(locationName) {
		if location, err := time.LoadLocation(locationName); err != nil {
			glog.Fatalln(err)
		} else {
			thls.xEngine.DatabaseTZ = location
			thls.xEngine.TZLocation = location
		}
	}
	if err = thls.xEngine.CreateTables(&KeyValue{}, &CommonAtosDataNode{}); err != nil { //应该是:只要存在这个tablename,就跳过它.
		glog.Fatalln(err)
	}
	if err = thls.xEngine.Sync2(&KeyValue{}, &CommonAtosDataNode{}); err != nil { //同步数据库结构
		glog.Fatalln(err)
	}
}

func (thls *businessNode) checkCachedDatabase() {
	//程序第一次启动后,可能接收并缓存了数据,然后关闭了程序,然后可能有人修改了缓存数据库里的配置,然后又启动程序,
	//程序启动时,需要检查,缓存数据库里的数据和配置是否冲突,有冲突的话,则拒绝启动.
	var err error
	//(CommonAtosDataNode.UserID)必须等于(txdata.ConnectionInfo.UserID)
	var rowData CommonAtosDataNode
	var affected1, affected2 int64
	if affected1, err = thls.xEngine.Count(&rowData); err != nil {
		glog.Fatalln(err)
	}
	rowData.UserID = thls.ownInfo.UserID
	if affected2, err = thls.xEngine.Count(&rowData); err != nil {
		glog.Fatalln(err)
	}
	if affected1 != affected2 {
		glog.Fatalln(affected1, affected2)
	}
}

func (thls *businessNode) backgroundWork() {
	CommonAtosDataNode2CommonNtosReq := func(src *CommonAtosDataNode) *txdata.CommonNtosReq {
		//负数的RequestID表示背景工作在做事情.
		req := &txdata.CommonNtosReq{RequestID: -1, UserID: src.UserID, SeqNo: src.SeqNo, DataType: src.ReqDataType, Data: src.ReqData, ReqTime: nil}
		req.ReqTime, _ = ptypes.TimestampProto(src.ReqTime)
		return req
	}
	data4qry := &CommonAtosDataNode{}
	fnReqTime := zxxorm.GuessColName(thls.xEngine, data4qry, unsafe.Offsetof(data4qry.ReqTime), true)
	fnmFinish := zxxorm.GuessColName(thls.xEngine, data4qry, unsafe.Offsetof(data4qry.Finish), true)
	go func() {
		for {
			time.Sleep(time.Second * 2)
			thls.workChan <- -1
		}
	}()
	//查询单条数据使用Get方法，在调用Get方法时需要传入一个对应结构体的指针，同时结构体中的非空field自动成为查询的条件和前面的方法条件组合在一起查询.
	var result CommonAtosDataNode
	var has bool
	var err error
	secRecover := float64(30)
	for {
		result = CommonAtosDataNode{}
		data4qry.ReqTime = time.Now().Add(-1 * time.Duration(secRecover) * time.Second) //查询secRecover之前的数据(可能刚执行了一个上报操作,刚插入数据库,所以要有一个缓存时段).
		if has, err = thls.xEngine.Where(builder.Eq{fnmFinish: true}.And(builder.Lt{fnReqTime: data4qry.ReqTime})).Get(&result); err != nil {
			glog.Fatalf("xorm.Get with has=%v, err=%v", has, err)
		} else if has {
			err = thls.sendDataToParent(txdata.MsgType_ID_CommonNtosReq, CommonAtosDataNode2CommonNtosReq(&result))
			//如果没有东西要发送(has == false),也是等待secRecover,然后再查询一下数据库.
			glog.Infof("background report data with SeqNo=%v and err=%v", result.SeqNo, err)
		}
		for looping := true; looping; {
			select {
			case tmpSeqNo, isOk := <-thls.workChan:
				if !isOk {
					glog.Fatalf("recv data from chan with tmpSeqNo=%v, isOk=%v", tmpSeqNo, isOk)
				}
				if tmpSeqNo < 0 { //负数是超时协程发送的数据.
					if secRecover*2 < time.Now().Sub(data4qry.ReqTime).Seconds() { //超时secRecover了,就跳出循环
						looping = false
					}
					continue
				}
				if tmpSeqNo != result.SeqNo {
					glog.Warningf("val=%v, result.SeqNo=%v", tmpSeqNo, result.SeqNo)
				}
				looping = false //上报给SERVER并且收到正确的回复了,就跳出循环.
			default:
			}
		}
	}
}

func (thls *businessNode) onConnected(msgConn *wsnet.WsSocket, isAccepted bool) {
	glog.Warningf("[   onConnected] msgConn=%p, isAccepted=%v, LocalAddr=%v, RemoteAddr=%v", msgConn, isAccepted, msgConn.LocalAddr(), msgConn.RemoteAddr())
	if !thls.cacheSock.insertData(msgConn, isAccepted) {
		glog.Fatalf("already exists, msgConn=%p", msgConn)
	}
	if !isAccepted {
		tmpTxData := txdata.ConnectedData{Info: &thls.ownInfo, Pathway: []string{thls.ownInfo.UserID}}
		msgConn.Send(msg2slice(txdata.MsgType_ID_ConnectedData, &tmpTxData))
	}
}

func (thls *businessNode) onDisconnected(msgConn *wsnet.WsSocket, err error) {
	checkSunWhenDisconnected := func(dataSlice []*connInfoEx) {
		sonNum := 0
		for _, node := range dataSlice {
			if len(node.Pathway) == 1 { //步长为1的是儿子.
				sonNum++
			}
		}
		if sonNum != 1 {
			glog.Fatalf("one msgConn with sonNum=%v", sonNum)
		}
	}
	glog.Warningf("[onDisconnected] msgConn=%p, err=%v", msgConn, err)
	if thls.parentData.conn == msgConn {
		//如果与父亲断开连接,就清理父亲的数据,这样就不用sendDataToParent了.
		glog.Infof("disconnected with father, msgConn=%p", msgConn)
		thls.parentData = connInfoEx{}
	}
	if dataSlice := thls.cacheNode.deleteDataByConn(msgConn); dataSlice != nil { //儿子和我断开连接,我要清理掉儿子和孙子的缓存.
		checkSunWhenDisconnected(dataSlice)
		for _, data := range dataSlice { //发给父亲,让父亲也清理掉对应的缓存.
			tmpTxData := txdata.DisconnectedData{Info: &data.Info}
			thls.sendDataToParent(txdata.MsgType_ID_DisconnectedData, &tmpTxData)
		}
	}
	thls.deleteConnectionFromAll(msgConn, false)
}

func (thls *businessNode) onMessage(msgConn *wsnet.WsSocket, msgData []byte, msgType int) {
	txMsgType, txMsgData, err := slice2msg(msgData)
	if err != nil {
		glog.Errorln(txMsgType, txMsgData, err)
		return
	}
	switch txMsgType {
	case txdata.MsgType_ID_ConnectedData:
		thls.handle_MsgType_ID_ConnectedData(txMsgData.(*txdata.ConnectedData), msgConn)
	case txdata.MsgType_ID_DisconnectedData:
		thls.handle_MsgType_ID_DisconnectedData(txMsgData.(*txdata.DisconnectedData), msgConn)
	case txdata.MsgType_ID_CommonNtosReq:
		thls.handle_MsgType_ID_CommonNtosReq(txMsgData.(*txdata.CommonNtosReq), msgConn)
	case txdata.MsgType_ID_CommonNtosRsp:
		thls.handle_MsgType_ID_CommonNtosRsp(txMsgData.(*txdata.CommonNtosRsp), msgConn)
	case txdata.MsgType_ID_ExecuteCommandReq:
		thls.handle_MsgType_ID_ExecuteCommandReq(txMsgData.(*txdata.ExecuteCommandReq), msgConn)
	case txdata.MsgType_ID_ExecuteCommandRsp:
		thls.handle_MsgType_ID_ExecuteCommandRsp(txMsgData.(*txdata.ExecuteCommandRsp), msgConn)
	case txdata.MsgType_ID_ParentDataReq:
		thls.handle_MsgType_ID_ParentDataReq(txMsgData.(*txdata.ParentDataReq), msgConn)
	default:
		glog.Errorf("unknown txdata.MsgType, msgConn=%p, txMsgType=%v", msgConn, txMsgType)
	}
}

func (thls *businessNode) handle_MsgType_ID_ConnectedData(msgData *txdata.ConnectedData, msgConn *wsnet.WsSocket) {
	sendToParent := false

	if (msgData.Pathway == nil) || (len(msgData.Pathway) == 0) {
		glog.Errorf("empty Pathway, will disconnect with it, msgConn=%p, msgData=%v", msgConn, msgData)
		thls.deleteConnectionFromAll(msgConn, true)
		sendToParent = false
	} else if len(msgData.Pathway) == 1 {
		sendToParent = thls.zxTestDeal4stepOne(msgData, msgConn)
	} else {
		// 1<len(msgData.Pathway)时,这个消息是[孙代]发送给[子代],[子代]转发给[父代](我这一代)的消息.
		sendToParent = thls.zxTestDeal4stepMulti(msgData, msgConn)
	}

	if sendToParent {
		msgData.Pathway = append(msgData.Pathway, thls.ownInfo.UserID)
		thls.sendDataToParent(txdata.MsgType_ID_ConnectedData, msgData)
	}
}

func (thls *businessNode) handle_MsgType_ID_DisconnectedData(msgData *txdata.DisconnectedData, msgConn *wsnet.WsSocket) {
	if thls.parentData.conn == msgConn {
		glog.Errorf("the data must not be from my father, msgConn=%p, msgData=%v", msgConn, msgData)
		return
	}
	if msgData.Info.UserKey.ExecType == txdata.ProgramType_NODE {
		if thls.cacheNode.deleteData(msgData.Info.UserID) == false {
			glog.Fatalf("cache data error, msgConn=%p, msgData=%v", msgConn, msgData)
		}
	} else {
		glog.Errorf("unmanageable, msgConn=%p, msgData=%v", msgConn, msgData)
		return
	}
	thls.sendDataToParent(txdata.MsgType_ID_DisconnectedData, msgData)
}

func (thls *businessNode) handle_MsgType_ID_CommonNtosReq(msgData *txdata.CommonNtosReq, msgConn *wsnet.WsSocket) {
	if thls.parentData.conn == msgConn { //上报请求,肯定要发给父亲,所以,一定不是父亲发过来的.
		glog.Errorf("the data must not be from my father, msgConn=%p, msgData=%v", msgConn, msgData)
		return
	}

	if err := thls.sendDataToParent(txdata.MsgType_ID_CommonNtosReq, msgData); err != nil {
		if _, qau, qas, r := CommonNtosReq_flag(msgData); qau || qas || r {
			if connInfoEx, isExist := thls.cacheNode.queryData(msgData.UserID); isExist {
				rspData := CommonNtosReq2CommonNtosRsp4Err(msgData, -1, err.Error(), false)
				rspData.Pathway = connInfoEx.Pathway
				connInfoEx.conn.Send(msg2slice(txdata.MsgType_ID_CommonNtosRsp, rspData))
			} else {
				//儿子刚发过来数据,我还没处理呢,结果儿子和我断开了,缓存也清理掉了,然后我才开始处理.
				glog.Warningf("user not found, msgConn=%p, msgData=%v", msgConn, msgData)
			}
		}
	}
}

func (thls *businessNode) handle_MsgType_ID_CommonNtosRsp(msgData *txdata.CommonNtosRsp, msgConn *wsnet.WsSocket) {
	if thls.parentData.conn != msgConn {
		glog.Errorf("the data must come from my father, msgConn=%p, msgData=%v", msgConn, msgData)
		return
	}
	var pathwayLen int
	if pathwayLen = len(msgData.Pathway); pathwayLen == 0 {
		glog.Errorf("empty Pathway, msgConn=%p, msgData=%v", msgConn, msgData)
		return
	}
	if thls.ownInfo.UserID != msgData.Pathway[pathwayLen-1] {
		glog.Errorf("illegal Pathway, msgConn=%p, msgData=%v", msgConn, msgData)
		return
	}
	msgData.Pathway = msgData.Pathway[:pathwayLen-1]
	if pathwayLen = len(msgData.Pathway); pathwayLen != 0 {
		nextUID := msgData.Pathway[pathwayLen-1]
		if nextConnInfoEx, isExist := thls.cacheNode.queryData(nextUID); isExist {
			nextConnInfoEx.conn.Send(msg2slice(txdata.MsgType_ID_CommonNtosRsp, msgData))
		} else {
			glog.Warningf("user not found, msgConn=%p, msgData=%v", msgConn, msgData)
		}
	} else {
		isPush, isReqRspUnsafe, isReqRspSafe, isRetry := CommonNtosRsp_flag(msgData)
		assert4true(isPush == false)
		if isReqRspUnsafe || isReqRspSafe { //请求响应相关.
			if node, isExist := thls.cacheReqRsp.deleteElement(msgData.RequestID); isExist {
				node.rspType = txdata.MsgType_ID_CommonNtosRsp
				node.rspData = msgData
				node.condVar.notifyAll()
			} else {
				glog.Infof("data not found in cache, RequestID=%v", msgData.RequestID)
			}
		}
		if isReqRspSafe || isRetry { //数据库相关.
			if msgData.FromServer { //ErrNo为0,表示SERVER处理成功,NODE可以删除自己的缓存了.
				if _, err := thls.xEngine.ID(core.PK{msgData.SeqNo}).Update(&CommonAtosDataNode{Finish: msgData.FromServer, ErrNo: msgData.ErrNo, ErrMsg: msgData.ErrMsg}); err != nil {
					glog.Fatalf("Engine.Update with err=%v", err)
				}
				//可能NODE短时间内发送了两个相同的请求,此时,第一个响应已经删除了数据,第二个响应会执行成功,同时删除零行(猜测//TODO:).
				//所以,可能存在(err == nil && affected == 0)的情况.
			}
		}
		if isRetry {
			thls.workChan <- msgData.SeqNo
		}
	}
}

func (thls *businessNode) handle_MsgType_ID_ExecuteCommandReq(msgData *txdata.ExecuteCommandReq, msgConn *wsnet.WsSocket) {
	if thls.parentData.conn != msgConn {
		glog.Errorf("the data must be from the father, msgConn=%p, msgData=%v", msgConn, msgData)
		return
	}
	if len(msgData.Pathway) == 0 {
		glog.Errorf("empty Pathway, msgConn=%p, msgData=%v", msgConn, msgData)
		return
	}
	if thls.ownInfo.UserID != msgData.Pathway[len(msgData.Pathway)-1] {
		glog.Errorf("illegal Pathway, msgConn=%p, msgData=%v", msgConn, msgData)
		return
	}
	msgData.Pathway = msgData.Pathway[:len(msgData.Pathway)-1]
	if length := len(msgData.Pathway); length != 0 {
		nextUID := msgData.Pathway[length-1]
		if nextConnInfo, isExist := thls.cacheNode.queryData(nextUID); isExist {
			nextConnInfo.conn.Send(msg2slice(txdata.MsgType_ID_ExecuteCommandReq, msgData))
		} else {
			tempTxData := txdata.ExecuteCommandRsp{RequestID: msgData.RequestID, ErrMsg: fmt.Sprintf("next step is unreachable, nextUID=%v", nextUID)}
			thls.sendDataToParent(txdata.MsgType_ID_ExecuteCommandRsp, &tempTxData)
		}
	} else {
		glog.Warningln("ExecuteCommand:", msgData.Command) //TODO:待添加真正的执行代码.

		tempTxData := txdata.ExecuteCommandRsp{RequestID: msgData.RequestID, UserID: thls.ownInfo.UserID, Result: "OK, Now=" + time.Now().Format("2006-01-02_15:04:05")}
		thls.sendDataToParent(txdata.MsgType_ID_ExecuteCommandRsp, &tempTxData)
	}
}

func (thls *businessNode) handle_MsgType_ID_ExecuteCommandRsp(msgData *txdata.ExecuteCommandRsp, msgConn *wsnet.WsSocket) {
	if thls.parentData.conn == msgConn {
		glog.Errorf("the data must not be from my father, msgConn=%p, msgData=%v", msgConn, msgData)
		return
	}
	thls.sendDataToParent(txdata.MsgType_ID_ExecuteCommandRsp, msgData)
}

func (thls *businessNode) handle_MsgType_ID_ParentDataReq(msgData *txdata.ParentDataReq, msgConn *wsnet.WsSocket) {
	rspData := &txdata.ParentDataRsp{}
	rspData.Data = make([]*txdata.ConnectedData, 0)
	thls.cacheNode.Lock()
	for _, node := range thls.cacheNode.M {
		rspData.Data = append(rspData.Data, &txdata.ConnectedData{Info: &node.Info, Pathway: node.Pathway})
	}
	thls.cacheNode.Unlock()
	msgConn.Send(msg2slice(txdata.MsgType_ID_ParentDataRsp, rspData))
}

func (thls *businessNode) sendDataToParent(msgType txdata.MsgType, msgData proto.Message) error {
	conn := thls.parentData.conn
	if conn == nil {
		return errors.New("parent is offline")
	}
	return conn.Send(msg2slice(msgType, msgData))
}

//func (thls *businessNode) isParentConnection(data *txdata.ConnectedData) bool {
//	var isParent bool
//	for range "1" {
//		assert4true(thls.ownInfo.UserKey.ExecType == txdata.ProgramType_NODE)
//		if data.Info.UserID != thls.ownInfo.BelongID {
//			break
//		}
//		if data.Info.UserKey.ExecType != txdata.ProgramType_NODE && data.Info.UserKey.ExecType != txdata.ProgramType_SERVER {
//			break
//		}
//		isParent = true
//	}
//	return isParent
//}

func (thls *businessNode) zxTestDeal4son(msgData *txdata.ConnectedData, msgConn *wsnet.WsSocket) (sendToParent bool) {
	assert4true(len(msgData.Pathway) == 1)
	assert4true(msgData.Info.BelongID == thls.ownInfo.UserID)

	var isAccepted bool
	var isExist bool
	if isAccepted, isExist = thls.cacheSock.deleteData(msgConn); !isExist {
		glog.Errorf("msgConn not found in cache, msgConn=%p, msgData=%v", msgConn, msgData)
		thls.deleteConnectionFromAll(msgConn, true)
		sendToParent = false
		return
	}

	curData := new(connInfoEx)
	curData.conn = msgConn
	curData.Info = *msgData.Info
	curData.Pathway = msgData.Pathway

	if isSuccess := thls.cacheNode.insertData(curData); !isSuccess {
		glog.Errorf("agent already exists, msgConn=%p, msgData=%v", msgConn, msgData)
		thls.deleteConnectionFromAll(msgConn, true)
		sendToParent = false
		return
	}

	if isAccepted {
		tmpTxData := txdata.ConnectedData{Info: &thls.ownInfo, Pathway: []string{thls.ownInfo.UserID}}
		msgConn.Send(msg2slice(txdata.MsgType_ID_ConnectedData, &tmpTxData))
	}

	sendToParent = true
	return
}

func (thls *businessNode) zxTestDeal4parent(msgData *txdata.ConnectedData, msgConn *wsnet.WsSocket) (sendToParent bool) {
	assert4true(len(msgData.Pathway) == 1)
	assert4true(msgData.Info.UserID == thls.ownInfo.BelongID)

	var isAccepted bool
	var isExist bool
	if isAccepted, isExist = thls.cacheSock.deleteData(msgConn); !isExist {
		glog.Errorf("msgConn not found in cache, msgConn=%p, msgData=%v", msgConn, msgData)
		thls.deleteConnectionFromAll(msgConn, true)
		return
	}

	if thls.parentData.conn != nil {
		glog.Errorf("father already exists, msgConn=%p, msgData=%v", msgConn, msgData)
		thls.deleteConnectionFromAll(msgConn, true)
		return
	}

	thls.parentData.conn = msgConn
	thls.parentData.Info = *msgData.Info
	if isAccepted {
		tmpTxData := txdata.ConnectedData{Info: &thls.ownInfo, Pathway: []string{thls.ownInfo.UserID}}
		msgConn.Send(msg2slice(txdata.MsgType_ID_ConnectedData, &tmpTxData))
	}

	if true {
		//和父亲建立连接了,要把自己的缓存发送给父亲,更新父亲的缓存.
		thls.cacheNode.Lock()
		for _, node := range thls.cacheNode.M {
			tmpTxData := txdata.ConnectedData{Info: &node.Info, Pathway: append(node.Pathway, thls.ownInfo.UserID)}
			thls.sendDataToParent(txdata.MsgType_ID_ConnectedData, &tmpTxData)
		}
		thls.cacheNode.Unlock()
	}

	return
}

func (thls *businessNode) zxTestDeal4stepOne(msgData *txdata.ConnectedData, msgConn *wsnet.WsSocket) (sendToParent bool) {
	assert4true(len(msgData.Pathway) == 1)

	if msgData.Info.UserID != msgData.Pathway[0] {
		glog.Errorf("illegal Pathway, msgConn=%p, msgData=%v", msgConn, msgData)
		thls.deleteConnectionFromAll(msgConn, true)
		sendToParent = false
		return
	}

	if msgData.Info.UserID == thls.ownInfo.UserID {
		glog.Errorf("maybe i connected myself, msgConn=%p, msgData=%v", msgConn, msgData)
		thls.deleteConnectionFromAll(msgConn, true)
		sendToParent = false
		return
	}
	if (msgData.Info.UserID != thls.ownInfo.BelongID) &&
		(msgData.Info.BelongID != thls.ownInfo.UserID) {
		glog.Errorf("he is not my father, i am not his father, msgConn=%p, msgData=%v", msgConn, msgData)
		thls.deleteConnectionFromAll(msgConn, true)
		sendToParent = false
		return
	}
	if atomicKey2Str(msgData.Info.UserKey) != msgData.Info.UserID {
		glog.Errorf("msgData.Info.UserID error, msgConn=%p, msgData=%v", msgConn, msgData)
		thls.deleteConnectionFromAll(msgConn, true)
		sendToParent = false
		return
	}
	if atomicKey2Str(msgData.Info.BelongKey) != msgData.Info.BelongID {
		glog.Errorf("BelongID error, msgConn=%p, msgData=%v", msgConn, msgData)
		thls.deleteConnectionFromAll(msgConn, true)
		sendToParent = false
		return
	}
	if msgData.Info.UserKey.ExecType == txdata.ProgramType_POINT {
		if msgData.Info.BelongID == thls.ownInfo.UserID {
			sendToParent = thls.zxTestDeal4son(msgData, msgConn)
			return
		} else {
			glog.Errorf("POINT and i'm not it's father, msgConn=%p, msgData=%v", msgConn, msgData)
			thls.deleteConnectionFromAll(msgConn, true)
			sendToParent = false
			return
		}
	} else if msgData.Info.UserKey.ExecType == txdata.ProgramType_SERVER {
		if msgData.Info.UserID == thls.ownInfo.BelongID {
			sendToParent = thls.zxTestDeal4parent(msgData, msgConn)
			return
		} else {
			glog.Errorf("SERVER and it's not my father, msgConn=%p, msgData=%v", msgConn, msgData)
			thls.deleteConnectionFromAll(msgConn, true)
			sendToParent = false
			return
		}
	} else if msgData.Info.UserKey.ExecType == txdata.ProgramType_NODE {
		if msgData.Info.BelongID == thls.ownInfo.UserID {
			sendToParent = thls.zxTestDeal4son(msgData, msgConn)
			return
		} else if msgData.Info.UserID == thls.ownInfo.BelongID {
			sendToParent = thls.zxTestDeal4parent(msgData, msgConn)
			return
		} else {
			assert4true(false)
			sendToParent = false
			return
		}
	} else {
		glog.Errorf("disable to connect, msgConn=%p, msgData=%v", msgConn, msgData)
		thls.deleteConnectionFromAll(msgConn, true)
		sendToParent = false
		return
	}
}

func (thls *businessNode) zxTestDeal4stepMulti(msgData *txdata.ConnectedData, msgConn *wsnet.WsSocket) (sendToParent bool) {
	assert4true(len(msgData.Pathway) > 1)

	sendToParent = false

	curData := new(connInfoEx)
	curData.conn = msgConn
	curData.Info = *msgData.Info
	curData.Pathway = msgData.Pathway

	if isSuccess := thls.cacheNode.insertData(curData); !isSuccess {
		glog.Errorf("agent already exists, msgConn=%p, msgData=%v", msgConn, msgData)
		thls.deleteConnectionFromAll(msgConn, true)
		sendToParent = false
		return
	}

	sendToParent = true
	return
}

//func (thls *businessNode) doDeal4parent(msgData *txdata.ConnectedData, msgConn *wsnet.WsSocket) (sendToParent bool) {
//	var isAccepted bool
//	var isExist bool
//	if isAccepted, isExist = thls.cacheSock.deleteData(msgConn); !isExist {
//		glog.Errorf("msgConn not found in cache, msgConn=%p, msgData=%v", msgConn, msgData)
//		thls.deleteConnectionFromAll(msgConn, true)
//		return
//	}
//
//	if thls.parentData.conn != nil {
//		glog.Errorf("father already exists, msgConn=%p, msgData=%v", msgConn, msgData)
//		thls.deleteConnectionFromAll(msgConn, true)
//		return
//	}
//
//	thls.parentData.conn = msgConn
//	thls.parentData.Info = *msgData.Info
//	if isAccepted {
//		//if thls.parentData.info.LinkDir != txdata.ConnectionInfo_CONNECT {
//		//	log.Panicln("parent info is abnormal", thls.parentData.info)
//		//}
//		tmpTxData := txdata.ConnectedData{Info: &thls.ownInfo, Pathway: []string{thls.ownInfo.UserID}}
//		msgConn.Send(msg2slice(txdata.MsgType_ID_ConnectedData, &tmpTxData))
//	} else {
//		//if thls.parentData.info.LinkDir != txdata.ConnectionInfo_ACCEPT {
//		//	log.Panicln("parent info is abnormal", thls.parentData.info)
//		//}
//	}
//
//	if true {
//		//和父亲建立连接了,要把自己的缓存发送给父亲,更新父亲的缓存.
//		thls.cacheNode.Lock()
//		for _, node := range thls.cacheNode.M {
//			tmpTxData := txdata.ConnectedData{Info: &node.Info, Pathway: append(node.Pathway, thls.ownInfo.UserID)}
//			thls.sendDataToParent(txdata.MsgType_ID_ConnectedData, &tmpTxData)
//		}
//		thls.cacheNode.Unlock()
//	}
//
//	return
//}

//func (thls *businessNode) doDeal4agent(msgData *txdata.ConnectedData, msgConn *wsnet.WsSocket) (sendToParent bool) {
//	var isAccepted bool
//	isSon := len(msgData.Pathway) == 1
//	if isSon {
//		var isExist bool
//		if isAccepted, isExist = thls.cacheSock.deleteData(msgConn); !isExist {
//			glog.Errorf("msgConn not found in cache, msgConn=%p, msgData=%v", msgConn, msgData)
//			thls.deleteConnectionFromAll(msgConn, true)
//			return
//		}
//	}
//
//	curData := new(connInfoEx)
//	curData.conn = msgConn
//	curData.Info = *msgData.Info
//	curData.Pathway = msgData.Pathway
//
//	if isSuccess := thls.cacheNode.insertData(curData); !isSuccess {
//		glog.Errorf("agent already exists, msgConn=%p, msgData=%v", msgConn, msgData)
//		thls.deleteConnectionFromAll(msgConn, true)
//		return
//	}
//
//	if isSon && isAccepted {
//		tmpTxData := txdata.ConnectedData{Info: &thls.ownInfo, Pathway: []string{thls.ownInfo.UserID}}
//		msgConn.Send(msg2slice(txdata.MsgType_ID_ConnectedData, &tmpTxData))
//	}
//
//	sendToParent = true
//
//	return
//}

func (thls *businessNode) deleteConnectionFromAll(conn *wsnet.WsSocket, closeIt bool) {
	if closeIt {
		conn.Close()
	}
	if thls.parentData.conn == conn {
		thls.parentData = connInfoEx{}
	}
	thls.cacheSock.deleteData(conn)
	thls.cacheNode.deleteDataByConn(conn)
}

func (thls *businessNode) commonAtos(reqInOut *txdata.CommonNtosReq, saveDB bool, d time.Duration) (rspOut *txdata.CommonNtosRsp) {
	if true { //修复请求结构体的相关字段.
		reqInOut.RequestID = 0
		reqInOut.UserID = thls.ownInfo.UserID
		reqInOut.SeqNo = 0
		//reqInOut.Endeavour
		//reqInOut.DataType
		//reqInOut.Data
		reqInOut.ReqTime, _ = ptypes.TimestampProto(time.Now())
	}
	for range "1" {
		var err error
		if saveDB { //要缓存到数据库.
			rowData := CommonNtosReq2CommonAtosDataNode(reqInOut)
			var affected int64
			if affected, err = thls.xEngine.InsertOne(rowData); err != nil {
				rspOut = CommonNtosReq2CommonNtosRsp4Err(reqInOut, -1, err.Error(), false)
				break
			}
			if affected != 1 {
				glog.Fatalf("Engine.InsertOne with affected=%v, err=%v", affected, err) //我就是想知道,成功的话,除了1,还有其他值吗.
			}
			reqInOut.SeqNo = rowData.SeqNo //利用xorm的特性.
		}
		//
		if d <= 0 {
			if err = thls.sendDataToParent(txdata.MsgType_ID_CommonNtosReq, reqInOut); err != nil {
				rspOut = CommonNtosReq2CommonNtosRsp4Err(reqInOut, -1, err.Error(), false)
			} else {
				rspOut = CommonNtosReq2CommonNtosRsp4Err(reqInOut, 0, "send success and no wait.", false)
			}
		} else { //根据RequestID可知它是否在等待.
			node := thls.cacheReqRsp.generateElement()
			if true {
				reqInOut.RequestID = node.requestID
				//
				node.reqType = txdata.MsgType_ID_CommonNtosReq
				node.reqData = reqInOut
			}
			//
			if err = thls.sendDataToParent(node.reqType, node.reqData); err != nil {
				rspOut = CommonNtosReq2CommonNtosRsp4Err(reqInOut, -1, err.Error(), false)
				break
			}
			//
			if isTimeout := node.condVar.waitFor(d); isTimeout {
				rspOut = CommonNtosReq2CommonNtosRsp4Err(reqInOut, -1, "timeout", false)
				break
			}
			rspOut = node.rspData.(*txdata.CommonNtosRsp)
			if (rspOut.RequestID != reqInOut.RequestID) || (rspOut.SeqNo != reqInOut.SeqNo) {
				glog.Fatalf("unmanageable, reqInOut=%v, rspOut=%v", reqInOut, rspOut)
			}
		}
	}
	if reqInOut.RequestID != 0 { //为0的话,还没有使用缓存呢,所以无需清理.
		thls.cacheReqRsp.deleteElement(reqInOut.RequestID)
	}
	return
}

func (thls *businessNode) reportData(dataIn *txdata.ReportDataItem, d time.Duration, isEndeavour bool) *CommRspData {
	reqInOut := Message2CommonNtosReq(dataIn, time.Now(), thls.ownInfo.UserID)
	rspOut := thls.commonAtos(reqInOut, isEndeavour, d)
	return CommonNtosReqRsp2CommRspData(reqInOut, rspOut)
}

func (thls *businessNode) sendMail(dataIn *txdata.SendMailItem, d time.Duration, isEndeavour bool) *CommRspData {
	reqInOut := Message2CommonNtosReq(dataIn, time.Now(), thls.ownInfo.UserID)
	rspOut := thls.commonAtos(reqInOut, isEndeavour, d)
	return CommonNtosReqRsp2CommRspData(reqInOut, rspOut)
}
