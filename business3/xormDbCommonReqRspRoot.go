package main

import (
	"fmt"
	"strings"
	"time"
	"unsafe"

	"github.com/golang/protobuf/ptypes"
	"github.com/zx9229/myexe/txdata"
)

//DbCommonReqRspRoot omit
type DbCommonReqRspRoot struct {
	MsgType    int32  `xorm:"   notnull"`
	UserID     string `xorm:"pk notnull"`
	MsgNo      int64  `xorm:"pk notnull"`
	SeqNo      int32  `xorm:"pk notnull"`
	BatchNo    int64
	RefNum     int64
	RefText    string
	SenderID   string
	RecverID   string
	ToRoot     bool
	IsLog      bool
	IsSafe     bool
	IsPush     bool
	UpCache    bool
	TxType     int32
	TxData     []byte
	TxTime     time.Time
	IsLast     bool
	InsertTime time.Time `xorm:"created"` //这个Field将在Insert时自动赋值为当前时间
}

func (thls *DbCommonReqRspRoot) insertOneResult(affected int64, err error) (isExist, isInsert bool) {
	if affected == 1 && err == nil {
		isExist = false
		isInsert = true
		return
	}
	if affected == 0 && err != nil {
		isInsert = false
		//我们先故意构造一个重复数据,然后就知道报错详情了,然后进行字符串匹配.
		if strings.Index(err.Error(), "UNIQUE constraint failed") == -1 {
			isExist = false
		} else {
			isExist = true
		}
		return
	}
	panic(fmt.Sprintf("affected=%v, err=%v", affected, err))
}

//FromC2Req omit
func (thls *DbCommonReqRspRoot) FromC2Req(src *txdata.Common2Req) {
	thls.MsgType = CalcMessageIndex(src)
	thls.UserID = src.Key.UserID
	thls.MsgNo = src.Key.MsgNo
	thls.SeqNo = src.Key.SeqNo
	thls.BatchNo = src.BatchNo
	thls.RefNum = src.RefNum
	thls.RefText = src.RefText
	thls.SenderID = src.SenderID
	thls.RecverID = src.RecverID
	thls.ToRoot = src.ToRoot
	thls.IsLog = src.IsLog
	thls.IsSafe = src.IsSafe
	thls.IsPush = src.IsPush
	thls.UpCache = src.UpCache
	thls.TxType = *(*int32)(unsafe.Pointer(&src.ReqType))
	thls.TxData = src.ReqData
	thls.TxTime, _ = ptypes.Timestamp(src.ReqTime)
	//thls.IsLast = src.IsLast
	//thls.InsertTime = src.InsertTime
}

//FromC2Rsp omit
func (thls *DbCommonReqRspRoot) FromC2Rsp(src *txdata.Common2Rsp) {
	thls.MsgType = CalcMessageIndex(src)
	thls.UserID = src.Key.UserID
	thls.MsgNo = src.Key.MsgNo
	thls.SeqNo = src.Key.SeqNo
	thls.BatchNo = src.BatchNo
	thls.RefNum = src.RefNum
	thls.RefText = src.RefText
	thls.SenderID = src.SenderID
	thls.RecverID = src.RecverID
	thls.ToRoot = src.ToRoot
	thls.IsLog = src.IsLog
	thls.IsSafe = src.IsSafe
	thls.IsPush = src.IsPush
	thls.UpCache = src.UpCache
	thls.TxType = *(*int32)(unsafe.Pointer(&src.RspType))
	thls.TxData = src.RspData
	thls.TxTime, _ = ptypes.Timestamp(src.RspTime)
	thls.IsLast = src.IsLast
	//thls.InsertTime = src.InsertTime
}

//To omit
func (thls *DbCommonReqRspRoot) To() ProtoMessage {
	curType := (*txdata.MsgType)(unsafe.Pointer(&thls.MsgType))
	if *curType == txdata.MsgType_ID_Common2Req {
		return thls.ToC2Req()
	}
	if *curType == txdata.MsgType_ID_Common2Rsp {
		return thls.ToC2Rsp()
	}
	return nil
}

//ToC2Req omit
func (thls *DbCommonReqRspRoot) ToC2Req() (dst *txdata.Common2Req) {
	dst = &txdata.Common2Req{Key: &txdata.UniKey{}}
	dst.Key.UserID = thls.UserID
	dst.Key.MsgNo = thls.MsgNo
	dst.Key.SeqNo = thls.SeqNo
	dst.BatchNo = thls.BatchNo
	dst.RefNum = thls.RefNum
	dst.RefText = thls.RefText
	dst.SenderID = thls.SenderID
	dst.RecverID = thls.RecverID
	dst.ToRoot = thls.ToRoot
	dst.IsLog = thls.IsLog
	dst.IsSafe = thls.IsSafe
	dst.IsPush = thls.IsPush
	dst.UpCache = thls.UpCache
	dst.ReqType = txdata.MsgType(thls.TxType)
	dst.ReqData = thls.TxData
	dst.ReqTime, _ = ptypes.TimestampProto(thls.TxTime)
	return
}

//ToC2Rsp omit
func (thls *DbCommonReqRspRoot) ToC2Rsp() (dst *txdata.Common2Rsp) {
	dst = &txdata.Common2Rsp{Key: &txdata.UniKey{}}
	dst.Key.UserID = thls.UserID
	dst.Key.MsgNo = thls.MsgNo
	dst.Key.SeqNo = thls.SeqNo
	dst.BatchNo = thls.BatchNo
	dst.RefNum = thls.RefNum
	dst.RefText = thls.RefText
	dst.SenderID = thls.SenderID
	dst.RecverID = thls.RecverID
	dst.ToRoot = thls.ToRoot
	dst.IsLog = thls.IsLog
	dst.IsSafe = thls.IsSafe
	dst.IsPush = thls.IsPush
	dst.UpCache = thls.UpCache
	dst.RspType = txdata.MsgType(thls.TxType)
	dst.RspData = thls.TxData
	dst.RspTime, _ = ptypes.TimestampProto(thls.TxTime)
	dst.IsLast = thls.IsLast
	return
}
