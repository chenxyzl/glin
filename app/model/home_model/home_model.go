package home_model

import (
	"context"
	"encoding/json"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/chenxyzl/glin/model"
	"github.com/chenxyzl/glin/slog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"laiya/proto/common"
	"laiya/proto/outer"
	"laiya/share/mongo_helper"
)

var _ model.IModel = new(Player)

type Friend struct {
	uid     string `bson:"uid"`      //id
	AddTime uint64 `bson:"add_time"` //添加时间
}

type Group struct {
	Gid     uint64 `bson:"gid"`      //id
	IsSelf  bool   `bson:"is_self"`  //是否是拥有者
	AddTime uint64 `bson:"add_time"` //加入时间
}

type LastInviteBroadcastGroup struct {
	Gid                  uint64 `bson:"gid"`                     //id
	TempVoiceRoomVersion uint64 `bson:"temp_voice_room_version"` //临时的语音房间版本,用于过滤重复的摇人--垃圾需求的临时实现
	Time                 int64  `bson:"time"`
}
type Player struct {
	model.Dirty                  `bson:"-" json:"model_._dirty"`         //不入库
	Uid                          uint64                                  `bson:"_id,omitempty"`
	Account                      string                                  `bson:"account"`    //不要修改
	IsDeleted                    bool                                    `bson:"is_deleted"` //玩家数据已删除
	Name                         string                                  `bson:"name"`
	Icon                         string                                  `bson:"icon"`
	Des                          string                                  `bson:"des"`
	Friends                      map[uint64]*Friend                      `bson:"friends"`
	Groups                       map[uint64]*Group                       `bson:"groups"`                 //群id-群信息
	LastOnlineTime               int64                                   `bson:"last_online_time"`       //单位秒
	LastOfflineTime              int64                                   `bson:"last_offline_time"`      //单位秒
	LastChatIdMap                map[uint64]uint64                       `bson:"last_chatId_map"`        //最后的聊天消息
	ChatMsgMap                   map[uint64][]*outer.ChatMsg             `bson:"chat_msg_map"`           //聊天消息--最多缓存n条
	NowInGroup                   uint64                                  `bson:"-"`                      //不入库--当前所在的群--下线自动离开
	Sess                         map[common.LoginPlatformType]*actor.PID `bson:"-"`                      //不入库--如果要支持多个链接这里需要修改
	ThirdPlatformAccount         map[int32]string                        `bson:"third_platform_account"` //第三方平台账号
	LastSyncUserStateTime        int64                                   `bson:"-"`                      //不入库--上次同步用户状态时间
	LastCheckInviteBroadcastTime int64                                   `bson:"-"`                      //不入库--上次检查摇人功能的时间
	LastInviteBroadcastGroup     []LastInviteBroadcastGroup              `bson:"-"`                      //不入库--最近n次推送
	OnlineDl                     OnlineDl                                `bson:"-"`                      //不入库--首次在线的代理函数
	OfflineDl                    OfflineDl                               `bson:"-"`                      //不入库--全部离线时候的代理函数
	AddGroupDl                   AddGroupDl                              `bson:"-"`                      //不入库--添加群后的代理函数
	RemoveGroupDl                RemoveGroupDl                           `bson:"-"`                      //不入库--移除群后的代理函数
}

func (mod *Player) CoName() string {
	return "player"
}

// Load 一定存在,因为创建是在登录流程里
func (mod *Player) Load() error {
	filter := bson.D{{"_id", mod.Uid}}
	col := mongo_helper.GetCol(mod)
	err := col.FindOne(context.TODO(), filter).Decode(mod)
	return err
}

func (mod *Player) Save() error {
	if mod == nil || mod.Uid <= 0 {
		return nil
	}
	filter := bson.D{{"_id", mod.Uid}}
	col := mongo_helper.GetCol(mod)
	_, err := col.UpdateOne(context.Background(), filter, bson.M{"$set": mod}, options.Update().SetUpsert(true))
	if err != nil {
		tmp, _ := json.Marshal(mod)
		slog.Errorf("save error, save to log file, _id:%v|data:%s", mod.Uid, tmp)
	} else {
		mod.CleanDirty()
	}
	return err
}

func (mod *Player) Delete() error {
	return nil
}
