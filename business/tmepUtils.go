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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/smtp"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/zx9229/myexe/txdata"
	"github.com/zx9229/myexe/wsnet"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
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
		if len(data.Info.UserID) == 0 {
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
	if _, isSuccess = thls.M[data.Info.UserID]; !isSuccess {
		thls.M[data.Info.UserID] = data
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

func msg2slice(msgData ProtoMessage) (dst []byte) {
	var err error
	if dst, err = proto.Marshal(msgData); err != nil {
		glog.Fatalln(err, msgData)
	}
	msgType := CalcMessageIndex(msgData)
	dst = append((*byte4type)(unsafe.Pointer(&msgType))[:2], dst...)
	return
}

func slice2msg(src []byte) (msgType txdata.MsgType, msgData ProtoMessage, err error) {
	b4 := (*byte4type)(unsafe.Pointer(&msgType))
	b4[0] = src[0]
	b4[1] = src[1]

	switch msgType {
	case txdata.MsgType_ID_ConnectedData:
		msgData = new(txdata.ConnectedData)
	case txdata.MsgType_ID_DisconnectedData:
		msgData = new(txdata.DisconnectedData)
	case txdata.MsgType_ID_CommonNtosReq:
		msgData = new(txdata.CommonNtosReq)
	case txdata.MsgType_ID_CommonNtosRsp:
		msgData = new(txdata.CommonNtosRsp)
	case txdata.MsgType_ID_ParentDataReq:
		msgData = new(txdata.ParentDataReq)
	case txdata.MsgType_ID_ParentDataRsp:
		msgData = new(txdata.ParentDataRsp)
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
	assert4true(CalcMessageType(msgData) == msgType)

	return
}

func tryToUTF8(src []byte) (utf8data string) {
	var dst []byte
	var err error
	if dst, err = ioutil.ReadAll(transform.NewReader(bytes.NewReader(src), simplifiedchinese.GBK.NewDecoder())); err == nil {
		return string(dst)
	}
	if dst, err = ioutil.ReadAll(transform.NewReader(bytes.NewReader(src), traditionalchinese.Big5.NewDecoder())); err == nil {
		return string(dst)
	}
	if dst, err = ioutil.ReadAll(transform.NewReader(bytes.NewReader(src), simplifiedchinese.GB18030.NewDecoder())); err == nil {
		return string(dst)
	}
	return string(src)
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
	mailMsg += tryToUTF8([]byte(content))
	return smtp.SendMail(smtpAddr, currAuth, username, strings.Split(to, ";"), []byte(mailMsg))
}

//ReportDataNode 上报的数据(存储到Node)
type ReportDataNode struct {
	SeqNo       int64     `xorm:"pk autoincr notnull unique"`
	UserID      string    //
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
	UserID     string    `xorm:"pk notnull"`
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

//CommonNtosReqDbN omit
type CommonNtosReqDbN struct {
	RequestID  int64
	UserID     string `xorm:"   notnull"`
	SeqNo      int64  `xorm:"pk notnull"`
	ReqType    txdata.MsgType
	ReqData    []byte
	ReqTime    time.Time
	RefNum     int64
	CreateTime time.Time `xorm:"created"` //这个Field将在Insert时自动赋值为当前时间
	State      int32     `xorm:"notnull"` //请求消息的状态(目前将int用作bool;0=>false)
}

func CommonNtosReqDbN2CommonNtosReq(src *CommonNtosReqDbN) (dst *txdata.CommonNtosReq) {
	dst = &txdata.CommonNtosReq{}
	dst.RequestID = src.RequestID
	dst.UserID = src.UserID
	dst.SeqNo = src.SeqNo
	dst.ReqType = src.ReqType
	dst.ReqData = src.ReqData
	dst.ReqTime, _ = ptypes.TimestampProto(src.ReqTime)
	dst.RefNum = src.RefNum
	return
}

func CommonNtosReq2CommonNtosReqDbN(src *txdata.CommonNtosReq, dst *CommonNtosReqDbN) {
	dst.RequestID = src.RequestID
	dst.UserID = src.UserID
	dst.SeqNo = src.SeqNo
	dst.ReqType = src.ReqType
	dst.ReqData = src.ReqData
	dst.ReqTime, _ = ptypes.Timestamp(src.ReqTime)
	dst.RefNum = src.RefNum
}

//CommonNtosReqDbS omit
type CommonNtosReqDbS struct {
	RequestID  int64
	UserID     string `xorm:"pk notnull"`
	SeqNo      int64  `xorm:"pk notnull"`
	ReqType    txdata.MsgType
	ReqData    []byte
	ReqTime    time.Time
	RefNum     int64
	CreateTime time.Time `xorm:"created"` //这个Field将在Insert时自动赋值为当前时间
}

func CommonNtosReq2CommonNtosReqDbS(src *txdata.CommonNtosReq, dst *CommonNtosReqDbS) {
	dst.RequestID = src.RequestID
	dst.UserID = src.UserID
	dst.SeqNo = src.SeqNo
	dst.ReqType = src.ReqType
	dst.ReqData = src.ReqData
	dst.ReqTime, _ = ptypes.Timestamp(src.ReqTime)
	dst.RefNum = src.RefNum
}

//CommonNtosRspDb omit
type CommonNtosRspDb struct {
	RequestID  int64
	Pathway    []string
	SeqNo      int64
	RspType    txdata.MsgType
	RspData    []byte
	RspTime    time.Time
	FromServer bool
	ErrNo      int32
	ErrMsg     string
	RefNum     int64
	CreateTime time.Time `xorm:"created"` //这个Field将在Insert时自动赋值为当前时间
}

func CommonNtosRsp2CommonNtosRspDb(src *txdata.CommonNtosRsp) (dst *CommonNtosRspDb) {
	dst = &CommonNtosRspDb{}
	dst.RequestID = src.RequestID
	dst.Pathway = src.Pathway
	dst.SeqNo = src.SeqNo
	dst.RspType = src.RspType
	dst.RspData = src.RspData
	dst.RspTime, _ = ptypes.Timestamp(src.RspTime)
	dst.FromServer = src.FromServer
	dst.ErrNo = src.ErrNo
	dst.ErrMsg = src.ErrMsg
	dst.RefNum = src.RefNum
	return
}

//CommonStonReqDb omit
type CommonStonReqDb struct {
	RequestID  int64          //
	Pathway    string         //
	SeqNo      int64          `xorm:"pk notnull"`
	ReqType    txdata.MsgType //
	ReqData    []byte         //
	ReqTime    time.Time      //
	RefNum     int64          //
	CreateTime time.Time      `xorm:"created"` //这个Field将在Insert时自动赋值为当前时间.
	State      int32          `xorm:"notnull"` //请求消息的状态(目前将int用作bool;0=>false)
}

func CommonStonReq2CommonStonReqDb(src *txdata.CommonStonReq, dst *CommonStonReqDb) {
	dst.RequestID = src.RequestID
	dst.Pathway = strings.Join(src.Pathway, ".")
	dst.SeqNo = src.SeqNo
	dst.ReqType = src.ReqType
	dst.ReqData = src.ReqData
	dst.ReqTime, _ = ptypes.Timestamp(src.ReqTime)
	dst.RefNum = src.RefNum
}

//CommonStonRspDb omit
type CommonStonRspDb struct {
	RequestID  int64          //
	UserID     string         //
	SeqNo      int64          //
	RspType    txdata.MsgType //
	RspData    []byte         //
	RspTime    time.Time      //
	RefNum     int64          //
	FromTarget bool           //
	State      int32          //
	ErrNo      int32          //
	ErrMsg     string         //
	CreateTime time.Time      `xorm:"created"` //这个Field将在Insert时自动赋值为当前时间.
}

func CommonStonRsp2CommonStonRspDb(src *txdata.CommonStonRsp) (dst *CommonStonRspDb) {
	dst.RequestID = src.RequestID
	dst.UserID = src.UserID
	dst.SeqNo = src.SeqNo
	dst.RspType = src.RspType
	dst.RspData = src.RspData
	dst.RspTime, _ = ptypes.Timestamp(src.RspTime)
	dst.RefNum = src.RefNum
	dst.FromTarget = src.FromTarget
	dst.State = src.State
	dst.ErrNo = src.ErrNo
	dst.ErrMsg = src.ErrMsg
	return
}

type CommRspData struct {
	UserID string
	SeqNo  int64
	ErrNo  int32
	ErrMsg string
}

func fillCommonNtosRspByCommonNtosReq(req *txdata.CommonNtosReq, rsp *txdata.CommonNtosRsp) {
	rsp.RequestID = req.RequestID
	rsp.SeqNo = req.SeqNo
	rsp.RefNum = req.RefNum
}

func fillCommonStonRspByCommonStonReq(req *txdata.CommonStonReq, rsp *txdata.CommonStonRsp) {
	rsp.RequestID = req.RequestID
	rsp.SeqNo = req.SeqNo
	rsp.RefNum = req.RefNum
}

func CommonNtosReq2CommonNtosRsp4Err(reqObj *txdata.CommonNtosReq, eNo int32, eMsg string, fromS bool) *txdata.CommonNtosRsp {
	rspObj := &txdata.CommonNtosRsp{FromServer: fromS, ErrNo: eNo, ErrMsg: eMsg}
	fillCommonNtosRspByCommonNtosReq(reqObj, rspObj)
	return rspObj
}

func CommonStonReq2CommonStonRsp4Err(reqObj *txdata.CommonStonReq, eNo int32, eMsg string, fromT bool, uId string) *txdata.CommonStonRsp {
	rspObj := &txdata.CommonStonRsp{UserID: uId, FromTarget: fromT, ErrNo: eNo, ErrMsg: eMsg}
	fillCommonStonRspByCommonStonReq(reqObj, rspObj)
	return rspObj
}

func CommonNtosReq2CommonNtosRsp4Rsp(reqObj *txdata.CommonNtosReq, eNo int32, eMsg string, fromS bool, rspType txdata.MsgType, rspData []byte) *txdata.CommonNtosRsp {
	rspObj := &txdata.CommonNtosRsp{FromServer: fromS, ErrNo: eNo, ErrMsg: eMsg, RspType: rspType, RspData: rspData}
	fillCommonNtosRspByCommonNtosReq(reqObj, rspObj)
	return rspObj
}

func Message2CommonNtosReq(src ProtoMessage, reportTime time.Time, userID string) *txdata.CommonNtosReq {
	dst := &txdata.CommonNtosReq{RequestID: 0, UserID: userID, SeqNo: 0, ReqType: CalcMessageType(src), ReqData: nil, ReqTime: nil}
	var err error
	if dst.ReqData, err = proto.Marshal(src); err != nil {
		glog.Fatalln(err, src)
	}
	if dst.ReqTime, err = ptypes.TimestampProto(reportTime); err != nil {
		glog.Fatalln(err, reportTime)
	}
	return dst
}

func CommonNtosReqRsp2CommRspData(req *txdata.CommonNtosReq, rsp *txdata.CommonNtosRsp) *CommRspData {
	return &CommRspData{UserID: req.UserID, SeqNo: req.SeqNo, ErrNo: rsp.ErrNo, ErrMsg: rsp.ErrMsg}
}

func atomicKey2Str(src *txdata.AtomicKey) string {
	execType := (*int32)(unsafe.Pointer(&src.ExecType))
	return fmt.Sprintf("/%v/%v/%v/%v", src.ZoneName, src.NodeName, *execType, src.ExecName)
}

func str2ProgramType(src string) txdata.ProgramType {
	if dst, ok := txdata.ProgramType_value[src]; ok {
		return txdata.ProgramType(dst)
	}
	return txdata.ProgramType_Zero2
}

func atomicKeyIsValid(src *txdata.AtomicKey) bool {
	isValidChar := func(c byte) bool {
		if 33 <= c && c <= 126 && c != '/' {
			return true
		}
		return false
	}
	for _, c := range []byte(src.ZoneName) {
		if !isValidChar(c) {
			return false
		}
	}
	for _, c := range []byte(src.NodeName) {
		if !isValidChar(c) {
			return false
		}
	}
	if src.ExecType == txdata.ProgramType_Zero2 {
		return false
	}
	for _, c := range []byte(src.ExecName) {
		if !isValidChar(c) {
			return false
		}
	}
	return true
}

func assert4true(cond bool) {
	if !cond {
		panic(cond)
	}
}

func assert4false(cond bool) {
	if cond {
		panic(cond)
	}
}

func calc_flag_RequestID_SeqNo(RequestID, SeqNo int64) (p, qau, qas, r bool) {
	p = false   //推送数据,发出去之后就不管了.(isPush)
	qau = false //请求响应,中途丢包就丢了,(question-answer-unsafe)
	qas = false //请求响应,中途丢包会重试.(question-answer-safe)
	r = false   //请求响应,中途丢包后重试消息(retransmit)
	assert4false(SeqNo < 0)
	assert4false(RequestID < 0 && SeqNo == 0)
	p = RequestID == 0 && SeqNo == 0
	qau = RequestID > 0 && SeqNo == 0
	assert4false(RequestID == 0 && SeqNo > 0)
	r = RequestID < 0 && SeqNo > 0
	qas = RequestID > 0 && SeqNo > 0
	return
}

func CommonNtosReq_flag(data *txdata.CommonNtosReq) (p, qau, qas, r bool) {
	return calc_flag_RequestID_SeqNo(data.RequestID, data.SeqNo)
}

func CommonNtosRsp_flag(data *txdata.CommonNtosRsp) (p, qau, qas, r bool) {
	return calc_flag_RequestID_SeqNo(data.RequestID, data.SeqNo)
}

func CommonStonReq_flag(data *txdata.CommonStonReq) (p, qau, qas, r bool) {
	return calc_flag_RequestID_SeqNo(data.RequestID, data.SeqNo)
}

func CommonStonRsp_flag(data *txdata.CommonStonRsp) (p, qau, qas, r bool) {
	return calc_flag_RequestID_SeqNo(data.RequestID, data.SeqNo)
}

//ProtoMessage omit
type ProtoMessage interface {
	Reset()                      //pb.Message
	String() string              //pb.Message
	ProtoMessage()               //pb.Message
	Descriptor() ([]byte, []int) //自动生成的结构体,全都包含该成员函数.
}

//CalcMessageIndex 用法示例:CalcMessageIndex(&txdata.CommonNtosReq{})
func CalcMessageIndex(protoMessage ProtoMessage) int32 {
	var data []int
	_, data = protoMessage.Descriptor()
	return int32(data[0])
}

//CalcMessageType 用法示例:CalcMessageType(&txdata.CommonNtosReq{})
func CalcMessageType(protoMessage ProtoMessage) txdata.MsgType {
	var data []int
	_, data = protoMessage.Descriptor()
	return txdata.MsgType(data[0])
}
