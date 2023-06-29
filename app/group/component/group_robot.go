package component

import (
	"golang.org/x/exp/slices"
	"laiya/common_service"
	"laiya/config"
	"laiya/proto/inner"
	"laiya/proto/outer"
	"laiya/share/call"
	"laiya/share/datarangers"
	"laiya/share/getui_helper"
)

func (a *GroupComponent) HandleRo2GrRobotSayChatMsg(req *inner.Ro2GrRobotSayChatMsg_Notify) {
	//
	a.NotifyRobotSay(req.GetRobotId(), req.GetReplyToUid(), req.GetContent())
	//
	datarangers.SendEvent("robot_say_chat", map[string]interface{}{
		"group_id":   a.GetModel().Gid,
		"group_name": a.GetModel().Name,
		"reply_to":   req.GetReplyToUid(),
	}, nil)

	a.GetLogger().Infof("robot say msg, robotId:%v|content:%v", req.GetRobotId(), req.GetContent())
}

func (a *GroupComponent) checkPush(offlineUids []uint64, msg *outer.ChatMsg_NormalMsg) {
	withAtMsg := common_service.ParseChatMsgFormat(msg.GetText())
	//atList := common_service.ParseChatMsgAt(msg.GetText())
	//
	//var pushToUids []uint64
	////检查推送目标
	//if slices.Contains(atList, 0) {
	//	//如果at的是所有人则全体推送
	//	pushToUids = offlineUids
	//} else {
	//	//如果被@的是离线的人则需要加入推送
	//	for _, uid := range atList {
	//		if slices.Contains(offlineUids, uid) {
	//			pushToUids = append(pushToUids, uid)
	//		}
	//	}
	//}
	//注意拼接的的中文的：
	pushContent := msg.SenderName + "：" + withAtMsg
	getui_helper.PushMsg(offlineUids, a.GetModel().Name, pushContent)
}
func (a *GroupComponent) checkAtRobotMsg(replyToUid1 uint64, msg *outer.ChatMsg_NormalMsg) {
	//解析消息内容
	content1 := common_service.ParseChatMsgPure(msg.GetText())
	atList1 := common_service.ParseChatMsgAt(msg.GetText())
	if len(atList1) == 0 {
		return
	}
	//下一帧发送
	a.Next(func(params ...any) {
		replyToUid := params[0].(uint64)
		atList := params[1].([]uint64)
		content := params[2].(string)
		robotUids := config.Get().RobotConfig.Robots
		for robotId := range robotUids {
			//是否有at自身,at所有人的情况不算
			if !slices.Contains(atList, robotId) {
				continue
			}
			//at的机器人没有激活
			robot := a.GetRobotByUid(robotId)
			if robot == nil {
				a.NotifyRobotSay(robotId, replyToUid, config.Get().GameCopywritingConfig.RobotNotActive)
				a.GetLogger().Warnf("robot not active, robotId:%v|nowGameType:%v", robotId, a.GetModel().GameType)
				continue
			}
			//
			call.Send(a.GetCtx(), robot, &inner.Gr2RoAtRobotChatMsg_Notify{SendUid: replyToUid, Content: content})
			//给机器人发送消息
			a.GetLogger().Infof("at robot msg, robotId:%v|content:%v|realContent:%v", robotId, content, content)
		}
	}, replyToUid1, atList1, content1)
}
