package component

import (
	"github.com/livekit/protocol/livekit"
	"laiya/model/group_model"
	"laiya/model/sub_model"
	"laiya/proto/outer"
	"laiya/share/call"
	"laiya/share/global"
	"laiya/share/livekit_helper"
	"strconv"

	"github.com/asynkron/protoactor-go/actor"
	"google.golang.org/protobuf/proto"
)

//	NotifyGroupMsg 广播
//
// @onlyInGroupParam[0] true: 只通知当前在群内的
func (a *GroupComponent) NotifyGroupMsg(msg proto.Message, onlyNowInGroupParam ...bool) {
	pushMsg := outer.BuildPushPack(msg)
	onlyInGroup := false
	if len(onlyNowInGroupParam) > 0 {
		onlyInGroup = onlyNowInGroupParam[0]
	}
	le := a.GetModel().OnlineSort.Length()
	sendTimes := 0
	a.GetModel().OnlineSort.ForRange(func(v *group_model.OnlineSortT) bool {
		//如果只发送在群内的,那需要判断是否在群中
		if onlyInGroup && !v.IsNowInGroup() {
			return true
		}
		//闭包避免异步的sess覆盖
		//每个推送都分批次发送,避免单条消息处理过长导致超时
		a.Next(func(params ...any) {
			call.Send(a.GetCtx(), params[0].(*actor.PID), pushMsg)
		}, v.GetSess())
		sendTimes++
		return true
	})
	a.GetLogger().Infof("notify group, push msg:%v, player has:%v|send:%v|msg:%v", msg.ProtoReflect().Descriptor().Name(), le, sendTimes, msg)
}

// NotifyChatMsg 广播聊天消息
func (a *GroupComponent) NotifyChatMsg(msg proto.Message, senderId uint64) {
	//必须不为空
	if msg == nil {
		return
	}
	//先构造+存消息
	chatMsg, err := a.GetModel().AddChatMsg(msg, senderId)
	if err != nil {
		a.GetLogger().Error(err)
		return
	}
	//准备push
	pushMsg := outer.BuildPushPack(&outer.PushChatMsg_Push{MsgList: []*outer.ChatMsg{chatMsg}})
	//
	var offlineUids []uint64
	var onlineUids []uint64
	var notGroupOnlineUids []uint64
	//群内用户
	for uid := range a.GetModel().Players {
		//机器人不用管
		if global.IsRobot(uid) {
			continue
		}
		sess, ok := a.GetModel().OnlineSort.Get(uid)
		//离线的加入推送列表
		if !ok || sess.GetSess() == nil {
			offlineUids = append(offlineUids, uid)
			continue
		}
		//在线的走在线消息
		onlineUids = append(onlineUids, uid)
		//每个推送都分批次发送,避免单条消息处理过长导致超时
		a.Next(func(params ...any) {
			call.Send(a.GetCtx(), params[0].(*actor.PID), pushMsg)
		}, sess.GetSess())
	}
	//非群内用户
	a.GetModel().OnlineSort.ForRange(func(v *group_model.OnlineSortT) bool {
		//群内用户在前面已处理,这里不用管
		if a.GetModel().GetPlayer(v.Key()) != nil {
			return true
		}
		notGroupOnlineUids = append(notGroupOnlineUids, v.Key())
		//每个推送都分批次发送,避免单条消息处理过长导致超时
		a.Next(func(params ...any) {
			call.Send(a.GetCtx(), params[0].(*actor.PID), pushMsg)
		}, v.GetSess())
		return true
	})

	//普通聊天消息处理后续的推送和机器人at检查
	if chatMsg.MsgType == outer.ChatMsg_Normal {
		normalMsg, ok := msg.(*outer.ChatMsg_NormalMsg)
		if ok && normalMsg != nil {
			//push相关
			a.checkPush(offlineUids, normalMsg)
			//机器人聊天处理
			a.checkAtRobotMsg(senderId, normalMsg)
		}
	}
	//
	a.GetLogger().Infof("group push msg:%v, offlineUids:%v|onlineUids:%v|notGroupOnlineUids:%v|msg:%v", msg.ProtoReflect().Descriptor().Name(), offlineUids, onlineUids, len(notGroupOnlineUids), chatMsg)
}

// NotifyRobotSay 机器人说话
func (a *GroupComponent) NotifyRobotSay(robotId uint64, toUid uint64, content string) {
	if content == "" {
		a.GetLogger().Errorf("robot say empty, ignore, robot:%v", robotId)
		return
	}
	//获取机器人头像
	robotInfo := sub_model.PlayerHead{Uid: robotId}
	err := robotInfo.Load()
	if err != nil {
		a.GetLogger().Errorf("load robot name and icon err, robot:%v|err:%v", robotId, err)
		return
	}
	//发送聊天
	//警告:一定不能改为Normal类型的聊天,回导致死循环
	a.NotifyChatMsg(&outer.ChatMsg_RobotSayMsg{
		SenderName: robotInfo.Name,
		SenderIcon: robotInfo.Icon,
		Text:       content,
	}, robotId)
}

func (a *GroupComponent) publishActive() {
	//_, err := glin.GetSystem().GetCluster().Publisher().Publish(context.Background(), global.GetGroupActiveTopic(a.GetModel().Gid), &inner.GroupActiveEvent_Notify{Gid: a.GetModel().Gid})
	//if err != nil {
	//	a.GetLogger().Errorf("publish group active event err, err:%v", err)
	//	return
	//}
	//a.GetLogger().Infof("publish group active event success")
}

// NotifyToOnMicUids 只通知给麦上用户
// @onlyInGroupParam[0] true: 只通知当前在群内的
func (a *GroupComponent) NotifyToOnMicUids(msg proto.Message) {
	pushMsg := outer.BuildPushPack(msg)
	if len(a.GetModel().OnMicPlayers) <= 0 {
		a.GetLogger().Info("notify, no onMic users, ignore")
		return
	}
	var offlineUids []uint64
	var onlineUids []uint64
	for _, uid := range a.GetModel().OnMicPlayers {
		sess, ok := a.GetModel().OnlineSort.Get(uid)
		//只发送给在线的
		if !ok || sess.GetSess() == nil {
			offlineUids = append(offlineUids, uid)
			continue
		}
		onlineUids = append(onlineUids, uid)
		a.Next(func(params ...any) {
			call.Send(a.GetCtx(), params[0].(*actor.PID), pushMsg)
		}, sess)
	}
	a.GetLogger().Infof("notify onMic ids, push msg, success:%v|fail:%v|msg:%v", onlineUids, offlineUids, msg)
}

func (a *GroupComponent) syncLivekitRoomState() ([]uint64, []uint64, *livekit.ListParticipantsResponse, error) {
	var afterUids []uint64
	var afterScreenSharingPlayers []uint64
	rsp, err := livekit_helper.GetVoicePlayers(a.GetModel().Gid)
	if err != nil {
		a.GetLogger().Errorf("sync voice players err, will check at next tick, err:%v", err)
		return nil, nil, nil, err
	}
	for _, player := range rsp.GetParticipants() {
		uid, err := strconv.ParseUint(player.Identity, 10, 64)
		if err != nil {
			a.GetLogger().Errorf("sync voice data, identity convert to uint64 err, err:%v", err)
			continue
		}
		afterUids = append(afterUids, uid)
		//查询屏幕分享的用户
		for _, track := range player.Tracks {
			if track.GetSource() == livekit.TrackSource_SCREEN_SHARE || track.GetSource() == livekit.TrackSource_SCREEN_SHARE_AUDIO {
				afterScreenSharingPlayers = append(afterScreenSharingPlayers, uid)
				break
			}
		}
	}
	//取消语音脏标记
	a.GetModel().CleanVoiceDirty()
	//
	beforeUids := a.GetModel().OnMicPlayers
	a.GetModel().OnMicPlayers = afterUids
	a.GetModel().OnScreenSharingPlayers = afterScreenSharingPlayers
	//
	a.GetLogger().Infof("sync voice player success, beforeUids:%v|afterUids:%v|screenUids:%v", beforeUids, afterUids, afterScreenSharingPlayers)
	return afterUids, beforeUids, rsp, nil
}
