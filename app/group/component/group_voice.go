package component

import (
	"golang.org/x/exp/slices"
	"laiya/config"
	"laiya/proto/code"
	"laiya/proto/inner"
	"laiya/proto/outer"
	"laiya/share/call"
	"laiya/share/datarangers"
	"laiya/share/livekit_helper"
)

func (a *GroupComponent) checkSyncLivekitRoomState(forceFlag ...bool) {
	//是否强制检查
	force := false
	if len(forceFlag) > 0 && forceFlag[0] {
		force = true
	}
	//语音同步检查
	if !force && !a.GetModel().IsVoiceDirty() {
		return
	}
	//获取语音数据
	afterUids, beforeUids, _, err := a.syncLivekitRoomState()
	if err != nil {
		a.GetLogger().Errorf("sync voice players err, will check at next tick, err:%v", err)
		return
	}
	//更新数据
	beforeOnMicCounts := len(beforeUids)
	afterOnMicCounts := len(afterUids)
	//人数有变化则需要通知到大厅
	needUpdateToHall := beforeOnMicCounts != afterOnMicCounts
	//如果麦上没人则需要恢复语音房间名字
	if afterOnMicCounts == 0 {
		if a.resetVoiceRoomName() {
			needUpdateToHall = true
		}
	}
	//如果人数不同则标记大厅更新排序
	if needUpdateToHall {
		a.UpdateToHall()
	}
	//火山埋点
	datarangers.SendEvent("group_players_mic", map[string]interface{}{"group_id": a.GetModel().Gid, "group_name": a.GetModel().Name, "before": beforeOnMicCounts, "after": afterOnMicCounts}, nil)
	//同步机器人人数变化
	robot := a.GetCurrentRobot()
	if robot != nil {
		call.Send(a.GetCtx(), robot, &inner.Gr2RoOnMicPlayerCountChanged_Notify{OnMicPlayerCount: uint32(len(a.GetModel().OnMicPlayers))})
	}
	//通知客户端
	a.NotifyGroupMsg(&outer.PushLivekitRoomChange_Push{OnMicPlayers: a.GetModel().OnMicPlayers, ScreenSharingUids: a.GetModel().OnScreenSharingPlayers}, true)
	//
	a.GetLogger().Infof("sync voice players data, afterUids:%v", afterUids)
}

// 老的语音变化同步接口-需要触发机器人邀请
func (a *GroupComponent) HandleGroupVoiceStateChange(req *outer.GroupVoiceStateChange_Request) (*outer.GroupVoiceStateChange_Reply, code.Code) {
	//强制同步语音
	a.checkSyncLivekitRoomState(true)
	//获取机器人
	robot := a.GetCurrentRobot()
	if robot == nil {
		a.GetLogger().Warnf("语音人数变更,触发微信群邀请,但对应游戏类型的机器人未激活, gameType:%v", a.GetModel().GameType)
		return &outer.GroupVoiceStateChange_Reply{OnmicPlayers: a.GetModel().OnMicPlayers}, code.Code_Ok
	}
	if len(a.GetModel().OnMicPlayers) == 1 {
		//发送给机器人找人通知
		call.Send(a.GetCtx(), robot, &inner.Gr2RoLookingPlayer_Notify{
			Uid:            req.GetUid(),
			GroupName:      a.GetModel().Name,
			GroupVoiceName: a.GetModel().GetVoiceRoomName(),
			GameType:       a.GetModel().GameType,
		})
	}
	a.GetLogger().Infof("群语音人数变更, uid:%v|count:%v", req.GetUid(), len(a.GetModel().OnMicPlayers))
	return &outer.GroupVoiceStateChange_Reply{OnmicPlayers: a.GetModel().OnMicPlayers}, code.Code_Ok
}

// 语音变化通知
// V1版不要触发机器人邀请
func (a *GroupComponent) HandleGroupVoiceStateChangeV1(req *outer.GroupVoiceStateChangeV1_Request) (*outer.GroupVoiceStateChangeV1_Reply, code.Code) {
	//强制同步语音
	a.checkSyncLivekitRoomState(true)
	//
	a.GetLogger().Infof("群语音人数变更V1, uid:%v|after:%v|screen:%v", req.GetUid(), a.GetModel().OnMicPlayers, a.GetModel().OnScreenSharingPlayers)
	return &outer.GroupVoiceStateChangeV1_Reply{OnmicPlayers: a.GetModel().OnMicPlayers}, code.Code_Ok
}

func (a *GroupComponent) HandleGetGroupVoiceToken(req *outer.GetGroupVoiceToken_Request) (*outer.GetGroupVoiceToken_Reply, code.Code) {
	voiceToken, err := livekit_helper.GetToken(req.GetUid(), req.GetGid())
	if err != nil {
		a.GetLogger().Errorf("req get group voice token err, gid:%v|err:%v", req.GetGid(), err)
		return nil, code.Code_InnerError
	}
	a.GetLogger().Infof("get group voice token success, gid:%v|voiceToken:%v", req.GetGid(), voiceToken)
	return &outer.GetGroupVoiceToken_Reply{VoiceToken: voiceToken}, code.Code_Ok
}

// todo 临时解决多端登录时候, 用户杀了在语音的链接, 语音放人不同步问题
// checkVoicePlayer 麦上有人时候没10s检查一次麦上用户是否在线
func (a *GroupComponent) checkVoicePlayer() {
	if len(a.GetModel().OnMicPlayers) > 0 {
		a.GetModel().TempLastCheckVoicePlayerTick++
	} else {
		a.GetModel().TempLastCheckVoicePlayerTick = 0
	}
	if a.GetModel().TempLastCheckVoicePlayerTick >= config.Get().ParamsConfig.VoiceSyncCheckInterval {
		a.GetModel().MarkVoiceDirty()
		a.GetModel().TempLastCheckVoicePlayerTick = 0
	}
}

func (a *GroupComponent) resetVoiceRoomName() bool {
	beforeVoiceName := a.GetModel().GetVoiceRoomName()
	afterName := config.Get().ConstConfig.DefaultVoiceRoomName
	if beforeVoiceName == afterName {
		return false
	}
	//
	a.GetModel().SetVoiceRoomName(afterName, false)
	//更新到大厅
	a.UpdateToHall()
	//恢复名字
	a.NotifyGroupMsg(&outer.PushGroupVoiceNameChange_Push{
		Gid:           a.GetModel().Gid,
		VoiceRoomName: afterName, //这里获取
	})
	a.GetLogger().Infof("reset voice room name, beforeName:%v|afterName:%v", beforeVoiceName, afterName)
	return true
}

func (a *GroupComponent) HandleMuteAllMic(req *outer.MuteAllMic_Request) (*outer.MuteAllMic_Reply, code.Code) {
	if !a.GetModel().CanMuteOther(req.GetUid()) {
		a.GetLogger().Errorf("mute all mic player err, no access, op:%v", req.GetUid())
		return nil, code.Code_NoAccess
	}
	_, _, rsp, err := a.syncLivekitRoomState()
	if err != nil {
		a.GetLogger().Errorf("mute all mic player err, op:%v|uids:%v|err:%v", req.GetUid(), a.GetModel().OnMicPlayers, err)
		return nil, code.Code_InnerError
	}
	if len(a.GetModel().OnMicPlayers) == 0 {
		a.GetLogger().Infof("mute all mic player ignore, no uids on mic, op:%v", req.GetUid())
		return nil, code.Code_Ok
	}
	livekit_helper.MuteGidVoicePlayerByGid(a.GetModel().Gid, rsp, []uint64{req.GetUid()}, true)
	a.GetLogger().Infof("mute all mic player success, op:%v|uids:%v", req.GetUid(), a.GetModel().OnMicPlayers)
	return &outer.MuteAllMic_Reply{}, code.Code_Ok
}

func (a *GroupComponent) HandleMuteUserMic(req *outer.MuteUserMic_Request) (*outer.MuteUserMic_Reply, code.Code) {
	if !a.GetModel().CanMuteOther(req.GetUid()) {
		a.GetLogger().Errorf("mute target mic player err, no access, op:%v", req.GetUid())
		return nil, code.Code_NoAccess
	}
	_, _, rsp, err := a.syncLivekitRoomState()
	if err != nil {
		a.GetLogger().Errorf("mute target mic player err, op:%v|target:%v|uids:%v|err:%v", req.GetUid(), req.GetTargetUid(), a.GetModel().OnMicPlayers, err)
		return nil, code.Code_InnerError
	}
	if !slices.Contains(a.GetModel().OnMicPlayers, req.GetTargetUid()) {
		a.GetLogger().Infof("mute target mic player ignore, target not on mic, op:%v|target:%v|uids:%v", req.GetUid(), req.GetTargetUid(), a.GetModel().OnMicPlayers)
		return nil, code.Code_Ok
	}
	livekit_helper.MuteGidVoicePlayerByUid(a.GetModel().Gid, rsp, []uint64{req.GetTargetUid()}, true)
	a.GetLogger().Infof("mute target mic player success, op:%v|target:%v|uids:%v", req.GetUid(), req.GetTargetUid(), a.GetModel().OnMicPlayers)
	return &outer.MuteUserMic_Reply{}, code.Code_Ok
}
