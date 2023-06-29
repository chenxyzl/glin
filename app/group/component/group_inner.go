package component

import (
	"fmt"
	"golang.org/x/exp/slices"
	"laiya/config"
	"laiya/model/sub_model"
	"laiya/proto/inner"
	"laiya/proto/outer"
	"laiya/share/getui_helper"
)

func (a *GroupComponent) HandleNotifyGroupChatMsg(req *inner.NotifyGroupChatMsg_Notify) {
	a.NotifyChatMsg(req.GetChatMsg(), req.GetSenderUid())
	a.GetLogger().Infof("recv group chat msg, req:%v", req)
}

func (a *GroupComponent) HandleNotifyGroupActivityMsg(req *inner.NotifyGroupActivityInfo_Notify) {
	a.NotifyChatMsg(req.GetActivityInfo(), req.GetSenderUid())
	a.GetLogger().Infof("recv group chat msg, activity info, req:%v", req)
}

func (a *GroupComponent) HandleNotifyGroupInviteOnMic(req *inner.NotifyGroupInviteOnMic_Notify) {
	var uids []uint64
	for _, uid := range req.GetUids() {
		if !a.GetModel().IsPlayerInGroup(uid) {
			continue
		}
		uids = append(uids, uid)
	}
	var notOnMicUids []uint64
	for _, uid := range uids {
		if slices.Contains(a.GetModel().OnMicPlayers, uid) {
			continue
		}
		notOnMicUids = append(notOnMicUids, uid)
	}
	if len(notOnMicUids) <= 0 {
		a.GetLogger().Infof("HandleNotifyGroupInviteOnMic all uids on mic, uid:%v", uids)
		return
	}
	owner, err := sub_model.LoadPlayerHead(req.GetSenderUid())
	if err != nil {
		a.GetLogger().Errorf("HandleNotifyGroupInviteOnMic Load Owner Data err, uid:%v|err:%v", req.GetSenderUid(), err)
		return
	}
	playerData, err := sub_model.BatchLoadPlayerHead(notOnMicUids, a.GetLogger())
	if err != nil {
		a.GetLogger().Errorf("HandleNotifyGroupInviteOnMic BatchLoadPlayerHead err,uid:%v|err:%v", req.GetSenderUid(), err)
		return
	}
	//@[abc#123,de1f#12387132698,哈1哈哈#7314698734] 大#1家好呀
	pushChatText := fmt.Sprintf(config.Get().GroupPluginActivityConfig.NotOnMicCopywriting, len(notOnMicUids))
	pushChatText += "@[" //开头
	for i, datum := range playerData {
		if i != 0 {
			pushChatText += ","
		}
		pushChatText += fmt.Sprintf("%s#%d", datum.Name, datum.Uid)
	}
	pushChatText += "]" //尾巴
	//推送到聊天框
	a.NotifyChatMsg(&outer.ChatMsg_NormalMsg{Text: pushChatText, SenderName: owner.Name, SenderIcon: owner.Icon}, owner.Uid)
	//推送到个推
	pushGeTuiText := fmt.Sprintf(config.Get().GroupPluginActivityConfig.CallOnMicCopywriting, req.GetActivityTitle())
	getui_helper.PushMsg(notOnMicUids, a.GetModel().Name, pushGeTuiText)
	a.GetLogger().Infof("recv group chat msg, invite onMic, req:%v", req)
}
