package main

import (
	"fmt"
	"unsafe"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/zx9229/myexe/txdata"
)

//ProtoMessage 对protobufMessage进行了抽象.
type ProtoMessage interface {
	Reset()                      //pb.Message
	String() string              //pb.Message
	ProtoMessage()               //pb.Message
	Descriptor() ([]byte, []int) //自动生成的结构体,全都包含该成员函数.
}

//CalcMessageIndex 用法示例:CalcMessageIndex(&txdata.ConnectionInfo{})
func CalcMessageIndex(protoMessage ProtoMessage) int32 {
	var data []int
	_, data = protoMessage.Descriptor()
	return int32(data[0])
}

//CalcMessageType 用法示例:CalcMessageType(&txdata.ConnectionInfo{})
func CalcMessageType(protoMessage ProtoMessage) txdata.MsgType {
	var data []int
	_, data = protoMessage.Descriptor()
	return txdata.MsgType(data[0])
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
	// 二进制数据的前2个字节标识了后面数据的类型.
	b4 := (*byte4type)(unsafe.Pointer(&msgType))
	b4[0] = src[0]
	b4[1] = src[1]
	// 需要在shell下,先创建ff函数,再执行ff函数.
	// ff(){ sed -n '/^enum MsgType/,/}/p' "$1" | sed 's/[ \t]*\?\(ID_\)\([^ \t]\+\).*/case txdata.MsgType_\1\2: \n msgData = new(txdata.\2)/g' ; }
	// ff  txdata.proto
	switch msgType {
	case txdata.MsgType_ID_DataPsh:
		msgData = new(txdata.DataPsh)
	case txdata.MsgType_ID_DataAck:
		msgData = new(txdata.DataAck)
	case txdata.MsgType_ID_CommonReq:
		msgData = new(txdata.CommonReq)
	case txdata.MsgType_ID_CommonRsp:
		msgData = new(txdata.CommonRsp)
	case txdata.MsgType_ID_ConnectionInfo:
		msgData = new(txdata.ConnectionInfo)
	case txdata.MsgType_ID_ConnectReq:
		msgData = new(txdata.ConnectReq)
	case txdata.MsgType_ID_ConnectRsp:
		msgData = new(txdata.ConnectRsp)
	case txdata.MsgType_ID_DisconnectedData:
		msgData = new(txdata.DisconnectedData)
	case txdata.MsgType_ID_ParentDataReq:
		msgData = new(txdata.ParentDataReq)
	case txdata.MsgType_ID_ParentDataRsp:
		msgData = new(txdata.ParentDataRsp)
	case txdata.MsgType_ID_EchoItem:
		msgData = new(txdata.EchoItem)
	case txdata.MsgType_ID_SendMailItem:
		msgData = new(txdata.SendMailItem)
	case txdata.MsgType_ID_ReportDataItem:
		msgData = new(txdata.ReportDataItem)
	default:
		msgData = nil
		err = fmt.Errorf("unknown txdata.MsgType(%v)", msgType)
	}
	if msgData != nil {
		if err = proto.Unmarshal(src[2:], msgData); err != nil {
			msgData = nil
			err = fmt.Errorf("Unmarshal failure(err=%v, msgType=%v)", err, msgType)
		}
	}
	assert4true(CalcMessageType(msgData) == msgType)

	return
}
