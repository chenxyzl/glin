package group_model

import (
	"context"
	"encoding/json"
	"github.com/chenxyzl/glin/model"
	"github.com/chenxyzl/glin/slog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"laiya/proto/outer"
	"laiya/share/algorithm/skiplist"
	"laiya/share/global"
	"laiya/share/mongo_helper"
)

var _ model.IModel = new(Group)

type Room struct {
	RoomId      uint64           `bson:"_id,omitempty"`
	RoomName    string           `bson:"name"`
	RoomDes     string           `bson:"des"` //房间介绍
	CreateTime  int64            `bson:"create_time"`
	GroupTag    uint64           `bson:"group_tag"`
	Chats       []*outer.ChatMsg `bson:"-"` //不入库-聊天消息不重要
	LasIncAgUid uint32           `bson:"-"` //不入库-声网自增id
	RoomPlayers map[uint64]int64 `bson:"-"` //不入库-进入时间
}

type Player struct {
	Uid         uint64 `bson:"_id,omitempty"`
	Permissions int32  `bson:"permissions"` //权限类型--common.GroupDef.Permissions
	//Sess *actor.PID `bson:"-"` //不入库- 在线的session
}

type Group struct {
	model.Dirty                  `bson:"-"`                       //不入库
	Version                      uint64                           `bson:"version"`
	Gid                          uint64                           `bson:"_id,omitempty"`
	Name                         string                           `bson:"name"`
	Icon                         string                           `bson:"icon"`
	Des                          string                           `bson:"des"`
	CreateTime                   int64                            `bson:"create_time"`
	OwnerUid                     uint64                           `bson:"owner_uid"`               //所属玩家
	GameType                     string                           `bson:"game_type"`               //游戏类型
	RunningVersion               uint64                           `bson:"running_version"`         //运行版本号
	IsDisband                    bool                             `bson:"is_disband"`              //删除群只是标记为删除。不会真的删除
	Room                         map[uint64]*Room                 `bson:"rooms"`                   //房间id-房间
	TempVoiceRoomVersion         uint64                           `bson:"temp_voice_room_version"` //临时的语音房间版本,用于过滤重复的摇人--垃圾需求的临时实现
	Players                      map[uint64]*Player               `bson:"players"`                 //群玩家列表
	LastActivityTime             int64                            `bson:"last_activity_time"`      //最后活跃时间-已离开房间为准
	IsPrivate                    bool                             `bson:"is_private"`              //是否是私有房间
	ChatMsg                      []*outer.ChatMsg                 `bson:"chat_msg"`                //聊天消息--最多缓存n条
	LastChatId                   uint64                           `bson:"last_chat_id"`            //不要手动修改-保存时计算
	JoinAccess                   bool                             `bson:"join_access"`             //加入权限
	OnlineSort                   *skiplist.SkipList[*OnlineSortT] `bson:"-"`                       //不入库-玩家在线时间排序的跳表
	VoiceDirty                   bool                             `bson:"-"`                       //不入库-语音脏标记
	OnMicPlayers                 []uint64                         `bson:"-"`                       //不入库-在麦上的玩家
	OnScreenSharingPlayers       []uint64                         `bson:"-"`                       //不入库--在屏幕分享的用户
	TempLastCheckVoicePlayerTick uint64                           `bson:"-"`                       //不入库--临时方案解决语音问题
	NeedUpdateToHall             bool                             `bson:"-"`                       //不入库 更新到大厅
}

func (mod *Group) CoName() string {
	return "group"
}

func (mod *Group) Load() error {
	filter := bson.D{{"_id", mod.Gid}}
	col := mongo_helper.GetCol(mod)
	err := col.FindOne(context.TODO(), filter).Decode(mod)
	return err
}

func (mod *Group) Save() error {
	//db已删除
	if mod == nil {
		return nil
	}
	//机器人群不保存
	if global.IsRobot(mod.OwnerUid) {
		return nil
	}
	//save
	filter := bson.D{{"_id", mod.Gid}}
	col := mongo_helper.GetCol(mod)
	_, err := col.UpdateOne(context.Background(), filter, bson.M{"$set": mod}, options.Update().SetUpsert(true))
	if err != nil {
		tmp, _ := json.Marshal(mod)
		slog.Errorf("save error, save to log file, _id:%v|data:%s", mod.Gid, tmp)
	} else {
		mod.CleanDirty()
	}
	return err
}

func (mod *Group) Delete() error {
	mod.IsDisband = true
	return mod.Save()
}
