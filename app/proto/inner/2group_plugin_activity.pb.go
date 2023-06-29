// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.5.1
// source: inner/2group_plugin_activity.proto

package inner

import (
	_ "github.com/asynkron/protoactor-go/actor"
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

//用户离开群
type Gr2AcPlayerUnfavoriteGroup struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Gr2AcPlayerUnfavoriteGroup) Reset() {
	*x = Gr2AcPlayerUnfavoriteGroup{}
	if protoimpl.UnsafeEnabled {
		mi := &file_inner_2group_plugin_activity_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Gr2AcPlayerUnfavoriteGroup) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Gr2AcPlayerUnfavoriteGroup) ProtoMessage() {}

func (x *Gr2AcPlayerUnfavoriteGroup) ProtoReflect() protoreflect.Message {
	mi := &file_inner_2group_plugin_activity_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Gr2AcPlayerUnfavoriteGroup.ProtoReflect.Descriptor instead.
func (*Gr2AcPlayerUnfavoriteGroup) Descriptor() ([]byte, []int) {
	return file_inner_2group_plugin_activity_proto_rawDescGZIP(), []int{0}
}

type Gr2AcPlayerUnfavoriteGroup_Notify struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uid uint64 `protobuf:"varint,1,opt,name=uid,proto3" json:"uid,omitempty"`
}

func (x *Gr2AcPlayerUnfavoriteGroup_Notify) Reset() {
	*x = Gr2AcPlayerUnfavoriteGroup_Notify{}
	if protoimpl.UnsafeEnabled {
		mi := &file_inner_2group_plugin_activity_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Gr2AcPlayerUnfavoriteGroup_Notify) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Gr2AcPlayerUnfavoriteGroup_Notify) ProtoMessage() {}

func (x *Gr2AcPlayerUnfavoriteGroup_Notify) ProtoReflect() protoreflect.Message {
	mi := &file_inner_2group_plugin_activity_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Gr2AcPlayerUnfavoriteGroup_Notify.ProtoReflect.Descriptor instead.
func (*Gr2AcPlayerUnfavoriteGroup_Notify) Descriptor() ([]byte, []int) {
	return file_inner_2group_plugin_activity_proto_rawDescGZIP(), []int{0, 0}
}

func (x *Gr2AcPlayerUnfavoriteGroup_Notify) GetUid() uint64 {
	if x != nil {
		return x.Uid
	}
	return 0
}

var File_inner_2group_plugin_activity_proto protoreflect.FileDescriptor

var file_inner_2group_plugin_activity_proto_rawDesc = []byte{
	0x0a, 0x22, 0x69, 0x6e, 0x6e, 0x65, 0x72, 0x2f, 0x32, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x70,
	0x6c, 0x75, 0x67, 0x69, 0x6e, 0x5f, 0x61, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x69, 0x6e, 0x6e, 0x65, 0x72, 0x1a, 0x11, 0x61, 0x63, 0x74,
	0x6f, 0x72, 0x2f, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x38,
	0x0a, 0x1a, 0x47, 0x72, 0x32, 0x41, 0x63, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x55, 0x6e, 0x66,
	0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x1a, 0x1a, 0x0a, 0x06,
	0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x03, 0x75, 0x69, 0x64, 0x42, 0x13, 0x5a, 0x11, 0x6c, 0x61, 0x69, 0x79,
	0x61, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x69, 0x6e, 0x6e, 0x65, 0x72, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_inner_2group_plugin_activity_proto_rawDescOnce sync.Once
	file_inner_2group_plugin_activity_proto_rawDescData = file_inner_2group_plugin_activity_proto_rawDesc
)

func file_inner_2group_plugin_activity_proto_rawDescGZIP() []byte {
	file_inner_2group_plugin_activity_proto_rawDescOnce.Do(func() {
		file_inner_2group_plugin_activity_proto_rawDescData = protoimpl.X.CompressGZIP(file_inner_2group_plugin_activity_proto_rawDescData)
	})
	return file_inner_2group_plugin_activity_proto_rawDescData
}

var file_inner_2group_plugin_activity_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_inner_2group_plugin_activity_proto_goTypes = []interface{}{
	(*Gr2AcPlayerUnfavoriteGroup)(nil),        // 0: inner.Gr2AcPlayerUnfavoriteGroup
	(*Gr2AcPlayerUnfavoriteGroup_Notify)(nil), // 1: inner.Gr2AcPlayerUnfavoriteGroup.Notify
}
var file_inner_2group_plugin_activity_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_inner_2group_plugin_activity_proto_init() }
func file_inner_2group_plugin_activity_proto_init() {
	if File_inner_2group_plugin_activity_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_inner_2group_plugin_activity_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Gr2AcPlayerUnfavoriteGroup); i {
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
		file_inner_2group_plugin_activity_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Gr2AcPlayerUnfavoriteGroup_Notify); i {
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
			RawDescriptor: file_inner_2group_plugin_activity_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_inner_2group_plugin_activity_proto_goTypes,
		DependencyIndexes: file_inner_2group_plugin_activity_proto_depIdxs,
		MessageInfos:      file_inner_2group_plugin_activity_proto_msgTypes,
	}.Build()
	File_inner_2group_plugin_activity_proto = out.File
	file_inner_2group_plugin_activity_proto_rawDesc = nil
	file_inner_2group_plugin_activity_proto_goTypes = nil
	file_inner_2group_plugin_activity_proto_depIdxs = nil
}
