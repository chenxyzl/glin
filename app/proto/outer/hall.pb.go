// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.5.1
// source: outer/hall.proto

package outer

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

//获取推荐群
type GetRecommendGroupList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RpcId int32 `protobuf:"varint,40000,opt,name=rpcId,proto3" json:"rpcId,omitempty"`
}

func (x *GetRecommendGroupList) Reset() {
	*x = GetRecommendGroupList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_outer_hall_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetRecommendGroupList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRecommendGroupList) ProtoMessage() {}

func (x *GetRecommendGroupList) ProtoReflect() protoreflect.Message {
	mi := &file_outer_hall_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRecommendGroupList.ProtoReflect.Descriptor instead.
func (*GetRecommendGroupList) Descriptor() ([]byte, []int) {
	return file_outer_hall_proto_rawDescGZIP(), []int{0}
}

func (x *GetRecommendGroupList) GetRpcId() int32 {
	if x != nil {
		return x.RpcId
	}
	return 0
}

//搜索群
type SearchGroup struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RpcId int32 `protobuf:"varint,40001,opt,name=rpcId,proto3" json:"rpcId,omitempty"`
}

func (x *SearchGroup) Reset() {
	*x = SearchGroup{}
	if protoimpl.UnsafeEnabled {
		mi := &file_outer_hall_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SearchGroup) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchGroup) ProtoMessage() {}

func (x *SearchGroup) ProtoReflect() protoreflect.Message {
	mi := &file_outer_hall_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchGroup.ProtoReflect.Descriptor instead.
func (*SearchGroup) Descriptor() ([]byte, []int) {
	return file_outer_hall_proto_rawDescGZIP(), []int{1}
}

func (x *SearchGroup) GetRpcId() int32 {
	if x != nil {
		return x.RpcId
	}
	return 0
}

type GetRecommendGroupList_Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Page uint32 `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"` //页
}

func (x *GetRecommendGroupList_Request) Reset() {
	*x = GetRecommendGroupList_Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_outer_hall_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetRecommendGroupList_Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRecommendGroupList_Request) ProtoMessage() {}

func (x *GetRecommendGroupList_Request) ProtoReflect() protoreflect.Message {
	mi := &file_outer_hall_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRecommendGroupList_Request.ProtoReflect.Descriptor instead.
func (*GetRecommendGroupList_Request) Descriptor() ([]byte, []int) {
	return file_outer_hall_proto_rawDescGZIP(), []int{0, 0}
}

func (x *GetRecommendGroupList_Request) GetPage() uint32 {
	if x != nil {
		return x.Page
	}
	return 0
}

type GetRecommendGroupList_Reply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GroupIds []uint64 `protobuf:"varint,1,rep,packed,name=groupIds,proto3" json:"groupIds,omitempty"`
}

func (x *GetRecommendGroupList_Reply) Reset() {
	*x = GetRecommendGroupList_Reply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_outer_hall_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetRecommendGroupList_Reply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRecommendGroupList_Reply) ProtoMessage() {}

func (x *GetRecommendGroupList_Reply) ProtoReflect() protoreflect.Message {
	mi := &file_outer_hall_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRecommendGroupList_Reply.ProtoReflect.Descriptor instead.
func (*GetRecommendGroupList_Reply) Descriptor() ([]byte, []int) {
	return file_outer_hall_proto_rawDescGZIP(), []int{0, 1}
}

func (x *GetRecommendGroupList_Reply) GetGroupIds() []uint64 {
	if x != nil {
		return x.GroupIds
	}
	return nil
}

type SearchGroup_Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StrVal string `protobuf:"bytes,1,opt,name=strVal,proto3" json:"strVal,omitempty"` //搜索--(全数字则为id搜索)
	Page   uint32 `protobuf:"varint,2,opt,name=page,proto3" json:"page,omitempty"`    //页
}

func (x *SearchGroup_Request) Reset() {
	*x = SearchGroup_Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_outer_hall_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SearchGroup_Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchGroup_Request) ProtoMessage() {}

func (x *SearchGroup_Request) ProtoReflect() protoreflect.Message {
	mi := &file_outer_hall_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchGroup_Request.ProtoReflect.Descriptor instead.
func (*SearchGroup_Request) Descriptor() ([]byte, []int) {
	return file_outer_hall_proto_rawDescGZIP(), []int{1, 0}
}

func (x *SearchGroup_Request) GetStrVal() string {
	if x != nil {
		return x.StrVal
	}
	return ""
}

func (x *SearchGroup_Request) GetPage() uint32 {
	if x != nil {
		return x.Page
	}
	return 0
}

type SearchGroup_Reply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GroupIds []uint64 `protobuf:"varint,1,rep,packed,name=groupIds,proto3" json:"groupIds,omitempty"`
}

func (x *SearchGroup_Reply) Reset() {
	*x = SearchGroup_Reply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_outer_hall_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SearchGroup_Reply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchGroup_Reply) ProtoMessage() {}

func (x *SearchGroup_Reply) ProtoReflect() protoreflect.Message {
	mi := &file_outer_hall_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchGroup_Reply.ProtoReflect.Descriptor instead.
func (*SearchGroup_Reply) Descriptor() ([]byte, []int) {
	return file_outer_hall_proto_rawDescGZIP(), []int{1, 1}
}

func (x *SearchGroup_Reply) GetGroupIds() []uint64 {
	if x != nil {
		return x.GroupIds
	}
	return nil
}

var File_outer_hall_proto protoreflect.FileDescriptor

var file_outer_hall_proto_rawDesc = []byte{
	0x0a, 0x10, 0x6f, 0x75, 0x74, 0x65, 0x72, 0x2f, 0x68, 0x61, 0x6c, 0x6c, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x05, 0x6f, 0x75, 0x74, 0x65, 0x72, 0x22, 0x73, 0x0a, 0x15, 0x47, 0x65, 0x74,
	0x52, 0x65, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x64, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x4c, 0x69,
	0x73, 0x74, 0x12, 0x16, 0x0a, 0x05, 0x72, 0x70, 0x63, 0x49, 0x64, 0x18, 0xc0, 0xb8, 0x02, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x05, 0x72, 0x70, 0x63, 0x49, 0x64, 0x1a, 0x1d, 0x0a, 0x07, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x1a, 0x23, 0x0a, 0x05, 0x52, 0x65, 0x70,
	0x6c, 0x79, 0x12, 0x1a, 0x0a, 0x08, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x64, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x04, 0x52, 0x08, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x64, 0x73, 0x22, 0x81,
	0x01, 0x0a, 0x0b, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x12, 0x16,
	0x0a, 0x05, 0x72, 0x70, 0x63, 0x49, 0x64, 0x18, 0xc1, 0xb8, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x05, 0x72, 0x70, 0x63, 0x49, 0x64, 0x1a, 0x35, 0x0a, 0x07, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x72, 0x56, 0x61, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x73, 0x74, 0x72, 0x56, 0x61, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x1a, 0x23, 0x0a,
	0x05, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x1a, 0x0a, 0x08, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x49,
	0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x04, 0x52, 0x08, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x49,
	0x64, 0x73, 0x42, 0x13, 0x5a, 0x11, 0x6c, 0x61, 0x69, 0x79, 0x61, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x6f, 0x75, 0x74, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_outer_hall_proto_rawDescOnce sync.Once
	file_outer_hall_proto_rawDescData = file_outer_hall_proto_rawDesc
)

func file_outer_hall_proto_rawDescGZIP() []byte {
	file_outer_hall_proto_rawDescOnce.Do(func() {
		file_outer_hall_proto_rawDescData = protoimpl.X.CompressGZIP(file_outer_hall_proto_rawDescData)
	})
	return file_outer_hall_proto_rawDescData
}

var file_outer_hall_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_outer_hall_proto_goTypes = []interface{}{
	(*GetRecommendGroupList)(nil),         // 0: outer.GetRecommendGroupList
	(*SearchGroup)(nil),                   // 1: outer.SearchGroup
	(*GetRecommendGroupList_Request)(nil), // 2: outer.GetRecommendGroupList.Request
	(*GetRecommendGroupList_Reply)(nil),   // 3: outer.GetRecommendGroupList.Reply
	(*SearchGroup_Request)(nil),           // 4: outer.SearchGroup.Request
	(*SearchGroup_Reply)(nil),             // 5: outer.SearchGroup.Reply
}
var file_outer_hall_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_outer_hall_proto_init() }
func file_outer_hall_proto_init() {
	if File_outer_hall_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_outer_hall_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetRecommendGroupList); i {
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
		file_outer_hall_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SearchGroup); i {
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
		file_outer_hall_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetRecommendGroupList_Request); i {
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
		file_outer_hall_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetRecommendGroupList_Reply); i {
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
		file_outer_hall_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SearchGroup_Request); i {
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
		file_outer_hall_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SearchGroup_Reply); i {
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
			RawDescriptor: file_outer_hall_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_outer_hall_proto_goTypes,
		DependencyIndexes: file_outer_hall_proto_depIdxs,
		MessageInfos:      file_outer_hall_proto_msgTypes,
	}.Build()
	File_outer_hall_proto = out.File
	file_outer_hall_proto_rawDesc = nil
	file_outer_hall_proto_goTypes = nil
	file_outer_hall_proto_depIdxs = nil
}