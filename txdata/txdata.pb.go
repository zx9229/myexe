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
	return fileDescriptor_txdata_259ee25c9cd68073, []int{0}
}

type ConnectionInfo_AppType int32

const (
	ConnectionInfo_Zero2  ConnectionInfo_AppType = 0
	ConnectionInfo_SERVER ConnectionInfo_AppType = 1
	ConnectionInfo_NODE   ConnectionInfo_AppType = 2
	ConnectionInfo_CLIENT ConnectionInfo_AppType = 3
)

var ConnectionInfo_AppType_name = map[int32]string{
	0: "Zero2",
	1: "SERVER",
	2: "NODE",
	3: "CLIENT",
}
var ConnectionInfo_AppType_value = map[string]int32{
	"Zero2":  0,
	"SERVER": 1,
	"NODE":   2,
	"CLIENT": 3,
}

func (x ConnectionInfo_AppType) String() string {
	return proto.EnumName(ConnectionInfo_AppType_name, int32(x))
}
func (ConnectionInfo_AppType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_txdata_259ee25c9cd68073, []int{0, 0}
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
	return fileDescriptor_txdata_259ee25c9cd68073, []int{0, 1}
}

type ConnectionInfo struct {
	UniqueID             string                  `protobuf:"bytes,1,opt,name=UniqueID,proto3" json:"UniqueID,omitempty"`
	BelongID             string                  `protobuf:"bytes,2,opt,name=BelongID,proto3" json:"BelongID,omitempty"`
	Version              string                  `protobuf:"bytes,3,opt,name=Version,proto3" json:"Version,omitempty"`
	ExeType              ConnectionInfo_AppType  `protobuf:"varint,4,opt,name=ExeType,proto3,enum=txdata.ConnectionInfo_AppType" json:"ExeType,omitempty"`
	LinkDir              ConnectionInfo_LinkType `protobuf:"varint,5,opt,name=LinkDir,proto3,enum=txdata.ConnectionInfo_LinkType" json:"LinkDir,omitempty"`
	ExePid               int32                   `protobuf:"varint,6,opt,name=ExePid,proto3" json:"ExePid,omitempty"`
	ExePath              string                  `protobuf:"bytes,7,opt,name=ExePath,proto3" json:"ExePath,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *ConnectionInfo) Reset()         { *m = ConnectionInfo{} }
func (m *ConnectionInfo) String() string { return proto.CompactTextString(m) }
func (*ConnectionInfo) ProtoMessage()    {}
func (*ConnectionInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_txdata_259ee25c9cd68073, []int{0}
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

func (m *ConnectionInfo) GetUniqueID() string {
	if m != nil {
		return m.UniqueID
	}
	return ""
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

func (m *ConnectionInfo) GetExeType() ConnectionInfo_AppType {
	if m != nil {
		return m.ExeType
	}
	return ConnectionInfo_Zero2
}

func (m *ConnectionInfo) GetLinkDir() ConnectionInfo_LinkType {
	if m != nil {
		return m.LinkDir
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
	return fileDescriptor_txdata_259ee25c9cd68073, []int{1}
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
	return fileDescriptor_txdata_259ee25c9cd68073, []int{2}
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
	return fileDescriptor_txdata_259ee25c9cd68073, []int{3}
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
	UniqueID             string   `protobuf:"bytes,2,opt,name=UniqueID,proto3" json:"UniqueID,omitempty"`
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
	return fileDescriptor_txdata_259ee25c9cd68073, []int{4}
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

func (m *ExecuteCommandRsp) GetUniqueID() string {
	if m != nil {
		return m.UniqueID
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

type CommonNtosReq struct {
	RequestID            int64                `protobuf:"varint,1,opt,name=RequestID,proto3" json:"RequestID,omitempty"`
	UniqueID             string               `protobuf:"bytes,2,opt,name=UniqueID,proto3" json:"UniqueID,omitempty"`
	SeqNo                int64                `protobuf:"varint,3,opt,name=SeqNo,proto3" json:"SeqNo,omitempty"`
	Endeavour            bool                 `protobuf:"varint,4,opt,name=Endeavour,proto3" json:"Endeavour,omitempty"`
	DataType             string               `protobuf:"bytes,5,opt,name=DataType,proto3" json:"DataType,omitempty"`
	Data                 []byte               `protobuf:"bytes,6,opt,name=Data,proto3" json:"Data,omitempty"`
	ReportTime           *timestamp.Timestamp `protobuf:"bytes,7,opt,name=ReportTime,proto3" json:"ReportTime,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *CommonNtosReq) Reset()         { *m = CommonNtosReq{} }
func (m *CommonNtosReq) String() string { return proto.CompactTextString(m) }
func (*CommonNtosReq) ProtoMessage()    {}
func (*CommonNtosReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_txdata_259ee25c9cd68073, []int{5}
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

func (m *CommonNtosReq) GetUniqueID() string {
	if m != nil {
		return m.UniqueID
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

func (m *CommonNtosReq) GetReportTime() *timestamp.Timestamp {
	if m != nil {
		return m.ReportTime
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
	return fileDescriptor_txdata_259ee25c9cd68073, []int{6}
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
	return fileDescriptor_txdata_259ee25c9cd68073, []int{7}
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
	return fileDescriptor_txdata_259ee25c9cd68073, []int{8}
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

func init() {
	proto.RegisterType((*ConnectionInfo)(nil), "txdata.ConnectionInfo")
	proto.RegisterType((*ConnectedData)(nil), "txdata.ConnectedData")
	proto.RegisterType((*DisconnectedData)(nil), "txdata.DisconnectedData")
	proto.RegisterType((*ExecuteCommandReq)(nil), "txdata.ExecuteCommandReq")
	proto.RegisterType((*ExecuteCommandRsp)(nil), "txdata.ExecuteCommandRsp")
	proto.RegisterType((*CommonNtosReq)(nil), "txdata.CommonNtosReq")
	proto.RegisterType((*CommonNtosRsp)(nil), "txdata.CommonNtosRsp")
	proto.RegisterType((*ReportDataItem)(nil), "txdata.ReportDataItem")
	proto.RegisterType((*SendMailItem)(nil), "txdata.SendMailItem")
	proto.RegisterEnum("txdata.MsgType", MsgType_name, MsgType_value)
	proto.RegisterEnum("txdata.ConnectionInfo_AppType", ConnectionInfo_AppType_name, ConnectionInfo_AppType_value)
	proto.RegisterEnum("txdata.ConnectionInfo_LinkType", ConnectionInfo_LinkType_name, ConnectionInfo_LinkType_value)
}

func init() { proto.RegisterFile("txdata.proto", fileDescriptor_txdata_259ee25c9cd68073) }

var fileDescriptor_txdata_259ee25c9cd68073 = []byte{
	// 747 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x54, 0xcd, 0x6e, 0xe3, 0x36,
	0x10, 0xb6, 0xe4, 0xff, 0x71, 0x6a, 0xb0, 0xac, 0x91, 0x0a, 0x41, 0xd1, 0xb8, 0x6a, 0x0f, 0x41,
	0x0e, 0x0e, 0xea, 0x00, 0x45, 0x9b, 0x43, 0x81, 0x54, 0xd2, 0x41, 0x40, 0xe2, 0x18, 0xb4, 0x92,
	0x43, 0x2f, 0x81, 0x62, 0x31, 0xae, 0x5a, 0x8b, 0x94, 0x25, 0xba, 0x71, 0x9e, 0x61, 0x8f, 0xfb,
	0x52, 0x7b, 0xdb, 0xe7, 0xd8, 0xe3, 0xbe, 0xc1, 0x82, 0x14, 0xe5, 0xd8, 0xf9, 0xd9, 0x60, 0xf7,
	0xc6, 0x6f, 0x7e, 0xc8, 0x6f, 0x38, 0xdf, 0x0c, 0xec, 0x88, 0x55, 0x14, 0x8a, 0x70, 0x90, 0x66,
	0x5c, 0x70, 0xdc, 0x28, 0xd0, 0xde, 0xfe, 0x8c, 0xf3, 0xd9, 0x9c, 0x1e, 0x29, 0xeb, 0xcd, 0xf2,
	0xf6, 0x48, 0xc4, 0x09, 0xcd, 0x45, 0x98, 0xa4, 0x45, 0xa0, 0xfd, 0xd1, 0x84, 0xae, 0xc3, 0x19,
	0xa3, 0x53, 0x11, 0x73, 0xe6, 0xb3, 0x5b, 0x8e, 0xf7, 0xa0, 0x75, 0xc9, 0xe2, 0xc5, 0x92, 0xfa,
	0xae, 0x65, 0xf4, 0x8d, 0x83, 0x36, 0x59, 0x63, 0xe9, 0xfb, 0x8b, 0xce, 0x39, 0x9b, 0xf9, 0xae,
	0x65, 0x16, 0xbe, 0x12, 0x63, 0x0b, 0x9a, 0x57, 0x34, 0xcb, 0x63, 0xce, 0xac, 0xaa, 0x72, 0x95,
	0x10, 0xff, 0x0e, 0x4d, 0x6f, 0x45, 0x83, 0xfb, 0x94, 0x5a, 0xb5, 0xbe, 0x71, 0xd0, 0x1d, 0xfe,
	0x38, 0xd0, 0x6c, 0xb7, 0x9f, 0x1e, 0x9c, 0xa6, 0xa9, 0x8c, 0x22, 0x65, 0x38, 0xfe, 0x03, 0x9a,
	0x67, 0x31, 0xfb, 0xcf, 0x8d, 0x33, 0xab, 0xae, 0x32, 0xf7, 0x5f, 0xc8, 0x94, 0x51, 0x45, 0xaa,
	0x8e, 0xc7, 0xbb, 0xd0, 0xf0, 0x56, 0x74, 0x1c, 0x47, 0x56, 0xa3, 0x6f, 0x1c, 0xd4, 0x89, 0x46,
	0x92, 0xa6, 0x3c, 0x85, 0xe2, 0x1f, 0xab, 0x59, 0xd0, 0xd4, 0xd0, 0xfe, 0x0d, 0x9a, 0x9a, 0x00,
	0x6e, 0x43, 0xfd, 0x6f, 0x9a, 0xf1, 0x21, 0xaa, 0x60, 0x80, 0xc6, 0xc4, 0x23, 0x57, 0x1e, 0x41,
	0x06, 0x6e, 0x41, 0x6d, 0x74, 0xe1, 0x7a, 0xc8, 0x94, 0x56, 0xe7, 0xcc, 0xf7, 0x46, 0x01, 0xaa,
	0xda, 0x03, 0x68, 0x95, 0xcf, 0x97, 0x89, 0xc7, 0xa8, 0x82, 0x3b, 0xd0, 0x74, 0x2e, 0x46, 0x23,
	0xcf, 0x09, 0x90, 0x21, 0xe3, 0x4f, 0x1d, 0xc7, 0x1b, 0x07, 0xc8, 0xb4, 0x2f, 0xe1, 0x1b, 0xcd,
	0x9e, 0x46, 0x6e, 0x28, 0x42, 0x7c, 0x08, 0x35, 0x59, 0x84, 0xfa, 0xed, 0xce, 0x70, 0xf7, 0xf9,
	0x12, 0x89, 0x8a, 0x91, 0xf4, 0x25, 0xd9, 0xbb, 0xf0, 0xde, 0x32, 0xfb, 0x55, 0x49, 0x5f, 0x43,
	0xfb, 0x4f, 0x40, 0x6e, 0x9c, 0x4f, 0xbf, 0xf6, 0x66, 0x9b, 0xc2, 0xb7, 0xde, 0x8a, 0x4e, 0x97,
	0x82, 0x3a, 0x3c, 0x49, 0x42, 0x16, 0x11, 0xba, 0xc0, 0x3f, 0x40, 0x9b, 0xd0, 0xc5, 0x92, 0xe6,
	0x42, 0xab, 0xa1, 0x4a, 0x1e, 0x0c, 0x2f, 0x93, 0x91, 0x1e, 0x7d, 0x4b, 0x29, 0x06, 0x0d, 0xed,
	0xb7, 0xc6, 0x93, 0x77, 0xf2, 0xf4, 0x95, 0x77, 0x36, 0x25, 0x69, 0x3e, 0x92, 0xe4, 0x2e, 0x34,
	0x08, 0xcd, 0x97, 0x73, 0xa1, 0x1f, 0xd2, 0x08, 0xf7, 0xa0, 0xee, 0x65, 0xd9, 0x88, 0x2b, 0xc9,
	0xd5, 0x49, 0x01, 0x94, 0x2a, 0xb2, 0xec, 0x3c, 0x9f, 0x29, 0x3d, 0xb5, 0x89, 0x46, 0xf6, 0x07,
	0x43, 0x36, 0x25, 0x49, 0x38, 0x1b, 0x09, 0x9e, 0xbf, 0x5e, 0xf9, 0xe7, 0x18, 0xf5, 0xa0, 0x3e,
	0xa1, 0x8b, 0x11, 0x57, 0x84, 0xaa, 0xa4, 0x00, 0xf2, 0x3e, 0x8f, 0x45, 0x34, 0xfc, 0x9f, 0x2f,
	0x33, 0xc5, 0xa9, 0x45, 0x1e, 0x0c, 0xf2, 0x3e, 0xd9, 0x30, 0x35, 0x23, 0x05, 0xb3, 0x35, 0xc6,
	0x18, 0x6a, 0xf2, 0xac, 0x74, 0xbc, 0x43, 0xd4, 0x19, 0x9f, 0x00, 0x10, 0x9a, 0xf2, 0x4c, 0x04,
	0x71, 0x42, 0x95, 0x90, 0x3b, 0xc3, 0xbd, 0x41, 0x31, 0xed, 0x83, 0x72, 0xda, 0x07, 0x41, 0x39,
	0xed, 0x64, 0x23, 0xda, 0x7e, 0xb3, 0x5d, 0xeb, 0xab, 0xbf, 0xff, 0x72, 0x97, 0x9f, 0xaf, 0xf4,
	0xcb, 0x7e, 0xfe, 0x04, 0xba, 0x05, 0x37, 0x59, 0x97, 0x2f, 0x68, 0x22, 0xf3, 0x03, 0x9e, 0xc6,
	0x53, 0xbd, 0x7d, 0x0a, 0xb0, 0xfe, 0x85, 0xe2, 0xb7, 0xd5, 0xd9, 0x7e, 0x67, 0xc0, 0xce, 0x84,
	0xb2, 0xe8, 0x3c, 0x8c, 0xe7, 0x2a, 0x55, 0xb6, 0x25, 0xa7, 0x19, 0x0b, 0x13, 0xba, 0xde, 0x5d,
	0x1a, 0x4b, 0xdf, 0x38, 0xcc, 0xf3, 0x3b, 0x9e, 0x45, 0x65, 0xcb, 0x4a, 0x2c, 0x7d, 0x93, 0x44,
	0xa4, 0xa7, 0x51, 0x94, 0x69, 0x19, 0xad, 0x31, 0xee, 0x82, 0x19, 0x14, 0xb5, 0xb4, 0x89, 0x19,
	0xa8, 0x09, 0x9c, 0x2c, 0x6f, 0xfe, 0xa5, 0x53, 0xa1, 0x2b, 0x29, 0x21, 0xee, 0x43, 0xc7, 0xe1,
	0x4c, 0x50, 0x26, 0x54, 0x1f, 0x1b, 0xca, 0xbb, 0x69, 0x2a, 0xc6, 0x42, 0xc1, 0x72, 0xf9, 0x68,
	0x78, 0xf8, 0xde, 0x80, 0xe6, 0x79, 0x3e, 0xdb, 0x5c, 0x22, 0xbf, 0xa2, 0x0a, 0xee, 0x01, 0xf2,
	0xdd, 0xeb, 0xad, 0x75, 0x81, 0x0c, 0xfc, 0x3d, 0x7c, 0xe7, 0xbb, 0xd7, 0x8f, 0xa7, 0x1d, 0x99,
	0xeb, 0xf0, 0x0d, 0x21, 0xa3, 0xfd, 0xa7, 0xd6, 0x3c, 0x45, 0xfd, 0x2d, 0xeb, 0x44, 0x70, 0x26,
	0x63, 0x7f, 0x7a, 0x6c, 0x0d, 0x65, 0xac, 0x8d, 0x2d, 0xe8, 0xf9, 0xee, 0xf5, 0x93, 0xf5, 0x80,
	0x7e, 0x7e, 0xde, 0x93, 0xa7, 0xe8, 0x97, 0x71, 0xe5, 0xa6, 0xa1, 0x84, 0x78, 0xfc, 0x29, 0x00,
	0x00, 0xff, 0xff, 0x5f, 0xf9, 0x13, 0xbc, 0x9c, 0x06, 0x00, 0x00,
}
