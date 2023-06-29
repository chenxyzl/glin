// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.5.1
// source: outer/pack.proto

package outer

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	code "laiya/proto/code"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

//请求
type ReqPack struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sn    uint64 `protobuf:"varint,1,opt,name=sn,proto3" json:"sn,omitempty"`       //流水号 request/reply
	RpcId int32  `protobuf:"varint,2,opt,name=rpcId,proto3" json:"rpcId,omitempty"` //rpcId all
	Data  []byte `protobuf:"bytes,4,opt,name=data,proto3" json:"data,omitempty"`    //消息体 all
}

func (x *ReqPack) Reset() {
	*x = ReqPack{}
	if protoimpl.UnsafeEnabled {
		mi := &file_outer_pack_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReqPack) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReqPack) ProtoMessage() {}

func (x *ReqPack) ProtoReflect() protoreflect.Message {
	mi := &file_outer_pack_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReqPack.ProtoReflect.Descriptor instead.
func (*ReqPack) Descriptor() ([]byte, []int) {
	return file_outer_pack_proto_rawDescGZIP(), []int{0}
}

func (x *ReqPack) GetSn() uint64 {
	if x != nil {
		return x.Sn
	}
	return 0
}

func (x *ReqPack) GetRpcId() int32 {
	if x != nil {
		return x.RpcId
	}
	return 0
}

func (x *ReqPack) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

//请求--json版本
type ReqJsonPack struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sn    uint64 `protobuf:"varint,1,opt,name=sn,proto3" json:"sn,omitempty"`       //流水号 request/reply
	RpcId int32  `protobuf:"varint,2,opt,name=rpcId,proto3" json:"rpcId,omitempty"` //rpcId all
	Data  string `protobuf:"bytes,4,opt,name=data,proto3" json:"data,omitempty"`    //消息体 all
}

func (x *ReqJsonPack) Reset() {
	*x = ReqJsonPack{}
	if protoimpl.UnsafeEnabled {
		mi := &file_outer_pack_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReqJsonPack) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReqJsonPack) ProtoMessage() {}

func (x *ReqJsonPack) ProtoReflect() protoreflect.Message {
	mi := &file_outer_pack_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReqJsonPack.ProtoReflect.Descriptor instead.
func (*ReqJsonPack) Descriptor() ([]byte, []int) {
	return file_outer_pack_proto_rawDescGZIP(), []int{1}
}

func (x *ReqJsonPack) GetSn() uint64 {
	if x != nil {
		return x.Sn
	}
	return 0
}

func (x *ReqJsonPack) GetRpcId() int32 {
	if x != nil {
		return x.RpcId
	}
	return 0
}

func (x *ReqJsonPack) GetData() string {
	if x != nil {
		return x.Data
	}
	return ""
}

//应答
type RspPack struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sn    uint64    `protobuf:"varint,1,opt,name=sn,proto3" json:"sn,omitempty"`                    //流水号 request/reply
	RpcId int32     `protobuf:"varint,2,opt,name=rpcId,proto3" json:"rpcId,omitempty"`              //rpcId all
	Code  code.Code `protobuf:"varint,3,opt,name=code,proto3,enum=code.Code" json:"code,omitempty"` //错误码 reply
	Data  []byte    `protobuf:"bytes,4,opt,name=data,proto3" json:"data,omitempty"`                 //消息体 all
	Msg   string    `protobuf:"bytes,5,opt,name=msg,proto3" json:"msg,omitempty"`                   //错误消息,可能有
}

func (x *RspPack) Reset() {
	*x = RspPack{}
	if protoimpl.UnsafeEnabled {
		mi := &file_outer_pack_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RspPack) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RspPack) ProtoMessage() {}

func (x *RspPack) ProtoReflect() protoreflect.Message {
	mi := &file_outer_pack_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RspPack.ProtoReflect.Descriptor instead.
func (*RspPack) Descriptor() ([]byte, []int) {
	return file_outer_pack_proto_rawDescGZIP(), []int{2}
}

func (x *RspPack) GetSn() uint64 {
	if x != nil {
		return x.Sn
	}
	return 0
}

func (x *RspPack) GetRpcId() int32 {
	if x != nil {
		return x.RpcId
	}
	return 0
}

func (x *RspPack) GetCode() code.Code {
	if x != nil {
		return x.Code
	}
	return code.Code(0)
}

func (x *RspPack) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *RspPack) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

//应答--json版本
type RspJsonPack struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sn    uint64    `protobuf:"varint,1,opt,name=sn,proto3" json:"sn,omitempty"`                    //流水号 request/reply
	RpcId int32     `protobuf:"varint,2,opt,name=rpcId,proto3" json:"rpcId,omitempty"`              //rpcId all
	Code  code.Code `protobuf:"varint,3,opt,name=code,proto3,enum=code.Code" json:"code,omitempty"` //错误码 reply
	Data  string    `protobuf:"bytes,4,opt,name=data,proto3" json:"data,omitempty"`                 //消息体 all
	Msg   string    `protobuf:"bytes,5,opt,name=msg,proto3" json:"msg,omitempty"`                   //错误消息,可能有
}

func (x *RspJsonPack) Reset() {
	*x = RspJsonPack{}
	if protoimpl.UnsafeEnabled {
		mi := &file_outer_pack_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RspJsonPack) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RspJsonPack) ProtoMessage() {}

func (x *RspJsonPack) ProtoReflect() protoreflect.Message {
	mi := &file_outer_pack_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RspJsonPack.ProtoReflect.Descriptor instead.
func (*RspJsonPack) Descriptor() ([]byte, []int) {
	return file_outer_pack_proto_rawDescGZIP(), []int{3}
}

func (x *RspJsonPack) GetSn() uint64 {
	if x != nil {
		return x.Sn
	}
	return 0
}

func (x *RspJsonPack) GetRpcId() int32 {
	if x != nil {
		return x.RpcId
	}
	return 0
}

func (x *RspJsonPack) GetCode() code.Code {
	if x != nil {
		return x.Code
	}
	return code.Code(0)
}

func (x *RspJsonPack) GetData() string {
	if x != nil {
		return x.Data
	}
	return ""
}

func (x *RspJsonPack) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

//http回复
type HttpReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code code.Code `protobuf:"varint,1,opt,name=code,proto3,enum=code.Code" json:"code,omitempty"` //错误码 reply
	Data []byte    `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`                 // 正确回复时候的消息体
	Msg  string    `protobuf:"bytes,3,opt,name=msg,proto3" json:"msg,omitempty"`                   // 错误回复时候的提示
}

func (x *HttpReply) Reset() {
	*x = HttpReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_outer_pack_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HttpReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HttpReply) ProtoMessage() {}

func (x *HttpReply) ProtoReflect() protoreflect.Message {
	mi := &file_outer_pack_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HttpReply.ProtoReflect.Descriptor instead.
func (*HttpReply) Descriptor() ([]byte, []int) {
	return file_outer_pack_proto_rawDescGZIP(), []int{4}
}

func (x *HttpReply) GetCode() code.Code {
	if x != nil {
		return x.Code
	}
	return code.Code(0)
}

func (x *HttpReply) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *HttpReply) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

//推送--服务器通知gateway给session的push消息
type PushPack struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sn    uint64 `protobuf:"varint,1,opt,name=sn,proto3" json:"sn,omitempty"`       //流水号 request/reply
	RpcId int32  `protobuf:"varint,2,opt,name=rpcId,proto3" json:"rpcId,omitempty"` //rpcId all
	Data  []byte `protobuf:"bytes,4,opt,name=data,proto3" json:"data,omitempty"`    //消息体--
}

func (x *PushPack) Reset() {
	*x = PushPack{}
	if protoimpl.UnsafeEnabled {
		mi := &file_outer_pack_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PushPack) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushPack) ProtoMessage() {}

func (x *PushPack) ProtoReflect() protoreflect.Message {
	mi := &file_outer_pack_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushPack.ProtoReflect.Descriptor instead.
func (*PushPack) Descriptor() ([]byte, []int) {
	return file_outer_pack_proto_rawDescGZIP(), []int{5}
}

func (x *PushPack) GetSn() uint64 {
	if x != nil {
		return x.Sn
	}
	return 0
}

func (x *PushPack) GetRpcId() int32 {
	if x != nil {
		return x.RpcId
	}
	return 0
}

func (x *PushPack) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_outer_pack_proto protoreflect.FileDescriptor

var file_outer_pack_proto_rawDesc = []byte{
	0x0a, 0x10, 0x6f, 0x75, 0x74, 0x65, 0x72, 0x2f, 0x70, 0x61, 0x63, 0x6b, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x05, 0x6f, 0x75, 0x74, 0x65, 0x72, 0x1a, 0x0f, 0x63, 0x6f, 0x64, 0x65, 0x2f,
	0x63, 0x6f, 0x64, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x43, 0x0a, 0x07, 0x52, 0x65,
	0x71, 0x50, 0x61, 0x63, 0x6b, 0x12, 0x0e, 0x0a, 0x02, 0x73, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x02, 0x73, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x70, 0x63, 0x49, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x72, 0x70, 0x63, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x64,
	0x61, 0x74, 0x61, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22,
	0x47, 0x0a, 0x0b, 0x52, 0x65, 0x71, 0x4a, 0x73, 0x6f, 0x6e, 0x50, 0x61, 0x63, 0x6b, 0x12, 0x0e,
	0x0a, 0x02, 0x73, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x73, 0x6e, 0x12, 0x14,
	0x0a, 0x05, 0x72, 0x70, 0x63, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x72,
	0x70, 0x63, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x75, 0x0a, 0x07, 0x52, 0x73, 0x70, 0x50,
	0x61, 0x63, 0x6b, 0x12, 0x0e, 0x0a, 0x02, 0x73, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x02, 0x73, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x70, 0x63, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x05, 0x72, 0x70, 0x63, 0x49, 0x64, 0x12, 0x1e, 0x0a, 0x04, 0x63, 0x6f, 0x64,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0a, 0x2e, 0x63, 0x6f, 0x64, 0x65, 0x2e, 0x43,
	0x6f, 0x64, 0x65, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x10, 0x0a,
	0x03, 0x6d, 0x73, 0x67, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x22,
	0x79, 0x0a, 0x0b, 0x52, 0x73, 0x70, 0x4a, 0x73, 0x6f, 0x6e, 0x50, 0x61, 0x63, 0x6b, 0x12, 0x0e,
	0x0a, 0x02, 0x73, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x73, 0x6e, 0x12, 0x14,
	0x0a, 0x05, 0x72, 0x70, 0x63, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x72,
	0x70, 0x63, 0x49, 0x64, 0x12, 0x1e, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x0a, 0x2e, 0x63, 0x6f, 0x64, 0x65, 0x2e, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x04,
	0x63, 0x6f, 0x64, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x22, 0x51, 0x0a, 0x09, 0x48, 0x74,
	0x74, 0x70, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x1e, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0a, 0x2e, 0x63, 0x6f, 0x64, 0x65, 0x2e, 0x43, 0x6f, 0x64,
	0x65, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x10, 0x0a, 0x03, 0x6d,
	0x73, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x22, 0x44, 0x0a,
	0x08, 0x50, 0x75, 0x73, 0x68, 0x50, 0x61, 0x63, 0x6b, 0x12, 0x0e, 0x0a, 0x02, 0x73, 0x6e, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x73, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x70, 0x63,
	0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x72, 0x70, 0x63, 0x49, 0x64, 0x12,
	0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64,
	0x61, 0x74, 0x61, 0x42, 0x13, 0x5a, 0x11, 0x6c, 0x61, 0x69, 0x79, 0x61, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2f, 0x6f, 0x75, 0x74, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_outer_pack_proto_rawDescOnce sync.Once
	file_outer_pack_proto_rawDescData = file_outer_pack_proto_rawDesc
)

func file_outer_pack_proto_rawDescGZIP() []byte {
	file_outer_pack_proto_rawDescOnce.Do(func() {
		file_outer_pack_proto_rawDescData = protoimpl.X.CompressGZIP(file_outer_pack_proto_rawDescData)
	})
	return file_outer_pack_proto_rawDescData
}

var file_outer_pack_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_outer_pack_proto_goTypes = []interface{}{
	(*ReqPack)(nil),     // 0: outer.ReqPack
	(*ReqJsonPack)(nil), // 1: outer.ReqJsonPack
	(*RspPack)(nil),     // 2: outer.RspPack
	(*RspJsonPack)(nil), // 3: outer.RspJsonPack
	(*HttpReply)(nil),   // 4: outer.HttpReply
	(*PushPack)(nil),    // 5: outer.PushPack
	(code.Code)(0),      // 6: code.Code
}
var file_outer_pack_proto_depIdxs = []int32{
	6, // 0: outer.RspPack.code:type_name -> code.Code
	6, // 1: outer.RspJsonPack.code:type_name -> code.Code
	6, // 2: outer.HttpReply.code:type_name -> code.Code
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_outer_pack_proto_init() }
func file_outer_pack_proto_init() {
	if File_outer_pack_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_outer_pack_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReqPack); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_outer_pack_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReqJsonPack); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_outer_pack_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RspPack); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_outer_pack_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RspJsonPack); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_outer_pack_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HttpReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_outer_pack_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PushPack); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_outer_pack_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_outer_pack_proto_goTypes,
		DependencyIndexes: file_outer_pack_proto_depIdxs,
		MessageInfos:      file_outer_pack_proto_msgTypes,
	}.Build()
	File_outer_pack_proto = out.File
	file_outer_pack_proto_rawDesc = nil
	file_outer_pack_proto_goTypes = nil
	file_outer_pack_proto_depIdxs = nil
}