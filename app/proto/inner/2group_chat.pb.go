// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.5.1
// source: inner/2group_chat.proto

package inner

import (
	_ "github.com/asynkron/protoactor-go/actor"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	outer "laiya/proto/outer"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

//通知群聊天
type NotifyGroupChatMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *NotifyGroupChatMsg) Reset() {
	*x = NotifyGroupChatMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_inner_2group_chat_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NotifyGroupChatMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NotifyGroupChatMsg) ProtoMessage() {}

func (x *NotifyGroupChatMsg) ProtoReflect() protoreflect.Message {
	mi := &file_inner_2group_chat_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NotifyGroupChatMsg.ProtoReflect.Descriptor instead.
func (*NotifyGroupChatMsg) Descriptor() ([]byte, []int) {
	return file_inner_2group_chat_proto_rawDescGZIP(), []int{0}
}

//通知群活动聊天
type NotifyGroupActivityInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *NotifyGroupActivityInfo) Reset() {
	*x = NotifyGroupActivityInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_inner_2group_chat_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NotifyGroupActivityInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NotifyGroupActivityInfo) ProtoMessage() {}

func (x *NotifyGroupActivityInfo) ProtoReflect() protoreflect.Message {
	mi := &file_inner_2group_chat_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NotifyGroupActivityInfo.ProtoReflect.Descriptor instead.
func (*NotifyGroupActivityInfo) Descriptor() ([]byte, []int) {
	return file_inner_2group_chat_proto_rawDescGZIP(), []int{1}
}

//活动插件通知群邀请上麦
type NotifyGroupInviteOnMic struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *NotifyGroupInviteOnMic) Reset() {
	*x = NotifyGroupInviteOnMic{}
	if protoimpl.UnsafeEnabled {
		mi := &file_inner_2group_chat_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NotifyGroupInviteOnMic) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NotifyGroupInviteOnMic) ProtoMessage() {}

func (x *NotifyGroupInviteOnMic) ProtoReflect() protoreflect.Message {
	mi := &file_inner_2group_chat_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NotifyGroupInviteOnMic.ProtoReflect.Descriptor instead.
func (*NotifyGroupInviteOnMic) Descriptor() ([]byte, []int) {
	return file_inner_2group_chat_proto_rawDescGZIP(), []int{2}
}

type NotifyGroupChatMsg_Notify struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SenderUid uint64                   `protobuf:"varint,1,opt,name=senderUid,proto3" json:"senderUid,omitempty"` //默认房主
	ChatMsg   *outer.ChatMsg_NormalMsg `protobuf:"bytes,2,opt,name=chatMsg,proto3" json:"chatMsg,omitempty"`
}

func (x *NotifyGroupChatMsg_Notify) Reset() {
	*x = NotifyGroupChatMsg_Notify{}
	if protoimpl.UnsafeEnabled {
		mi := &file_inner_2group_chat_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NotifyGroupChatMsg_Notify) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NotifyGroupChatMsg_Notify) ProtoMessage() {}

func (x *NotifyGroupChatMsg_Notify) ProtoReflect() protoreflect.Message {
	mi := &file_inner_2group_chat_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NotifyGroupChatMsg_Notify.ProtoReflect.Descriptor instead.
func (*NotifyGroupChatMsg_Notify) Descriptor() ([]byte, []int) {
	return file_inner_2group_chat_proto_rawDescGZIP(), []int{0, 0}
}

func (x *NotifyGroupChatMsg_Notify) GetSenderUid() uint64 {
	if x != nil {
		return x.SenderUid
	}
	return 0
}

func (x *NotifyGroupChatMsg_Notify) GetChatMsg() *outer.ChatMsg_NormalMsg {
	if x != nil {
		return x.ChatMsg
	}
	return nil
}

type NotifyGroupActivityInfo_Notify struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SenderUid    uint64                          `protobuf:"varint,1,opt,name=senderUid,proto3" json:"senderUid,omitempty"` //默认房主
	ActivityInfo *outer.ChatMsg_GroupActivityMsg `protobuf:"bytes,2,opt,name=activityInfo,proto3" json:"activityInfo,omitempty"`
}

func (x *NotifyGroupActivityInfo_Notify) Reset() {
	*x = NotifyGroupActivityInfo_Notify{}
	if protoimpl.UnsafeEnabled {
		mi := &file_inner_2group_chat_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NotifyGroupActivityInfo_Notify) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NotifyGroupActivityInfo_Notify) ProtoMessage() {}

func (x *NotifyGroupActivityInfo_Notify) ProtoReflect() protoreflect.Message {
	mi := &file_inner_2group_chat_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NotifyGroupActivityInfo_Notify.ProtoReflect.Descriptor instead.
func (*NotifyGroupActivityInfo_Notify) Descriptor() ([]byte, []int) {
	return file_inner_2group_chat_proto_rawDescGZIP(), []int{1, 0}
}

func (x *NotifyGroupActivityInfo_Notify) GetSenderUid() uint64 {
	if x != nil {
		return x.SenderUid
	}
	return 0
}

func (x *NotifyGroupActivityInfo_Notify) GetActivityInfo() *outer.ChatMsg_GroupActivityMsg {
	if x != nil {
		return x.ActivityInfo
	}
	return nil
}

type NotifyGroupInviteOnMic_Notify struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SenderUid     uint64   `protobuf:"varint,1,opt,name=senderUid,proto3" json:"senderUid,omitempty"` //默认房主
	ActivityTitle string   `protobuf:"bytes,2,opt,name=activityTitle,proto3" json:"activityTitle,omitempty"`
	Uids          []uint64 `protobuf:"varint,3,rep,packed,name=uids,proto3" json:"uids,omitempty"`
}

func (x *NotifyGroupInviteOnMic_Notify) Reset() {
	*x = NotifyGroupInviteOnMic_Notify{}
	if protoimpl.UnsafeEnabled {
		mi := &file_inner_2group_chat_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NotifyGroupInviteOnMic_Notify) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NotifyGroupInviteOnMic_Notify) ProtoMessage() {}

func (x *NotifyGroupInviteOnMic_Notify) ProtoReflect() protoreflect.Message {
	mi := &file_inner_2group_chat_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NotifyGroupInviteOnMic_Notify.ProtoReflect.Descriptor instead.
func (*NotifyGroupInviteOnMic_Notify) Descriptor() ([]byte, []int) {
	return file_inner_2group_chat_proto_rawDescGZIP(), []int{2, 0}
}

func (x *NotifyGroupInviteOnMic_Notify) GetSenderUid() uint64 {
	if x != nil {
		return x.SenderUid
	}
	return 0
}

func (x *NotifyGroupInviteOnMic_Notify) GetActivityTitle() string {
	if x != nil {
		return x.ActivityTitle
	}
	return ""
}

func (x *NotifyGroupInviteOnMic_Notify) GetUids() []uint64 {
	if x != nil {
		return x.Uids
	}
	return nil
}

var File_inner_2group_chat_proto protoreflect.FileDescriptor

var file_inner_2group_chat_proto_rawDesc = []byte{
	0x0a, 0x17, 0x69, 0x6e, 0x6e, 0x65, 0x72, 0x2f, 0x32, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x63,
	0x68, 0x61, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x69, 0x6e, 0x6e, 0x65, 0x72,
	0x1a, 0x11, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x2f, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x18, 0x6f, 0x75, 0x74, 0x65, 0x72, 0x2f, 0x6f, 0x75, 0x74, 0x65, 0x72,
	0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x70, 0x0a,
	0x12, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x43, 0x68, 0x61, 0x74,
	0x4d, 0x73, 0x67, 0x1a, 0x5a, 0x0a, 0x06, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x12, 0x1c, 0x0a,
	0x09, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x55, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x09, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x55, 0x69, 0x64, 0x12, 0x32, 0x0a, 0x07, 0x63,
	0x68, 0x61, 0x74, 0x4d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x6f,
	0x75, 0x74, 0x65, 0x72, 0x2e, 0x43, 0x68, 0x61, 0x74, 0x4d, 0x73, 0x67, 0x2e, 0x4e, 0x6f, 0x72,
	0x6d, 0x61, 0x6c, 0x4d, 0x73, 0x67, 0x52, 0x07, 0x63, 0x68, 0x61, 0x74, 0x4d, 0x73, 0x67, 0x22,
	0x86, 0x01, 0x0a, 0x17, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x41,
	0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x49, 0x6e, 0x66, 0x6f, 0x1a, 0x6b, 0x0a, 0x06, 0x4e,
	0x6f, 0x74, 0x69, 0x66, 0x79, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x55,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72,
	0x55, 0x69, 0x64, 0x12, 0x43, 0x0a, 0x0c, 0x61, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x49,
	0x6e, 0x66, 0x6f, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x6f, 0x75, 0x74, 0x65,
	0x72, 0x2e, 0x43, 0x68, 0x61, 0x74, 0x4d, 0x73, 0x67, 0x2e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x41,
	0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x4d, 0x73, 0x67, 0x52, 0x0c, 0x61, 0x63, 0x74, 0x69,
	0x76, 0x69, 0x74, 0x79, 0x49, 0x6e, 0x66, 0x6f, 0x22, 0x7a, 0x0a, 0x16, 0x4e, 0x6f, 0x74, 0x69,
	0x66, 0x79, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x6e, 0x76, 0x69, 0x74, 0x65, 0x4f, 0x6e, 0x4d,
	0x69, 0x63, 0x1a, 0x60, 0x0a, 0x06, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x12, 0x1c, 0x0a, 0x09,
	0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x55, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x09, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x55, 0x69, 0x64, 0x12, 0x24, 0x0a, 0x0d, 0x61, 0x63,
	0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x54, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0d, 0x61, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x54, 0x69, 0x74, 0x6c, 0x65,
	0x12, 0x12, 0x0a, 0x04, 0x75, 0x69, 0x64, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x04, 0x52, 0x04,
	0x75, 0x69, 0x64, 0x73, 0x42, 0x13, 0x5a, 0x11, 0x6c, 0x61, 0x69, 0x79, 0x61, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x69, 0x6e, 0x6e, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_inner_2group_chat_proto_rawDescOnce sync.Once
	file_inner_2group_chat_proto_rawDescData = file_inner_2group_chat_proto_rawDesc
)

func file_inner_2group_chat_proto_rawDescGZIP() []byte {
	file_inner_2group_chat_proto_rawDescOnce.Do(func() {
		file_inner_2group_chat_proto_rawDescData = protoimpl.X.CompressGZIP(file_inner_2group_chat_proto_rawDescData)
	})
	return file_inner_2group_chat_proto_rawDescData
}

var file_inner_2group_chat_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_inner_2group_chat_proto_goTypes = []interface{}{
	(*NotifyGroupChatMsg)(nil),             // 0: inner.NotifyGroupChatMsg
	(*NotifyGroupActivityInfo)(nil),        // 1: inner.NotifyGroupActivityInfo
	(*NotifyGroupInviteOnMic)(nil),         // 2: inner.NotifyGroupInviteOnMic
	(*NotifyGroupChatMsg_Notify)(nil),      // 3: inner.NotifyGroupChatMsg.Notify
	(*NotifyGroupActivityInfo_Notify)(nil), // 4: inner.NotifyGroupActivityInfo.Notify
	(*NotifyGroupInviteOnMic_Notify)(nil),  // 5: inner.NotifyGroupInviteOnMic.Notify
	(*outer.ChatMsg_NormalMsg)(nil),        // 6: outer.ChatMsg.NormalMsg
	(*outer.ChatMsg_GroupActivityMsg)(nil), // 7: outer.ChatMsg.GroupActivityMsg
}
var file_inner_2group_chat_proto_depIdxs = []int32{
	6, // 0: inner.NotifyGroupChatMsg.Notify.chatMsg:type_name -> outer.ChatMsg.NormalMsg
	7, // 1: inner.NotifyGroupActivityInfo.Notify.activityInfo:type_name -> outer.ChatMsg.GroupActivityMsg
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_inner_2group_chat_proto_init() }
func file_inner_2group_chat_proto_init() {
	if File_inner_2group_chat_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_inner_2group_chat_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NotifyGroupChatMsg); i {
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
		file_inner_2group_chat_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NotifyGroupActivityInfo); i {
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
		file_inner_2group_chat_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NotifyGroupInviteOnMic); i {
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
		file_inner_2group_chat_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NotifyGroupChatMsg_Notify); i {
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
		file_inner_2group_chat_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NotifyGroupActivityInfo_Notify); i {
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
		file_inner_2group_chat_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NotifyGroupInviteOnMic_Notify); i {
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
			RawDescriptor: file_inner_2group_chat_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_inner_2group_chat_proto_goTypes,
		DependencyIndexes: file_inner_2group_chat_proto_depIdxs,
		MessageInfos:      file_inner_2group_chat_proto_msgTypes,
	}.Build()
	File_inner_2group_chat_proto = out.File
	file_inner_2group_chat_proto_rawDesc = nil
	file_inner_2group_chat_proto_goTypes = nil
	file_inner_2group_chat_proto_depIdxs = nil
}