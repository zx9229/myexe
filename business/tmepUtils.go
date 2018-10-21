/*
关于golang的日志库,
https://github.com/golang/glog
https://github.com/uber-go/zap
https://godoc.org/go.uber.org/zap
深度 | 从Go高性能日志库zap看如何实现高性能Go组件
https://studygolang.com/articles/14220
在Github中stars数最多的Go日志库集合
https://studygolang.com/articles/11995
*/
package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
	"unsafe"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
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
	case txdata.MsgType_ID_PushData:
		msgData = new(txdata.PushData)
	case txdata.MsgType_ID_ConnectedData:
		msgData = new(txdata.ConnectedData)
	case txdata.MsgType_ID_DisconnectedData:
		msgData = new(txdata.DisconnectedData)
	case txdata.MsgType_ID_ExecuteCommandReq:
		msgData = new(txdata.ExecuteCommandReq)
	case txdata.MsgType_ID_ExecuteCommandRsp:
		msgData = new(txdata.ExecuteCommandRsp)
	case txdata.MsgType_ID_ReportDataReq:
		msgData = new(txdata.ReportDataReq)
	case txdata.MsgType_ID_ReportDataRsp:
		msgData = new(txdata.ReportDataRsp)
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

//ReportDataAgent 上报的数据(存储到Agent)
type ReportDataAgent struct {
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
