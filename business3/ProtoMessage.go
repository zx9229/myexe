package main

import (
	"fmt"
	"reflect"
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

//DeepCopy omit
func DeepCopy(src ProtoMessage) (dst ProtoMessage) {
	if src == nil {
		return
	}
	dst = reflect.New(reflect.TypeOf(src).Elem()).Interface().(ProtoMessage)
	var err error
	var byteSlice []byte
	if byteSlice, err = proto.Marshal(src); err != nil {
		panic(err)
	}
	if err = proto.Unmarshal(byteSlice, dst); err != nil {
		panic(err)
	}
	return
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
	return
}

func msg2package(msgData ProtoMessage) (dst []byte) {
	dst = msg2slice(msgData)
	msgType := CalcMessageIndex(msgData)
	dst = append((*byte4type)(unsafe.Pointer(&msgType))[:2], dst...)
	return
}

func package2msg(src []byte) (msgType txdata.MsgType, msgData ProtoMessage, err error) {
	// 二进制数据的前2个字节标识了后面数据的类型.
	b4 := (*byte4type)(unsafe.Pointer(&msgType))
	b4[0] = src[0]
	b4[1] = src[1]
	msgData, err = slice2msg(msgType, src[2:])
	if err == nil {
		assert4true(CalcMessageType(msgData) == msgType)
	}
	return
}

func slice2msg(msgType txdata.MsgType, src []byte) (msgData ProtoMessage, err error) {
	// 需要在shell下,先创建ff函数,再执行ff函数.
	// ff(){ sed -n '/^enum MsgType/,/}/p' "$1" | sed 's/[ \t]*\?\(ID_\)\([^ \t]\+\).*/case txdata.MsgType_\1\2: \n msgData = new(txdata.\2)/g' ; }
	// ff  txdata.proto
	switch msgType {
	case txdata.MsgType_ID_Common1Req:
		msgData = new(txdata.Common1Req)
	case txdata.MsgType_ID_Common1Rsp:
		msgData = new(txdata.Common1Rsp)
	case txdata.MsgType_ID_Common2Req:
		msgData = new(txdata.Common2Req)
	case txdata.MsgType_ID_Common2Rsp:
		msgData = new(txdata.Common2Rsp)
	case txdata.MsgType_ID_Common2Ack:
		msgData = new(txdata.Common2Ack)
	case txdata.MsgType_ID_CommonErr:
		msgData = new(txdata.CommonErr)
	case txdata.MsgType_ID_ConnectionInfo:
		msgData = new(txdata.ConnectionInfo)
	case txdata.MsgType_ID_DisconnectedData:
		msgData = new(txdata.DisconnectedData)
	case txdata.MsgType_ID_ConnectReq:
		msgData = new(txdata.ConnectReq)
	case txdata.MsgType_ID_ConnectRsp:
		msgData = new(txdata.ConnectRsp)
	case txdata.MsgType_ID_OnlineNotice:
		msgData = new(txdata.OnlineNotice)
	case txdata.MsgType_ID_SystemReport:
		msgData = new(txdata.SystemReport)
	case txdata.MsgType_ID_EchoItem:
		msgData = new(txdata.EchoItem)
	case txdata.MsgType_ID_BinItem:
		msgData = new(txdata.BinItem)
	case txdata.MsgType_ID_EmailItem:
		msgData = new(txdata.EmailItem)
	case txdata.MsgType_ID_QryConnInfoReq:
		msgData = new(txdata.QryConnInfoReq)
	case txdata.MsgType_ID_QryConnInfoRsp:
		msgData = new(txdata.QryConnInfoRsp)
	case txdata.MsgType_ID_QueryRecordReq:
		msgData = new(txdata.QueryRecordReq)
	case txdata.MsgType_ID_QueryRecordRsp:
		msgData = new(txdata.QueryRecordRsp)
	case txdata.MsgType_ID_ExecCmdReq:
		msgData = new(txdata.ExecCmdReq)
	case txdata.MsgType_ID_ExecCmdRsp:
		msgData = new(txdata.ExecCmdRsp)
	case txdata.MsgType_ID_PushWrap:
		msgData = new(txdata.PushWrap)
	case txdata.MsgType_ID_PushItem:
		msgData = new(txdata.PushItem)
	case txdata.MsgType_ID_SubscribeReq:
		msgData = new(txdata.SubscribeReq)
	case txdata.MsgType_ID_SubscribeRsp:
		msgData = new(txdata.SubscribeRsp)
	default:
		msgData = nil
		err = fmt.Errorf("unknown txdata.MsgType(%v)", msgType)
	}
	if msgData != nil {
		if err = proto.Unmarshal(src, msgData); err != nil {
			msgData = nil
			err = fmt.Errorf("Unmarshal failure(err=%v, msgType=%v)", err, msgType)
		}
	}
	return
}
