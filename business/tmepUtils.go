/*
关于golang的日志库,
https://github.com/golang/glog
https://github.com/uber-go/zap
https://godoc.org/go.uber.org/zap
深度 | 从Go高性能日志库zap看如何实现高性能Go组件
https://studygolang.com/articles/14220
在Github中stars数最多的Go日志库集合
https://studygolang.com/articles/11995

关于orm,
有Golang的ORM框架推荐么？或者相互比较的文章?
https://www.zhihu.com/question/55072439
在Github中stars数最多的Go数据库框架库集合
https://my.oschina.net/u/168737/blog/1531834
golang orm对比
https://segmentfault.com/a/1190000015606291
最终，我选择在gorm和xorm里面进行挑选。
因为因为gorm不用gorm.Model的话，update时会在控制台打印警告信息，我还不知道怎么关掉它，另外xorm有gorm.Model的功能，所以我选择xorm。

XORM - 官方博客
http://blog.xorm.io/
使用手册 - xorm: 简单而强大的 Go 语言ORM框架
http://xorm.io/docs/

sqlite 不允许非主键的自增属性
https://www.sqlite.org/autoinc.html
Because AUTOINCREMENT keyword changes the behavior of the ROWID selection algorithm, AUTOINCREMENT is not allowed on WITHOUT ROWID tables or on any table column other than INTEGER PRIMARY KEY. Any attempt to use AUTOINCREMENT on a WITHOUT ROWID table or on a column other than the INTEGER PRIMARY KEY column results in an error.

sqlite中 是不是可以设置非主键的自动增长列-CSDN论坛
https://bbs.csdn.net/topics/390603321
*/

package main

import (
	"encoding/json"
	"fmt"
	"net/smtp"
	"reflect"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/zx9229/myexe/txdata"
	"github.com/zx9229/myexe/wsnet"
)

type safeWsSocketMap struct {
	sync.Mutex
	M map[*wsnet.WsSocket]bool
}

func newSafeWsSocketMap() *safeWsSocketMap {
	return &safeWsSocketMap{M: make(map[*wsnet.WsSocket]bool)}
}

func (thls *safeWsSocketMap) insertData(k *wsnet.WsSocket, v bool) (isSuccess bool) {
	thls.Lock()
	if _, isSuccess = thls.M[k]; !isSuccess {
		thls.M[k] = v
	}
	thls.Unlock()
	isSuccess = !isSuccess
	return
}

func (thls *safeWsSocketMap) deleteData(k *wsnet.WsSocket) (v bool, isSuccess bool) {
	thls.Lock()
	if v, isSuccess = thls.M[k]; isSuccess {
		delete(thls.M, k)
	}
	thls.Unlock()
	return
}

type connInfoEx struct {
	conn    *wsnet.WsSocket
	Info    txdata.ConnectionInfo
	Pathway []string
}

type safeConnInfoMap struct {
	sync.Mutex
	M map[string]*connInfoEx
}

func newSafeConnInfoMap() *safeConnInfoMap {
	return &safeConnInfoMap{M: make(map[string]*connInfoEx)}
}

func (thls *safeConnInfoMap) humanReadable() (jsonContent string) {
	thls.Lock()
	if byteSlice, err := json.Marshal(thls.M); err != nil {
		glog.Fatalln(err, thls.M)
	} else {
		jsonContent = string(byteSlice)
	}
	thls.Unlock()
	return
}

func (thls *safeConnInfoMap) queryData(key string) (connEx *connInfoEx, isExist bool) {
	thls.Lock()
	connEx, isExist = thls.M[key]
	thls.Unlock()
	return
}

func (thls *safeConnInfoMap) isValidData(data *connInfoEx) bool {
	var isOk bool
	for range "1" {
		if data.conn == nil {
			break
		}
		if len(data.Info.UniqueID) == 0 {
			break
		}
		if data.Pathway == nil {
			break
		}
		isOk = true
	}
	return isOk
}

func (thls *safeConnInfoMap) insertData(data *connInfoEx) (isSuccess bool) {
	thls.Lock()
	if _, isSuccess = thls.M[data.Info.UniqueID]; !isSuccess {
		thls.M[data.Info.UniqueID] = data
	}
	thls.Unlock()
	isSuccess = !isSuccess
	return
}

func (thls *safeConnInfoMap) deleteData(key string) (isSuccess bool) {
	thls.Lock()
	if _, isSuccess = thls.M[key]; isSuccess {
		delete(thls.M, key)
	}
	thls.Unlock()
	return
}

func (thls *safeConnInfoMap) deleteDataByConn(conn *wsnet.WsSocket) []*connInfoEx {
	var dataSlice []*connInfoEx
	thls.Lock()
	for key, val := range thls.M {
		if val.conn == conn {
			if dataSlice == nil {
				dataSlice = make([]*connInfoEx, 0)
			}
			dataSlice = append(dataSlice, val)
			delete(thls.M, key)
		}
	}
	thls.Unlock()
	return dataSlice
}

type byte4type [4]byte //用于int32相关

func msg2slice(msgType txdata.MsgType, msgData proto.Message) (dst []byte) {
	var err error
	if dst, err = proto.Marshal(msgData); err != nil {
		glog.Fatalln(err, msgData)
	}
	dst = append((*byte4type)(unsafe.Pointer(&msgType))[:2], dst...)
	return
}

func slice2msg(src []byte) (msgType txdata.MsgType, msgData proto.Message, err error) {
	b4 := (*byte4type)(unsafe.Pointer(&msgType))
	b4[0] = src[0]
	b4[1] = src[1]

	switch msgType {
	case txdata.MsgType_ID_ConnectedData:
		msgData = new(txdata.ConnectedData)
	case txdata.MsgType_ID_DisconnectedData:
		msgData = new(txdata.DisconnectedData)
	case txdata.MsgType_ID_ExecuteCommandReq:
		msgData = new(txdata.ExecuteCommandReq)
	case txdata.MsgType_ID_ExecuteCommandRsp:
		msgData = new(txdata.ExecuteCommandRsp)
	case txdata.MsgType_ID_CommonAtosReq:
		msgData = new(txdata.CommonAtosReq)
	case txdata.MsgType_ID_CommonAtosRsp:
		msgData = new(txdata.CommonAtosRsp)
	default:
		msgData = nil
		err = fmt.Errorf("unknown txdata.MsgType=%v", msgType)
	}
	if msgData != nil {
		if err = proto.Unmarshal(src[2:], msgData); err != nil {
			msgData = nil
			err = fmt.Errorf("Unmarshal with err=%v, msgType=%v", err, msgType)
		}
	}

	return
}

func sendMail(username, password, smtpAddr, to, subject, contentType, content string) error {
	/*
		username := "sender@163.com"
		password := "senderPassword"
		smtpAddr := "smtp.163.com:25"
		to := "receiver1@163.com;receiver2@126.com;receiver3@hotmail.com"
		subject := "测试邮件"
		bodyType := "plain"
		content := "这是一封测试邮件，用于测试自动发送。"
		它只负责发送信息到服务器,至于,收件人是否正确,是否被退信,之类的详细提示,需要登录邮箱查看.
	*/
	if contentType != "html" {
		contentType = "plain"
	}
	currAuth := smtp.PlainAuth("", username, password, strings.Split(smtpAddr, ":")[0])
	var mailMsg string
	mailMsg += fmt.Sprintf("From: %s\r\n", username)
	mailMsg += fmt.Sprintf("To: %s\r\n", to)
	mailMsg += fmt.Sprintf("Subject: %s\r\n", subject)
	mailMsg += fmt.Sprintf("Content-Type: text/%s; charset=UTF-8\r\n", contentType)
	mailMsg += fmt.Sprintf("\r\n")
	mailMsg += content
	return smtp.SendMail(smtpAddr, currAuth, username, strings.Split(to, ";"), []byte(mailMsg))
}

//ReportDataNode 上报的数据(存储到Node)
type ReportDataNode struct {
	SeqNo       int64     `xorm:"pk autoincr notnull unique"`
	UniqueID    string    //
	Topic       string    //
	Data        string    //
	ReportTime  time.Time //报告时间
	CreateTime  time.Time `xorm:"created"` //这个Field将在Insert时自动赋值为当前时间
	FatalErrNo  int32     `xorm:"notnull"` //不为0,表示这一条数据,SERVER处理不了(比如:主键冲突等原因,插数据库失败),防止无限循环.
	FatalErrMsg string    //错误的具体原因.
}

//ReportDataServer 上报的数据(存储到Server)(sqlite不是主键的话,ID不能自增,所以没有添加ID)
//如果field名称为Id而且类型为int64并且没有定义tag，则会被xorm视为主键，并且拥有自增属性。
//如果想用Id以外的名字或非int64类型做为主键名，必须在对应的Tag上加上xorm:"pk"来定义主键，加上xorm:"autoincr"作为自增。
//这里需要注意的是，有些数据库并不允许非主键的自增属性。
type ReportDataServer struct {
	SeqNo      int64     `xorm:"pk notnull"`
	UniqueID   string    `xorm:"pk notnull"`
	Topic      string    //
	Data       string    //
	ReportTime time.Time //报告时间
	CreateTime time.Time `xorm:"created"` //这个Field将在Insert时自动赋值为当前时间
}

//KeyValue omit
type KeyValue struct {
	Key   string `xorm:"notnull pk"`
	Value string `xorm:"notnull"`
}

//CommonAtosDataNode omit
type CommonAtosDataNode struct {
	SeqNo       int64     `xorm:"pk autoincr notnull unique"`
	UniqueID    string    //
	DataType    string    //
	Data        []byte    //
	ReportTime  time.Time //报告时间
	CreateTime  time.Time `xorm:"created"` //这个Field将在Insert时自动赋值为当前时间
	FatalErrNo  int32     `xorm:"notnull"` //不为0,表示这一条数据,SERVER处理不了(比如:主键冲突等原因,插数据库失败),防止无限循环.
	FatalErrMsg string    //错误的具体原因.
}

//CommonAtosDataServer omit
type CommonAtosDataServer struct {
	SeqNo      int64     `xorm:"pk notnull"`
	UniqueID   string    `xorm:"pk notnull"`
	DataType   string    //
	Data       []byte    //
	ReportTime time.Time //报告时间
	CreateTime time.Time `xorm:"created"` //这个Field将在Insert时自动赋值为当前时间
}

func needSendRsp_CommonAtos_RequestID(requestID int64) bool {
	//(正:超时等待,要回响应);(零:不等待,不用回复响应);(负:背景上报,要回响应)
	return (requestID != 0)
}

func reqrspRelated_RequestID(requestID int64) bool {
	//req&rsp相关(//从safeNodeReqRspCache出来的RequestID都是正数)
	return (0 < requestID)
}

func backgroundRelated_RequestID(requestID int64) bool {
	//背景执行相关的请求ID. //从background出来的RequestID都是负数
	return (requestID < 0)
}

func dbRelated_CommonAtos_SeqNo(seqNo int64) bool {
	//(正:缓存数据,发不过去要重试)(零:未缓存数据,发不过去就算了)(负:绝无可能)//(SeqNo非0,表示插入了数据库)
	return (seqNo != 0)
}

type CommRspData struct {
	UniqueID string
	SeqNo    int64
	ErrNo    int32
	ErrMsg   string
}

func CommonAtosReq2CommonAtosRsp4Err(reqIn *txdata.CommonAtosReq, errNo int32, errMsg string) *txdata.CommonAtosRsp {
	return &txdata.CommonAtosRsp{RequestID: reqIn.RequestID, Pathway: nil, SeqNo: reqIn.SeqNo, ErrNo: errNo, ErrMsg: errMsg}
}

func Message2CommonAtosReq(src proto.Message, reportTime time.Time, uniqueID string, isEndeavour bool) *txdata.CommonAtosReq {
	dst := &txdata.CommonAtosReq{RequestID: 0, UniqueID: uniqueID, SeqNo: 0, Endeavour: isEndeavour, DataType: reflect.TypeOf(src).String(), Data: nil, ReportTime: nil}
	var err error
	if dst.Data, err = proto.Marshal(src); err != nil {
		glog.Fatalln(err, src)
	}
	if dst.ReportTime, err = ptypes.TimestampProto(reportTime); err != nil {
		glog.Fatalln(err, reportTime)
	}
	return dst
}

func CommonAtosReqRsp2CommRspData(req *txdata.CommonAtosReq, rsp *txdata.CommonAtosRsp) *CommRspData {
	return &CommRspData{UniqueID: req.UniqueID, SeqNo: req.SeqNo, ErrNo: rsp.ErrNo, ErrMsg: rsp.ErrMsg}
}

func CommonAtosReq2CommonAtosDataNode(reqIn *txdata.CommonAtosReq) *CommonAtosDataNode {
	var err error
	cada := &CommonAtosDataNode{SeqNo: 0, UniqueID: reqIn.UniqueID, DataType: reqIn.DataType, Data: reqIn.Data, ReportTime: time.Time{}}
	if cada.ReportTime, err = ptypes.Timestamp(reqIn.ReportTime); err != nil {
		glog.Fatalln(err)
	}
	return cada
}
