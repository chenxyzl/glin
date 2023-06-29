package group_model

import (
	"context"
	"errors"
	"github.com/chenxyzl/glin/slog"
	"github.com/chenxyzl/glin/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/protobuf/proto"
	"laiya/config"
	"laiya/proto/code"
	"laiya/proto/common"
	"laiya/proto/outer"
	"laiya/share/global"
	"laiya/share/mongo_helper"
	"time"
)

func (mod *Group) GetPlayer(playerId uint64) *Player {
	player, ok := mod.Players[playerId]
	if !ok {
		return nil
	}
	return player
}
func (mod *Group) AddPlayer(player *Player) {
	mod.Players[player.Uid] = player
	mod.MarkDirty()
}

func (mod *Group) RemovePlayer(uid uint64) {
	_, ok := mod.Players[uid]
	if !ok {
		return
	}
	delete(mod.Players, uid)
	mod.MarkDirty()
}
func (mod *Group) IsPlayerInGroup(uid uint64) bool {
	return mod.Players[uid] != nil
}
func (mod *Group) IsPlayerNowInGroup(uid uint64) bool {
	p, ok := mod.OnlineSort.Get(uid)
	if !ok {
		return false
	}
	return p.IsNowInGroup()
}
func (mod *Group) IsPlayerOnline(uid uint64) bool {
	if global.IsRobot(uid) {
		return true
	}
	_, ok := mod.OnlineSort.Get(uid)
	return ok
}

func (mod *Group) PlayerCount() int {
	return len(mod.Players)
}

// GetOnlinePlayerCount 获取在线人数(只获取群内玩家-即已收藏用户)
func (mod *Group) GetOnlinePlayerCount() uint32 {
	count := uint32(0)
	mod.OnlineSort.ForRange(func(v *OnlineSortT) bool {
		if mod.GetPlayer(v.Key()) != nil {
			count++
		}
		return true
	})
	return count
}

// CreateRoom 创建房间
// @param imSave 是否立即保存
func (mod *Group) CreateRoom(name string, roomDes string) (*Room, code.Code) {
	roomId, err := uuid.GenerateIntervalSubId(mod.Gid, func(u uint64) bool {
		for id := range mod.Room {
			if id == u {
				return true
			}
		}
		return false
	})
	if err != nil {
		if errors.Is(err, uuid.FormatErr) {
			return nil, code.Code_GroupIdFormatErr
		}
		if errors.Is(err, uuid.AllUsedErr) {
			return nil, code.Code_GroupSubIdFull
		}
		return nil, code.Code_InnerError
	}

	room := &Room{
		RoomId:      roomId,
		RoomName:    name,
		RoomDes:     roomDes,
		CreateTime:  time.Now().UnixMilli(),
		GroupTag:    0,
		RoomPlayers: map[uint64]int64{},
	}
	return room, code.Code_Ok
}

func (mod *Group) RemoveRoom(roomId uint64) {
	if _, ok := mod.Room[roomId]; !ok {
		return
	}
	delete(mod.Room, roomId)
	mod.MarkDirty()
}

func (r *Room) AddChat(chat *outer.ChatMsg) {
	r.Chats = append(r.Chats, chat)
	l := len(r.Chats)
	max := int(config.Get().ConstConfig.MaxChatHistory)
	if l > max {
		r.Chats = r.Chats[l-max:]
	}
}

func (mod *Group) GetOnMicPlayerCount() int {
	return len(mod.OnMicPlayers)
}

func (mod *Group) GetHallSortScore() uint64 {
	//官方标记
	var officialGroupTag = uint64(0)
	if config.Get().GmConfig.IsOfficialGroup(mod.Gid) {
		officialGroupTag = 1
	}
	//时间偏移量
	const offset = 1672502400 //2023-01-01:00:00:00 GMT
	//位运算
	// | 1 Bit Empty| 8 Bit Free | 16 Bit OnMicPlayers | 1 Bit officialGroupTag | 32 Bit Timestamp |
	return (uint64(len(mod.OnMicPlayers)) << global.PlayerCountOffset) | officialGroupTag<<32 | uint64(time.Now().Unix()-offset)
}

func (mod *Group) GetRoom(roomId uint64) *Room {
	room, ok := mod.Room[roomId]
	if !ok {
		return nil
	}
	return room
}

func (mod *Group) SetGroupName(name string) {
	mod.Name = name
	mod.MarkDirty()
}
func (mod *Group) SetGroupIcon(icon string) {
	mod.Icon = icon
	mod.MarkDirty()
}
func (mod *Group) SetGroupDes(des string) {
	mod.Des = des
	mod.MarkDirty()
}
func (mod *Group) SetGroupGameType(gameType string) {
	mod.GameType = gameType
	mod.MarkDirty()
}
func (mod *Group) GetOnlineCount(onlyInGroupOps ...bool) int {
	//actor因mod不存在时候销毁自己,此时mod为nil
	if mod == nil {
		return 0
	}
	//
	onlyInGroup := false
	if len(onlyInGroupOps) > 0 && onlyInGroupOps[0] == true {
		onlyInGroup = true
	}
	if onlyInGroup {
		onlineCount := 0
		mod.OnlineSort.ForRange(func(v *OnlineSortT) bool {
			if mod.IsPlayerInGroup(v.Key()) {
				onlineCount++
			}
			return true
		})
		return onlineCount
	} else {
		return mod.OnlineSort.Length()
	}
}
func (mod *Group) GetPlayers() []uint64 {
	var ret []uint64
	for _, player := range mod.Players {
		ret = append(ret, player.Uid)
	}
	return ret
}

func (mod *Group) SetJoinAccess(enable bool) {
	mod.JoinAccess = enable
	mod.MarkDirty()
}

// GetVoiceRoomName 获取语音房间的名字
// todo 目前只有1个语音房间，临时这么处理
func (mod *Group) GetVoiceRoomName() string {
	voiceRoomName := ""
	for _, v := range mod.Room {
		voiceRoomName = v.RoomName
		break
	}
	return voiceRoomName
}

// SetVoiceRoomName 设置语音房间的名字
// todo 目前只有1个语音房间，临时这么处理
func (mod *Group) SetVoiceRoomName(name string, newVoiceRoomVersion bool) {
	for _, v := range mod.Room {
		v.RoomName = name
	}
	if newVoiceRoomVersion {
		mod.TempVoiceRoomVersion++
	}
	mod.MarkDirty()
}

func (mod *Group) MarkVoiceDirty() {
	mod.VoiceDirty = true
}
func (mod *Group) CleanVoiceDirty() {
	mod.VoiceDirty = false
}
func (mod *Group) IsVoiceDirty() bool {
	return mod.VoiceDirty
}

func (mod *Group) AddChatMsg(msg proto.Message, senderId uint64) (*outer.ChatMsg, error) {
	chatMsg, err := outer.BuildChatMsg(msg)
	if err != nil {
		return nil, err
	}
	//自增
	mod.LastChatId++
	//立即更新last_chat_id(用户登录批量查询这个字段时候需要使用)
	filter := bson.M{"_id": mod.Gid}
	update := bson.M{"$set": bson.M{"last_chat_id": mod.LastChatId}}
	_, err = mongo_helper.GetCol(mod).UpdateOne(context.Background(), filter, update)
	if err != nil {
		//恢复之前自增的消息id
		mod.LastChatId--
		return nil, err
	}
	//赋值
	chatMsg.RId = mod.Gid
	chatMsg.MsgId = mod.LastChatId
	chatMsg.FromType = outer.ChatMsg_Group
	chatMsg.FromId = mod.Gid
	chatMsg.SenderId = senderId
	chatMsg.ToId = mod.Gid
	chatMsg.Time = time.Now().Unix()
	chatMsg.Uuid = uuid.Generate()
	//添加到db
	mod.ChatMsg = append(mod.ChatMsg, chatMsg)
	l := len(mod.ChatMsg)
	max := int(config.Get().ConstConfig.ChatCacheSize)
	if l > max {
		mod.ChatMsg = mod.ChatMsg[l-max:]
	}
	//
	mod.MarkDirty()
	return chatMsg, nil
}

func (mod *Group) GetChatMsgByBeginId(beginId uint64, logger slog.Logger) ([]*outer.ChatMsg, uint64) {
	//清楚老数据---和私聊不一样-群聊先检查再返回
	mod.CleanOldChatMsg(logger)
	//
	chatMsgList := mod.ChatMsg
	l := len(chatMsgList)
	if l == 0 {
		return nil, 0
	}
	last := chatMsgList[l-1]
	//已经到最新就不返回了
	if beginId >= last.MsgId {
		return nil, last.MsgId
	}
	//0则获取最新的10条
	if beginId == 0 {
		if l > global.PageChatCount {
			beginId = uint64(l - global.PageChatCount)
		}
		return chatMsgList[beginId:], last.MsgId
	}
	//todo 后续再换更高效的二分法
	var msgs []*outer.ChatMsg
	for _, msg := range chatMsgList {
		if msg.MsgId <= beginId {
			continue
		}
		msgs = append(msgs, msg)
		if len(msgs) >= global.PageChatCount {
			break
		}
	}
	return msgs, last.MsgId
}

func (mod *Group) CleanOldChatMsg(logger slog.Logger) {
	if len(mod.ChatMsg) == 0 {
		return
	}
	keepSec := int64(config.Get().ConstConfig.GroupChatKeepTime / time.Second)
	now := time.Now().Unix()
	var deleteMsgIds []uint64
	for {
		//保留最小n条
		l := len(mod.ChatMsg)
		if l == 0 {
			break
		}
		//检查首个--未超时者跳过
		if now < mod.ChatMsg[0].Time+keepSec {
			break
		}
		//删除老的
		deleteMsgIds = append(deleteMsgIds, mod.ChatMsg[0].MsgId)
		mod.ChatMsg = mod.ChatMsg[1:]
	}
	if len(deleteMsgIds) > 0 {
		mod.MarkDirty()
		logger.Infof("delete expire group chat msg, ids:%v", deleteMsgIds)
	}
}

func (mod *Group) SetUserPermissions(targetUid uint64, permissions common.GroupPermissionsDef_Type) code.Code {
	if mod.GetPlayer(targetUid) == nil {
		return code.Code_NotInGroup
	}
	mod.GetPlayer(targetUid).Permissions = int32(permissions.Number())
	return code.Code_Ok
}

func (mod *Group) GetUserPermissions(targetUid uint64) common.GroupPermissionsDef_Type {
	if mod.GetPlayer(targetUid) == nil {
		return common.GroupPermissionsDef_UnknownType
	}
	//所有者默认管理员权限
	if targetUid == mod.OwnerUid {
		return common.GroupPermissionsDef_Admin
	}
	return common.GroupPermissionsDef_Type(mod.GetPlayer(targetUid).Permissions)
}

// CanPermissionsConfig 是否可以设置权限
func (mod *Group) CanPermissionsConfig(opUid uint64) bool {
	return config.Get().GroupPermissionsConfig.HasAccess(mod.GetUserPermissions(opUid), common.GroupPermissionsDef_ChangePermissions)
}

// CanMuteOther 是否可以静音其他用户
func (mod *Group) CanMuteOther(opUid uint64) bool {
	return config.Get().GroupPermissionsConfig.HasAccess(mod.GetUserPermissions(opUid), common.GroupPermissionsDef_MuteOther)
}
