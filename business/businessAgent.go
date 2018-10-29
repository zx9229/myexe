package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
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

type businessAgent struct {
	cacheSock   *safeWsSocketMap
	cacheAgent  *safeConnInfoMap
	cacheReqRsp *safeNodeReqRspCache
	ownInfo     txdata.ConnectionInfo
	parentData  connInfoEx
	xEngine     *xorm.Engine
	workChan    chan int64
}

func newBusinessAgent(cfg *configAgent) *businessAgent {
	if len(cfg.UniqueID) == 0 || len(cfg.BelongID) == 0 {
		glog.Fatalf("must not be empty, UniqueID=%v, BelongID=%v", cfg.UniqueID, cfg.BelongID)
	}
	if cfg.UniqueID == cfg.BelongID {
		glog.Fatalf("must not be equal, UniqueID=%v, BelongID=%v", cfg.UniqueID, cfg.BelongID)
	}
	curData := new(businessAgent)
	//
	curData.cacheSock = newSafeWsSocketMap()
	curData.cacheAgent = newSafeConnInfoMap()
	curData.cacheReqRsp = newSafeNodeReqRspCache()
	//
	curData.ownInfo.UniqueID = cfg.UniqueID
	curData.ownInfo.BelongID = cfg.BelongID
	curData.ownInfo.Version = "Version20181020"
	curData.ownInfo.ExeType = txdata.ConnectionInfo_AGENT
	curData.ownInfo.LinkDir = txdata.ConnectionInfo_Zero3
	curData.ownInfo.ExePid = int32(os.Getpid())
	curData.ownInfo.ExePath, _ = filepath.Abs(os.Args[0])
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

func (thls *businessAgent) initEngine(dataSourceName string, locationName string) {
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
	if err = thls.xEngine.CreateTables(&KeyValue{}, &CommonAtosDataAgent{}); err != nil { //应该是:只要存在这个tablename,就跳过它.
		glog.Fatalln(err)
	}
	if err = thls.xEngine.Sync2(&KeyValue{}, &CommonAtosDataAgent{}); err != nil { //同步数据库结构
		glog.Fatalln(err)
	}
}

func (thls *businessAgent) checkCachedDatabase() {
	//程序第一次启动后,可能接收并缓存了数据,然后关闭了程序,然后可能有人修改了缓存数据库里的配置,然后又启动程序,
	//程序启动时,需要检查,缓存数据库里的数据和配置是否冲突,有冲突的话,则拒绝启动.
	var err error
	//(CommonAtosDataAgent.UniqueID)必须等于(txdata.ConnectionInfo.UniqueID)
	var rowData CommonAtosDataAgent
	var affected1, affected2 int64
	if affected1, err = thls.xEngine.Count(&rowData); err != nil {
		glog.Fatalln(err)
	}
	rowData.UniqueID = thls.ownInfo.UniqueID
	if affected2, err = thls.xEngine.Count(&rowData); err != nil {
		glog.Fatalln(err)
	}
	if affected1 != affected2 {
		glog.Fatalln(affected1, affected2)
	}
}

func (thls *businessAgent) backgroundWork() {
	CommonAtosDataAgent2CommonAtosReq := func(src *CommonAtosDataAgent) *txdata.CommonAtosReq {
		//负数的RequestID表示背景工作在做事情.
		req := &txdata.CommonAtosReq{RequestID: -1, UniqueID: src.UniqueID, SeqNo: src.SeqNo, Endeavour: true, DataType: src.DataType, Data: src.Data, ReportTime: nil}
		req.ReportTime, _ = ptypes.TimestampProto(src.ReportTime)
		return req
	}
	data4qry := &CommonAtosDataAgent{}
	fnReportTime := zxxorm.GuessColName(thls.xEngine, data4qry, unsafe.Offsetof(data4qry.ReportTime), true)
	fnFatalErrNo := zxxorm.GuessColName(thls.xEngine, data4qry, unsafe.Offsetof(data4qry.FatalErrNo), true)
	go func() {
		for {
			time.Sleep(time.Second * 5)
			thls.workChan <- -1
		}
	}()
	//查询单条数据使用Get方法，在调用Get方法时需要传入一个对应结构体的指针，同时结构体中的非空field自动成为查询的条件和前面的方法条件组合在一起查询.
	var result CommonAtosDataAgent
	var has bool
	var err error
	for true {
		result = CommonAtosDataAgent{}
		data4qry.ReportTime = time.Now().Add(-30 * time.Second) //查询30秒之前的数据(可能刚执行了一个上报操作,刚插入数据库,所以要有一个缓存时段).
		if has, err = thls.xEngine.Where(builder.Eq{fnFatalErrNo: 0}.And(builder.Lt{fnReportTime: data4qry.ReportTime})).Get(&result); err != nil {
			glog.Fatalf("xorm.Get with has=%v, err=%v", has, err)
		} else if has {
			err = thls.sendDataToParent(txdata.MsgType_ID_CommonAtosReq, CommonAtosDataAgent2CommonAtosReq(&result))
			//如果没有东西要发送(has == false),也是等待30秒,然后再查询一下数据库.
			glog.Infof("background report data with SeqNo=%v and err=%v", result.SeqNo, err)
		}
		for looping := true; looping; {
			select {
			case val, isOk := <-thls.workChan:
				if !isOk {
					glog.Fatalf("recv data from chan with val=%v, isOk=%v", val, isOk)
				}
				if val < 0 { //负数是超时协程发送的数据.
					if 60 < time.Now().Sub(data4qry.ReportTime).Seconds() { //超时30秒了,就跳出循环
						looping = false
					}
					continue
				}
				if val != result.SeqNo {
					glog.Warningf("val=%v, result.SeqNo=%v", val, result.SeqNo)
				}
				looping = false //上报给SERVER并且收到正确的回复了,就跳出循环.
			default:
			}
		}
	}
}

func (thls *businessAgent) onConnected(msgConn *wsnet.WsSocket, isAccepted bool) {
	glog.Warningf("[   onConnected] msgConn=%p, isAccepted=%v, LocalAddr=%v, RemoteAddr=%v", msgConn, isAccepted, msgConn.LocalAddr(), msgConn.RemoteAddr())
	if !thls.cacheSock.insertData(msgConn, isAccepted) {
		glog.Fatalf("already exists, msgConn=%p", msgConn)
	}
	if !isAccepted {
		tmpTxData := txdata.ConnectedData{Info: &thls.ownInfo, Pathway: []string{thls.ownInfo.UniqueID}}
		//tmpTxData.Info.LinkDir = txdata.ConnectionInfo_CONNECT
		msgConn.Send(msg2slice(txdata.MsgType_ID_ConnectedData, &tmpTxData))
		//tmpTxData.Info.LinkDir = txdata.ConnectionInfo_Zero3
	}
}

func (thls *businessAgent) onDisconnected(msgConn *wsnet.WsSocket, err error) {
	glog.Warningf("[onDisconnected] msgConn=%p, err=%v", msgConn, err)
	if thls.parentData.conn == msgConn {
		//如果与父亲断开连接,就清理父亲的数据,这样就不用sendDataToParent了.
		glog.Infof("disconnected with father, msgConn=%p", msgConn)
		thls.parentData = connInfoEx{}
	}
	if dataSlice := thls.cacheAgent.deleteDataByConn(msgConn); dataSlice != nil {
		//儿子和我断开连接,我要清理掉儿子和孙子的缓存.
		sonNum := 0
		for _, node := range dataSlice {
			if len(node.Pathway) == 1 { //步长为1的是儿子.
				sonNum++
			}
		}
		if 1 < sonNum {
			glog.Fatalf("one msgConn to multi son connInfoEx, msgConn=%p", msgConn)
		}
		for _, data := range dataSlice { //发给父亲,让父亲也清理掉对应的缓存.
			tmpTxData := txdata.DisconnectedData{Info: &data.Info}
			thls.sendDataToParent(txdata.MsgType_ID_DisconnectedData, &tmpTxData)
		}
	}
	thls.deleteConnectionFromAll(msgConn, false)
}

func (thls *businessAgent) onMessage(msgConn *wsnet.WsSocket, msgData []byte, msgType int) {
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
	case txdata.MsgType_ID_CommonAtosReq:
		thls.handle_MsgType_ID_CommonAtosReq(txMsgData.(*txdata.CommonAtosReq), msgConn)
	case txdata.MsgType_ID_CommonAtosRsp:
		thls.handle_MsgType_ID_CommonAtosRsp(txMsgData.(*txdata.CommonAtosRsp), msgConn)
	case txdata.MsgType_ID_ExecuteCommandReq:
		thls.handle_MsgType_ID_ExecuteCommandReq(txMsgData.(*txdata.ExecuteCommandReq), msgConn)
	case txdata.MsgType_ID_ExecuteCommandRsp:
		thls.handle_MsgType_ID_ExecuteCommandRsp(txMsgData.(*txdata.ExecuteCommandRsp), msgConn)
	default:
		glog.Errorf("unknown txdata.MsgType, msgConn=%p, txMsgType=%v", msgConn, txMsgType)
	}
}

func (thls *businessAgent) handle_MsgType_ID_ConnectedData(msgData *txdata.ConnectedData, msgConn *wsnet.WsSocket) {
	sendToParent := true
	for range "1" {
		if (msgData.Pathway == nil) || (len(msgData.Pathway) == 0) {
			glog.Errorf("empty Pathway, will disconnect with it, msgConn=%p, msgData=%v", msgConn, msgData)
			thls.deleteConnectionFromAll(msgConn, true)
			sendToParent = false
			break
		}
		if len(msgData.Pathway) == 1 {
			// 1<len(msgData.Pathway)时,这个消息是[孙代]发送给[子代],[子代]转发给[父代](我这一代)的消息.
			if (msgData.Info.ExeType == txdata.ConnectionInfo_AGENT) &&
				(thls.ownInfo.ExeType == txdata.ConnectionInfo_AGENT) &&
				(msgData.Info.UniqueID == thls.ownInfo.UniqueID) {
				glog.Errorf("maybe i connected myself, msgConn=%p, msgData=%v", msgConn, msgData)
				thls.deleteConnectionFromAll(msgConn, true)
				sendToParent = false
				break
			}
			if (msgData.Info.UniqueID != thls.ownInfo.BelongID) && (msgData.Info.BelongID != thls.ownInfo.UniqueID) {
				glog.Errorf("he is not my father, i am not his father, msgConn=%p, msgData=%v", msgConn, msgData)
				thls.deleteConnectionFromAll(msgConn, true)
				sendToParent = false
				break
			}
		}
		if thls.isParentConnection(msgData) {
			sendToParent = thls.doDeal4parent(msgData, msgConn)
			break
		} else if msgData.Info.ExeType == txdata.ConnectionInfo_AGENT {
			sendToParent = thls.doDeal4agent(msgData, msgConn)
			break
		} else {
			glog.Errorf("unknown message, msgConn=%p, msgData=%v", msgConn, msgData)
			thls.deleteConnectionFromAll(msgConn, true)
			sendToParent = false
			break
		}
	}

	if sendToParent {
		msgData.Pathway = append(msgData.Pathway, thls.ownInfo.UniqueID)
		thls.sendDataToParent(txdata.MsgType_ID_ConnectedData, msgData)
	}
}

func (thls *businessAgent) handle_MsgType_ID_DisconnectedData(msgData *txdata.DisconnectedData, msgConn *wsnet.WsSocket) {
	if thls.parentData.conn == msgConn {
		glog.Errorf("the data must not be from my father, msgConn=%p, msgData=%v", msgConn, msgData)
		return
	}
	if msgData.Info.ExeType == txdata.ConnectionInfo_AGENT {
		if thls.cacheAgent.deleteData(msgData.Info.UniqueID) == false {
			glog.Fatalf("cache data error, msgConn=%p, msgData=%v", msgConn, msgData)
		}
	} else {
		glog.Errorf("unmanageable, msgConn=%p, msgData=%v", msgConn, msgData)
		return
	}
	thls.sendDataToParent(txdata.MsgType_ID_DisconnectedData, msgData)
}

func (thls *businessAgent) handle_MsgType_ID_CommonAtosReq(msgData *txdata.CommonAtosReq, msgConn *wsnet.WsSocket) {
	if thls.parentData.conn == msgConn { //上报请求,肯定要发给父亲,所以,一定不是父亲发过来的.
		glog.Errorf("the data must not be from my father, msgConn=%p, msgData=%v", msgConn, msgData)
		return
	}

	if err := thls.sendDataToParent(txdata.MsgType_ID_CommonAtosReq, msgData); err != nil && needSendRsp_CommonAtos_RequestID(msgData.RequestID) {
		if connInfoEx, isExist := thls.cacheAgent.queryData(msgData.UniqueID); isExist {
			rspData := CommonAtosReq2CommonAtosRsp4Err(msgData, -1, err.Error())
			rspData.Pathway = connInfoEx.Pathway
			connInfoEx.conn.Send(msg2slice(txdata.MsgType_ID_CommonAtosRsp, rspData))
		} else {
			//儿子刚发过来数据,我还没处理呢,结果儿子和我断开了,缓存也清理掉了,然后我才开始处理.
			glog.Warningf("user not found, msgConn=%p, msgData=%v", msgConn, msgData)
		}
	}
}

func (thls *businessAgent) handle_MsgType_ID_CommonAtosRsp(msgData *txdata.CommonAtosRsp, msgConn *wsnet.WsSocket) {
	if thls.parentData.conn != msgConn {
		glog.Errorf("the data must come from my father, msgConn=%p, msgData=%v", msgConn, msgData)
		return
	}
	var pathwayLen int
	if pathwayLen = len(msgData.Pathway); pathwayLen == 0 {
		glog.Errorf("empty Pathway, msgConn=%p, msgData=%v", msgConn, msgData)
		return
	}
	if thls.ownInfo.UniqueID != msgData.Pathway[pathwayLen-1] {
		glog.Errorf("illegal Pathway, msgConn=%p, msgData=%v", msgConn, msgData)
		return
	}
	msgData.Pathway = msgData.Pathway[:pathwayLen-1]
	if pathwayLen = len(msgData.Pathway); pathwayLen != 0 {
		nextUID := msgData.Pathway[pathwayLen-1]
		if nextConnInfoEx, isExist := thls.cacheAgent.queryData(nextUID); isExist {
			nextConnInfoEx.conn.Send(msg2slice(txdata.MsgType_ID_CommonAtosRsp, msgData))
		} else {
			glog.Warningf("user not found, msgConn=%p, msgData=%v", msgConn, msgData)
		}
	} else {
		if 0 < msgData.RequestID { //从safeNodeReqRspCache出来的RequestID都是正数
			if node, isExist := thls.cacheReqRsp.deleteElement(msgData.RequestID); isExist {
				node.rspType = txdata.MsgType_ID_CommonAtosRsp
				node.rspData = msgData
				node.condVar.notifyAll()
			} else {
				glog.Infof("data not found in cache, RequestID=%v", msgData.RequestID)
			}
		}
		if dbRelated_CommonAtos_SeqNo(msgData.SeqNo) { //SeqNo非0,表示插入了数据库.
			if msgData.ErrNo == 0 { //ErrNo为0,表示SERVER处理成功,AGENT可以删除自己的缓存了.
				if affected, err := thls.xEngine.Delete(&CommonAtosDataAgent{SeqNo: msgData.SeqNo}); err != nil {
					glog.Fatalf("Engine.Delete with affected=%v, err=%v", affected, err)
				}
				//可能AGENT短时间内发送了两个相同的请求,此时,第一个响应已经删除了数据,第二个响应会执行成功,同时删除零行(猜测//TODO:).
				//所以,可能存在(err == nil && affected == 0)的情况.
			}
			if msgData.ErrNo == -83 { //为(-83)表示SERVER无法处理这个数据,此时AGENT不应当再上报它了,因为上报了也处理不了.
				if _, err := thls.xEngine.ID(core.PK{msgData.SeqNo}).Update(&CommonAtosDataAgent{FatalErrNo: msgData.ErrNo, FatalErrMsg: msgData.ErrMsg}); err != nil {
					glog.Fatalf("Engine.Update with err=%v", err)
				}
			}
		}
		if msgData.RequestID < 0 { //从background出来的RequestID都是负数
			thls.workChan <- msgData.SeqNo
		}
	}
}

func (thls *businessAgent) handle_MsgType_ID_ExecuteCommandReq(msgData *txdata.ExecuteCommandReq, msgConn *wsnet.WsSocket) {
	if thls.parentData.conn != msgConn {
		glog.Errorf("the data must be from the father, msgConn=%p, msgData=%v", msgConn, msgData)
		return
	}
	if len(msgData.Pathway) == 0 {
		glog.Errorf("empty Pathway, msgConn=%p, msgData=%v", msgConn, msgData)
		return
	}
	if thls.ownInfo.UniqueID != msgData.Pathway[len(msgData.Pathway)-1] {
		glog.Errorf("illegal Pathway, msgConn=%p, msgData=%v", msgConn, msgData)
		return
	}
	msgData.Pathway = msgData.Pathway[:len(msgData.Pathway)-1]
	if length := len(msgData.Pathway); length != 0 {
		nextUID := msgData.Pathway[length-1]
		if nextConnInfo, isExist := thls.cacheAgent.queryData(nextUID); isExist {
			nextConnInfo.conn.Send(msg2slice(txdata.MsgType_ID_ExecuteCommandReq, msgData))
		} else {
			tempTxData := txdata.ExecuteCommandRsp{RequestID: msgData.RequestID, ErrMsg: fmt.Sprintf("next step is unreachable, nextUID=%v", nextUID)}
			thls.sendDataToParent(txdata.MsgType_ID_ExecuteCommandRsp, &tempTxData)
		}
	} else {
		glog.Warningln("ExecuteCommand:", msgData.Command) //TODO:待添加真正的执行代码.

		tempTxData := txdata.ExecuteCommandRsp{RequestID: msgData.RequestID, UniqueID: thls.ownInfo.UniqueID, Result: "OK, Now=" + time.Now().Format("2006-01-02_15:04:05")}
		thls.sendDataToParent(txdata.MsgType_ID_ExecuteCommandRsp, &tempTxData)
	}
}

func (thls *businessAgent) handle_MsgType_ID_ExecuteCommandRsp(msgData *txdata.ExecuteCommandRsp, msgConn *wsnet.WsSocket) {
	if thls.parentData.conn == msgConn {
		glog.Errorf("the data must not be from my father, msgConn=%p, msgData=%v", msgConn, msgData)
		return
	}
	thls.sendDataToParent(txdata.MsgType_ID_ExecuteCommandRsp, msgData)
}

func (thls *businessAgent) sendDataToParent(msgType txdata.MsgType, msgData proto.Message) error {
	conn := thls.parentData.conn
	if conn == nil {
		return errors.New("parent is offline")
	}
	return conn.Send(msg2slice(msgType, msgData))
}

func (thls *businessAgent) isParentConnection(data *txdata.ConnectedData) bool {
	var isParent bool
	for range "1" {
		if thls.ownInfo.ExeType != txdata.ConnectionInfo_AGENT {
			glog.Fatalf("illegal data, ownInfo=%v", thls.ownInfo)
			break
		}
		if data.Info.ExeType != txdata.ConnectionInfo_AGENT && data.Info.ExeType != txdata.ConnectionInfo_SERVER {
			break
		}
		if data.Info.UniqueID != thls.ownInfo.BelongID {
			break
		}
		isParent = true
	}
	return isParent
}

func (thls *businessAgent) doDeal4parent(msgData *txdata.ConnectedData, msgConn *wsnet.WsSocket) (sendToParent bool) {
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
		//if thls.parentData.info.LinkDir != txdata.ConnectionInfo_CONNECT {
		//	log.Panicln("parent info is abnormal", thls.parentData.info)
		//}
		tmpTxData := txdata.ConnectedData{Info: &thls.ownInfo, Pathway: []string{thls.ownInfo.UniqueID}}
		msgConn.Send(msg2slice(txdata.MsgType_ID_ConnectedData, &tmpTxData))
	} else {
		//if thls.parentData.info.LinkDir != txdata.ConnectionInfo_ACCEPT {
		//	log.Panicln("parent info is abnormal", thls.parentData.info)
		//}
	}

	if true {
		//和父亲建立连接了,要把自己的缓存发送给父亲,更新父亲的缓存.
		thls.cacheAgent.Lock()
		for _, node := range thls.cacheAgent.M {
			tmpTxData := txdata.ConnectedData{Info: &node.Info, Pathway: append(node.Pathway, thls.ownInfo.UniqueID)}
			thls.sendDataToParent(txdata.MsgType_ID_ConnectedData, &tmpTxData)
		}
		thls.cacheAgent.Unlock()
	}

	return
}

func (thls *businessAgent) doDeal4agent(msgData *txdata.ConnectedData, msgConn *wsnet.WsSocket) (sendToParent bool) {
	var isAccepted bool
	isSon := len(msgData.Pathway) == 1
	if isSon {
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

	sendToParent = true

	return
}

func (thls *businessAgent) deleteConnectionFromAll(conn *wsnet.WsSocket, closeIt bool) {
	if closeIt {
		conn.Close()
	}
	if thls.parentData.conn == conn {
		thls.parentData = connInfoEx{}
	}
	thls.cacheSock.deleteData(conn)
	thls.cacheAgent.deleteDataByConn(conn)
}

func (thls *businessAgent) commonAtos(reqInOut *txdata.CommonAtosReq, d time.Duration) (rspOut *txdata.CommonAtosRsp) {
	if true { //修复请求结构体的相关字段.
		reqInOut.RequestID = 0
		reqInOut.UniqueID = thls.ownInfo.UniqueID
		reqInOut.SeqNo = 0
		//reqInOut.Endeavour
		//reqInOut.DataType
		//reqInOut.Data
		reqInOut.ReportTime, _ = ptypes.TimestampProto(time.Now())
	}
	CommonAtosReq2CommonAtosDataAgent := func(reqIn *txdata.CommonAtosReq) *CommonAtosDataAgent {
		cada := &CommonAtosDataAgent{SeqNo: 0, UniqueID: reqIn.UniqueID, DataType: reqIn.DataType, Data: reqIn.Data, ReportTime: time.Time{}}
		cada.ReportTime, _ = ptypes.Timestamp(reqIn.ReportTime)
		return cada
	}
	for range "1" {
		var err error
		if reqInOut.Endeavour { //根据SeqNo可知它是否缓存
			cada := CommonAtosReq2CommonAtosDataAgent(reqInOut)
			var affected int64
			if affected, err = thls.xEngine.InsertOne(cada); err != nil {
				rspOut = CommonAtosReq2CommonAtosRsp4Err(reqInOut, -1, fmt.Sprintf("insert to db with err=%v", err))
				break
			}
			if affected != 1 {
				glog.Fatalf("Engine.InsertOne with affected=%v, err=%v", affected, err) //我就是想知道,成功的话,除了1,还有其他值吗.
			}
			reqInOut.SeqNo = cada.SeqNo //利用xorm的特性.
		}
		//
		if d <= 0 {
			if err = thls.sendDataToParent(txdata.MsgType_ID_CommonAtosReq, reqInOut); err != nil {
				rspOut = CommonAtosReq2CommonAtosRsp4Err(reqInOut, -1, err.Error())
			} else {
				rspOut = CommonAtosReq2CommonAtosRsp4Err(reqInOut, 0, "send success and no wait.")
			}
		} else { //根据RequestID可知它是否在等待.
			node := thls.cacheReqRsp.generateElement()
			if true {
				reqInOut.RequestID = node.requestID
				//
				node.reqType = txdata.MsgType_ID_CommonAtosReq
				node.reqData = reqInOut
			}
			//
			if err = thls.sendDataToParent(node.reqType, node.reqData); err != nil {
				rspOut = CommonAtosReq2CommonAtosRsp4Err(reqInOut, -1, err.Error())
				break
			}
			//
			if isTimeout := node.condVar.waitFor(d); isTimeout {
				rspOut = CommonAtosReq2CommonAtosRsp4Err(reqInOut, -1, "timeout")
				break
			}
			rspOut = node.rspData.(*txdata.CommonAtosRsp)
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

func (thls *businessAgent) reportData(dataIn *txdata.ReportDataItem, d time.Duration, isEndeavour bool) *CommRspData {
	toCommonAtosReq := func(src *txdata.ReportDataItem, reportTime time.Time) *txdata.CommonAtosReq {
		dst := &txdata.CommonAtosReq{RequestID: 0, UniqueID: thls.ownInfo.UniqueID, SeqNo: 0, Endeavour: isEndeavour, DataType: reflect.TypeOf(src).String(), Data: nil, ReportTime: nil}
		var err error
		if dst.Data, err = proto.Marshal(src); err != nil {
			glog.Fatalln(err, src)
		}
		if dst.ReportTime, err = ptypes.TimestampProto(reportTime); err != nil {
			glog.Fatalln(err, reportTime)
		}
		return dst
	}
	toCommRspData := func(req *txdata.CommonAtosReq, rsp *txdata.CommonAtosRsp) *CommRspData {
		return &CommRspData{UniqueID: req.UniqueID, SeqNo: req.SeqNo, ErrNo: rsp.ErrNo, ErrMsg: rsp.ErrMsg}
	}
	reqInOut := toCommonAtosReq(dataIn, time.Now())
	rspOut := thls.commonAtos(reqInOut, d)
	return toCommRspData(reqInOut, rspOut)
}
