package main

import (
	"fmt"
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
	cacheAgent  *safeConnInfoMap
	cacheReqRsp *safeNodeReqRspCache
	ownInfo     txdata.ConnectionInfo
	xEngine     *xorm.Engine
}

func newBusinessServer(cfg *configServer) *businessServer {
	if len(cfg.UniqueID) == 0 {
		glog.Fatalf("must not be empty, UniqueID=%v", cfg.UniqueID)
	}
	curData := new(businessServer)
	//
	curData.cacheSock = newSafeWsSocketMap()
	curData.cacheAgent = newSafeConnInfoMap()
	curData.cacheReqRsp = newSafeNodeReqRspCache()
	//
	curData.ownInfo.UniqueID = cfg.UniqueID
	curData.ownInfo.BelongID = ""
	curData.ownInfo.Version = "Version20181021"
	curData.ownInfo.ExeType = txdata.ConnectionInfo_SERVER
	curData.ownInfo.LinkDir = txdata.ConnectionInfo_Zero3
	curData.ownInfo.ExePid = int32(os.Getpid())
	curData.ownInfo.ExePath, _ = filepath.Abs(os.Args[0])
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
	if err = thls.xEngine.CreateTables(&KeyValue{}, &CommonAtosDataAgent{}, &CommonAtosDataServer{}); err != nil { //应该是:只要存在这个tablename,就跳过它.
		glog.Fatalln(err)
	}
	if err = thls.xEngine.Sync2(&KeyValue{}, &CommonAtosDataAgent{}, &CommonAtosDataServer{}); err != nil { //同步数据库结构
		glog.Fatalln(err)
	}
}

func (thls *businessServer) onConnected(msgConn *wsnet.WsSocket, isAccepted bool) {
	glog.Warningf("[   onConnected] msgConn=%p, isAccepted=%v, LocalAddr=%v, RemoteAddr=%v", msgConn, isAccepted, msgConn.LocalAddr(), msgConn.RemoteAddr())
	if !thls.cacheSock.insertData(msgConn, isAccepted) {
		glog.Fatalf("already exists, msgConn=%v", msgConn)
	}
	if !isAccepted {
		tmpTxData := txdata.ConnectedData{Info: &thls.ownInfo, Pathway: []string{thls.ownInfo.UniqueID}}
		msgConn.Send(msg2slice(txdata.MsgType_ID_ConnectedData, &tmpTxData))
	}
}

func (thls *businessServer) onDisconnected(msgConn *wsnet.WsSocket, err error) {
	glog.Warningf("[onDisconnected] msgConn=%p, err=%v", msgConn, err)
	if dataSlice := thls.cacheAgent.deleteDataByConn(msgConn); dataSlice != nil {
		//儿子和我断开连接,我要清理掉儿子和孙子的缓存.
		sonNum := 0
		for _, node := range dataSlice {
			if len(node.Pathway) == 1 { //步长为1的是儿子.
				sonNum++
			}
		}
		if 1 < sonNum {
			glog.Fatalf("one msgConn to multi connInfoEx, msgConn=%p", msgConn)
		}
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
	case txdata.MsgType_ID_ExecuteCommandReq:
		thls.handle_MsgType_ID_ExecuteCommandReq(txMsgData.(*txdata.ExecuteCommandReq), msgConn)
	case txdata.MsgType_ID_ExecuteCommandRsp:
		thls.handle_MsgType_ID_ExecuteCommandRsp(txMsgData.(*txdata.ExecuteCommandRsp), msgConn)
	case txdata.MsgType_ID_CommonAtosReq:
		thls.handle_MsgType_ID_CommonAtosReq(txMsgData.(*txdata.CommonAtosReq), msgConn)
	case txdata.MsgType_ID_CommonAtosRsp:
		thls.handle_MsgType_ID_CommonAtosRsp(txMsgData.(*txdata.CommonAtosRsp), msgConn)
	default:
		glog.Errorf("unknown txdata.MsgType, msgConn=%p, txMsgType=%v", msgConn, txMsgType)
	}
}

func (thls *businessServer) handle_MsgType_ID_ConnectedData(msgData *txdata.ConnectedData, msgConn *wsnet.WsSocket) {
	for range "1" {
		if (msgData.Pathway == nil) || (len(msgData.Pathway) == 0) {
			glog.Errorf("empty Pathway, will disconnect with it, msgConn=%p, msgData=%v", msgConn, msgData)
			thls.deleteConnectionFromAll(msgConn, true)
			break
		}
		if msgData.Info.ExeType != txdata.ConnectionInfo_AGENT {
			glog.Errorf("unknown message, msgConn=%p, msgData=%v", msgConn, msgData)
			thls.deleteConnectionFromAll(msgConn, true)
			break
		}
		if msgData.Info.UniqueID != msgData.Pathway[0] {
			glog.Errorf("illegal Pathway, msgConn=%p, msgData=%v", msgConn, msgData)
			thls.deleteConnectionFromAll(msgConn, true)
			break
		}
		if (msgData.Info.ExeType == txdata.ConnectionInfo_AGENT) &&
			(len(msgData.Pathway) == 1) &&
			(msgData.Info.BelongID != thls.ownInfo.UniqueID) {
			glog.Errorf("i am not his father, msgConn=%p, msgData=%v", msgConn, msgData)
			thls.deleteConnectionFromAll(msgConn, true)
			break
		}
		thls.doDeal4agent(msgData, msgConn)
	}
}

func (thls *businessServer) handle_MsgType_ID_DisconnectedData(msgData *txdata.DisconnectedData, msgConn *wsnet.WsSocket) {
	if msgData.Info.ExeType == txdata.ConnectionInfo_AGENT {
		if thls.cacheAgent.deleteData(msgData.Info.UniqueID) == false {
			glog.Fatalf("cache data error, msgConn=%p, msgData=%v", msgConn, msgData)
		}
	} else {
		glog.Errorf("unmanageable, msgConn=%p, msgData=%v", msgConn, msgData)
		return
	}
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

func (thls *businessServer) doDeal4agent(msgData *txdata.ConnectedData, msgConn *wsnet.WsSocket) {
	var isAccepted bool
	isSon := len(msgData.Pathway) == 1
	if isSon {
		//步长为0的数据,在主调函数中已经进行了过滤.
		//步长等于1时,是儿子socket发送的数据.
		//步长大于1时,是孙子socket发送的数据.
		var isExist bool
		if isAccepted, isExist = thls.cacheSock.deleteData(msgConn); !isExist {
			glog.Errorf("msgConn not found in cache, msgConn=%p, msgData=%v", msgConn, msgData)
			thls.deleteConnectionFromAll(msgConn, true)
			return
		}
	}

	curData := new(connInfoEx)
	curData.conn = msgConn
	curData.Info = *msgData.Info
	curData.Pathway = msgData.Pathway

	if isSuccess := thls.cacheAgent.insertData(curData); !isSuccess {
		glog.Errorf("agent already exists, msgConn=%p, msgData=%v", msgConn, msgData)
		thls.deleteConnectionFromAll(msgConn, true)
		return
	}

	if isSon && isAccepted {
		tmpTxData := txdata.ConnectedData{Info: &thls.ownInfo, Pathway: []string{thls.ownInfo.UniqueID}}
		msgConn.Send(msg2slice(txdata.MsgType_ID_ConnectedData, &tmpTxData))
	}
}

func (thls *businessServer) deleteConnectionFromAll(conn *wsnet.WsSocket, closeIt bool) {
	if closeIt {
		conn.Close()
	}
	thls.cacheSock.deleteData(conn)
	thls.cacheAgent.deleteDataByConn(conn)
}

func (thls *businessServer) executeCommand(reqInOut *txdata.ExecuteCommandReq, d time.Duration) (rspOut *txdata.ExecuteCommandRsp) {
	//这里选择用(reqInOut.Pathway[0])承载(要控制哪个AGENT)信息.
	tmpUID := reqInOut.Pathway[0]
	if true { //修复请求结构体的相关字段.
		reqInOut.RequestID = 0
		reqInOut.Pathway = nil
		//reqIn.Command
	}
	ExecuteCommandReq2ExecuteCommandRsp := func(reqIn *txdata.ExecuteCommandReq, errNo int32, errMsg string, uid string) *txdata.ExecuteCommandRsp {
		return &txdata.ExecuteCommandRsp{RequestID: reqIn.RequestID, UniqueID: uid, Result: "", ErrNo: errNo, ErrMsg: errMsg}
	}
	for range "1" {
		connInfoEx, isExist := thls.cacheAgent.queryData(tmpUID)
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

func (thls *businessServer) handle_MsgType_ID_CommonAtosReq_txdata_ReportDataItem(commReq *txdata.CommonAtosReq, item *txdata.ReportDataItem) (errNo int32, errMsg string) {
	glog.Infoln(commReq.DataType, item)
	return
}

func (thls *businessServer) handle_MsgType_ID_CommonAtosReq_process(msgData *txdata.CommonAtosReq) (rspData *txdata.CommonAtosRsp) {
	CommonAtosReq2CommonAtosRsp4Err := func(reqIn *txdata.CommonAtosReq, errNo int32, errMsg string) *txdata.CommonAtosRsp {
		return &txdata.CommonAtosRsp{RequestID: reqIn.RequestID, Pathway: nil, SeqNo: reqIn.SeqNo, ErrNo: errNo, ErrMsg: errMsg}
	}
	var errNo int32
	var errMsg string
	switch msgData.DataType {
	case "txdata.ReportDataItem":
		curData := &txdata.ReportDataItem{}
		if err := proto.Unmarshal(msgData.Data, curData); err != nil {
			glog.Fatalln(msgData)
		}
		errNo, errMsg = thls.handle_MsgType_ID_CommonAtosReq_txdata_ReportDataItem(msgData, curData)
	default:
		errNo = -1
		errMsg = "unknown data type"
	}
	return CommonAtosReq2CommonAtosRsp4Err(msgData, errNo, errMsg)
}

func (thls *businessServer) handle_MsgType_ID_CommonAtosReq_inner(msgData *txdata.CommonAtosReq) (rspData *txdata.CommonAtosRsp) {
	CommonAtosReq2CommonAtosRsp4Err := func(reqIn *txdata.CommonAtosReq, errNo int32, errMsg string) *txdata.CommonAtosRsp {
		return &txdata.CommonAtosRsp{RequestID: reqIn.RequestID, Pathway: nil, SeqNo: reqIn.SeqNo, ErrNo: errNo, ErrMsg: errMsg}
	}
	for range "1" {
		//以(UniqueID+SeqNo)唯一定位一条数据.
		cads := &CommonAtosDataServer{SeqNo: msgData.SeqNo, UniqueID: msgData.UniqueID}
		if has, err := thls.xEngine.Get(cads); err != nil {
			glog.Fatalf("Engine.Get with has=%v, err=%v, rds=%v", has, err, cads)
			rspData = CommonAtosReq2CommonAtosRsp4Err(msgData, -83, fmt.Sprintf("query from db with err=%v", err))
			break
		} else if has {
			//已经存在这一条数据了,就(ErrNo=0)然后AGENT会从缓存中移除这一条数据.
			rspData = CommonAtosReq2CommonAtosRsp4Err(msgData, 0, "already existed")
			break
		}
		if true {
			//cads.SeqNo
			//cads.UniqueID
			cads.DataType = msgData.DataType
			cads.Data = msgData.Data
			cads.ReportTime, _ = ptypes.Timestamp(msgData.ReportTime)
		}
		if affected, err := thls.xEngine.InsertOne(cads); err != nil {
			glog.Fatalf("Engine.InsertOne with affected=%v, err=%v, cads=%v", affected, err, cads)
			rspData = CommonAtosReq2CommonAtosRsp4Err(msgData, -83, fmt.Sprintf("insert to db with err=%v", err))
			break
		} else if affected != 1 {
			glog.Fatalf("Engine.InsertOne with affected=%v, err=%v, cads=%v", affected, err, cads) //我就是想知道,成功的话,除了1,还有其他值吗.
		}
		rspData = thls.handle_MsgType_ID_CommonAtosReq_process(msgData)
	}
	return
}

func (thls *businessServer) handle_MsgType_ID_CommonAtosReq(msgData *txdata.CommonAtosReq, msgConn *wsnet.WsSocket) {
	rspData := thls.handle_MsgType_ID_CommonAtosReq_inner(msgData)
	//
	if needSendRsp_CommonAtos_RequestID(msgData.RequestID) {
		if connInfoEx, isExist := thls.cacheAgent.queryData(msgData.UniqueID); isExist {
			rspData.Pathway = connInfoEx.Pathway
			connInfoEx.conn.Send(msg2slice(txdata.MsgType_ID_CommonAtosRsp, rspData))
		} else {
			glog.Infof("user not found, msgConn=%p, msgData=%v", msgConn, msgData)
		}
	}
}

func (thls *businessServer) handle_MsgType_ID_CommonAtosRsp(msgData *txdata.CommonAtosRsp, msgConn *wsnet.WsSocket) {
	glog.Errorf("the data must not be from my father, msgConn=%p, msgData=%v", msgConn, msgData)
	return
}
