package main

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/golang/protobuf/ptypes"
	"github.com/zx9229/myexe/txdata"
)

//DbPushWrap omit
type DbPushWrap struct {
	MsgNo      int64  `xorm:"notnull pk autoincr"`
	UserID     string `xorm:"notnull"`
	MsgTime    time.Time
	MsgType    int32
	MsgData    []byte
	InsertTime time.Time `xorm:"created"`
}

func (thls *DbPushWrap) insertOneResult(affected int64, err error) (isExist, isInsert bool) {
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

func (thls *DbPushWrap) From(src *txdata.PushWrap) {
	thls.MsgNo = src.MsgNo
	thls.UserID = src.UserID
	thls.MsgTime, _ = ptypes.Timestamp(src.MsgTime)
	thls.MsgType = int32(src.MsgType)
	thls.MsgData = src.MsgData
}
func (thls *DbPushWrap) To() (dst *txdata.PushWrap) {
	dst = &txdata.PushWrap{}
	dst.MsgNo = thls.MsgNo
	dst.UserID = thls.UserID
	dst.MsgTime, _ = ptypes.TimestampProto(thls.MsgTime)
	dst.MsgType = txdata.MsgType(thls.MsgType)
	dst.MsgData = thls.MsgData
	return dst
}

//DbOp omit
type DbOp struct {
	wg       sync.WaitGroup
	pm       ProtoMessage
	handler  func(session *xorm.Session, data *DbOp)
	affected int64
	err      error
}
