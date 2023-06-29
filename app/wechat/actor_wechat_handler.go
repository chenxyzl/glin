package wechat

import (
	"fmt"
	"github.com/chenxyzl/glin/slog"
	"github.com/eatmoreapple/openwechat"
	"laiya/config"
	"laiya/model/wechat_model"
	"laiya/proto/code"
	"laiya/proto/inner"
	"laiya/share/call"
	"laiya/share/datarangers"
	"laiya/share/global"
	"math/rand"
)

func (a *WechatActor) LookingGroup(wechatCtx *openwechat.MessageContext, group *openwechat.Group) {
	gameType := wechat_model.GetGameType(group)
	if gameType == global.DefaultGameType {
		a.GetLogger().Errorf("group not config game type, group:%v", group.NickName)
		delayReply(wechatCtx, config.Get().GameCopywritingConfig.GroupNotConfigGameType, group.NickName)
		return
	}
	reply, cod := call.RequestRemote[*inner.We2HaLookingGroupAvailable_Reply](global.HallUID, &inner.We2HaLookingGroupAvailable_Request{GameType: gameType})
	if cod != code.Code_Ok {
		a.GetLogger().Errorf("req hall looking available group  err, group:%v|cod:%v", group.NickName, cod)
		return
	}
	copywritingConfig := config.Get().GameCopywritingConfig.Get(gameType)
	if len(reply.Groups) == 0 {
		delayReply(wechatCtx, copywritingConfig.NotFoundGroup, group.NickName)
		a.GetLogger().Infof("req hall looking available group, not found, group:%v", group.NickName)
		return
	}
	var items string
	var itemT = copywritingConfig.LookingGroupItem
	for i, availableGroup := range reply.Groups {
		if i != 0 {
			items += "\n"
		}
		items += fmt.Sprintf(itemT, availableGroup.GroupName, availableGroup.GetVoiceRoomName(), config.Get().WebConfig.GetShorUrl(availableGroup.GetGid()), availableGroup.OnMicPlayerCount)
	}
	var total = fmt.Sprintf(copywritingConfig.LookingGroup, len(reply.Groups), items)
	delayReply(wechatCtx, total, group.NickName)
	a.GetLogger().Infof("looking group success, group:%v|reply:%v", group.NickName, reply.GetGroups())
}

func (a *WechatActor) NotifyToWechat(gameType string, content string, filter func(group *wechat_model.WechatGroup) bool) {
	if content == "" {
		a.GetLogger().Errorf("notify to wechat, but content is empty, gameType:%v", gameType)
		return
	}
	var sendGroups []*wechat_model.WechatGroup
	//过滤可发送的
	for _, group := range a.GetModel().Groups {
		if filter != nil && !filter(group) {
			continue
		}
		if group.OpenWechatGroup == nil {
			continue
		}
		groupGameType := group.GetGameType()
		if groupGameType != gameType {
			continue
		}
		if group.OpenWechatGroup == nil {
			continue
		}
		if group.RecentNotifyTimes >= config.Get().ParamsConfig.ContinuityNotifyTimes {
			continue
		}
		sendGroups = append(sendGroups, group)
	}
	//
	if len(sendGroups) > config.Get().ParamsConfig.NotifyGroupCount {
		//打乱顺序
		rand.Shuffle(len(sendGroups), func(i, j int) {
			sendGroups[i], sendGroups[j] = sendGroups[j], sendGroups[i]
		})
		//获取前面n个
		sendGroups = sendGroups[:config.Get().ParamsConfig.NotifyGroupCount]
	}
	//
	for _, group := range sendGroups {
		//发送消息
		group.RecentNotifyTimes++
		go func(temp *openwechat.Group) {
			temp.SendText(content)
			datarangers.SendEvent("wechat_robot_say", map[string]interface{}{
				"wechat_group_name": group.OpenWechatGroup.NickName,
				"reply_to":          0,
			}, nil)
		}(group.OpenWechatGroup)
		slog.Infof("wechat_robot_say, groupName:%v|content:%v", group.OpenWechatGroup.NickName, content)
	}
	//
	a.GetModel().MarkDirty()
	//
	a.GetLogger().Infof("notify to wechat, allCount:%v|sendCount:%v|gameType:%v|groups:%v|content:%v", len(a.GetModel().Groups), len(sendGroups), gameType, sendGroups, content)
}

func (a *WechatActor) UpdateOpenWechatGroupInfo(openWechatGroup *openwechat.Group) {
	//todo 目前group不存库,所以暂时不用管
	group, ok := a.GetModel().Groups[openWechatGroup.UserName]
	if !ok {
		group = &wechat_model.WechatGroup{
			RecentNotifyTimes: 0,
		}
		a.GetModel().Groups[openWechatGroup.UserName] = group
	}
	group.OpenWechatGroup = openWechatGroup
	group.RecentNotifyTimes = 0
	//
	a.GetModel().MarkDirty()
	//data, _ := json.Marshal(openWechatGroup)
	//
	gameType := group.GetGameType()
	slog.Infof("update group info, nickName:%v|userName:%v|remarkName:%v|gameType:%v", openWechatGroup.NickName, openWechatGroup.UserName, openWechatGroup.RemarkName, gameType)
}
