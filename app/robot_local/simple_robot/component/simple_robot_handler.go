package component

import (
	"github.com/chenxyzl/glin/grain"
	"laiya/config"
	"laiya/config_impl"
	"laiya/model/robot_local_model"
	"laiya/proto/inner"
	"laiya/robot_local/simple_robot/iface"
	"laiya/share/call"
	"laiya/share/global"
	"time"
)

type SimpleRobotComponent[T robot_local_model.ILocalSimpleRobotModel] struct {
	grain.BaseComponent
	iface.ISimpleRobotActor[T]
}

func (a *SimpleRobotComponent[T]) BeforeInit() error {
	a.ISimpleRobotActor = a.GetEntity().(iface.ISimpleRobotActor[T])
	return nil
}
func (a *SimpleRobotComponent[T]) AfterInit() error {
	return nil
}                                                   //actor Init完成前
func (a *SimpleRobotComponent[T]) BeforeTerminate() {} //actor Terminate完成前
func (a *SimpleRobotComponent[T]) AfterTerminate()  {} //actor Terminate完成前
func (a *SimpleRobotComponent[T]) Tick()            {}
func (a *SimpleRobotComponent[T]) HandleGr2RoAddRobot(req *inner.Gr2RoAddRobot_Notify) {
	gameType := config.Get().RobotConfig.GetRobotGameType(a.ISimpleRobotActor.GetRobotId())
	copywritingConfig := config.Get().GameCopywritingConfig.Get(gameType)
	call.Send(a.GetCtx(), a.GetHost(), &inner.Ro2GrRobotSayChatMsg_Notify{
		RobotId:    a.ISimpleRobotActor.GetRobotId(),
		ReplyToUid: 0,
		Content:    copywritingConfig.GroupHelper,
	})
}

func (a *SimpleRobotComponent[T]) HandleGr2RoAtRobotChatMsg(req *inner.Gr2RoAtRobotChatMsg_Notify) {

	command := config.Get().RobotConfig.GetCommand(req.GetContent())
	a.GetLogger().Infof("get robot reply msg, command:%v|contend:%v", command, req.GetContent())
	switch command {
	case config_impl.GameCommandUnknown:
		a.doCommandHelper(req.GetSendUid(), req.GetContent())
	case config_impl.GameCommandLookingGroup:
		a.doCommandLookingGroup(req.GetSendUid(), req.GetContent())
	default:
		a.doCommandHelper(req.GetSendUid(), req.GetContent())
	}
}

func (a *SimpleRobotComponent[T]) HandleGr2RoGr2RoNeedInviteBroadcast(req *inner.Gr2RoNeedInviteBroadcast_Notify) {
	//群通知
	now := time.Now().UnixNano()
	inviteBroadcastImCopywritingInterval := 100 * time.Millisecond //默认初始延时100毫秒
	canWechatNotify := now > int64(config.Get().ParamsConfig.WechatInviteBroadcastInterval)+a.GetModel().GetLastWechatInviteBroadcastTime()
	//检查是否可微信群发送
	if canWechatNotify { //能微信群发送
		//设置本次发送时间
		a.GetModel().SetLastWechatInviteBroadcastTime(now)
		//本次的回复0
		call.Send(a.GetCtx(), a.GetHost(), &inner.Ro2GrRobotSayChatMsg_Notify{
			RobotId:    a.ISimpleRobotActor.GetRobotId(),
			ReplyToUid: 0,
			Content:    config.Get().GameCopywritingConfig.WaitingInviteBroadcast0,
		})
		inviteBroadcastImCopywritingInterval += config.Get().ParamsConfig.InviteBroadcastImCopywritingInterval //发送完成加间隔
		//通知微信机器人
		call.NotifyRemote(global.WechatUID, &inner.Ro2WechatLookingMsg_Notify{
			GameType: req.GetGameType(),
			Group: &inner.Ro2WechatLookingMsg_Group{
				Gid:            a.GetHostGid(),
				GroupName:      req.GetGroupName(),
				GroupVoiceName: req.GetGroupVoiceName(),
				PlayerCount:    1,
			},
		})
	}
	//本次的回复1
	a.Delay(inviteBroadcastImCopywritingInterval, func(i ...interface{}) {
		call.Send(a.GetCtx(), a.GetHost(), &inner.Ro2GrRobotSayChatMsg_Notify{
			RobotId:    a.ISimpleRobotActor.GetRobotId(),
			ReplyToUid: 0,
			Content:    config.Get().GameCopywritingConfig.WaitingInviteBroadcast1,
		})
	})
	inviteBroadcastImCopywritingInterval += config.Get().ParamsConfig.InviteBroadcastImCopywritingInterval //发送完成加间隔
	//本次的回复2
	a.Delay(inviteBroadcastImCopywritingInterval, func(i ...interface{}) {
		call.Send(a.GetCtx(), a.GetHost(), &inner.Ro2GrRobotSayChatMsg_Notify{
			RobotId:    a.ISimpleRobotActor.GetRobotId(),
			ReplyToUid: 0,
			Content:    config.Get().GameCopywritingConfig.WaitingInviteBroadcast2,
		})
	})
	inviteBroadcastImCopywritingInterval += config.Get().ParamsConfig.InviteBroadcastImCopywritingInterval //发送完成加间隔
	a.GetLogger().Infof("invite broadcast publish success, trigger by uid:%v|wechat:%v", req.GetUid(), canWechatNotify)
}
func (a *SimpleRobotComponent[T]) HandleGr2RoLookingPlayer(req *inner.Gr2RoLookingPlayer_Notify) {
	//群通知
	now := time.Now().UnixNano()
	canWechatNotify := now > int64(config.Get().ParamsConfig.WechatInviteBroadcastInterval)+a.GetModel().GetLastWechatInviteBroadcastTime()
	//检查是否可微信群发送
	if !canWechatNotify { //不能微信群发送
		//不可微信邀请的提示
		call.Send(a.GetCtx(), a.GetHost(), &inner.Ro2GrRobotSayChatMsg_Notify{
			RobotId:    a.ISimpleRobotActor.GetRobotId(),
			ReplyToUid: 0,
			Content:    config.Get().GameCopywritingConfig.LookingNotifyInterval,
		})
		a.GetLogger().Infof("微信广播cd中, last:%v", a.GetModel().GetLastWechatInviteBroadcastTime())
		return
	}
	//可微信邀请的提示
	call.Send(a.GetCtx(), a.GetHost(), &inner.Ro2GrRobotSayChatMsg_Notify{
		RobotId:    a.ISimpleRobotActor.GetRobotId(),
		ReplyToUid: 0,
		Content:    config.Get().GameCopywritingConfig.LookingNotifyWaiting,
	})
	//延时30秒后发送
	a.Delay(config.Get().ParamsConfig.LookingNotifyWaitingTime, func(i ...interface{}) {
		//二次检查是否可以发送~因为这个人数可能重复触发
		secondNow := time.Now().UnixNano()
		secondCheck := secondNow > int64(config.Get().ParamsConfig.WechatInviteBroadcastInterval)+a.GetModel().GetLastWechatInviteBroadcastTime()
		if !secondCheck { //不能微信群发送
			return
		}
		//设置本次发送时间
		a.GetModel().SetLastWechatInviteBroadcastTime(now)
		//发送到群聊天框
		call.Send(a.GetCtx(), a.GetHost(), &inner.Ro2GrRobotSayChatMsg_Notify{
			RobotId:    a.ISimpleRobotActor.GetRobotId(),
			ReplyToUid: 0,
			Content:    config.Get().GameCopywritingConfig.Get(req.GetGameType()).LookingNotifySuccess,
		})
		//通知微信机器人
		call.NotifyRemote(global.WechatUID, &inner.Ro2WechatLookingMsg_Notify{
			GameType: req.GetGameType(),
			Group: &inner.Ro2WechatLookingMsg_Group{
				Gid:            a.GetHostGid(),
				GroupName:      req.GetGroupName(),
				GroupVoiceName: req.GetGroupVoiceName(),
				PlayerCount:    1,
			},
		})
	})
	a.GetLogger().Infof("微信广播准备中,将在间隔:%v后发送, trigger by uid:%v", config.Get().ParamsConfig.LookingNotifyWaitingTime, req.GetUid())
}

// HandleOnMicPlayerCountChanged 在麦上的用户人数变化
func (a *SimpleRobotComponent[T]) HandleOnMicPlayerCountChanged(req *inner.Gr2RoOnMicPlayerCountChanged_Notify) {
	//如果没有了在麦上的语音人数,要取消通知
	if req.GetOnMicPlayerCount() == 0 {
		actor := a.GetEntity().(grain.IActor)
		actor.CleanAllDelayIds()
		a.GetLogger().Infof("群人数变更,当前无人在麦,清空所有延时定时器, count:%v", req.GetOnMicPlayerCount())
		return
	}
	a.GetLogger().Infof("群人数变更,当前在麦人数, count:%v", req.GetOnMicPlayerCount())
}
