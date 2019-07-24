// Code generated by protoc-gen-go. DO NOT EDIT.
// source: wire.proto

package wire

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type InternalTransaction struct {
	Nonce                uint64   `protobuf:"varint,1,opt,name=Nonce,proto3" json:"Nonce,omitempty"`
	Amount               uint64   `protobuf:"varint,2,opt,name=Amount,proto3" json:"Amount,omitempty"`
	Receiver             string   `protobuf:"bytes,3,opt,name=Receiver,proto3" json:"Receiver,omitempty"`
	UntilBlock           uint64   `protobuf:"varint,4,opt,name=UntilBlock,proto3" json:"UntilBlock,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InternalTransaction) Reset()         { *m = InternalTransaction{} }
func (m *InternalTransaction) String() string { return proto.CompactTextString(m) }
func (*InternalTransaction) ProtoMessage()    {}
func (*InternalTransaction) Descriptor() ([]byte, []int) {
	return fileDescriptor_f2dcdddcdf68d8e0, []int{0}
}

func (m *InternalTransaction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InternalTransaction.Unmarshal(m, b)
}
func (m *InternalTransaction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InternalTransaction.Marshal(b, m, deterministic)
}
func (m *InternalTransaction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InternalTransaction.Merge(m, src)
}
func (m *InternalTransaction) XXX_Size() int {
	return xxx_messageInfo_InternalTransaction.Size(m)
}
func (m *InternalTransaction) XXX_DiscardUnknown() {
	xxx_messageInfo_InternalTransaction.DiscardUnknown(m)
}

var xxx_messageInfo_InternalTransaction proto.InternalMessageInfo

func (m *InternalTransaction) GetNonce() uint64 {
	if m != nil {
		return m.Nonce
	}
	return 0
}

func (m *InternalTransaction) GetAmount() uint64 {
	if m != nil {
		return m.Amount
	}
	return 0
}

func (m *InternalTransaction) GetReceiver() string {
	if m != nil {
		return m.Receiver
	}
	return ""
}

func (m *InternalTransaction) GetUntilBlock() uint64 {
	if m != nil {
		return m.UntilBlock
	}
	return 0
}

type ExtTxns struct {
	List                 [][]byte `protobuf:"bytes,2,rep,name=List,proto3" json:"List,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ExtTxns) Reset()         { *m = ExtTxns{} }
func (m *ExtTxns) String() string { return proto.CompactTextString(m) }
func (*ExtTxns) ProtoMessage()    {}
func (*ExtTxns) Descriptor() ([]byte, []int) {
	return fileDescriptor_f2dcdddcdf68d8e0, []int{1}
}

func (m *ExtTxns) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ExtTxns.Unmarshal(m, b)
}
func (m *ExtTxns) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ExtTxns.Marshal(b, m, deterministic)
}
func (m *ExtTxns) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExtTxns.Merge(m, src)
}
func (m *ExtTxns) XXX_Size() int {
	return xxx_messageInfo_ExtTxns.Size(m)
}
func (m *ExtTxns) XXX_DiscardUnknown() {
	xxx_messageInfo_ExtTxns.DiscardUnknown(m)
}

var xxx_messageInfo_ExtTxns proto.InternalMessageInfo

func (m *ExtTxns) GetList() [][]byte {
	if m != nil {
		return m.List
	}
	return nil
}

type Event struct {
	SfNum                uint64                 `protobuf:"varint,1,opt,name=SfNum,proto3" json:"SfNum,omitempty"`
	Seq                  uint64                 `protobuf:"varint,2,opt,name=Seq,proto3" json:"Seq,omitempty"`
	Creator              string                 `protobuf:"bytes,3,opt,name=Creator,proto3" json:"Creator,omitempty"`
	Parents              [][]byte               `protobuf:"bytes,4,rep,name=Parents,proto3" json:"Parents,omitempty"`
	LamportTime          uint64                 `protobuf:"varint,5,opt,name=LamportTime,proto3" json:"LamportTime,omitempty"`
	InternalTransactions []*InternalTransaction `protobuf:"bytes,6,rep,name=InternalTransactions,proto3" json:"InternalTransactions,omitempty"`
	// Types that are valid to be assigned to ExternalTransactions:
	//	*Event_ExtTxnsValue
	//	*Event_ExtTxnsHash
	ExternalTransactions isEvent_ExternalTransactions `protobuf_oneof:"ExternalTransactions"`
	Sign                 []byte                       `protobuf:"bytes,9,opt,name=Sign,proto3" json:"Sign,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                     `json:"-"`
	XXX_unrecognized     []byte                       `json:"-"`
	XXX_sizecache        int32                        `json:"-"`
}

func (m *Event) Reset()         { *m = Event{} }
func (m *Event) String() string { return proto.CompactTextString(m) }
func (*Event) ProtoMessage()    {}
func (*Event) Descriptor() ([]byte, []int) {
	return fileDescriptor_f2dcdddcdf68d8e0, []int{2}
}

func (m *Event) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Event.Unmarshal(m, b)
}
func (m *Event) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Event.Marshal(b, m, deterministic)
}
func (m *Event) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Event.Merge(m, src)
}
func (m *Event) XXX_Size() int {
	return xxx_messageInfo_Event.Size(m)
}
func (m *Event) XXX_DiscardUnknown() {
	xxx_messageInfo_Event.DiscardUnknown(m)
}

var xxx_messageInfo_Event proto.InternalMessageInfo

func (m *Event) GetSfNum() uint64 {
	if m != nil {
		return m.SfNum
	}
	return 0
}

func (m *Event) GetSeq() uint64 {
	if m != nil {
		return m.Seq
	}
	return 0
}

func (m *Event) GetCreator() string {
	if m != nil {
		return m.Creator
	}
	return ""
}

func (m *Event) GetParents() [][]byte {
	if m != nil {
		return m.Parents
	}
	return nil
}

func (m *Event) GetLamportTime() uint64 {
	if m != nil {
		return m.LamportTime
	}
	return 0
}

func (m *Event) GetInternalTransactions() []*InternalTransaction {
	if m != nil {
		return m.InternalTransactions
	}
	return nil
}

type isEvent_ExternalTransactions interface {
	isEvent_ExternalTransactions()
}

type Event_ExtTxnsValue struct {
	ExtTxnsValue *ExtTxns `protobuf:"bytes,7,opt,name=ExtTxnsValue,proto3,oneof"`
}

type Event_ExtTxnsHash struct {
	ExtTxnsHash []byte `protobuf:"bytes,8,opt,name=ExtTxnsHash,proto3,oneof"`
}

func (*Event_ExtTxnsValue) isEvent_ExternalTransactions() {}

func (*Event_ExtTxnsHash) isEvent_ExternalTransactions() {}

func (m *Event) GetExternalTransactions() isEvent_ExternalTransactions {
	if m != nil {
		return m.ExternalTransactions
	}
	return nil
}

func (m *Event) GetExtTxnsValue() *ExtTxns {
	if x, ok := m.GetExternalTransactions().(*Event_ExtTxnsValue); ok {
		return x.ExtTxnsValue
	}
	return nil
}

func (m *Event) GetExtTxnsHash() []byte {
	if x, ok := m.GetExternalTransactions().(*Event_ExtTxnsHash); ok {
		return x.ExtTxnsHash
	}
	return nil
}

func (m *Event) GetSign() []byte {
	if m != nil {
		return m.Sign
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*Event) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*Event_ExtTxnsValue)(nil),
		(*Event_ExtTxnsHash)(nil),
	}
}

type Block struct {
	Index                uint64   `protobuf:"varint,1,opt,name=Index,proto3" json:"Index,omitempty"`
	Events               [][]byte `protobuf:"bytes,2,rep,name=Events,proto3" json:"Events,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Block) Reset()         { *m = Block{} }
func (m *Block) String() string { return proto.CompactTextString(m) }
func (*Block) ProtoMessage()    {}
func (*Block) Descriptor() ([]byte, []int) {
	return fileDescriptor_f2dcdddcdf68d8e0, []int{3}
}

func (m *Block) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Block.Unmarshal(m, b)
}
func (m *Block) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Block.Marshal(b, m, deterministic)
}
func (m *Block) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Block.Merge(m, src)
}
func (m *Block) XXX_Size() int {
	return xxx_messageInfo_Block.Size(m)
}
func (m *Block) XXX_DiscardUnknown() {
	xxx_messageInfo_Block.DiscardUnknown(m)
}

var xxx_messageInfo_Block proto.InternalMessageInfo

func (m *Block) GetIndex() uint64 {
	if m != nil {
		return m.Index
	}
	return 0
}

func (m *Block) GetEvents() [][]byte {
	if m != nil {
		return m.Events
	}
	return nil
}

func init() {
	proto.RegisterType((*InternalTransaction)(nil), "wire.InternalTransaction")
	proto.RegisterType((*ExtTxns)(nil), "wire.ExtTxns")
	proto.RegisterType((*Event)(nil), "wire.Event")
	proto.RegisterType((*Block)(nil), "wire.Block")
}

func init() { proto.RegisterFile("wire.proto", fileDescriptor_f2dcdddcdf68d8e0) }

var fileDescriptor_f2dcdddcdf68d8e0 = []byte{
	// 395 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x92, 0xd1, 0xaf, 0xd2, 0x30,
	0x14, 0xc6, 0xef, 0xee, 0x36, 0xb8, 0xf7, 0x0c, 0x13, 0x53, 0x09, 0xa9, 0x26, 0x9a, 0x65, 0x4f,
	0x7b, 0x81, 0x25, 0x10, 0x63, 0x7c, 0x14, 0x83, 0x81, 0x04, 0x89, 0x29, 0xe8, 0x83, 0x6f, 0x65,
	0x14, 0x68, 0xdc, 0x5a, 0x6c, 0x3b, 0xdc, 0x9b, 0x7f, 0xa7, 0xff, 0x8d, 0x69, 0xd9, 0x0c, 0xe6,
	0xf2, 0x76, 0x7e, 0xe7, 0xb4, 0xdf, 0xbe, 0x7e, 0x67, 0x00, 0xbf, 0xb8, 0x62, 0xa3, 0x93, 0x92,
	0x46, 0xa2, 0xc0, 0xd6, 0xc9, 0x6f, 0x78, 0xb1, 0x10, 0x86, 0x29, 0x41, 0x8b, 0x8d, 0xa2, 0x42,
	0xd3, 0xdc, 0x70, 0x29, 0x50, 0x1f, 0xc2, 0x95, 0x14, 0x39, 0xc3, 0x5e, 0xec, 0xa5, 0x01, 0xb9,
	0x00, 0x1a, 0x40, 0xe7, 0x43, 0x29, 0x2b, 0x61, 0xf0, 0xbd, 0x6b, 0x37, 0x84, 0x5e, 0xc1, 0x03,
	0x61, 0x39, 0xe3, 0x67, 0xa6, 0xb0, 0x1f, 0x7b, 0xe9, 0x23, 0xf9, 0xc7, 0xe8, 0x0d, 0xc0, 0x57,
	0x61, 0x78, 0x31, 0x2d, 0x64, 0xfe, 0x03, 0x07, 0xee, 0xde, 0x55, 0x27, 0x79, 0x0d, 0xdd, 0x59,
	0x6d, 0x36, 0xb5, 0xd0, 0x08, 0x41, 0xb0, 0xe4, 0xda, 0x8a, 0xfb, 0x69, 0x8f, 0xb8, 0x3a, 0xf9,
	0x73, 0x0f, 0xe1, 0xec, 0xcc, 0x84, 0xb1, 0x96, 0xd6, 0xfb, 0x55, 0x55, 0xb6, 0x96, 0x1c, 0xa0,
	0xe7, 0xe0, 0xaf, 0xd9, 0xcf, 0xc6, 0x8f, 0x2d, 0x11, 0x86, 0xee, 0x47, 0xc5, 0xa8, 0x91, 0xad,
	0x97, 0x16, 0xed, 0xe4, 0x0b, 0x55, 0x4c, 0x18, 0x8d, 0x03, 0xf7, 0x89, 0x16, 0x51, 0x0c, 0xd1,
	0x92, 0x96, 0x27, 0xa9, 0xcc, 0x86, 0x97, 0x0c, 0x87, 0x4e, 0xed, 0xba, 0x85, 0x3e, 0x43, 0xff,
	0x46, 0x4e, 0x1a, 0x77, 0x62, 0x3f, 0x8d, 0xc6, 0x2f, 0x47, 0x2e, 0xd8, 0x1b, 0x27, 0xc8, 0xcd,
	0x6b, 0x68, 0x02, 0xbd, 0xe6, 0xd5, 0xdf, 0x68, 0x51, 0x31, 0xdc, 0x8d, 0xbd, 0x34, 0x1a, 0x3f,
	0xbb, 0xc8, 0x34, 0x93, 0xf9, 0x1d, 0xf9, 0xef, 0x10, 0x4a, 0x20, 0x6a, 0x47, 0x54, 0x1f, 0xf1,
	0x43, 0xec, 0xa5, 0xbd, 0xf9, 0x1d, 0xb9, 0x6e, 0xda, 0x0c, 0xd7, 0xfc, 0x20, 0xf0, 0xa3, 0x1d,
	0x12, 0x57, 0x4f, 0x07, 0xd0, 0x9f, 0xd5, 0x4f, 0x4d, 0x24, 0x6f, 0x21, 0x74, 0x3b, 0xb0, 0xd1,
	0x2e, 0xc4, 0x8e, 0xd5, 0x6d, 0xb4, 0x0e, 0xec, 0xb6, 0x5d, 0xf2, 0xba, 0x59, 0x48, 0x43, 0xd3,
	0xf7, 0xdf, 0xdf, 0x1d, 0xb8, 0x39, 0x56, 0xdb, 0x51, 0x2e, 0xcb, 0xec, 0x13, 0x15, 0x46, 0x96,
	0xc3, 0xbd, 0xac, 0xc4, 0x8e, 0x5a, 0xdd, 0xec, 0x20, 0x87, 0x05, 0xcd, 0x8f, 0x4c, 0x73, 0x9d,
	0x69, 0x95, 0x67, 0xdc, 0x06, 0x90, 0xd9, 0x97, 0x6d, 0x3b, 0xee, 0xd7, 0x9b, 0xfc, 0x0d, 0x00,
	0x00, 0xff, 0xff, 0xb1, 0xa0, 0x2e, 0xc4, 0x88, 0x02, 0x00, 0x00,
}
