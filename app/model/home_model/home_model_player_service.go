package home_model

import (
	"fmt"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/chenxyzl/glin/slog"
	"google.golang.org/protobuf/proto"
	"laiya/config"
	"laiya/proto/common"
	"laiya/proto/outer"
	"laiya/share/global"
	"time"
)

// OnlineUpdateSess sess在线
// @return 返回是否首个session登录
func (mod *Player) OnlineUpdateSess(loginPlatform common.LoginPlatformType, newSess *actor.PID, repeated func(oldSess *actor.PID)) bool {
	firstSessOnline := len(mod.Sess) == 0
	if mod.Sess[loginPlatform] != nil && !mod.Sess[loginPlatform].Equal(newSess) {
		//
		repeated(mod.Sess[loginPlatform])
		//删除老sess
		delete(mod.Sess, loginPlatform)
	}
	mod.Sess[loginPlatform] = newSess
	mod.MarkDirty()
	//
	if firstSessOnline {
		mod.LastOnlineTime = time.Now().Unix()
		//
		if mod.OnlineDl != nil {
			mod.OnlineDl()
		}
	}
	return firstSessOnline
}

// OfflineCleanSess sess离线
// @return 返回是否全部sess离线
func (mod *Player) OfflineCleanSess(loginPlatform common.LoginPlatformType, newSess *actor.PID) bool {
	if len(mod.Sess) == 0 {
		return true
	}
	//
	if mod.Sess[loginPlatform] != nil && mod.Sess[loginPlatform].Equal(newSess) {
		delete(mod.Sess, loginPlatform)
	}
	allOffline := len(mod.Sess) == 0
	if allOffline {
		//处理下线业务
		mod.LastOfflineTime = time.Now().Unix()
		mod.MarkDirty()
		//
		if mod.OfflineDl != nil {
			mod.OfflineDl()
		}
	}
	return allOffline
}

func (mod *Player) SetUserName(name string) {
	mod.Name = name
	mod.MarkDirty()
}
func (mod *Player) SetUserIcon(icon string) {
	mod.Icon = icon
	mod.MarkDirty()
}
func (mod *Player) SetUserDes(des string) {
	mod.Des = des
	mod.MarkDirty()
}

func (mod *Player) AddGroup(group *Group) {
	mod.Groups[group.Gid] = group
	mod.MarkDirty()
	if mod.RemoveGroupDl != nil {
		mod.RemoveGroupDl(group.Gid)
	}
}

func (mod *Player) RemoveGroup(gid uint64) {
	if _, ok := mod.Groups[gid]; !ok {
		return
	}
	delete(mod.Groups, gid)
	mod.MarkDirty()
	if mod.RemoveGroupDl != nil {
		mod.RemoveGroupDl(gid)
	}
}

func (mod *Player) IsOnline() bool {
	//actor因mod不存在时候销毁自己,此时mod为nil
	if mod == nil {
		return false
	}
	//
	return len(mod.Sess) > 0
}

func (mod *Player) DeletePlayerData() {
	mod.IsDeleted = true
	mod.MarkDirty()
}

func (mod *Player) SetThirdPlatformAccount(id int32, code string) {
	mod.ThirdPlatformAccount[id] = code
	mod.MarkDirty()
}

func (mod *Player) AddLastInviteGroup(gid uint64, tempVoiceRoomVersion uint64) {
	//避免脏数据
	if gid == 0 {
		return
	}
	//检查是否存在
	foundIdx := -1
	for idx := range mod.LastInviteBroadcastGroup {
		if gid == mod.LastInviteBroadcastGroup[idx].Gid {
			foundIdx = idx
			break
		}
	}
	//不存在则追加
	if foundIdx < 0 {
		mod.LastInviteBroadcastGroup = append(mod.LastInviteBroadcastGroup, LastInviteBroadcastGroup{Gid: gid})
		foundIdx = len(mod.LastInviteBroadcastGroup) - 1
	}
	//更新数据
	mod.LastInviteBroadcastGroup[foundIdx].TempVoiceRoomVersion = tempVoiceRoomVersion
	mod.LastInviteBroadcastGroup[foundIdx].Time = time.Now().UnixNano()
	//
	m := 10
	l := len(mod.LastInviteBroadcastGroup)
	//只保留最近的n次
	if l > m {
		mod.LastInviteBroadcastGroup = mod.LastInviteBroadcastGroup[l-m:]
	}
}

func (mod *Player) SelfGroupCount() int {
	count := 0
	for _, group := range mod.Groups {
		if group.IsSelf {
			count++
		}
	}
	return count
}

// GetLastPrivateChatId 获取私聊的聊天Id最大值
func (mod *Player) GetLastPrivateChatId() map[uint64]uint64 {
	return mod.LastChatIdMap
}

func (mod *Player) AddPrivateTextChat(content string, senderId uint64, toUid uint64, uuid uint64, logger slog.Logger) (*outer.ChatMsg, error) {
	msg := &outer.ChatMsg_NormalMsg{Text: content, SenderName: "", SenderIcon: ""}
	return mod.addPrivateChat(msg, senderId, toUid, uuid, logger)
}

func (mod *Player) AddPrivatePicChat(url string, senderId uint64, toUid uint64, uuid uint64, logger slog.Logger) (*outer.ChatMsg, error) {
	msg := &outer.ChatMsg_PicMsg{Url: url}
	return mod.addPrivateChat(msg, senderId, toUid, uuid, logger)
}

func (mod *Player) addPrivateChat(content proto.Message, senderId uint64, toUid uint64, uuid uint64, logger slog.Logger) (*outer.ChatMsg, error) {
	chatMsg, err := outer.BuildChatMsg(content)
	if err != nil {
		return nil, err
	}
	//key
	var rid uint64
	if mod.Uid == senderId {
		//如果自己是发送方--则房间id为目标方
		rid = toUid
	} else if mod.Uid == toUid {
		//如果自己是接收方--则房间id为发送方
		rid = senderId
	} else {
		return nil, fmt.Errorf("not found rid, senderId:%v|toUiend:%v", senderId, toUid)
	}
	//自增
	mod.LastChatIdMap[rid]++
	//赋值
	chatMsg.RId = rid
	chatMsg.MsgId = mod.LastChatIdMap[rid]
	chatMsg.FromType = outer.ChatMsg_Friend
	chatMsg.FromId = senderId
	chatMsg.SenderId = senderId
	chatMsg.ToId = toUid
	chatMsg.Time = time.Now().Unix()
	chatMsg.Uuid = uuid
	//添加到db
	mod.ChatMsgMap[rid] = append(mod.ChatMsgMap[rid], chatMsg)
	l := len(mod.ChatMsgMap[rid])
	max := int(config.Get().ConstConfig.PrivateChatCacheSize)
	if l > max {
		mod.ChatMsgMap[rid] = mod.ChatMsgMap[rid][l-max:]
	}
	//
	mod.MarkDirty()
	return chatMsg, nil
}

func (mod *Player) GetChatMsgByRidBeginId(rid uint64, beginId uint64, logger slog.Logger) ([]*outer.ChatMsg, uint64) {
	//清楚老数据---和群聊不一样-私聊先返回再检查
	defer mod.CleanOldChatMsg(rid, logger)
	//
	chatMsgList := mod.ChatMsgMap[rid]
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
func (mod *Player) CleanOldChatMsg(rid uint64, logger slog.Logger) {
	if len(mod.ChatMsgMap[rid]) == 0 {
		return
	}
	keepSec := int64(config.Get().ConstConfig.PrivateChatKeepTime / time.Second)
	now := time.Now().Unix()
	var deleteMsgIds []uint64
	for {
		//保留最小n条
		l := len(mod.ChatMsgMap[rid])
		if l == 0 || l <= config.Get().ConstConfig.PrivateChatKeepSize {
			break
		}
		//检查首个--未超时者跳过
		if now < mod.ChatMsgMap[rid][0].Time+keepSec {
			break
		}
		//删除老的
		deleteMsgIds = append(deleteMsgIds, mod.ChatMsgMap[rid][0].MsgId)
		mod.ChatMsgMap[rid] = mod.ChatMsgMap[rid][1:]
	}
	if len(deleteMsgIds) > 0 {
		mod.MarkDirty()
		logger.Infof("delete expire private chat msg, rid:%v|ids:%v", rid, deleteMsgIds)
	}
}
