package global

import (
	"github.com/chenxyzl/glin/share"
	"time"
)

// 集群actor类型定义
const (
	WebKind                 share.EntityKind = "web"              //web--本地句柄(web节点,非actor)
	SessionKind             share.EntityKind = "session"          //session链接--本地actor(gate节点)
	PlayerKind              share.EntityKind = "player"           //玩家--远程actor(home节点)
	GroupKind               share.EntityKind = "group"            //群--远程actor(group节点)
	HallKind                share.EntityKind = "hall"             //大厅--远程actor(hall节点)
	WechatKind              share.EntityKind = "wechat"           //微信--远程actor(wechat节点)
	LocalOvercookedKind     share.EntityKind = "local_overcooked" //分手厨房--本地actor(group的子actor)
	LocalAnimalPartyKind    share.EntityKind = "local_animal_party"
	GroupPluginActivityKind share.EntityKind = ParamPlugin + "_" + GroupKind + "_" + ParamGroupPluginType //群活动插件--本地actor(group节点)
)

// 特殊actor的一id
const (
	HallUID uint64 = 10000 //目前大厅就开1个

	WechatUID uint64 = 200000 //微信

	RobotBeginUid       uint64 = 300000 //机器人的uid起始范围
	RobotUidChatGpt     uint64 = 300001 //机器人-ChatGpt
	RobotUidOvercooked  uint64 = 300002 //机器人-分手厨房机器人id
	RobotUidAnimalParty uint64 = 300003 //机器人-动物排队
	RobotMaxUid         uint64 = 999999 //机器人的最大id
)

// 一些常用参数
const (
	RpcRequestTimeout        = time.Second * 10 //rpc的超时时间
	PageGroupCount    uint32 = 20               //每页获得群数量
	PageChatCount     int    = 10               //每页最多获取的聊天数量
	PlayerCountOffset uint64 = 33               //群在麦上人数偏移量
	DefaultGameType   string = "其他"             //默认游戏名字
	MaxLoop           int    = 1000             //循环保护
	PrivateChatId            = 1                //私聊的所属的uid
)

// IsRobot 是否是机器人
func IsRobot(uid uint64) bool {
	return uid >= RobotBeginUid && uid <= RobotMaxUid
}

//// 发布订阅的topic
//const (
//	TopicInviteBroadcast string = "TopicInviteBroadcast"
//)
