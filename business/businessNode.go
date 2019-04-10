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
	cacheUser   *safeConnInfoMap
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
	curData.cacheUser = newSafeConnInfoMap()
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
	go curData.backgroundWork2()
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
	//
	beanSlice := make([]interface{}, 0)
	beanSlice = append(beanSlice, &KeyValue{})
	beanSlice = append(beanSlice, &CommonReqDbN{})
	beanSlice = append(beanSlice, &CommonReqDbS{})
	//beanSlice = append(beanSlice, &CommonRspDbN{})
	//
	if err = thls.xEngine.CreateTables(beanSlice...); err != nil { //应该是:只要存在这个tablename,就跳过它.
		glog.Fatalln(err)
	}
	if err = thls.xEngine.Sync2(beanSlice...); err != nil { //同步数据库结构
		glog.Fatalln(err)
	}
}

func (thls *businessNode) checkCachedDatabase() {
	//程序第一次启动后,可能接收并缓存了数据,然后关闭了程序,然后可能有人修改了缓存数据库里的配置,然后又启动程序,
	//程序启动时,需要检查,缓存数据库里的数据和配置是否冲突,有冲突的话,则拒绝启动.
	var err error
	//(CommonNtosReqDbN.UserID)必须等于(txdata.ConnectionInfo.UserID)
	var rowData CommonReqDbN
	var affected1, affected2 int64
	if affected1, err = thls.xEngine.Count(&rowData); err != nil {
		glog.Fatalln(err)
	}
	rowData.SenderID = thls.ownInfo.UserID
	if affected2, err = thls.xEngine.Count(&rowData); err != nil {
		glog.Fatalln(err)
	}
	if affected1 != affected2 {
		glog.Fatalln(affected1, affected2)
	}
}

//func (thls *businessNode) backgroundWork() {
//	CommonAtosDataNode2CommonNtosReq := func(src *CommonAtosDataNode) *txdata.CommonNtosReq {
//		//(RequestID<0)表示背景工作在做事情.
//		req := &txdata.CommonNtosReq{RequestID: -1, UserID: src.UserID, SeqNo: src.SeqNo, ReqType: src.ReqType, ReqData: src.ReqData, ReqTime: nil}
//		req.ReqTime, _ = ptypes.TimestampProto(src.ReqTime)
//		return req
//	}
//	data4qry := &CommonAtosDataNode{}
//	fnReqTime := zxxorm.GuessColName(thls.xEngine, data4qry, unsafe.Offsetof(data4qry.ReqTime), true)
//	fnmFinish := zxxorm.GuessColName(thls.xEngine, data4qry, unsafe.Offsetof(data4qry.Finish), true)
//	go func() {
//		for {
//			time.Sleep(time.Second * 2)
//			thls.workChan <- -1
//		}
//	}()
//	//查询单条数据使用Get方法，在调用Get方法时需要传入一个对应结构体的指针，同时结构体中的非空field自动成为查询的条件和前面的方法条件组合在一起查询.
//	var result CommonAtosDataNode
//	var has bool
//	var err error
//	secRetransmit := float64(30)
//	for {
//		result = CommonAtosDataNode{}
//		data4qry.ReqTime = time.Now().Add(-1 * time.Duration(secRetransmit) * time.Second) //查询secRetransmit之前的数据(可能刚执行了一个上报操作,刚插入数据库,所以要有一个缓存时段).
//		if has, err = thls.xEngine.Where(builder.Eq{fnmFinish: false}.And(builder.Lt{fnReqTime: data4qry.ReqTime})).Get(&result); err != nil {
//			glog.Fatalf("xorm.Get with has=%v, err=%v", has, err)
//		} else if has {
//			err = thls.sendDataToParent(txdata.MsgType_ID_CommonNtosReq, CommonAtosDataNode2CommonNtosReq(&result))
//			//如果没有东西要发送(has == false),也是等待secRecover,然后再查询一下数据库.
//			glog.Infof("background report data with SeqNo=%v and err=%v", result.SeqNo, err)
//		}
//		for looping := true; looping; {
//			select {
//			case tmpSeqNo, isOk := <-thls.workChan:
//				if !isOk {
//					glog.Fatalf("recv data from chan with tmpSeqNo=%v, isOk=%v", tmpSeqNo, isOk)
//				}
//				if tmpSeqNo < 0 { //负数是超时协程发送的数据.
//					if secRetransmit*2 < time.Now().Sub(data4qry.ReqTime).Seconds() { //超时secRetransmit了,就跳出循环
//						looping = false
//					}
//					continue
//				}
//				if tmpSeqNo != result.SeqNo {
//					glog.Warningf("val=%v, result.SeqNo=%v", tmpSeqNo, result.SeqNo)
//				}
//				looping = false //上报给SERVER并且收到正确的回复了,就跳出循环.此时另一个协程已经修改数据库了,无需这边再次修改.
//			default:
//			}
//		}
//	}
//}

func (thls *businessNode) backgroundWork2() {
	CommonNtosReqDbN2CommonNtosReq4Retransmit := func(src *CommonReqDbN) (dst *txdata.CommonReq) {
		//(RequestID<0)表示背景工作在做事情.
		dst = CommonReqDbN2CommonReq(src)
		dst.RequestID = -1
		return
	}

	var qryResult CommonReqDbN
	funCreateTime := zxxorm.GuessColName(thls.xEngine, &qryResult, unsafe.Offsetof(qryResult.CreateTime), true)
	funcNameState := zxxorm.GuessColName(thls.xEngine, &qryResult, unsafe.Offsetof(qryResult.State), true)

	go func() {
		for {
			time.Sleep(time.Second * 2)
			thls.workChan <- -1
		}
	}()

	secRetransmit := float64(30)
	var qryValueState int //默认值是0
	var qryCreateTime time.Time

	var has bool
	var err error
	for {
		qryResult = CommonReqDbN{}
		qryCreateTime = time.Now().Add(-1 * time.Duration(secRetransmit) * time.Second) //查询secRetransmit之前的数据(可能刚执行了一个上报操作,刚插入数据库,所以要有一个缓存时段).

		if has, err = thls.xEngine.Where(builder.Neq{funcNameState: qryValueState}.And(builder.Lt{funCreateTime: qryCreateTime})).Get(&qryResult); err != nil {
			glog.Fatalf("xorm.Get with has=%v, err=%v", has, err)
		} else if has {
			err = thls.sendDataToParent(CommonNtosReqDbN2CommonNtosReq4Retransmit(&qryResult))
			//如果没有东西要发送(has == false),也是等待(secRetransmit)然后再查询一下数据库.
			glog.Infof("background report data with SeqNo=%v and err=%v", qryResult.SeqNo, err)
		}
		for looping := true; looping; {
			select {
			case tmpSeqNo, isOk := <-thls.workChan:
				if !isOk {
					glog.Fatalf("recv data from chan with tmpSeqNo=%v, isOk=%v", tmpSeqNo, isOk)
				}
				if tmpSeqNo < 0 { //负数是超时协程发送的数据.
					if secRetransmit*2 < time.Now().Sub(qryCreateTime).Seconds() { //超时(secRetransmit)就跳出循环.
						looping = false
					}
					continue
				}
				if tmpSeqNo != qryResult.SeqNo {
					glog.Warningf("val=%v, result.SeqNo=%v", tmpSeqNo, qryResult.SeqNo)
				}
				looping = false //上报给SERVER并且收到正确的回复了,就跳出循环.此时另一个协程已经修改数据库了,无需这边再次修改.
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
		msgConn.Send(msg2slice(&tmpTxData))
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
	if dataSlice := thls.cacheUser.deleteDataByConn(msgConn); dataSlice != nil { //儿子和我断开连接,我要清理掉儿子和孙子的缓存.
		checkSunWhenDisconnected(dataSlice)
		for _, data := range dataSlice { //发给父亲,让父亲也清理掉对应的缓存.
			tmpTxData := txdata.DisconnectedData{Info: &data.Info}
			thls.sendDataToParent(&tmpTxData)
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
	case txdata.MsgType_ID_CommonReq:
		thls.handle_MsgType_ID_CommonReq(txMsgData.(*txdata.CommonReq), msgConn)
	case txdata.MsgType_ID_CommonRsp:
		thls.handle_MsgType_ID_CommonRsp(txMsgData.(*txdata.CommonRsp), msgConn)
	/*case txdata.MsgType_ID_CommonStonReq:
		thls.handle_MsgType_ID_CommonStonReq(txMsgData.(*txdata.CommonStonReq), msgConn)
	case txdata.MsgType_ID_CommonStonRsp:
		thls.handle_MsgType_ID_CommonStonRsp(txMsgData.(*txdata.CommonStonRsp), msgConn)*/
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
		thls.sendDataToParent(msgData)
	}
}

func (thls *businessNode) handle_MsgType_ID_DisconnectedData(msgData *txdata.DisconnectedData, msgConn *wsnet.WsSocket) {
	if thls.parentData.conn == msgConn {
		glog.Errorf("the data must not be from my father, msgConn=%p, msgData=%v", msgConn, msgData)
		return
	}
	if msgData.Info.UserKey.ExecType == txdata.ProgramType_NODE {
		if thls.cacheUser.deleteData(msgData.Info.UserID) == false {
			glog.Fatalf("cache data error, msgConn=%p, msgData=%v", msgConn, msgData)
		}
	} else {
		glog.Errorf("unmanageable, msgConn=%p, msgData=%v", msgConn, msgData)
		return
	}
	thls.sendDataToParent(msgData)
}

func (thls *businessNode) handle_MsgType_ID_CommonReq_process(msgData *txdata.CommonReq, msgConn *wsnet.WsSocket) (rspData *txdata.CommonRsp) {
	var errNo int32
	var errMsg string
	switch txdata.MsgType(msgData.ReqType) {
	case txdata.MsgType_ID_EchoItem:
		curData := &txdata.EchoItem{}
		if err := proto.Unmarshal(msgData.ReqData, curData); err != nil {
			glog.Fatalln(msgData)
			assert4true(err == nil)
		}
		rspData = thls.handle_MsgType_ID_CommonReq_process_txdata_EchoItem(msgData, curData)
	default:
		errNo = -1
		errMsg = "unknown data type"
		rspData = CommonReq2CommonRsp4Err(msgData, errNo, errMsg, false, true)
	}
	return rspData
}

func (thls *businessNode) handle_MsgType_ID_CommonReq_process_txdata_EchoItem(msgData *txdata.CommonReq, innerData *txdata.EchoItem) (rspData *txdata.CommonRsp) {
	innerData.Data += "_rsp"
	bytes, err := proto.Marshal(innerData)
	assert4true(err == nil)
	return CommonReq2CommonRsp4Rsp(msgData, 0, "", false, true, txdata.MsgType_ID_EchoItem, bytes)
}

func (thls *businessNode) handle_MsgType_ID_CommonReq(msgData *txdata.CommonReq, msgConn *wsnet.WsSocket) {
	fmt.Println("ZX_REQ:", msgData)
	if msgData.CrossServer != (thls.parentData.conn == msgConn) {
		//crossServer了，那么Req一定是server发过来的，就一定是父亲发过来的。
		//未crossSERVER，那么Req一定是发往server的，，
		glog.Errorf("data transmission direction error, msgConn=%p, msgData=%v", msgConn, msgData)
		return
	}

	if msgData.CrossServer && thls.ownInfo.UserID == msgData.RecverID {
		//到目的地了
		fmt.Println("ZX_REQ_END:", msgData)
		rspData := thls.handle_MsgType_ID_CommonReq_process(msgData, msgConn)
		if _, qau, qas, r := CommonReq_flag(msgData); qau || qas || r {
			fmt.Println("ZX_REQ_RSP:", rspData)
			thls.sendDataToParent(rspData)
		}
		return
	}

	var err error

	if msgData.CrossServer {
		if connInfoEx, isExist := thls.cacheUser.queryData(msgData.RecverID); isExist {
			err = connInfoEx.conn.Send(msg2slice(msgData))
		} else {
			err = errors.New("can not send to recver")
		}
	} else {
		err = thls.sendDataToParent(msgData)
	}

	if err != nil {
		if _, qau, qas, r := CommonReq_flag(msgData); qau || qas || r {
			rspData := CommonReq2CommonRsp4Err(msgData, -1, err.Error(), msgData.CrossServer, false)
			if msgData.CrossServer {
				thls.sendDataToParent(rspData)
			} else {
				if connInfoEx, isExist := thls.cacheUser.queryData(msgData.SenderID); isExist {
					connInfoEx.conn.Send(msg2slice(rspData))
				}
			}
		}
	}
}

func (thls *businessNode) handle_MsgType_ID_CommonRsp(msgData *txdata.CommonRsp, msgConn *wsnet.WsSocket) {
	fmt.Println("ZX_RSP:", msgData)
	if msgData.CrossServer != (thls.parentData.conn == msgConn) {
		//crossServer了，那么Req一定是server发过来的，就一定是父亲发过来的。
		//未crossSERVER，那么Req一定是发往server的，，
		glog.Errorf("data transmission direction error, msgConn=%p, msgData=%v", msgConn, msgData)
		return
	}

	if msgData.CrossServer && thls.ownInfo.UserID == msgData.RecverID {
		//TODO:未完成
		isPush, isReqRspUnsafe, isReqRspSafe, isRetransmit := CommonRsp_flag(msgData)
		assert4true(isPush == false)
		if isReqRspUnsafe || isReqRspSafe { //请求响应相关.
			if node, isExist := thls.cacheReqRsp.deleteElement(msgData.RequestID); isExist {
				node.rspData = msgData
				node.condVar.notifyAll()
			} else {
				glog.Infof("data not found in cache, RequestID=%v", msgData.RequestID)
			}
		}
		if isRetransmit {
			//TODO:
		}
		return
	}

	var err error

	if msgData.CrossServer {
		if connInfoEx, isExist := thls.cacheUser.queryData(msgData.RecverID); isExist {
			err = connInfoEx.conn.Send(msg2slice(msgData))
		} else {
			err = errors.New("can not send to recver")
		}
	} else {
		err = thls.sendDataToParent(msgData)
	}

	if err != nil {
		glog.Warningf("send fail with err=%v", err)
	}
}

/*
func (thls *businessNode) handle_MsgType_ID_CommonStonReq(msgData *txdata.CommonStonReq, msgConn *wsnet.WsSocket) {
	assert4true(thls.parentData.conn == msgConn)
	pathwayLen := len(msgData.Pathway)
	assert4true(0 < pathwayLen)
	assert4true(thls.ownInfo.UserID == msgData.Pathway[pathwayLen-1])

	msgData.Pathway = msgData.Pathway[:pathwayLen-1]

	if pathwayLen = len(msgData.Pathway); pathwayLen != 0 {
		nextUID := msgData.Pathway[pathwayLen-1]
		if nextConnInfoEx, isExist := thls.cacheUser.queryData(nextUID); isExist {
			nextConnInfoEx.conn.Send(msg2slice(msgData))
		} else {
			//TODO:貌似应当回一个响应.
			glog.Warningf("user not found, msgConn=%p, msgData=%v", msgConn, msgData)
		}
	} else {
		var needSave bool
		var needResp bool
		if isPush, isReqRspUnsafe, isReqRspSafe, isRetransmit := CommonStonReq_flag(msgData); true || isPush {
			needSave = isRetransmit || isReqRspSafe
			needResp = isRetransmit || isReqRspSafe || isReqRspUnsafe
		}
		rspData := thls.handle_MsgType_ID_CommonStonReq_inner(msgData, needSave)
		if needResp {
			thls.sendDataToParent(rspData)
		}
	}
}

func (thls *businessNode) handle_MsgType_ID_CommonStonReq_inner(msgData *txdata.CommonStonReq, neesSave bool) (rspData *txdata.CommonStonRsp) {
	for range "1" {
		if neesSave {
			stonReqDb := &CommonStonReqDb{SeqNo: msgData.SeqNo}
			if has, err := thls.xEngine.Get(stonReqDb); err != nil {
				glog.Fatalf("Engine.Get with has=%v, err=%v, stonReqDb=%v", has, err, stonReqDb)
				rspData = CommonStonReq2CommonStonRsp4Err(msgData, -1, err.Error(), true, thls.ownInfo.UserID)
				break
			} else if has {
				rspData = CommonStonReq2CommonStonRsp4Err(msgData, -1, "key already exists", true, thls.ownInfo.UserID)
				break
			}
			CommonStonReq2CommonStonReqDb(msgData, stonReqDb)
			if affected, err := thls.xEngine.InsertOne(stonReqDb); err != nil {
				glog.Fatalf("Engine.InsertOne with affected=%v, err=%v, stonReqDb=%v", affected, err, stonReqDb)
				rspData = CommonStonReq2CommonStonRsp4Err(msgData, -1, err.Error(), true, thls.ownInfo.UserID)
				break
			} else if affected != 1 {
				glog.Fatalf("Engine.InsertOne with affected=%v, err=%v, stonReqDb=%v", affected, err, stonReqDb) //我就是想知道,成功的话,除了1,还有其他值吗.
				assert4true(affected == 1)
			}
		}
		rspData = thls.handle_MsgType_ID_CommonStonReq_process(msgData)
		if true {
			fillCommonStonRspByCommonStonReq(msgData, rspData)
			rspData.RspTime, _ = ptypes.TimestampProto(time.Now())
			rspData.FromTarget = true
		}
		if neesSave {
			stonRspDb := CommonStonRsp2CommonStonRspDb(rspData)
			if affected, err := thls.xEngine.InsertOne(stonRspDb); (err != nil) || (affected != 1) {
				glog.Fatalf("Engine.InsertOne with affected=%v, err=%v, stonRspDb=%v", affected, err, stonRspDb)
				assert4true(err == nil && affected == 1)
			}
		}
	}
	return
}

//handle_MsgType_ID_CommonStonReq_process 只要(ErrNo,ErrMsg,RspType,RspData)这几个字段的值正确就OK了.
func (thls *businessNode) handle_MsgType_ID_CommonStonReq_process(msgData *txdata.CommonStonReq) (rspObj *txdata.CommonStonRsp) {
	var rspType txdata.MsgType
	var rspData []byte
	if true {
		rdt := &txdata.ReportDataItem{Topic: "ston_req_rsp", Data: "test"}
		var err error
		rspData, err = proto.Marshal(rdt)
		assert4true(err == nil)
		rspType = CalcMessageType(rdt)
	}
	rspObj = &txdata.CommonStonRsp{ErrNo: 0, ErrMsg: "", RspType: rspType, RspData: rspData}
	return
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
		if nextConnInfoEx, isExist := thls.cacheUser.queryData(nextUID); isExist {
			nextConnInfoEx.conn.Send(msg2slice(msgData))
		} else {
			glog.Warningf("user not found, msgConn=%p, msgData=%v", msgConn, msgData)
		}
	} else {
		isPush, isReqRspUnsafe, isReqRspSafe, isRetransmit := CommonNtosRsp_flag(msgData)
		assert4true(isPush == false)
		if isReqRspUnsafe || isReqRspSafe { //请求响应相关.
			if node, isExist := thls.cacheReqRsp.deleteElement(msgData.RequestID); isExist {
				node.rspData = msgData
				node.condVar.notifyAll()
			} else {
				glog.Infof("data not found in cache, RequestID=%v", msgData.RequestID)
			}
		}
		if isReqRspSafe || isRetransmit { //数据库相关.
			if msgData.FromServer { //SERVER已处理本条数据,本条数据已结束,不用重传它了.
				ntosRspDb := CommonNtosRsp2CommonNtosRspDb(msgData)
				if affected, err := thls.xEngine.InsertOne(ntosRspDb); (err != nil) || (affected != 1) {
					glog.Fatalf("Engine.InsertOne with affected=%v, err=%v, ntosRspDb=%v", affected, err, ntosRspDb)
					assert4true(err == nil && affected == 1)
				}
				if _, err := thls.xEngine.ID(core.PK{msgData.SeqNo}).Update(&CommonNtosReqDbN{State: 1}); err != nil {
					glog.Fatalf("Engine.Update with err=%v", err)
					assert4true(err == nil)
				}
			}
		}
		if isRetransmit {
			thls.workChan <- msgData.SeqNo
			//给重传线程回一个消息,重传线程就会结束等待,开始处理下一条数据.
		}
	}
}
*/

// func (thls *businessNode) handle_MsgType_ID_ExecuteCommandReq(msgData *txdata.ExecuteCommandReq, msgConn *wsnet.WsSocket) {
// 	if thls.parentData.conn != msgConn {
// 		glog.Errorf("the data must be from the father, msgConn=%p, msgData=%v", msgConn, msgData)
// 		return
// 	}
// 	if len(msgData.Pathway) == 0 {
// 		glog.Errorf("empty Pathway, msgConn=%p, msgData=%v", msgConn, msgData)
// 		return
// 	}
// 	if thls.ownInfo.UserID != msgData.Pathway[len(msgData.Pathway)-1] {
// 		glog.Errorf("illegal Pathway, msgConn=%p, msgData=%v", msgConn, msgData)
// 		return
// 	}
// 	msgData.Pathway = msgData.Pathway[:len(msgData.Pathway)-1]
// 	if length := len(msgData.Pathway); length != 0 {
// 		nextUID := msgData.Pathway[length-1]
// 		if nextConnInfo, isExist := thls.cacheUser.queryData(nextUID); isExist {
// 			nextConnInfo.conn.Send(msg2slice(msgData))
// 		} else {
// 			tempTxData := txdata.ExecuteCommandRsp{RequestID: msgData.RequestID, ErrMsg: fmt.Sprintf("next step is unreachable, nextUID=%v", nextUID)}
// 			thls.sendDataToParent(&tempTxData)
// 		}
// 	} else {
// 		glog.Warningln("ExecuteCommand:", msgData.Command) //TODO:待添加真正的执行代码.

// 		tempTxData := txdata.ExecuteCommandRsp{RequestID: msgData.RequestID, UserID: thls.ownInfo.UserID, Result: "OK, Now=" + time.Now().Format("2006-01-02_15:04:05")}
// 		thls.sendDataToParent(&tempTxData)
// 	}
// }

// func (thls *businessNode) handle_MsgType_ID_ExecuteCommandRsp(msgData *txdata.ExecuteCommandRsp, msgConn *wsnet.WsSocket) {
// 	if thls.parentData.conn == msgConn {
// 		glog.Errorf("the data must not be from my father, msgConn=%p, msgData=%v", msgConn, msgData)
// 		return
// 	}
// 	thls.sendDataToParent(msgData)
// }

func (thls *businessNode) handle_MsgType_ID_ParentDataReq(msgData *txdata.ParentDataReq, msgConn *wsnet.WsSocket) {
	rspData := &txdata.ParentDataRsp{}
	rspData.Data = make([]*txdata.ConnectedData, 0)
	thls.cacheUser.Lock()
	for _, node := range thls.cacheUser.M {
		rspData.Data = append(rspData.Data, &txdata.ConnectedData{Info: &node.Info, Pathway: node.Pathway})
	}
	thls.cacheUser.Unlock()
	msgConn.Send(msg2slice(rspData))
}

func (thls *businessNode) sendDataToParent(msgData ProtoMessage) error {
	conn := thls.parentData.conn
	if conn == nil {
		return errors.New("parent is offline")
	}
	return conn.Send(msg2slice(msgData))
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

	if isSuccess := thls.cacheUser.insertData(curData); !isSuccess {
		glog.Errorf("agent already exists, msgConn=%p, msgData=%v", msgConn, msgData)
		thls.deleteConnectionFromAll(msgConn, true)
		sendToParent = false
		return
	}

	if isAccepted {
		tmpTxData := txdata.ConnectedData{Info: &thls.ownInfo, Pathway: []string{thls.ownInfo.UserID}}
		msgConn.Send(msg2slice(&tmpTxData))
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
		msgConn.Send(msg2slice(&tmpTxData))
	}

	if true {
		//和父亲建立连接了,要把自己的缓存发送给父亲,更新父亲的缓存.
		thls.cacheUser.Lock()
		for _, node := range thls.cacheUser.M {
			tmpTxData := txdata.ConnectedData{Info: &node.Info, Pathway: append(node.Pathway, thls.ownInfo.UserID)}
			thls.sendDataToParent(&tmpTxData)
		}
		thls.cacheUser.Unlock()
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

	if isSuccess := thls.cacheUser.insertData(curData); !isSuccess {
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
	thls.cacheUser.deleteDataByConn(conn)
}

func (thls *businessNode) commonAtos(reqInOut *txdata.CommonReq, saveDB bool, d time.Duration) (rspOut *txdata.CommonRsp) {
	if true { //修复请求结构体的相关字段.
		reqInOut.SenderID = thls.ownInfo.UserID
		//reqInOut.RecverID
		reqInOut.CrossServer = false
		reqInOut.RequestID = 0
		reqInOut.SeqNo = 0
		//reqInOut.ReqType
		//reqInOut.ReqData
		reqInOut.ReqTime, _ = ptypes.TimestampProto(time.Now())
		reqInOut.RefNum = 0
	}
	for range "1" {
		var err error
		var affected int64
		if saveDB { //要缓存到数据库.
			cReqDbN := &CommonReqDbN{}
			CommonReq2CommonReqDbN(reqInOut, cReqDbN)
			if affected, err = thls.xEngine.InsertOne(cReqDbN); err != nil {
				rspOut = CommonReq2CommonRsp4Err(reqInOut, -1, err.Error(), false, false)
				break
			}
			if affected != 1 {
				glog.Fatalf("Engine.InsertOne with affected=%v, err=%v", affected, err) //我就是想知道,成功的话,除了1,还有其他值吗.
			}
			reqInOut.SeqNo = cReqDbN.SeqNo //利用xorm的特性.
		}
		//
		if d <= 0 {
			if err = thls.sendDataToParent(reqInOut); err != nil {
				rspOut = CommonReq2CommonRsp4Err(reqInOut, -1, err.Error(), false, false)
			} else {
				rspOut = CommonReq2CommonRsp4Err(reqInOut, 0, "send success and no wait.", false, false)
			}
		} else { //根据RequestID可知它是否在等待.
			node := thls.cacheReqRsp.generateElement()
			if true {
				reqInOut.RequestID = node.requestID
				//
				node.reqData = reqInOut
			}
			//
			if err = thls.sendDataToParent(node.reqData); err != nil {
				rspOut = CommonReq2CommonRsp4Err(reqInOut, -1, err.Error(), false, false)
				break
			}
			//
			if isTimeout := node.condVar.waitFor(d); isTimeout {
				rspOut = CommonReq2CommonRsp4Err(reqInOut, -1, "timeout", false, false)
				break
			}
			rspOut = node.rspData.(*txdata.CommonRsp)
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

/*
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
*/
