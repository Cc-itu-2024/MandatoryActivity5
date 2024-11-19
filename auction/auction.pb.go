// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.28.2
// source: auction/auction.proto

package gRPC

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Amount struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Amount int32 `protobuf:"varint,1,opt,name=amount,proto3" json:"amount,omitempty"`
}

func (x *Amount) Reset() {
	*x = Amount{}
	mi := &file_auction_auction_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Amount) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Amount) ProtoMessage() {}

func (x *Amount) ProtoReflect() protoreflect.Message {
	mi := &file_auction_auction_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Amount.ProtoReflect.Descriptor instead.
func (*Amount) Descriptor() ([]byte, []int) {
	return file_auction_auction_proto_rawDescGZIP(), []int{0}
}

func (x *Amount) GetAmount() int32 {
	if x != nil {
		return x.Amount
	}
	return 0
}

type Outcome struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Outcome int32 `protobuf:"varint,1,opt,name=outcome,proto3" json:"outcome,omitempty"`
}

func (x *Outcome) Reset() {
	*x = Outcome{}
	mi := &file_auction_auction_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Outcome) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Outcome) ProtoMessage() {}

func (x *Outcome) ProtoReflect() protoreflect.Message {
	mi := &file_auction_auction_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Outcome.ProtoReflect.Descriptor instead.
func (*Outcome) Descriptor() ([]byte, []int) {
	return file_auction_auction_proto_rawDescGZIP(), []int{1}
}

func (x *Outcome) GetOutcome() int32 {
	if x != nil {
		return x.Outcome
	}
	return 0
}

type Ack struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *Ack) Reset() {
	*x = Ack{}
	mi := &file_auction_auction_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Ack) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Ack) ProtoMessage() {}

func (x *Ack) ProtoReflect() protoreflect.Message {
	mi := &file_auction_auction_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Ack.ProtoReflect.Descriptor instead.
func (*Ack) Descriptor() ([]byte, []int) {
	return file_auction_auction_proto_rawDescGZIP(), []int{2}
}

func (x *Ack) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type Void struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Void) Reset() {
	*x = Void{}
	mi := &file_auction_auction_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Void) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Void) ProtoMessage() {}

func (x *Void) ProtoReflect() protoreflect.Message {
	mi := &file_auction_auction_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Void.ProtoReflect.Descriptor instead.
func (*Void) Descriptor() ([]byte, []int) {
	return file_auction_auction_proto_rawDescGZIP(), []int{3}
}

var File_auction_auction_proto protoreflect.FileDescriptor

var file_auction_auction_proto_rawDesc = []byte{
	0x0a, 0x15, 0x61, 0x75, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x61, 0x75, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x67, 0x52, 0x50, 0x43, 0x1a, 0x1f, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x20,
	0x0a, 0x06, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x6d, 0x6f, 0x75,
	0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74,
	0x22, 0x23, 0x0a, 0x07, 0x4f, 0x75, 0x74, 0x63, 0x6f, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6f,
	0x75, 0x74, 0x63, 0x6f, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x6f, 0x75,
	0x74, 0x63, 0x6f, 0x6d, 0x65, 0x22, 0x1f, 0x0a, 0x03, 0x41, 0x63, 0x6b, 0x12, 0x18, 0x0a, 0x07,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x06, 0x0a, 0x04, 0x56, 0x6f, 0x69, 0x64, 0x32, 0x4e,
	0x0a, 0x07, 0x41, 0x75, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1e, 0x0a, 0x03, 0x42, 0x69, 0x64,
	0x12, 0x0c, 0x2e, 0x67, 0x52, 0x50, 0x43, 0x2e, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x1a, 0x09,
	0x2e, 0x67, 0x52, 0x50, 0x43, 0x2e, 0x41, 0x63, 0x6b, 0x12, 0x23, 0x0a, 0x06, 0x52, 0x65, 0x73,
	0x75, 0x6c, 0x74, 0x12, 0x0a, 0x2e, 0x67, 0x52, 0x50, 0x43, 0x2e, 0x56, 0x6f, 0x69, 0x64, 0x1a,
	0x0d, 0x2e, 0x67, 0x52, 0x50, 0x43, 0x2e, 0x4f, 0x75, 0x74, 0x63, 0x6f, 0x6d, 0x65, 0x42, 0x07,
	0x5a, 0x05, 0x2f, 0x67, 0x52, 0x50, 0x43, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_auction_auction_proto_rawDescOnce sync.Once
	file_auction_auction_proto_rawDescData = file_auction_auction_proto_rawDesc
)

func file_auction_auction_proto_rawDescGZIP() []byte {
	file_auction_auction_proto_rawDescOnce.Do(func() {
		file_auction_auction_proto_rawDescData = protoimpl.X.CompressGZIP(file_auction_auction_proto_rawDescData)
	})
	return file_auction_auction_proto_rawDescData
}

var file_auction_auction_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_auction_auction_proto_goTypes = []any{
	(*Amount)(nil),  // 0: gRPC.Amount
	(*Outcome)(nil), // 1: gRPC.Outcome
	(*Ack)(nil),     // 2: gRPC.Ack
	(*Void)(nil),    // 3: gRPC.Void
}
var file_auction_auction_proto_depIdxs = []int32{
	0, // 0: gRPC.Auction.Bid:input_type -> gRPC.Amount
	3, // 1: gRPC.Auction.Result:input_type -> gRPC.Void
	2, // 2: gRPC.Auction.Bid:output_type -> gRPC.Ack
	1, // 3: gRPC.Auction.Result:output_type -> gRPC.Outcome
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_auction_auction_proto_init() }
func file_auction_auction_proto_init() {
	if File_auction_auction_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_auction_auction_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_auction_auction_proto_goTypes,
		DependencyIndexes: file_auction_auction_proto_depIdxs,
		MessageInfos:      file_auction_auction_proto_msgTypes,
	}.Build()
	File_auction_auction_proto = out.File
	file_auction_auction_proto_rawDesc = nil
	file_auction_auction_proto_goTypes = nil
	file_auction_auction_proto_depIdxs = nil
}