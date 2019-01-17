package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/zx9229/myexe/txdata"
	"github.com/zx9229/myexe/wsnet"
)

type businessServer struct {
	cacheSock   *safeWsSocketMap
	cacheUser   *safeConnInfoMap
	cacheReqRsp *safeNodeReqRspCache
	ownInfo     txdata.ConnectionInfo
	mailCfg     config4mail
	xEngine     *xorm.Engine
}

func newBusinessServer(cfg *configServer) *businessServer {
	if false ||
		atomicKeyIsValid(cfg.UserKey.toTxAtomicKey()) == false ||
		len(cfg.UserKey.NodeName) == 0 ||
		cfg.UserKey.toTxAtomicKey().ExecType != txdata.ProgramType_SERVER ||
		len(cfg.UserKey.ExecName) != 0 {
		glog.Fatalf("newBusinessServer fail")
	}

	curData := new(businessServer)
	//
	curData.cacheSock = newSafeWsSocketMap()
	curData.cacheUser = newSafeConnInfoMap()
	curData.cacheReqRsp = newSafeNodeReqRspCache()
	//
	curData.ownInfo.UserKey = cfg.UserKey.toTxAtomicKey()
	curData.ownInfo.UserID = atomicKey2Str(curData.ownInfo.UserKey)
	curData.ownInfo.BelongKey = &txdata.AtomicKey{}
	curData.ownInfo.BelongID = atomicKey2Str(curData.ownInfo.BelongKey)
	curData.ownInfo.Version = "Version20190107"
	curData.ownInfo.LinkMode = txdata.ConnectionInfo_Zero3
	curData.ownInfo.ExePid = int32(os.Getpid())
	curData.ownInfo.ExePath, _ = filepath.Abs(os.Args[0])
	curData.ownInfo.Remark = ""
	//
	curData.mailCfg = cfg.MailCfg
	//
	curData.initEngine(cfg.DataSourceName, cfg.LocationName)
	//
	return curData
}

func (thls *businessServer) initEngine(dataSourceName string, locationName string) {
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
	if err = thls.xEngine.CreateTables(&KeyValue{}, &CommonAtosDataNode{}, &CommonAtosDataServer{}); err != nil { //应该是:只要存在这个tablename,就跳过它.
		glog.Fatalln(err)
	}
	if err = thls.xEngine.Sync2(&KeyValue{}, &CommonAtosDataNode{}, &CommonAtosDataServer{}); err != nil { //同步数据库结构
		glog.Fatalln(err)
	}
}

func (thls *businessServer) onConnected(msgConn *wsnet.WsSocket, isAccepted bool) {
	glog.Warningf("[   onConnected] msgConn=%p, isAccepted=%v, LocalAddr=%v, RemoteAddr=%v", msgConn, isAccepted, msgConn.LocalAddr(), msgConn.RemoteAddr())
	if thls.cacheSock.insertData(msgConn, isAccepted) == false {
		glog.Fatalf("already exists, msgConn=%v", msgConn)
	}
	if !isAccepted {
		tmpTxData := txdata.ConnectedData{Info: &thls.ownInfo, Pathway: []string{thls.ownInfo.UserID}}
		msgConn.Send(msg2slice(txdata.MsgType_ID_ConnectedData, &tmpTxData))
	}
}

func (thls *businessServer) onDisconnected(msgConn *wsnet.WsSocket, err error) {
	checkSonWhenDisconnected := func(dataSlice []*connInfoEx) {
		sonNum := 0
		for _, node := range dataSlice {
			if len(node.Pathway) == 1 { //步长为1的是儿子.
				sonNum++
			}
			assert4true(len(node.Pathway) != 0)
		}
		if sonNum != 1 {
			glog.Fatalf("one msgConn with sonNum=%v", sonNum)
		}
	}
	glog.Warningf("[onDisconnected] msgConn=%p, err=%v", msgConn, err)
	if dataSlice := thls.cacheUser.deleteDataByConn(msgConn); dataSlice != nil { //儿子和我断开连接,我要清理掉儿子和孙子的缓存.
		checkSonWhenDisconnected(dataSlice)
	}
	thls.deleteConnectionFromAll(msgConn, false)
}

func (thls *businessServer) onMessage(msgConn *wsnet.WsSocket, msgData []byte, msgType int) {
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

func (thls *businessServer) handle_MsgType_ID_ConnectedData(msgData *txdata.ConnectedData, msgConn *wsnet.WsSocket) {
	if (msgData.Pathway == nil) || (len(msgData.Pathway) == 0) {
		glog.Errorf("empty Pathway, will disconnect with it, msgConn=%p, msgData=%v", msgConn, msgData)
		thls.deleteConnectionFromAll(msgConn, true)
	} else if len(msgData.Pathway) == 1 {
		thls.zxTestDeal4stepOne(msgData, msgConn)
	} else {
		// 1<len(msgData.Pathway)时,这个消息是[孙代]发送给[子代],[子代]转发给[父代](我这一代)的消息.
		thls.zxTestDeal4stepMulti(msgData, msgConn)
	}
}

func (thls *businessServer) handle_MsgType_ID_DisconnectedData(msgData *txdata.DisconnectedData, msgConn *wsnet.WsSocket) {
	if msgData.Info.UserKey.ExecType == txdata.ProgramType_NODE {
		if thls.cacheUser.deleteData(msgData.Info.UserID) == false {
			glog.Fatalf("cache data error, msgConn=%p, msgData=%v", msgConn, msgData)
		}
	} else {
		glog.Errorf("unmanageable, msgConn=%p, msgData=%v", msgConn, msgData)
		return
	}
}

func (thls *businessServer) handle_MsgType_ID_CommonNtosReq(msgData *txdata.CommonNtosReq, msgConn *wsnet.WsSocket) {
	var needSave bool
	var needResp bool
	if isPush, isReqRspUnsafe, isReqRspSafe, isRetransmit := CommonNtosReq_flag(msgData); true || isPush {
		needSave = isRetransmit || isReqRspSafe
		needResp = isRetransmit || isReqRspSafe || isReqRspUnsafe
	}
	//
	rspData := thls.handle_MsgType_ID_CommonNtosReq_inner(msgData, needSave, needResp)
	//
	if needResp {
		if connInfoEx, isExist := thls.cacheUser.queryData(msgData.UserID); isExist {
			rspData.Pathway = connInfoEx.Pathway
			connInfoEx.conn.Send(msg2slice(txdata.MsgType_ID_CommonNtosRsp, rspData))
		} else {
			glog.Infof("user not found, msgConn=%p, msgData=%v", msgConn, msgData)
		}
	}
}

func (thls *businessServer) handle_MsgType_ID_CommonNtosReq_inner(msgData *txdata.CommonNtosReq, needSave, needResp bool) (rspData *txdata.CommonNtosRsp) {
	for range "1" {
		var cads *CommonAtosDataServer
		if needSave {
			//以(UserID+SeqNo)唯一定位一条数据.
			cads = &CommonAtosDataServer{SeqNo: msgData.SeqNo, UserID: msgData.UserID}
			if has, err := thls.xEngine.Get(cads); err != nil {
				glog.Fatalf("Engine.Get with has=%v, err=%v, rds=%v", has, err, cads)
				rspData = CommonNtosReq2CommonNtosRsp4Err(msgData, -1, err.Error(), true)
				break
			} else if has {
				rspData = CommonNtosReq2CommonNtosRsp4Err(msgData, -1, "key already exists", true)
				break
			}
			if true {
				//cads.SeqNo
				//cads.UserID
				cads.ReqDataType = msgData.DataType
				cads.ReqData = msgData.Data
				cads.ReqTime, _ = ptypes.Timestamp(msgData.ReqTime)
			}
			if affected, err := thls.xEngine.InsertOne(cads); err != nil {
				glog.Fatalf("Engine.InsertOne with affected=%v, err=%v, cads=%v", affected, err, cads)
				rspData = CommonNtosReq2CommonNtosRsp4Err(msgData, -1, err.Error(), true)
				break
			} else if affected != 1 {
				glog.Fatalf("Engine.InsertOne with affected=%v, err=%v, cads=%v", affected, err, cads) //我就是想知道,成功的话,除了1,还有其他值吗.
				assert4true(affected != 1)
			}
		}
		rspData = thls.handle_MsgType_ID_CommonNtosReq_process(msgData, needResp)
		if needSave {
			if true {
				cads.Finish = true
				cads.ErrNo = rspData.ErrNo
				cads.ErrMsg = rspData.ErrMsg
				cads.RspDataType = rspData.DataType
				cads.RspData = rspData.Data
			}
			if _, err := thls.xEngine.ID(core.PK{msgData.SeqNo}).Update(cads); err != nil {
				glog.Fatalf("Engine.Update with err=%v", err)
				assert4true(err == nil)
			}
		}
	}
	return
}

func (thls *businessServer) handle_MsgType_ID_CommonNtosReq_process(msgData *txdata.CommonNtosReq, needResp bool) (rspData *txdata.CommonNtosRsp) {
	var errNo int32
	var errMsg string
	switch msgData.DataType {
	case "txdata.ReportDataItem":
		curData := &txdata.ReportDataItem{}
		if err := proto.Unmarshal(msgData.Data, curData); err != nil {
			glog.Fatalln(msgData)
			assert4true(err == nil)
		}
		rspData = thls.handle_MsgType_ID_CommonNtosReq_txdata_ReportDataItem(msgData, curData, needResp)
	case "txdata.SendMailItem":
		curData := &txdata.SendMailItem{}
		if err := proto.Unmarshal(msgData.Data, curData); err != nil {
			glog.Fatalln(msgData)
			assert4true(err == nil)
		}
		rspData = thls.handle_MsgType_ID_CommonNtosReq_txdata_SendMailItem(msgData, curData, needResp)
	case "":
		rspData = CommonNtosReq2CommonNtosRsp4Rsp(msgData, true, -1, "emtpy_type", msgData.DataType, msgData.Data)
	default:
		errNo = -1
		errMsg = "unknown data type"
		rspData = CommonNtosReq2CommonNtosRsp4Err(msgData, errNo, errMsg, true)
	}
	return rspData
}

func (thls *businessServer) handle_MsgType_ID_CommonNtosReq_txdata_ReportDataItem(commReq *txdata.CommonNtosReq, item *txdata.ReportDataItem, needRsp bool) (rspOut *txdata.CommonNtosRsp) {
	glog.Infoln(commReq.DataType, item)
	var errNo int32
	var errMsg string
	var rspType string
	var rspData []byte
	if needRsp {
		rspOut = CommonNtosReq2CommonNtosRsp4Rsp(commReq, true, errNo, errMsg, rspType, rspData)
		rspOut.RspTime, _ = ptypes.TimestampProto(time.Now())
	}
	return
}

func (thls *businessServer) handle_MsgType_ID_CommonNtosReq_txdata_SendMailItem(commReq *txdata.CommonNtosReq, item *txdata.SendMailItem, needRsp bool) (rspOut *txdata.CommonNtosRsp) {
	var errNo int32
	var errMsg string
	var rspType string
	var rspData []byte

	if len(item.Username) == 0 {
		item.Username = thls.mailCfg.Username
		item.Password = thls.mailCfg.Password
		item.SmtpAddr = thls.mailCfg.SmtpAddr
	}
	if err := sendMail(item.Username, item.Password, item.SmtpAddr, item.To, item.Subject, item.ContentType, item.Content); err != nil {
		errNo = -1
		errMsg = err.Error()
	}

	if needRsp {
		rspOut = CommonNtosReq2CommonNtosRsp4Rsp(commReq, true, errNo, errMsg, rspType, rspData)
		rspOut.RspTime, _ = ptypes.TimestampProto(time.Now())
	}
	return
}

func (thls *businessServer) handle_MsgType_ID_CommonNtosRsp(msgData *txdata.CommonNtosRsp, msgConn *wsnet.WsSocket) {
	glog.Errorf("the data must not be from my father, msgConn=%p, msgData=%v", msgConn, msgData)
	return
}

func (thls *businessServer) handle_MsgType_ID_ExecuteCommandReq(msgData *txdata.ExecuteCommandReq, msgConn *wsnet.WsSocket) {
	glog.Errorf("the data must not be from my father, msgConn=%p, msgData=%v", msgConn, msgData)
	return
}

func (thls *businessServer) handle_MsgType_ID_ExecuteCommandRsp(msgData *txdata.ExecuteCommandRsp, msgConn *wsnet.WsSocket) {
	if node, isExist := thls.cacheReqRsp.deleteElement(msgData.RequestID); isExist {
		node.rspType = txdata.MsgType_ID_ExecuteCommandRsp
		node.rspData = msgData
		node.condVar.notifyAll()
	} else {
		glog.Infof("data not found in cache, RequestID=%v", msgData.RequestID)
	}
}

func (thls *businessServer) handle_MsgType_ID_ParentDataReq(msgData *txdata.ParentDataReq, msgConn *wsnet.WsSocket) {
	rspData := &txdata.ParentDataRsp{}
	rspData.Data = make([]*txdata.ConnectedData, 0)
	thls.cacheUser.Lock()
	for _, node := range thls.cacheUser.M {
		rspData.Data = append(rspData.Data, &txdata.ConnectedData{Info: &node.Info, Pathway: node.Pathway})
	}
	thls.cacheUser.Unlock()
	msgConn.Send(msg2slice(txdata.MsgType_ID_ParentDataRsp, rspData))
}

//func (thls *businessServer) doDeal4client(msgData *txdata.ConnectedData, msgConn *wsnet.WsSocket) {
//	if msgData.Info.UserKey.ExecType != txdata.ProgramType_CLIENT {
//		panic(msgData)
//	}
//
//	if isAccepted, isExist := thls.cacheSock.deleteData(msgConn); !(isExist && isAccepted) {
//		glog.Errorf("msgConn not found or not isAccepted, msgConn=%p, msgData=%v", msgConn, msgData)
//		thls.deleteConnectionFromAll(msgConn, true)
//		return
//	}
//
//	if msgData.Pathway != nil && 0 < len(msgData.Pathway) {
//		glog.Errorf("msgConn Pathway not empty, msgConn=%p, msgData=%v", msgConn, msgData)
//		thls.deleteConnectionFromAll(msgConn, true)
//		return
//	}
//
//	curData := new(connInfoEx)
//	curData.conn = msgConn
//	curData.Info = *msgData.Info
//	curData.Pathway = msgData.Pathway
//
//	if isSuccess := thls.cacheClient.insertData(curData); !isSuccess {
//		glog.Errorf("client already exists, msgConn=%p, msgData=%v", msgConn, msgData)
//		thls.deleteConnectionFromAll(msgConn, true)
//		return
//	}
//
//	tmpTxData := txdata.ConnectedData{Info: &thls.ownInfo, Pathway: []string{thls.ownInfo.UserID}}
//	msgConn.Send(msg2slice(txdata.MsgType_ID_ConnectedData, &tmpTxData))
//}

func (thls *businessServer) zxTestDeal4son(msgData *txdata.ConnectedData, msgConn *wsnet.WsSocket) {
	assert4true(len(msgData.Pathway) == 1)
	assert4true(msgData.Info.BelongID == thls.ownInfo.UserID)

	var isAccepted bool
	var isExist bool
	if isAccepted, isExist = thls.cacheSock.deleteData(msgConn); !isExist {
		glog.Errorf("msgConn not found in cache, msgConn=%p, msgData=%v", msgConn, msgData)
		thls.deleteConnectionFromAll(msgConn, true)
		return
	}

	curData := new(connInfoEx)
	curData.conn = msgConn
	curData.Info = *msgData.Info
	curData.Pathway = msgData.Pathway

	if isSuccess := thls.cacheUser.insertData(curData); !isSuccess {
		glog.Errorf("agent already exists, msgConn=%p, msgData=%v", msgConn, msgData)
		thls.deleteConnectionFromAll(msgConn, true)
		return
	}

	if isAccepted {
		tmpTxData := txdata.ConnectedData{Info: &thls.ownInfo, Pathway: []string{thls.ownInfo.UserID}}
		msgConn.Send(msg2slice(txdata.MsgType_ID_ConnectedData, &tmpTxData))
	}
}

func (thls *businessServer) zxTestDeal4stepOne(msgData *txdata.ConnectedData, msgConn *wsnet.WsSocket) {
	assert4true(len(msgData.Pathway) == 1)

	if msgData.Info.UserID != msgData.Pathway[0] {
		glog.Errorf("illegal Pathway, msgConn=%p, msgData=%v", msgConn, msgData)
		thls.deleteConnectionFromAll(msgConn, true)
		return
	}
	if msgData.Info.UserID == thls.ownInfo.UserID {
		glog.Errorf("maybe i connected myself, msgConn=%p, msgData=%v", msgConn, msgData)
		thls.deleteConnectionFromAll(msgConn, true)
		return
	}
	if msgData.Info.BelongID != thls.ownInfo.UserID {
		glog.Errorf("i am not his father, msgConn=%p, msgData=%v", msgConn, msgData)
		thls.deleteConnectionFromAll(msgConn, true)
		return
	}
	if atomicKey2Str(msgData.Info.UserKey) != msgData.Info.UserID {
		glog.Errorf("msgData.Info.UserID error, msgConn=%p, msgData=%v", msgConn, msgData)
		thls.deleteConnectionFromAll(msgConn, true)
		return
	}
	if atomicKey2Str(msgData.Info.BelongKey) != msgData.Info.BelongID {
		glog.Errorf("BelongID error, msgConn=%p, msgData=%v", msgConn, msgData)
		thls.deleteConnectionFromAll(msgConn, true)
		return
	}

	if (msgData.Info.UserKey.ExecType != txdata.ProgramType_CLIENT) &&
		(msgData.Info.UserKey.ExecType != txdata.ProgramType_NODE) {
		glog.Errorf("disable to connect, msgConn=%p, msgData=%v", msgConn, msgData)
		thls.deleteConnectionFromAll(msgConn, true)
		return
	}

	thls.zxTestDeal4son(msgData, msgConn)
}

func (thls *businessServer) zxTestDeal4stepMulti(msgData *txdata.ConnectedData, msgConn *wsnet.WsSocket) {
	assert4true(len(msgData.Pathway) > 1)

	if true &&
		(msgData.Info.UserKey.ExecType != txdata.ProgramType_NODE) &&
		(msgData.Info.UserKey.ExecType != txdata.ProgramType_POINT) {
		glog.Errorf("UserKey.ExecType abnormal, msgConn=%p, msgData=%v", msgConn, msgData)
		thls.deleteConnectionFromAll(msgConn, true)
		return
	}

	curData := new(connInfoEx)
	curData.conn = msgConn
	curData.Info = *msgData.Info
	curData.Pathway = msgData.Pathway

	if isSuccess := thls.cacheUser.insertData(curData); !isSuccess {
		glog.Errorf("agent already exists, msgConn=%p, msgData=%v", msgConn, msgData)
		thls.deleteConnectionFromAll(msgConn, true)
		return
	}
}

//func (thls *businessServer) doDeal4node(msgData *txdata.ConnectedData, msgConn *wsnet.WsSocket) {
//	if msgData.Info.UserKey.ExecType != txdata.ProgramType_NODE {
//		panic(msg2slice)
//	}
//
//	if msgData.Pathway == nil || len(msgData.Pathway) == 0 {
//		glog.Errorf("empty Pathway, will disconnect with it, msgConn=%p, msgData=%v", msgConn, msgData)
//		thls.deleteConnectionFromAll(msgConn, true)
//		return
//	}
//
//	if msgData.Info.UserID != msgData.Pathway[0] {
//		glog.Errorf("illegal Pathway, msgConn=%p, msgData=%v", msgConn, msgData)
//		thls.deleteConnectionFromAll(msgConn, true)
//		return
//	}
//
//	//我的儿子连过来了，我要自检，自己是不是它的父亲.
//	if (len(msgData.Pathway) == 1) && (msgData.Info.BelongID != thls.ownInfo.UserID) {
//		glog.Errorf("i am not his father, msgConn=%p, msgData=%v", msgConn, msgData)
//		thls.deleteConnectionFromAll(msgConn, true)
//		return
//	}
//
//	var isAccepted bool
//	isSon := len(msgData.Pathway) == 1
//	if isSon {
//		//步长为0的数据,在主调函数中已经进行了过滤.
//		//步长等于1时,是儿子socket发送的数据.
//		//步长大于1时,是孙子socket发送的数据.
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
//}

func (thls *businessServer) deleteConnectionFromAll(conn *wsnet.WsSocket, closeIt bool) {
	if closeIt {
		conn.Close()
	}
	thls.cacheSock.deleteData(conn)
	thls.cacheUser.deleteDataByConn(conn)
}

func (thls *businessServer) executeCommand(reqInOut *txdata.ExecuteCommandReq, d time.Duration) (rspOut *txdata.ExecuteCommandRsp) {
	//这里选择用(reqInOut.Pathway[0])承载(要控制哪个NODE)信息.
	tmpUID := reqInOut.Pathway[0]
	if true { //修复请求结构体的相关字段.
		reqInOut.RequestID = 0
		reqInOut.Pathway = nil
		//reqIn.Command
	}
	ExecuteCommandReq2ExecuteCommandRsp := func(reqIn *txdata.ExecuteCommandReq, errNo int32, errMsg string, uid string) *txdata.ExecuteCommandRsp {
		return &txdata.ExecuteCommandRsp{RequestID: reqIn.RequestID, UserID: uid, Result: "", ErrNo: errNo, ErrMsg: errMsg}
	}
	for range "1" {
		connInfoEx, isExist := thls.cacheUser.queryData(tmpUID)
		if !isExist {
			rspOut = ExecuteCommandReq2ExecuteCommandRsp(reqInOut, -1, "uid not found in cache", tmpUID)
			break
		}
		reqInOut.Pathway = connInfoEx.Pathway
		//
		node := thls.cacheReqRsp.generateElement()
		if true {
			reqInOut.RequestID = node.requestID
			//
			node.reqType = txdata.MsgType_ID_ExecuteCommandReq
			node.reqData = reqInOut
		}
		//
		if err := connInfoEx.conn.Send(msg2slice(node.reqType, node.reqData)); err != nil {
			rspOut = ExecuteCommandReq2ExecuteCommandRsp(reqInOut, -1, err.Error(), tmpUID)
			break
		}
		//
		if isTimeout := node.condVar.waitFor(d); isTimeout {
			rspOut = ExecuteCommandReq2ExecuteCommandRsp(reqInOut, -1, "timeout", tmpUID)
			break
		}
		rspOut = node.rspData.(*txdata.ExecuteCommandRsp)
		if rspOut.RequestID != reqInOut.RequestID {
			glog.Fatalf("unmanageable, reqInOut=%v, rspOut=%v", reqInOut, rspOut)
		}
	}
	if reqInOut.RequestID != 0 { //为0的话,还没有使用缓存呢,所以无需清理.
		thls.cacheReqRsp.deleteElement(reqInOut.RequestID)
	}
	return
}

func (thls *businessServer) commonAtos(reqInOut *txdata.CommonNtosReq, saveDB bool, d time.Duration) (rspOut *txdata.CommonNtosRsp) {
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
		var affected int64
		if saveDB { //要缓存到数据库.
			rowData := CommonNtosReq2CommonAtosDataNode(reqInOut)
			if affected, err = thls.xEngine.InsertOne(rowData); err != nil {
				rspOut = CommonNtosReq2CommonNtosRsp4Err(reqInOut, -1, err.Error(), false) //尚未进入SERVER处理请求的逻辑就出错了,故置false.
				break
			}
			if affected != 1 {
				glog.Fatalf("Engine.InsertOne with affected=%v, err=%v", affected, err) //我就是想知道,成功的话,除了1,还有其他值吗.
				assert4true(affected == 1)
			}
			reqInOut.SeqNo = rowData.SeqNo //利用xorm的特性.
		}
		rspOut = thls.handle_MsgType_ID_CommonNtosReq_inner(reqInOut, saveDB, true)
		if saveDB && rspOut.FromServer {
			if _, err := thls.xEngine.ID(core.PK{reqInOut.SeqNo}).Update(&CommonAtosDataNode{Finish: rspOut.FromServer, ErrNo: rspOut.ErrNo, ErrMsg: rspOut.ErrMsg}); err != nil {
				glog.Fatalf("Engine.Update with err=%v", err)
			}
		}
	}
	return
}

func (thls *businessServer) reportData(dataIn *txdata.ReportDataItem, d time.Duration, isEndeavour bool) *CommRspData {
	reqInOut := Message2CommonNtosReq(dataIn, time.Now(), thls.ownInfo.UserID)
	rspOut := thls.commonAtos(reqInOut, isEndeavour, d)
	return CommonNtosReqRsp2CommRspData(reqInOut, rspOut)
}

func (thls *businessServer) sendMail(dataIn *txdata.SendMailItem, d time.Duration, isEndeavour bool) *CommRspData {
	reqInOut := Message2CommonNtosReq(dataIn, time.Now(), thls.ownInfo.UserID)
	rspOut := thls.commonAtos(reqInOut, isEndeavour, d)
	return CommonNtosReqRsp2CommRspData(reqInOut, rspOut)
}
