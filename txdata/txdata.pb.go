// Code generated by protoc-gen-go. DO NOT EDIT.
// source: txdata.proto

package txdata

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import timestamp "github.com/golang/protobuf/ptypes/timestamp"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Timestamp from public import google/protobuf/timestamp.proto
type Timestamp = timestamp.Timestamp

type MsgType int32

const (
	MsgType_Zero1                MsgType = 0
	MsgType_ID_ConnectedData     MsgType = 1
	MsgType_ID_DisconnectedData  MsgType = 2
	MsgType_ID_CommonNtosReq     MsgType = 31
	MsgType_ID_CommonNtosRsp     MsgType = 32
	MsgType_ID_CommonStonReq     MsgType = 33
	MsgType_ID_CommonStoaRsp     MsgType = 34
	MsgType_ID_ExecuteCommandReq MsgType = 35
	MsgType_ID_ExecuteCommandRsp MsgType = 36
)

var MsgType_name = map[int32]string{
	0:  "Zero1",
	1:  "ID_ConnectedData",
	2:  "ID_DisconnectedData",
	31: "ID_CommonNtosReq",
	32: "ID_CommonNtosRsp",
	33: "ID_CommonStonReq",
	34: "ID_CommonStoaRsp",
	35: "ID_ExecuteCommandReq",
	36: "ID_ExecuteCommandRsp",
}
var MsgType_value = map[string]int32{
	"Zero1":                0,
	"ID_ConnectedData":     1,
	"ID_DisconnectedData":  2,
	"ID_CommonNtosReq":     31,
	"ID_CommonNtosRsp":     32,
	"ID_CommonStonReq":     33,
	"ID_CommonStoaRsp":     34,
	"ID_ExecuteCommandReq": 35,
	"ID_ExecuteCommandRsp": 36,
}

func (x MsgType) String() string {
	return proto.EnumName(MsgType_name, int32(x))
}
func (MsgType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_txdata_d3a92d0928e85c9c, []int{0}
}

type ProgramType int32

const (
	ProgramType_Zero2  ProgramType = 0
	ProgramType_CLIENT ProgramType = 1
	ProgramType_SERVER ProgramType = 2
	ProgramType_NODE   ProgramType = 3
	ProgramType_POINT  ProgramType = 4
)

var ProgramType_name = map[int32]string{
	0: "Zero2",
	1: "CLIENT",
	2: "SERVER",
	3: "NODE",
	4: "POINT",
}
var ProgramType_value = map[string]int32{
	"Zero2":  0,
	"CLIENT": 1,
	"SERVER": 2,
	"NODE":   3,
	"POINT":  4,
}

func (x ProgramType) String() string {
	return proto.EnumName(ProgramType_name, int32(x))
}
func (ProgramType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_txdata_d3a92d0928e85c9c, []int{1}
}

type ConnectionInfo_LinkType int32

const (
	ConnectionInfo_Zero3   ConnectionInfo_LinkType = 0
	ConnectionInfo_CONNECT ConnectionInfo_LinkType = 1
	ConnectionInfo_ACCEPT  ConnectionInfo_LinkType = 2
)

var ConnectionInfo_LinkType_name = map[int32]string{
	0: "Zero3",
	1: "CONNECT",
	2: "ACCEPT",
}
var ConnectionInfo_LinkType_value = map[string]int32{
	"Zero3":   0,
	"CONNECT": 1,
	"ACCEPT":  2,
}

func (x ConnectionInfo_LinkType) String() string {
	return proto.EnumName(ConnectionInfo_LinkType_name, int32(x))
}
func (ConnectionInfo_LinkType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_txdata_d3a92d0928e85c9c, []int{1, 0}
}

// 必须是可见的ASCII字符,且不能为(/),这样可以拼成(/zone/node/type/name)的样子.
type AtomicKey struct {
	ZoneName             string      `protobuf:"bytes,1,opt,name=ZoneName,proto3" json:"ZoneName,omitempty"`
	NodeName             string      `protobuf:"bytes,2,opt,name=NodeName,proto3" json:"NodeName,omitempty"`
	ExecType             ProgramType `protobuf:"varint,3,opt,name=ExecType,proto3,enum=txdata.ProgramType" json:"ExecType,omitempty"`
	ExecName             string      `protobuf:"bytes,4,opt,name=ExecName,proto3" json:"ExecName,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *AtomicKey) Reset()         { *m = AtomicKey{} }
func (m *AtomicKey) String() string { return proto.CompactTextString(m) }
func (*AtomicKey) ProtoMessage()    {}
func (*AtomicKey) Descriptor() ([]byte, []int) {
	return fileDescriptor_txdata_d3a92d0928e85c9c, []int{0}
}
func (m *AtomicKey) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AtomicKey.Unmarshal(m, b)
}
func (m *AtomicKey) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AtomicKey.Marshal(b, m, deterministic)
}
func (dst *AtomicKey) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AtomicKey.Merge(dst, src)
}
func (m *AtomicKey) XXX_Size() int {
	return xxx_messageInfo_AtomicKey.Size(m)
}
func (m *AtomicKey) XXX_DiscardUnknown() {
	xxx_messageInfo_AtomicKey.DiscardUnknown(m)
}

var xxx_messageInfo_AtomicKey proto.InternalMessageInfo

func (m *AtomicKey) GetZoneName() string {
	if m != nil {
		return m.ZoneName
	}
	return ""
}

func (m *AtomicKey) GetNodeName() string {
	if m != nil {
		return m.NodeName
	}
	return ""
}

func (m *AtomicKey) GetExecType() ProgramType {
	if m != nil {
		return m.ExecType
	}
	return ProgramType_Zero2
}

func (m *AtomicKey) GetExecName() string {
	if m != nil {
		return m.ExecName
	}
	return ""
}

type ConnectionInfo struct {
	UserKey              *AtomicKey              `protobuf:"bytes,1,opt,name=UserKey,proto3" json:"UserKey,omitempty"`
	UserID               string                  `protobuf:"bytes,2,opt,name=UserID,proto3" json:"UserID,omitempty"`
	BelongKey            *AtomicKey              `protobuf:"bytes,3,opt,name=BelongKey,proto3" json:"BelongKey,omitempty"`
	BelongID             string                  `protobuf:"bytes,4,opt,name=BelongID,proto3" json:"BelongID,omitempty"`
	Version              string                  `protobuf:"bytes,5,opt,name=Version,proto3" json:"Version,omitempty"`
	LinkMode             ConnectionInfo_LinkType `protobuf:"varint,6,opt,name=LinkMode,proto3,enum=txdata.ConnectionInfo_LinkType" json:"LinkMode,omitempty"`
	ExePid               int32                   `protobuf:"varint,7,opt,name=ExePid,proto3" json:"ExePid,omitempty"`
	ExePath              string                  `protobuf:"bytes,8,opt,name=ExePath,proto3" json:"ExePath,omitempty"`
	Remark               string                  `protobuf:"bytes,9,opt,name=Remark,proto3" json:"Remark,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *ConnectionInfo) Reset()         { *m = ConnectionInfo{} }
func (m *ConnectionInfo) String() string { return proto.CompactTextString(m) }
func (*ConnectionInfo) ProtoMessage()    {}
func (*ConnectionInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_txdata_d3a92d0928e85c9c, []int{1}
}
func (m *ConnectionInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConnectionInfo.Unmarshal(m, b)
}
func (m *ConnectionInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConnectionInfo.Marshal(b, m, deterministic)
}
func (dst *ConnectionInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConnectionInfo.Merge(dst, src)
}
func (m *ConnectionInfo) XXX_Size() int {
	return xxx_messageInfo_ConnectionInfo.Size(m)
}
func (m *ConnectionInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_ConnectionInfo.DiscardUnknown(m)
}

var xxx_messageInfo_ConnectionInfo proto.InternalMessageInfo

func (m *ConnectionInfo) GetUserKey() *AtomicKey {
	if m != nil {
		return m.UserKey
	}
	return nil
}

func (m *ConnectionInfo) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

func (m *ConnectionInfo) GetBelongKey() *AtomicKey {
	if m != nil {
		return m.BelongKey
	}
	return nil
}

func (m *ConnectionInfo) GetBelongID() string {
	if m != nil {
		return m.BelongID
	}
	return ""
}

func (m *ConnectionInfo) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *ConnectionInfo) GetLinkMode() ConnectionInfo_LinkType {
	if m != nil {
		return m.LinkMode
	}
	return ConnectionInfo_Zero3
}

func (m *ConnectionInfo) GetExePid() int32 {
	if m != nil {
		return m.ExePid
	}
	return 0
}

func (m *ConnectionInfo) GetExePath() string {
	if m != nil {
		return m.ExePath
	}
	return ""
}

func (m *ConnectionInfo) GetRemark() string {
	if m != nil {
		return m.Remark
	}
	return ""
}

// 某user维护的所有连接,打包成一个结构体,作为一个快照,发送出去.
type ConnInfoSnap struct {
	UserID               string            `protobuf:"bytes,1,opt,name=UserID,proto3" json:"UserID,omitempty"`
	Data                 []*ConnectionInfo `protobuf:"bytes,2,rep,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *ConnInfoSnap) Reset()         { *m = ConnInfoSnap{} }
func (m *ConnInfoSnap) String() string { return proto.CompactTextString(m) }
func (*ConnInfoSnap) ProtoMessage()    {}
func (*ConnInfoSnap) Descriptor() ([]byte, []int) {
	return fileDescriptor_txdata_d3a92d0928e85c9c, []int{2}
}
func (m *ConnInfoSnap) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConnInfoSnap.Unmarshal(m, b)
}
func (m *ConnInfoSnap) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConnInfoSnap.Marshal(b, m, deterministic)
}
func (dst *ConnInfoSnap) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConnInfoSnap.Merge(dst, src)
}
func (m *ConnInfoSnap) XXX_Size() int {
	return xxx_messageInfo_ConnInfoSnap.Size(m)
}
func (m *ConnInfoSnap) XXX_DiscardUnknown() {
	xxx_messageInfo_ConnInfoSnap.DiscardUnknown(m)
}

var xxx_messageInfo_ConnInfoSnap proto.InternalMessageInfo

func (m *ConnInfoSnap) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

func (m *ConnInfoSnap) GetData() []*ConnectionInfo {
	if m != nil {
		return m.Data
	}
	return nil
}

type ConnectedData struct {
	Info                 *ConnectionInfo `protobuf:"bytes,1,opt,name=Info,proto3" json:"Info,omitempty"`
	Pathway              []string        `protobuf:"bytes,2,rep,name=Pathway,proto3" json:"Pathway,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *ConnectedData) Reset()         { *m = ConnectedData{} }
func (m *ConnectedData) String() string { return proto.CompactTextString(m) }
func (*ConnectedData) ProtoMessage()    {}
func (*ConnectedData) Descriptor() ([]byte, []int) {
	return fileDescriptor_txdata_d3a92d0928e85c9c, []int{3}
}
func (m *ConnectedData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConnectedData.Unmarshal(m, b)
}
func (m *ConnectedData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConnectedData.Marshal(b, m, deterministic)
}
func (dst *ConnectedData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConnectedData.Merge(dst, src)
}
func (m *ConnectedData) XXX_Size() int {
	return xxx_messageInfo_ConnectedData.Size(m)
}
func (m *ConnectedData) XXX_DiscardUnknown() {
	xxx_messageInfo_ConnectedData.DiscardUnknown(m)
}

var xxx_messageInfo_ConnectedData proto.InternalMessageInfo

func (m *ConnectedData) GetInfo() *ConnectionInfo {
	if m != nil {
		return m.Info
	}
	return nil
}

func (m *ConnectedData) GetPathway() []string {
	if m != nil {
		return m.Pathway
	}
	return nil
}

type DisconnectedData struct {
	Info                 *ConnectionInfo `protobuf:"bytes,1,opt,name=Info,proto3" json:"Info,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *DisconnectedData) Reset()         { *m = DisconnectedData{} }
func (m *DisconnectedData) String() string { return proto.CompactTextString(m) }
func (*DisconnectedData) ProtoMessage()    {}
func (*DisconnectedData) Descriptor() ([]byte, []int) {
	return fileDescriptor_txdata_d3a92d0928e85c9c, []int{4}
}
func (m *DisconnectedData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DisconnectedData.Unmarshal(m, b)
}
func (m *DisconnectedData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DisconnectedData.Marshal(b, m, deterministic)
}
func (dst *DisconnectedData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DisconnectedData.Merge(dst, src)
}
func (m *DisconnectedData) XXX_Size() int {
	return xxx_messageInfo_DisconnectedData.Size(m)
}
func (m *DisconnectedData) XXX_DiscardUnknown() {
	xxx_messageInfo_DisconnectedData.DiscardUnknown(m)
}

var xxx_messageInfo_DisconnectedData proto.InternalMessageInfo

func (m *DisconnectedData) GetInfo() *ConnectionInfo {
	if m != nil {
		return m.Info
	}
	return nil
}

type CommonNtosReq struct {
	RequestID            int64                `protobuf:"varint,1,opt,name=RequestID,proto3" json:"RequestID,omitempty"`
	UserID               string               `protobuf:"bytes,2,opt,name=UserID,proto3" json:"UserID,omitempty"`
	SeqNo                int64                `protobuf:"varint,3,opt,name=SeqNo,proto3" json:"SeqNo,omitempty"`
	Endeavour            bool                 `protobuf:"varint,4,opt,name=Endeavour,proto3" json:"Endeavour,omitempty"`
	DataType             string               `protobuf:"bytes,5,opt,name=DataType,proto3" json:"DataType,omitempty"`
	Data                 []byte               `protobuf:"bytes,6,opt,name=Data,proto3" json:"Data,omitempty"`
	ReqTime              *timestamp.Timestamp `protobuf:"bytes,7,opt,name=ReqTime,proto3" json:"ReqTime,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *CommonNtosReq) Reset()         { *m = CommonNtosReq{} }
func (m *CommonNtosReq) String() string { return proto.CompactTextString(m) }
func (*CommonNtosReq) ProtoMessage()    {}
func (*CommonNtosReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_txdata_d3a92d0928e85c9c, []int{5}
}
func (m *CommonNtosReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommonNtosReq.Unmarshal(m, b)
}
func (m *CommonNtosReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommonNtosReq.Marshal(b, m, deterministic)
}
func (dst *CommonNtosReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommonNtosReq.Merge(dst, src)
}
func (m *CommonNtosReq) XXX_Size() int {
	return xxx_messageInfo_CommonNtosReq.Size(m)
}
func (m *CommonNtosReq) XXX_DiscardUnknown() {
	xxx_messageInfo_CommonNtosReq.DiscardUnknown(m)
}

var xxx_messageInfo_CommonNtosReq proto.InternalMessageInfo

func (m *CommonNtosReq) GetRequestID() int64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *CommonNtosReq) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

func (m *CommonNtosReq) GetSeqNo() int64 {
	if m != nil {
		return m.SeqNo
	}
	return 0
}

func (m *CommonNtosReq) GetEndeavour() bool {
	if m != nil {
		return m.Endeavour
	}
	return false
}

func (m *CommonNtosReq) GetDataType() string {
	if m != nil {
		return m.DataType
	}
	return ""
}

func (m *CommonNtosReq) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *CommonNtosReq) GetReqTime() *timestamp.Timestamp {
	if m != nil {
		return m.ReqTime
	}
	return nil
}

type CommonNtosRsp struct {
	RequestID            int64    `protobuf:"varint,1,opt,name=RequestID,proto3" json:"RequestID,omitempty"`
	Pathway              []string `protobuf:"bytes,2,rep,name=Pathway,proto3" json:"Pathway,omitempty"`
	SeqNo                int64    `protobuf:"varint,3,opt,name=SeqNo,proto3" json:"SeqNo,omitempty"`
	ErrNo                int32    `protobuf:"varint,4,opt,name=ErrNo,proto3" json:"ErrNo,omitempty"`
	ErrMsg               string   `protobuf:"bytes,5,opt,name=ErrMsg,proto3" json:"ErrMsg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CommonNtosRsp) Reset()         { *m = CommonNtosRsp{} }
func (m *CommonNtosRsp) String() string { return proto.CompactTextString(m) }
func (*CommonNtosRsp) ProtoMessage()    {}
func (*CommonNtosRsp) Descriptor() ([]byte, []int) {
	return fileDescriptor_txdata_d3a92d0928e85c9c, []int{6}
}
func (m *CommonNtosRsp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommonNtosRsp.Unmarshal(m, b)
}
func (m *CommonNtosRsp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommonNtosRsp.Marshal(b, m, deterministic)
}
func (dst *CommonNtosRsp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommonNtosRsp.Merge(dst, src)
}
func (m *CommonNtosRsp) XXX_Size() int {
	return xxx_messageInfo_CommonNtosRsp.Size(m)
}
func (m *CommonNtosRsp) XXX_DiscardUnknown() {
	xxx_messageInfo_CommonNtosRsp.DiscardUnknown(m)
}

var xxx_messageInfo_CommonNtosRsp proto.InternalMessageInfo

func (m *CommonNtosRsp) GetRequestID() int64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *CommonNtosRsp) GetPathway() []string {
	if m != nil {
		return m.Pathway
	}
	return nil
}

func (m *CommonNtosRsp) GetSeqNo() int64 {
	if m != nil {
		return m.SeqNo
	}
	return 0
}

func (m *CommonNtosRsp) GetErrNo() int32 {
	if m != nil {
		return m.ErrNo
	}
	return 0
}

func (m *CommonNtosRsp) GetErrMsg() string {
	if m != nil {
		return m.ErrMsg
	}
	return ""
}

type CommonStonReq struct {
	RequestID            int64                `protobuf:"varint,1,opt,name=RequestID,proto3" json:"RequestID,omitempty"`
	Pathway              []string             `protobuf:"bytes,2,rep,name=Pathway,proto3" json:"Pathway,omitempty"`
	DataType             string               `protobuf:"bytes,3,opt,name=DataType,proto3" json:"DataType,omitempty"`
	Data                 []byte               `protobuf:"bytes,4,opt,name=Data,proto3" json:"Data,omitempty"`
	ReqTime              *timestamp.Timestamp `protobuf:"bytes,5,opt,name=ReqTime,proto3" json:"ReqTime,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *CommonStonReq) Reset()         { *m = CommonStonReq{} }
func (m *CommonStonReq) String() string { return proto.CompactTextString(m) }
func (*CommonStonReq) ProtoMessage()    {}
func (*CommonStonReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_txdata_d3a92d0928e85c9c, []int{7}
}
func (m *CommonStonReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommonStonReq.Unmarshal(m, b)
}
func (m *CommonStonReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommonStonReq.Marshal(b, m, deterministic)
}
func (dst *CommonStonReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommonStonReq.Merge(dst, src)
}
func (m *CommonStonReq) XXX_Size() int {
	return xxx_messageInfo_CommonStonReq.Size(m)
}
func (m *CommonStonReq) XXX_DiscardUnknown() {
	xxx_messageInfo_CommonStonReq.DiscardUnknown(m)
}

var xxx_messageInfo_CommonStonReq proto.InternalMessageInfo

func (m *CommonStonReq) GetRequestID() int64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *CommonStonReq) GetPathway() []string {
	if m != nil {
		return m.Pathway
	}
	return nil
}

func (m *CommonStonReq) GetDataType() string {
	if m != nil {
		return m.DataType
	}
	return ""
}

func (m *CommonStonReq) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *CommonStonReq) GetReqTime() *timestamp.Timestamp {
	if m != nil {
		return m.ReqTime
	}
	return nil
}

type CommonStonRsp struct {
	RequestID            int64                `protobuf:"varint,1,opt,name=RequestID,proto3" json:"RequestID,omitempty"`
	UserID               string               `protobuf:"bytes,2,opt,name=UserID,proto3" json:"UserID,omitempty"`
	DataType             string               `protobuf:"bytes,3,opt,name=DataType,proto3" json:"DataType,omitempty"`
	Data                 []byte               `protobuf:"bytes,4,opt,name=Data,proto3" json:"Data,omitempty"`
	RspTime              *timestamp.Timestamp `protobuf:"bytes,5,opt,name=RspTime,proto3" json:"RspTime,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *CommonStonRsp) Reset()         { *m = CommonStonRsp{} }
func (m *CommonStonRsp) String() string { return proto.CompactTextString(m) }
func (*CommonStonRsp) ProtoMessage()    {}
func (*CommonStonRsp) Descriptor() ([]byte, []int) {
	return fileDescriptor_txdata_d3a92d0928e85c9c, []int{8}
}
func (m *CommonStonRsp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommonStonRsp.Unmarshal(m, b)
}
func (m *CommonStonRsp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommonStonRsp.Marshal(b, m, deterministic)
}
func (dst *CommonStonRsp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommonStonRsp.Merge(dst, src)
}
func (m *CommonStonRsp) XXX_Size() int {
	return xxx_messageInfo_CommonStonRsp.Size(m)
}
func (m *CommonStonRsp) XXX_DiscardUnknown() {
	xxx_messageInfo_CommonStonRsp.DiscardUnknown(m)
}

var xxx_messageInfo_CommonStonRsp proto.InternalMessageInfo

func (m *CommonStonRsp) GetRequestID() int64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *CommonStonRsp) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

func (m *CommonStonRsp) GetDataType() string {
	if m != nil {
		return m.DataType
	}
	return ""
}

func (m *CommonStonRsp) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *CommonStonRsp) GetRspTime() *timestamp.Timestamp {
	if m != nil {
		return m.RspTime
	}
	return nil
}

type ExecuteCommandReq struct {
	RequestID            int64    `protobuf:"varint,1,opt,name=RequestID,proto3" json:"RequestID,omitempty"`
	Pathway              []string `protobuf:"bytes,2,rep,name=Pathway,proto3" json:"Pathway,omitempty"`
	Command              string   `protobuf:"bytes,3,opt,name=Command,proto3" json:"Command,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ExecuteCommandReq) Reset()         { *m = ExecuteCommandReq{} }
func (m *ExecuteCommandReq) String() string { return proto.CompactTextString(m) }
func (*ExecuteCommandReq) ProtoMessage()    {}
func (*ExecuteCommandReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_txdata_d3a92d0928e85c9c, []int{9}
}
func (m *ExecuteCommandReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ExecuteCommandReq.Unmarshal(m, b)
}
func (m *ExecuteCommandReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ExecuteCommandReq.Marshal(b, m, deterministic)
}
func (dst *ExecuteCommandReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExecuteCommandReq.Merge(dst, src)
}
func (m *ExecuteCommandReq) XXX_Size() int {
	return xxx_messageInfo_ExecuteCommandReq.Size(m)
}
func (m *ExecuteCommandReq) XXX_DiscardUnknown() {
	xxx_messageInfo_ExecuteCommandReq.DiscardUnknown(m)
}

var xxx_messageInfo_ExecuteCommandReq proto.InternalMessageInfo

func (m *ExecuteCommandReq) GetRequestID() int64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *ExecuteCommandReq) GetPathway() []string {
	if m != nil {
		return m.Pathway
	}
	return nil
}

func (m *ExecuteCommandReq) GetCommand() string {
	if m != nil {
		return m.Command
	}
	return ""
}

type ExecuteCommandRsp struct {
	RequestID            int64    `protobuf:"varint,1,opt,name=RequestID,proto3" json:"RequestID,omitempty"`
	UserID               string   `protobuf:"bytes,2,opt,name=UserID,proto3" json:"UserID,omitempty"`
	Result               string   `protobuf:"bytes,3,opt,name=Result,proto3" json:"Result,omitempty"`
	ErrNo                int32    `protobuf:"varint,4,opt,name=ErrNo,proto3" json:"ErrNo,omitempty"`
	ErrMsg               string   `protobuf:"bytes,5,opt,name=ErrMsg,proto3" json:"ErrMsg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ExecuteCommandRsp) Reset()         { *m = ExecuteCommandRsp{} }
func (m *ExecuteCommandRsp) String() string { return proto.CompactTextString(m) }
func (*ExecuteCommandRsp) ProtoMessage()    {}
func (*ExecuteCommandRsp) Descriptor() ([]byte, []int) {
	return fileDescriptor_txdata_d3a92d0928e85c9c, []int{10}
}
func (m *ExecuteCommandRsp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ExecuteCommandRsp.Unmarshal(m, b)
}
func (m *ExecuteCommandRsp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ExecuteCommandRsp.Marshal(b, m, deterministic)
}
func (dst *ExecuteCommandRsp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExecuteCommandRsp.Merge(dst, src)
}
func (m *ExecuteCommandRsp) XXX_Size() int {
	return xxx_messageInfo_ExecuteCommandRsp.Size(m)
}
func (m *ExecuteCommandRsp) XXX_DiscardUnknown() {
	xxx_messageInfo_ExecuteCommandRsp.DiscardUnknown(m)
}

var xxx_messageInfo_ExecuteCommandRsp proto.InternalMessageInfo

func (m *ExecuteCommandRsp) GetRequestID() int64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *ExecuteCommandRsp) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

func (m *ExecuteCommandRsp) GetResult() string {
	if m != nil {
		return m.Result
	}
	return ""
}

func (m *ExecuteCommandRsp) GetErrNo() int32 {
	if m != nil {
		return m.ErrNo
	}
	return 0
}

func (m *ExecuteCommandRsp) GetErrMsg() string {
	if m != nil {
		return m.ErrMsg
	}
	return ""
}

type ReportDataItem struct {
	Topic                string   `protobuf:"bytes,1,opt,name=Topic,proto3" json:"Topic,omitempty"`
	Data                 string   `protobuf:"bytes,2,opt,name=Data,proto3" json:"Data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ReportDataItem) Reset()         { *m = ReportDataItem{} }
func (m *ReportDataItem) String() string { return proto.CompactTextString(m) }
func (*ReportDataItem) ProtoMessage()    {}
func (*ReportDataItem) Descriptor() ([]byte, []int) {
	return fileDescriptor_txdata_d3a92d0928e85c9c, []int{11}
}
func (m *ReportDataItem) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReportDataItem.Unmarshal(m, b)
}
func (m *ReportDataItem) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReportDataItem.Marshal(b, m, deterministic)
}
func (dst *ReportDataItem) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReportDataItem.Merge(dst, src)
}
func (m *ReportDataItem) XXX_Size() int {
	return xxx_messageInfo_ReportDataItem.Size(m)
}
func (m *ReportDataItem) XXX_DiscardUnknown() {
	xxx_messageInfo_ReportDataItem.DiscardUnknown(m)
}

var xxx_messageInfo_ReportDataItem proto.InternalMessageInfo

func (m *ReportDataItem) GetTopic() string {
	if m != nil {
		return m.Topic
	}
	return ""
}

func (m *ReportDataItem) GetData() string {
	if m != nil {
		return m.Data
	}
	return ""
}

type SendMailItem struct {
	Username             string   `protobuf:"bytes,1,opt,name=Username,proto3" json:"Username,omitempty"`
	Password             string   `protobuf:"bytes,2,opt,name=Password,proto3" json:"Password,omitempty"`
	SmtpAddr             string   `protobuf:"bytes,3,opt,name=SmtpAddr,proto3" json:"SmtpAddr,omitempty"`
	To                   string   `protobuf:"bytes,4,opt,name=To,proto3" json:"To,omitempty"`
	Subject              string   `protobuf:"bytes,5,opt,name=Subject,proto3" json:"Subject,omitempty"`
	ContentType          string   `protobuf:"bytes,6,opt,name=ContentType,proto3" json:"ContentType,omitempty"`
	Content              string   `protobuf:"bytes,7,opt,name=Content,proto3" json:"Content,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SendMailItem) Reset()         { *m = SendMailItem{} }
func (m *SendMailItem) String() string { return proto.CompactTextString(m) }
func (*SendMailItem) ProtoMessage()    {}
func (*SendMailItem) Descriptor() ([]byte, []int) {
	return fileDescriptor_txdata_d3a92d0928e85c9c, []int{12}
}
func (m *SendMailItem) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SendMailItem.Unmarshal(m, b)
}
func (m *SendMailItem) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SendMailItem.Marshal(b, m, deterministic)
}
func (dst *SendMailItem) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SendMailItem.Merge(dst, src)
}
func (m *SendMailItem) XXX_Size() int {
	return xxx_messageInfo_SendMailItem.Size(m)
}
func (m *SendMailItem) XXX_DiscardUnknown() {
	xxx_messageInfo_SendMailItem.DiscardUnknown(m)
}

var xxx_messageInfo_SendMailItem proto.InternalMessageInfo

func (m *SendMailItem) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *SendMailItem) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

func (m *SendMailItem) GetSmtpAddr() string {
	if m != nil {
		return m.SmtpAddr
	}
	return ""
}

func (m *SendMailItem) GetTo() string {
	if m != nil {
		return m.To
	}
	return ""
}

func (m *SendMailItem) GetSubject() string {
	if m != nil {
		return m.Subject
	}
	return ""
}

func (m *SendMailItem) GetContentType() string {
	if m != nil {
		return m.ContentType
	}
	return ""
}

func (m *SendMailItem) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

type ServerCacheItem struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ServerCacheItem) Reset()         { *m = ServerCacheItem{} }
func (m *ServerCacheItem) String() string { return proto.CompactTextString(m) }
func (*ServerCacheItem) ProtoMessage()    {}
func (*ServerCacheItem) Descriptor() ([]byte, []int) {
	return fileDescriptor_txdata_d3a92d0928e85c9c, []int{13}
}
func (m *ServerCacheItem) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ServerCacheItem.Unmarshal(m, b)
}
func (m *ServerCacheItem) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ServerCacheItem.Marshal(b, m, deterministic)
}
func (dst *ServerCacheItem) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ServerCacheItem.Merge(dst, src)
}
func (m *ServerCacheItem) XXX_Size() int {
	return xxx_messageInfo_ServerCacheItem.Size(m)
}
func (m *ServerCacheItem) XXX_DiscardUnknown() {
	xxx_messageInfo_ServerCacheItem.DiscardUnknown(m)
}

var xxx_messageInfo_ServerCacheItem proto.InternalMessageInfo

func init() {
	proto.RegisterType((*AtomicKey)(nil), "txdata.AtomicKey")
	proto.RegisterType((*ConnectionInfo)(nil), "txdata.ConnectionInfo")
	proto.RegisterType((*ConnInfoSnap)(nil), "txdata.ConnInfoSnap")
	proto.RegisterType((*ConnectedData)(nil), "txdata.ConnectedData")
	proto.RegisterType((*DisconnectedData)(nil), "txdata.DisconnectedData")
	proto.RegisterType((*CommonNtosReq)(nil), "txdata.CommonNtosReq")
	proto.RegisterType((*CommonNtosRsp)(nil), "txdata.CommonNtosRsp")
	proto.RegisterType((*CommonStonReq)(nil), "txdata.CommonStonReq")
	proto.RegisterType((*CommonStonRsp)(nil), "txdata.CommonStonRsp")
	proto.RegisterType((*ExecuteCommandReq)(nil), "txdata.ExecuteCommandReq")
	proto.RegisterType((*ExecuteCommandRsp)(nil), "txdata.ExecuteCommandRsp")
	proto.RegisterType((*ReportDataItem)(nil), "txdata.ReportDataItem")
	proto.RegisterType((*SendMailItem)(nil), "txdata.SendMailItem")
	proto.RegisterType((*ServerCacheItem)(nil), "txdata.ServerCacheItem")
	proto.RegisterEnum("txdata.MsgType", MsgType_name, MsgType_value)
	proto.RegisterEnum("txdata.ProgramType", ProgramType_name, ProgramType_value)
	proto.RegisterEnum("txdata.ConnectionInfo_LinkType", ConnectionInfo_LinkType_name, ConnectionInfo_LinkType_value)
}

func init() { proto.RegisterFile("txdata.proto", fileDescriptor_txdata_d3a92d0928e85c9c) }

var fileDescriptor_txdata_d3a92d0928e85c9c = []byte{
	// 914 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x56, 0xcf, 0x6e, 0xdb, 0xc6,
	0x13, 0x36, 0xf5, 0x5f, 0x23, 0xff, 0xfc, 0xa3, 0x37, 0x46, 0x4a, 0x18, 0x05, 0xac, 0xb2, 0x3d,
	0x18, 0x2e, 0x20, 0xa3, 0x4a, 0x4f, 0x2d, 0x50, 0xc0, 0x95, 0x78, 0x20, 0x12, 0xd3, 0xc2, 0x8a,
	0xc9, 0x21, 0x17, 0x83, 0x16, 0x37, 0x0a, 0x1b, 0x93, 0x4b, 0x2f, 0x57, 0x89, 0xfd, 0x0c, 0x3d,
	0xb4, 0xef, 0x51, 0xf4, 0x3d, 0x7a, 0xeb, 0x23, 0xf4, 0x21, 0xfa, 0x02, 0xc5, 0xec, 0x1f, 0x59,
	0x8a, 0x2d, 0x38, 0x75, 0x6f, 0xfb, 0xcd, 0xcc, 0xce, 0x7e, 0xf3, 0x71, 0x76, 0x96, 0xb0, 0x2d,
	0xaf, 0xd3, 0x44, 0x26, 0x83, 0x52, 0x70, 0xc9, 0x49, 0x4b, 0xa3, 0xfd, 0x83, 0x39, 0xe7, 0xf3,
	0x4b, 0x76, 0xac, 0xac, 0x17, 0x8b, 0x37, 0xc7, 0x32, 0xcb, 0x59, 0x25, 0x93, 0xbc, 0xd4, 0x81,
	0xfe, 0xaf, 0x0e, 0x74, 0x4f, 0x24, 0xcf, 0xb3, 0xd9, 0x73, 0x76, 0x43, 0xf6, 0xa1, 0xf3, 0x9a,
	0x17, 0x2c, 0x4a, 0x72, 0xe6, 0x39, 0x7d, 0xe7, 0xb0, 0x4b, 0x97, 0x18, 0x7d, 0x11, 0x4f, 0xb5,
	0xaf, 0xa6, 0x7d, 0x16, 0x93, 0x63, 0xe8, 0x04, 0xd7, 0x6c, 0x16, 0xdf, 0x94, 0xcc, 0xab, 0xf7,
	0x9d, 0xc3, 0x9d, 0xe1, 0x93, 0x81, 0xe1, 0x33, 0x11, 0x7c, 0x2e, 0x92, 0x1c, 0x5d, 0x74, 0x19,
	0x84, 0xc9, 0x70, 0xad, 0x92, 0x35, 0x74, 0x32, 0x8b, 0xfd, 0xbf, 0x6b, 0xb0, 0x33, 0xe2, 0x45,
	0xc1, 0x66, 0x32, 0xe3, 0x45, 0x58, 0xbc, 0xe1, 0xe4, 0x6b, 0x68, 0xbf, 0xac, 0x98, 0x78, 0xce,
	0x6e, 0x14, 0xad, 0xde, 0x70, 0xd7, 0xa6, 0x5f, 0x72, 0xa7, 0x36, 0x82, 0x3c, 0x85, 0x16, 0x2e,
	0xc3, 0xb1, 0xa1, 0x69, 0x10, 0x39, 0x86, 0xee, 0x8f, 0xec, 0x92, 0x17, 0x73, 0x4c, 0x53, 0xdf,
	0x94, 0xe6, 0x36, 0x06, 0x49, 0x6a, 0x10, 0x8e, 0x2d, 0x49, 0x8b, 0x89, 0x07, 0xed, 0x57, 0x4c,
	0x54, 0x19, 0x2f, 0xbc, 0xa6, 0x72, 0x59, 0x48, 0xbe, 0x87, 0xce, 0x8b, 0xac, 0x78, 0x77, 0xca,
	0x53, 0xe6, 0xb5, 0x94, 0x16, 0x07, 0xf6, 0x94, 0xf5, 0xaa, 0x06, 0x18, 0xa6, 0x75, 0xb1, 0x1b,
	0x90, 0x7b, 0x70, 0xcd, 0x26, 0x59, 0xea, 0xb5, 0xfb, 0xce, 0x61, 0x93, 0x1a, 0x84, 0xc7, 0xe1,
	0x2a, 0x91, 0x6f, 0xbd, 0x8e, 0x3e, 0xce, 0x40, 0xdc, 0x41, 0x59, 0x9e, 0x88, 0x77, 0x5e, 0x57,
	0x57, 0xab, 0x91, 0x3f, 0xd0, 0x34, 0x94, 0xda, 0x5d, 0x68, 0xbe, 0x66, 0x82, 0x3f, 0x73, 0xb7,
	0x48, 0x0f, 0xda, 0xa3, 0xb3, 0x28, 0x0a, 0x46, 0xb1, 0xeb, 0x10, 0x80, 0xd6, 0xc9, 0x68, 0x14,
	0x4c, 0x62, 0xb7, 0xe6, 0x53, 0xd8, 0x46, 0x7a, 0x48, 0x6c, 0x5a, 0x24, 0xe5, 0x8a, 0x8a, 0xce,
	0x9a, 0x8a, 0x47, 0xd0, 0xc0, 0x5a, 0xbc, 0x5a, 0xbf, 0x7e, 0xd8, 0x1b, 0x3e, 0xbd, 0xbf, 0x34,
	0xaa, 0x62, 0xfc, 0x97, 0xf0, 0x3f, 0x63, 0x67, 0xe9, 0x38, 0x91, 0x09, 0x6e, 0x46, 0xb7, 0xf9,
	0x88, 0x1b, 0x37, 0xab, 0x6f, 0xee, 0x41, 0x1b, 0x0b, 0xfc, 0x90, 0xdc, 0xa8, 0xb3, 0xba, 0xd4,
	0x42, 0xff, 0x07, 0x70, 0xc7, 0x59, 0x35, 0x7b, 0x6c, 0x66, 0xff, 0x2f, 0x07, 0x79, 0xe5, 0x39,
	0x2f, 0x22, 0xc9, 0x2b, 0xca, 0xae, 0xc8, 0xe7, 0xd0, 0xa5, 0xec, 0x6a, 0xc1, 0x2a, 0x69, 0xea,
	0xad, 0xd3, 0x5b, 0xc3, 0xc6, 0x86, 0xda, 0x83, 0xe6, 0x94, 0x5d, 0x45, 0x5c, 0x35, 0x53, 0x9d,
	0x6a, 0x80, 0xb9, 0x82, 0x22, 0x65, 0xc9, 0x7b, 0xbe, 0x10, 0xaa, 0x6d, 0x3a, 0xf4, 0xd6, 0x80,
	0x3d, 0x85, 0x7c, 0xd5, 0x4d, 0xd1, 0x8d, 0xb3, 0xc4, 0x84, 0x40, 0x03, 0xd7, 0xaa, 0x6b, 0xb6,
	0xa9, 0x5a, 0x93, 0x6f, 0xa1, 0x4d, 0xd9, 0x55, 0x9c, 0xe5, 0x4c, 0x75, 0x44, 0x6f, 0xb8, 0x3f,
	0xd0, 0x57, 0x7a, 0x60, 0xaf, 0xf4, 0x20, 0xb6, 0x57, 0x9a, 0xda, 0x50, 0xff, 0xe7, 0xf5, 0x0a,
	0xab, 0xf2, 0x81, 0x0a, 0x37, 0x6a, 0xbd, 0xa1, 0xc6, 0x3d, 0x68, 0x06, 0x42, 0x44, 0x5c, 0xd5,
	0xd7, 0xa4, 0x1a, 0xa8, 0xe6, 0x15, 0xe2, 0xb4, 0x9a, 0x9b, 0xca, 0x0c, 0xf2, 0x7f, 0x5f, 0xb2,
	0x99, 0x4a, 0x5e, 0x3c, 0xac, 0xf7, 0x66, 0x36, 0xab, 0xea, 0xd5, 0x37, 0xa8, 0xd7, 0xb8, 0x5f,
	0xbd, 0xe6, 0xa7, 0xab, 0xf7, 0xdb, 0x3a, 0xdf, 0x07, 0xd5, 0xdb, 0xd4, 0x1f, 0x8f, 0x61, 0x5b,
	0x95, 0x9f, 0xcc, 0x56, 0x87, 0xfa, 0x0c, 0x76, 0x71, 0x74, 0x2e, 0x24, 0x43, 0xce, 0x49, 0x91,
	0xfe, 0x17, 0x81, 0x3d, 0x68, 0x9b, 0x2c, 0x86, 0xb1, 0x85, 0xfe, 0x2f, 0xce, 0x9d, 0x73, 0x1e,
	0x2d, 0x8c, 0x9a, 0x59, 0xd5, 0xe2, 0x52, 0x9a, 0x43, 0x0c, 0xfa, 0x97, 0x6d, 0xf5, 0x1d, 0xec,
	0x50, 0x56, 0x72, 0x21, 0x51, 0xbc, 0x50, 0xb2, 0x1c, 0xf7, 0xc7, 0xbc, 0xcc, 0x66, 0x66, 0x64,
	0x69, 0xb0, 0x94, 0x5a, 0x73, 0x50, 0x6b, 0xff, 0x0f, 0x07, 0xb6, 0xa7, 0xac, 0x48, 0x4f, 0x93,
	0xec, 0x52, 0x6d, 0xdd, 0x87, 0x0e, 0x92, 0x2b, 0x56, 0x5e, 0x3e, 0x8b, 0xd1, 0x37, 0x49, 0xaa,
	0xea, 0x03, 0x17, 0xa9, 0x7d, 0xf9, 0x2c, 0x46, 0xdf, 0x34, 0x97, 0xe5, 0x49, 0x9a, 0x0a, 0xfb,
	0x8d, 0x2d, 0x26, 0x3b, 0x50, 0x8b, 0xb9, 0x79, 0x39, 0x6a, 0xb1, 0x9a, 0x68, 0xd3, 0xc5, 0xc5,
	0x4f, 0x6c, 0x26, 0xed, 0x9b, 0x61, 0x20, 0xe9, 0x43, 0x6f, 0xc4, 0x0b, 0xc9, 0x0a, 0xa9, 0x9a,
	0xa5, 0xa5, 0xbc, 0xab, 0x26, 0xfd, 0x61, 0x14, 0x54, 0x73, 0x40, 0x7d, 0x18, 0x05, 0xfd, 0x5d,
	0xf8, 0xff, 0x94, 0x89, 0xf7, 0x4c, 0x8c, 0x92, 0xd9, 0x5b, 0x86, 0xc5, 0x1c, 0xfd, 0xe9, 0x40,
	0xfb, 0xb4, 0x9a, 0xaf, 0xce, 0xfe, 0x6f, 0xdc, 0x2d, 0xb2, 0x07, 0x6e, 0x38, 0x3e, 0x5f, 0x9b,
	0xc8, 0xae, 0x43, 0x3e, 0x83, 0x27, 0xe1, 0xf8, 0xfc, 0xe3, 0x81, 0xea, 0xd6, 0x96, 0xe1, 0x2b,
	0x83, 0xd2, 0x3d, 0xb8, 0x6b, 0xad, 0x4a, 0xb7, 0xbf, 0x66, 0x35, 0x97, 0xdc, 0xfd, 0xe2, 0x63,
	0x6b, 0x82, 0xb1, 0x3e, 0xf1, 0x60, 0x2f, 0x1c, 0x9f, 0xdf, 0xe9, 0x59, 0xf7, 0xcb, 0xfb, 0x3d,
	0x55, 0xe9, 0x7e, 0x75, 0x14, 0x40, 0x6f, 0xe5, 0x47, 0xc2, 0x16, 0x35, 0x74, 0xb7, 0xf0, 0x0d,
	0x1b, 0xbd, 0x08, 0x83, 0xc8, 0xbc, 0x67, 0xd3, 0x80, 0xbe, 0x0a, 0xa8, 0x5b, 0x23, 0x1d, 0x68,
	0x44, 0x67, 0xe3, 0xc0, 0xad, 0x63, 0xf0, 0xe4, 0x2c, 0x8c, 0x62, 0xb7, 0x31, 0xd9, 0xba, 0x68,
	0xa9, 0xab, 0xf4, 0xec, 0x9f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x2d, 0x60, 0xb7, 0xa3, 0x2f, 0x09,
	0x00, 0x00,
}
