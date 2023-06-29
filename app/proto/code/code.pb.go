// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.5.1
// source: code/code.proto

package code

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

type Code int32

const (
	//base
	Code_UnknownErr           Code = 0 //未知错误--dart对于新增的不能识别的枚举都会认为是对应的0的枚举值
	Code_Ok                   Code = 1
	Code_Error                Code = 2  //公共错误
	Code_InnerError           Code = 3  //内部未知错误
	Code_InnerRpcError        Code = 4  //内部未知错误
	Code_MsgUnmarshalErr      Code = 5  //消息解析错误
	Code_UnsupportedMsgType   Code = 6  //不支持的消息类型
	Code_NotImpl              Code = 7  //未注册的消息
	Code_DataReadErr          Code = 8  //data读取错误
	Code_RpcTimeout           Code = 9  //消息处理超时
	Code_TokenNotExist        Code = 10 //token不存在
	Code_TokenFormatErr       Code = 11 //token格式错误
	Code_NotFoundMethod       Code = 12 //没有找到接口
	Code_NotFoundComponent    Code = 13 //没有找到接口
	Code_MsgDefineErr         Code = 14 //消息定义错误
	Code_UidNotFound          Code = 15 //没找到uid
	Code_NotFoundTargetId     Code = 16 //没找到目标id
	Code_EnumNotFoundImpl     Code = 17 //对应的枚举没找到实现
	Code_TokenErr             Code = 18 //token错误
	Code_NoAccess             Code = 19 //无权操作
	Code_RobotNotExist        Code = 20 //机器人不存在
	Code_AppVersionNeedUpdate Code = 21 //app版本需要升级
	Code_GmAuthError          Code = 22 //gm授权失败
	Code_ActorInitError       Code = 23 //actor初始化错误--一定不能改(23关联到InitErrReceive)
	//通用 10000~
	Code_UserNotExist           Code = 10000 //user不存在
	Code_ChatLenLimit           Code = 10008 //聊天内容过长
	Code_NowNotInGroup          Code = 10009 //当前没有进入群
	Code_MustQuitAllGroupBefore Code = 10010 //必须先退出所有的群
	Code_GameTypeNotExist       Code = 10011 //游戏类型不存在
	Code_PlatformTypeNotExist   Code = 10012 //平台类型不存在
	//20000 home
	Code_UserNameLenIllegal                 Code = 20001 //名字长度非法
	Code_UserDesLenIllegal                  Code = 20002 //描述长度非法
	Code_UpdateThirdPlatformInfoTypeNoExist Code = 20003 //更新第三方账户类型不存在
	//30000~ 群
	Code_GroupNotExist        Code = 30001 //群不存在
	Code_GroupCountFull       Code = 30002 //所有群总数量已满
	Code_CreateGroupCountFull Code = 30003 //自己创建的群数量已满
	Code_SelfGroup            Code = 30004 //自己的群
	Code_NotSelfGroup         Code = 30005 //不是自己的群
	Code_NotInGroup           Code = 30007 //不在群中
	Code_GroupIdFormatErr     Code = 30101 //群id格式错误
	Code_GroupSubIdFull       Code = 30102 //群的子id已用完
	Code_GroupPlayerFull      Code = 30103 //群的人数已满
	Code_GroupNameLenIllegal  Code = 30104 //名字长度非法
	Code_GroupDesLenIllegal   Code = 30105 //描述长度非法
	Code_GroupOnlineFull      Code = 30106 //群同时在线人数操作限制
	Code_JoinAccessDenied     Code = 30107 //加入权限-被拒绝
	Code_VoiceRoomNameIllegal Code = 30108 //语音名字不合法
	//35000~ 群插件-活动
	Code_GroupActivityNotExist               Code = 35001 //群活动不存在
	Code_GroupActivityFull                   Code = 35002 //群活动已满
	Code_GroupActivityTimeIllegal            Code = 35003 //群活动时间非法
	Code_GroupActivityExpired                Code = 35004 //群活动已过期
	Code_GroupActivityTitleLenIllegal        Code = 35005 //群活动标题长度不合法
	Code_GroupActivityDesLenIllegal          Code = 35006 //群活动描述长度不合法
	Code_GroupActivitySignUpRewardLenIllegal Code = 35007 //群活动报名备注长度不合法
	//90000 web错误
	Code_CaptchaTooFrequently        Code = 90000 //验证码请求过于频繁
	Code_CaptchaExpire               Code = 90001 //验证码过期
	Code_CaptchaNotEeq               Code = 90002 //验证码错误
	Code_PhNotIllegal                Code = 90003 //手机号非法
	Code_NotHasCaptcha               Code = 90004 //当前并没有验证码
	Code_UrlFormatErr                Code = 90005 //url格式错误
	Code_NotFoundUpdateFileNamedIcon Code = 90006 //没找到名为icon的上传文件
	Code_GeTuiCidIsEmpty             Code = 90007 //个推的cid没有设置
	Code_GeTuiOpenPushErr            Code = 90008 //个推推送开启失败
	Code_GeTuiClosePushErr           Code = 90009 //个推关闭失败
	//100000 微信
	Code_WechatWaitingLogin Code = 100000 //微信厨房还未启动完成
)

// Enum value maps for Code.
var (
	Code_name = map[int32]string{
		0:      "UnknownErr",
		1:      "Ok",
		2:      "Error",
		3:      "InnerError",
		4:      "InnerRpcError",
		5:      "MsgUnmarshalErr",
		6:      "UnsupportedMsgType",
		7:      "NotImpl",
		8:      "DataReadErr",
		9:      "RpcTimeout",
		10:     "TokenNotExist",
		11:     "TokenFormatErr",
		12:     "NotFoundMethod",
		13:     "NotFoundComponent",
		14:     "MsgDefineErr",
		15:     "UidNotFound",
		16:     "NotFoundTargetId",
		17:     "EnumNotFoundImpl",
		18:     "TokenErr",
		19:     "NoAccess",
		20:     "RobotNotExist",
		21:     "AppVersionNeedUpdate",
		22:     "GmAuthError",
		23:     "ActorInitError",
		10000:  "UserNotExist",
		10008:  "ChatLenLimit",
		10009:  "NowNotInGroup",
		10010:  "MustQuitAllGroupBefore",
		10011:  "GameTypeNotExist",
		10012:  "PlatformTypeNotExist",
		20001:  "UserNameLenIllegal",
		20002:  "UserDesLenIllegal",
		20003:  "UpdateThirdPlatformInfoTypeNoExist",
		30001:  "GroupNotExist",
		30002:  "GroupCountFull",
		30003:  "CreateGroupCountFull",
		30004:  "SelfGroup",
		30005:  "NotSelfGroup",
		30007:  "NotInGroup",
		30101:  "GroupIdFormatErr",
		30102:  "GroupSubIdFull",
		30103:  "GroupPlayerFull",
		30104:  "GroupNameLenIllegal",
		30105:  "GroupDesLenIllegal",
		30106:  "GroupOnlineFull",
		30107:  "JoinAccessDenied",
		30108:  "VoiceRoomNameIllegal",
		35001:  "GroupActivityNotExist",
		35002:  "GroupActivityFull",
		35003:  "GroupActivityTimeIllegal",
		35004:  "GroupActivityExpired",
		35005:  "GroupActivityTitleLenIllegal",
		35006:  "GroupActivityDesLenIllegal",
		35007:  "GroupActivitySignUpRewardLenIllegal",
		90000:  "CaptchaTooFrequently",
		90001:  "CaptchaExpire",
		90002:  "CaptchaNotEeq",
		90003:  "PhNotIllegal",
		90004:  "NotHasCaptcha",
		90005:  "UrlFormatErr",
		90006:  "NotFoundUpdateFileNamedIcon",
		90007:  "GeTuiCidIsEmpty",
		90008:  "GeTuiOpenPushErr",
		90009:  "GeTuiClosePushErr",
		100000: "WechatWaitingLogin",
	}
	Code_value = map[string]int32{
		"UnknownErr":                          0,
		"Ok":                                  1,
		"Error":                               2,
		"InnerError":                          3,
		"InnerRpcError":                       4,
		"MsgUnmarshalErr":                     5,
		"UnsupportedMsgType":                  6,
		"NotImpl":                             7,
		"DataReadErr":                         8,
		"RpcTimeout":                          9,
		"TokenNotExist":                       10,
		"TokenFormatErr":                      11,
		"NotFoundMethod":                      12,
		"NotFoundComponent":                   13,
		"MsgDefineErr":                        14,
		"UidNotFound":                         15,
		"NotFoundTargetId":                    16,
		"EnumNotFoundImpl":                    17,
		"TokenErr":                            18,
		"NoAccess":                            19,
		"RobotNotExist":                       20,
		"AppVersionNeedUpdate":                21,
		"GmAuthError":                         22,
		"ActorInitError":                      23,
		"UserNotExist":                        10000,
		"ChatLenLimit":                        10008,
		"NowNotInGroup":                       10009,
		"MustQuitAllGroupBefore":              10010,
		"GameTypeNotExist":                    10011,
		"PlatformTypeNotExist":                10012,
		"UserNameLenIllegal":                  20001,
		"UserDesLenIllegal":                   20002,
		"UpdateThirdPlatformInfoTypeNoExist":  20003,
		"GroupNotExist":                       30001,
		"GroupCountFull":                      30002,
		"CreateGroupCountFull":                30003,
		"SelfGroup":                           30004,
		"NotSelfGroup":                        30005,
		"NotInGroup":                          30007,
		"GroupIdFormatErr":                    30101,
		"GroupSubIdFull":                      30102,
		"GroupPlayerFull":                     30103,
		"GroupNameLenIllegal":                 30104,
		"GroupDesLenIllegal":                  30105,
		"GroupOnlineFull":                     30106,
		"JoinAccessDenied":                    30107,
		"VoiceRoomNameIllegal":                30108,
		"GroupActivityNotExist":               35001,
		"GroupActivityFull":                   35002,
		"GroupActivityTimeIllegal":            35003,
		"GroupActivityExpired":                35004,
		"GroupActivityTitleLenIllegal":        35005,
		"GroupActivityDesLenIllegal":          35006,
		"GroupActivitySignUpRewardLenIllegal": 35007,
		"CaptchaTooFrequently":                90000,
		"CaptchaExpire":                       90001,
		"CaptchaNotEeq":                       90002,
		"PhNotIllegal":                        90003,
		"NotHasCaptcha":                       90004,
		"UrlFormatErr":                        90005,
		"NotFoundUpdateFileNamedIcon":         90006,
		"GeTuiCidIsEmpty":                     90007,
		"GeTuiOpenPushErr":                    90008,
		"GeTuiClosePushErr":                   90009,
		"WechatWaitingLogin":                  100000,
	}
)

func (x Code) Enum() *Code {
	p := new(Code)
	*p = x
	return p
}

func (x Code) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Code) Descriptor() protoreflect.EnumDescriptor {
	return file_code_code_proto_enumTypes[0].Descriptor()
}

func (Code) Type() protoreflect.EnumType {
	return &file_code_code_proto_enumTypes[0]
}

func (x Code) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Code.Descriptor instead.
func (Code) EnumDescriptor() ([]byte, []int) {
	return file_code_code_proto_rawDescGZIP(), []int{0}
}

var File_code_code_proto protoreflect.FileDescriptor

var file_code_code_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x63, 0x6f, 0x64, 0x65, 0x2f, 0x63, 0x6f, 0x64, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x2a, 0xca, 0x0b, 0x0a, 0x04, 0x43, 0x6f, 0x64, 0x65,
	0x12, 0x0e, 0x0a, 0x0a, 0x55, 0x6e, 0x6b, 0x6e, 0x6f, 0x77, 0x6e, 0x45, 0x72, 0x72, 0x10, 0x00,
	0x12, 0x06, 0x0a, 0x02, 0x4f, 0x6b, 0x10, 0x01, 0x12, 0x09, 0x0a, 0x05, 0x45, 0x72, 0x72, 0x6f,
	0x72, 0x10, 0x02, 0x12, 0x0e, 0x0a, 0x0a, 0x49, 0x6e, 0x6e, 0x65, 0x72, 0x45, 0x72, 0x72, 0x6f,
	0x72, 0x10, 0x03, 0x12, 0x11, 0x0a, 0x0d, 0x49, 0x6e, 0x6e, 0x65, 0x72, 0x52, 0x70, 0x63, 0x45,
	0x72, 0x72, 0x6f, 0x72, 0x10, 0x04, 0x12, 0x13, 0x0a, 0x0f, 0x4d, 0x73, 0x67, 0x55, 0x6e, 0x6d,
	0x61, 0x72, 0x73, 0x68, 0x61, 0x6c, 0x45, 0x72, 0x72, 0x10, 0x05, 0x12, 0x16, 0x0a, 0x12, 0x55,
	0x6e, 0x73, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x4d, 0x73, 0x67, 0x54, 0x79, 0x70,
	0x65, 0x10, 0x06, 0x12, 0x0b, 0x0a, 0x07, 0x4e, 0x6f, 0x74, 0x49, 0x6d, 0x70, 0x6c, 0x10, 0x07,
	0x12, 0x0f, 0x0a, 0x0b, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x61, 0x64, 0x45, 0x72, 0x72, 0x10,
	0x08, 0x12, 0x0e, 0x0a, 0x0a, 0x52, 0x70, 0x63, 0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x10,
	0x09, 0x12, 0x11, 0x0a, 0x0d, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x4e, 0x6f, 0x74, 0x45, 0x78, 0x69,
	0x73, 0x74, 0x10, 0x0a, 0x12, 0x12, 0x0a, 0x0e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x46, 0x6f, 0x72,
	0x6d, 0x61, 0x74, 0x45, 0x72, 0x72, 0x10, 0x0b, 0x12, 0x12, 0x0a, 0x0e, 0x4e, 0x6f, 0x74, 0x46,
	0x6f, 0x75, 0x6e, 0x64, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x10, 0x0c, 0x12, 0x15, 0x0a, 0x11,
	0x4e, 0x6f, 0x74, 0x46, 0x6f, 0x75, 0x6e, 0x64, 0x43, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e,
	0x74, 0x10, 0x0d, 0x12, 0x10, 0x0a, 0x0c, 0x4d, 0x73, 0x67, 0x44, 0x65, 0x66, 0x69, 0x6e, 0x65,
	0x45, 0x72, 0x72, 0x10, 0x0e, 0x12, 0x0f, 0x0a, 0x0b, 0x55, 0x69, 0x64, 0x4e, 0x6f, 0x74, 0x46,
	0x6f, 0x75, 0x6e, 0x64, 0x10, 0x0f, 0x12, 0x14, 0x0a, 0x10, 0x4e, 0x6f, 0x74, 0x46, 0x6f, 0x75,
	0x6e, 0x64, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x49, 0x64, 0x10, 0x10, 0x12, 0x14, 0x0a, 0x10,
	0x45, 0x6e, 0x75, 0x6d, 0x4e, 0x6f, 0x74, 0x46, 0x6f, 0x75, 0x6e, 0x64, 0x49, 0x6d, 0x70, 0x6c,
	0x10, 0x11, 0x12, 0x0c, 0x0a, 0x08, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x45, 0x72, 0x72, 0x10, 0x12,
	0x12, 0x0c, 0x0a, 0x08, 0x4e, 0x6f, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x10, 0x13, 0x12, 0x11,
	0x0a, 0x0d, 0x52, 0x6f, 0x62, 0x6f, 0x74, 0x4e, 0x6f, 0x74, 0x45, 0x78, 0x69, 0x73, 0x74, 0x10,
	0x14, 0x12, 0x18, 0x0a, 0x14, 0x41, 0x70, 0x70, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x4e,
	0x65, 0x65, 0x64, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x10, 0x15, 0x12, 0x0f, 0x0a, 0x0b, 0x47,
	0x6d, 0x41, 0x75, 0x74, 0x68, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x10, 0x16, 0x12, 0x12, 0x0a, 0x0e,
	0x41, 0x63, 0x74, 0x6f, 0x72, 0x49, 0x6e, 0x69, 0x74, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x10, 0x17,
	0x12, 0x11, 0x0a, 0x0c, 0x55, 0x73, 0x65, 0x72, 0x4e, 0x6f, 0x74, 0x45, 0x78, 0x69, 0x73, 0x74,
	0x10, 0x90, 0x4e, 0x12, 0x11, 0x0a, 0x0c, 0x43, 0x68, 0x61, 0x74, 0x4c, 0x65, 0x6e, 0x4c, 0x69,
	0x6d, 0x69, 0x74, 0x10, 0x98, 0x4e, 0x12, 0x12, 0x0a, 0x0d, 0x4e, 0x6f, 0x77, 0x4e, 0x6f, 0x74,
	0x49, 0x6e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x10, 0x99, 0x4e, 0x12, 0x1b, 0x0a, 0x16, 0x4d, 0x75,
	0x73, 0x74, 0x51, 0x75, 0x69, 0x74, 0x41, 0x6c, 0x6c, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x42, 0x65,
	0x66, 0x6f, 0x72, 0x65, 0x10, 0x9a, 0x4e, 0x12, 0x15, 0x0a, 0x10, 0x47, 0x61, 0x6d, 0x65, 0x54,
	0x79, 0x70, 0x65, 0x4e, 0x6f, 0x74, 0x45, 0x78, 0x69, 0x73, 0x74, 0x10, 0x9b, 0x4e, 0x12, 0x19,
	0x0a, 0x14, 0x50, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x54, 0x79, 0x70, 0x65, 0x4e, 0x6f,
	0x74, 0x45, 0x78, 0x69, 0x73, 0x74, 0x10, 0x9c, 0x4e, 0x12, 0x18, 0x0a, 0x12, 0x55, 0x73, 0x65,
	0x72, 0x4e, 0x61, 0x6d, 0x65, 0x4c, 0x65, 0x6e, 0x49, 0x6c, 0x6c, 0x65, 0x67, 0x61, 0x6c, 0x10,
	0xa1, 0x9c, 0x01, 0x12, 0x17, 0x0a, 0x11, 0x55, 0x73, 0x65, 0x72, 0x44, 0x65, 0x73, 0x4c, 0x65,
	0x6e, 0x49, 0x6c, 0x6c, 0x65, 0x67, 0x61, 0x6c, 0x10, 0xa2, 0x9c, 0x01, 0x12, 0x28, 0x0a, 0x22,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x68, 0x69, 0x72, 0x64, 0x50, 0x6c, 0x61, 0x74, 0x66,
	0x6f, 0x72, 0x6d, 0x49, 0x6e, 0x66, 0x6f, 0x54, 0x79, 0x70, 0x65, 0x4e, 0x6f, 0x45, 0x78, 0x69,
	0x73, 0x74, 0x10, 0xa3, 0x9c, 0x01, 0x12, 0x13, 0x0a, 0x0d, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x4e,
	0x6f, 0x74, 0x45, 0x78, 0x69, 0x73, 0x74, 0x10, 0xb1, 0xea, 0x01, 0x12, 0x14, 0x0a, 0x0e, 0x47,
	0x72, 0x6f, 0x75, 0x70, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x46, 0x75, 0x6c, 0x6c, 0x10, 0xb2, 0xea,
	0x01, 0x12, 0x1a, 0x0a, 0x14, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x47, 0x72, 0x6f, 0x75, 0x70,
	0x43, 0x6f, 0x75, 0x6e, 0x74, 0x46, 0x75, 0x6c, 0x6c, 0x10, 0xb3, 0xea, 0x01, 0x12, 0x0f, 0x0a,
	0x09, 0x53, 0x65, 0x6c, 0x66, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x10, 0xb4, 0xea, 0x01, 0x12, 0x12,
	0x0a, 0x0c, 0x4e, 0x6f, 0x74, 0x53, 0x65, 0x6c, 0x66, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x10, 0xb5,
	0xea, 0x01, 0x12, 0x10, 0x0a, 0x0a, 0x4e, 0x6f, 0x74, 0x49, 0x6e, 0x47, 0x72, 0x6f, 0x75, 0x70,
	0x10, 0xb7, 0xea, 0x01, 0x12, 0x16, 0x0a, 0x10, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x64, 0x46,
	0x6f, 0x72, 0x6d, 0x61, 0x74, 0x45, 0x72, 0x72, 0x10, 0x95, 0xeb, 0x01, 0x12, 0x14, 0x0a, 0x0e,
	0x47, 0x72, 0x6f, 0x75, 0x70, 0x53, 0x75, 0x62, 0x49, 0x64, 0x46, 0x75, 0x6c, 0x6c, 0x10, 0x96,
	0xeb, 0x01, 0x12, 0x15, 0x0a, 0x0f, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x50, 0x6c, 0x61, 0x79, 0x65,
	0x72, 0x46, 0x75, 0x6c, 0x6c, 0x10, 0x97, 0xeb, 0x01, 0x12, 0x19, 0x0a, 0x13, 0x47, 0x72, 0x6f,
	0x75, 0x70, 0x4e, 0x61, 0x6d, 0x65, 0x4c, 0x65, 0x6e, 0x49, 0x6c, 0x6c, 0x65, 0x67, 0x61, 0x6c,
	0x10, 0x98, 0xeb, 0x01, 0x12, 0x18, 0x0a, 0x12, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x44, 0x65, 0x73,
	0x4c, 0x65, 0x6e, 0x49, 0x6c, 0x6c, 0x65, 0x67, 0x61, 0x6c, 0x10, 0x99, 0xeb, 0x01, 0x12, 0x15,
	0x0a, 0x0f, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x4f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x46, 0x75, 0x6c,
	0x6c, 0x10, 0x9a, 0xeb, 0x01, 0x12, 0x16, 0x0a, 0x10, 0x4a, 0x6f, 0x69, 0x6e, 0x41, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x44, 0x65, 0x6e, 0x69, 0x65, 0x64, 0x10, 0x9b, 0xeb, 0x01, 0x12, 0x1a, 0x0a,
	0x14, 0x56, 0x6f, 0x69, 0x63, 0x65, 0x52, 0x6f, 0x6f, 0x6d, 0x4e, 0x61, 0x6d, 0x65, 0x49, 0x6c,
	0x6c, 0x65, 0x67, 0x61, 0x6c, 0x10, 0x9c, 0xeb, 0x01, 0x12, 0x1b, 0x0a, 0x15, 0x47, 0x72, 0x6f,
	0x75, 0x70, 0x41, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x4e, 0x6f, 0x74, 0x45, 0x78, 0x69,
	0x73, 0x74, 0x10, 0xb9, 0x91, 0x02, 0x12, 0x17, 0x0a, 0x11, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x41,
	0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x46, 0x75, 0x6c, 0x6c, 0x10, 0xba, 0x91, 0x02, 0x12,
	0x1e, 0x0a, 0x18, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x41, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79,
	0x54, 0x69, 0x6d, 0x65, 0x49, 0x6c, 0x6c, 0x65, 0x67, 0x61, 0x6c, 0x10, 0xbb, 0x91, 0x02, 0x12,
	0x1a, 0x0a, 0x14, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x41, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79,
	0x45, 0x78, 0x70, 0x69, 0x72, 0x65, 0x64, 0x10, 0xbc, 0x91, 0x02, 0x12, 0x22, 0x0a, 0x1c, 0x47,
	0x72, 0x6f, 0x75, 0x70, 0x41, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x54, 0x69, 0x74, 0x6c,
	0x65, 0x4c, 0x65, 0x6e, 0x49, 0x6c, 0x6c, 0x65, 0x67, 0x61, 0x6c, 0x10, 0xbd, 0x91, 0x02, 0x12,
	0x20, 0x0a, 0x1a, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x41, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79,
	0x44, 0x65, 0x73, 0x4c, 0x65, 0x6e, 0x49, 0x6c, 0x6c, 0x65, 0x67, 0x61, 0x6c, 0x10, 0xbe, 0x91,
	0x02, 0x12, 0x29, 0x0a, 0x23, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x41, 0x63, 0x74, 0x69, 0x76, 0x69,
	0x74, 0x79, 0x53, 0x69, 0x67, 0x6e, 0x55, 0x70, 0x52, 0x65, 0x77, 0x61, 0x72, 0x64, 0x4c, 0x65,
	0x6e, 0x49, 0x6c, 0x6c, 0x65, 0x67, 0x61, 0x6c, 0x10, 0xbf, 0x91, 0x02, 0x12, 0x1a, 0x0a, 0x14,
	0x43, 0x61, 0x70, 0x74, 0x63, 0x68, 0x61, 0x54, 0x6f, 0x6f, 0x46, 0x72, 0x65, 0x71, 0x75, 0x65,
	0x6e, 0x74, 0x6c, 0x79, 0x10, 0x90, 0xbf, 0x05, 0x12, 0x13, 0x0a, 0x0d, 0x43, 0x61, 0x70, 0x74,
	0x63, 0x68, 0x61, 0x45, 0x78, 0x70, 0x69, 0x72, 0x65, 0x10, 0x91, 0xbf, 0x05, 0x12, 0x13, 0x0a,
	0x0d, 0x43, 0x61, 0x70, 0x74, 0x63, 0x68, 0x61, 0x4e, 0x6f, 0x74, 0x45, 0x65, 0x71, 0x10, 0x92,
	0xbf, 0x05, 0x12, 0x12, 0x0a, 0x0c, 0x50, 0x68, 0x4e, 0x6f, 0x74, 0x49, 0x6c, 0x6c, 0x65, 0x67,
	0x61, 0x6c, 0x10, 0x93, 0xbf, 0x05, 0x12, 0x13, 0x0a, 0x0d, 0x4e, 0x6f, 0x74, 0x48, 0x61, 0x73,
	0x43, 0x61, 0x70, 0x74, 0x63, 0x68, 0x61, 0x10, 0x94, 0xbf, 0x05, 0x12, 0x12, 0x0a, 0x0c, 0x55,
	0x72, 0x6c, 0x46, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x45, 0x72, 0x72, 0x10, 0x95, 0xbf, 0x05, 0x12,
	0x21, 0x0a, 0x1b, 0x4e, 0x6f, 0x74, 0x46, 0x6f, 0x75, 0x6e, 0x64, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x46, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x64, 0x49, 0x63, 0x6f, 0x6e, 0x10, 0x96,
	0xbf, 0x05, 0x12, 0x15, 0x0a, 0x0f, 0x47, 0x65, 0x54, 0x75, 0x69, 0x43, 0x69, 0x64, 0x49, 0x73,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x10, 0x97, 0xbf, 0x05, 0x12, 0x16, 0x0a, 0x10, 0x47, 0x65, 0x54,
	0x75, 0x69, 0x4f, 0x70, 0x65, 0x6e, 0x50, 0x75, 0x73, 0x68, 0x45, 0x72, 0x72, 0x10, 0x98, 0xbf,
	0x05, 0x12, 0x17, 0x0a, 0x11, 0x47, 0x65, 0x54, 0x75, 0x69, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x50,
	0x75, 0x73, 0x68, 0x45, 0x72, 0x72, 0x10, 0x99, 0xbf, 0x05, 0x12, 0x18, 0x0a, 0x12, 0x57, 0x65,
	0x63, 0x68, 0x61, 0x74, 0x57, 0x61, 0x69, 0x74, 0x69, 0x6e, 0x67, 0x4c, 0x6f, 0x67, 0x69, 0x6e,
	0x10, 0xa0, 0x8d, 0x06, 0x42, 0x12, 0x5a, 0x10, 0x6c, 0x61, 0x69, 0x79, 0x61, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x64, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_code_code_proto_rawDescOnce sync.Once
	file_code_code_proto_rawDescData = file_code_code_proto_rawDesc
)

func file_code_code_proto_rawDescGZIP() []byte {
	file_code_code_proto_rawDescOnce.Do(func() {
		file_code_code_proto_rawDescData = protoimpl.X.CompressGZIP(file_code_code_proto_rawDescData)
	})
	return file_code_code_proto_rawDescData
}

var file_code_code_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_code_code_proto_goTypes = []interface{}{
	(Code)(0), // 0: code.Code
}
var file_code_code_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_code_code_proto_init() }
func file_code_code_proto_init() {
	if File_code_code_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_code_code_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_code_code_proto_goTypes,
		DependencyIndexes: file_code_code_proto_depIdxs,
		EnumInfos:         file_code_code_proto_enumTypes,
	}.Build()
	File_code_code_proto = out.File
	file_code_code_proto_rawDesc = nil
	file_code_code_proto_goTypes = nil
	file_code_code_proto_depIdxs = nil
}
