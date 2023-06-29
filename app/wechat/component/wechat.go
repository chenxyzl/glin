package component

import (
	"fmt"
	"github.com/chenxyzl/glin/grain"
	"laiya/config"
	"laiya/model/wechat_model"
	"laiya/proto/code"
	"laiya/proto/inner"
	"laiya/wechat/iface"
	"strings"
)

type WechatComponent struct {
	grain.BaseComponent
	iface.IWechatActor
}

func (a *WechatComponent) BeforeInit() error {
	a.IWechatActor = a.GetEntity().(iface.IWechatActor)
	return nil
}
func (a *WechatComponent) AfterInit() error {
	return nil
}                                           //actor Init完成前
func (a *WechatComponent) BeforeTerminate() {} //actor Terminate完成前
func (a *WechatComponent) AfterTerminate()  {} //actor Terminate完成前
func (a *WechatComponent) Tick()            {}
func (a *WechatComponent) HandleAwakeWechat(notify *inner.AwakeWechat_Request) (*inner.AwakeWechat_Reply, code.Code) {
	a.GetLogger().Infof("wechat awake success")
	return &inner.AwakeWechat_Reply{}, code.Code_Ok
}

func (a *WechatComponent) HandleHa2WechatLookingMsg(notify *inner.Ha2WechatLookingMsg_Notify) {
	copywritingConfig := config.Get().GameCopywritingConfig.Get(notify.GetGameType())
	itemT := copywritingConfig.TimerNotifyWechatItem
	items := ""
	for i, group := range notify.GetGroups() {
		if i != 0 {
			items += "\n"
		}
		items += fmt.Sprintf(itemT, group.GetGroupName(), group.GetVoiceRoomName(), config.Get().WebConfig.GetShorUrl(group.Gid), group.GetPlayerCount())
	}
	var content = fmt.Sprintf(copywritingConfig.TimerNotifyWechat, items)

	a.GetLogger().Infof("wechat get a hall looking notify, gameType:%v|content:%v", notify.GetGameType(), content)
	a.NotifyToWechat(notify.GetGameType(), content, nil)
}

func (a *WechatComponent) HandleRo2WechatLookingMsg(notify *inner.Ro2WechatLookingMsg_Notify) {
	gameType := notify.GetGameType()
	copywritingConfig := config.Get().GameCopywritingConfig.Get(gameType)
	group := notify.GetGroup()
	content := copywritingConfig.LookingNotifyWechat
	content = fmt.Sprintf(content, group.GetGroupName(), group.GetGroupVoiceName(), config.Get().WebConfig.GetShorUrl(group.Gid), 1)
	a.GetLogger().Infof("wechat get a robot looking notify, gameType:%v|gid:%v|content:%v", gameType, group.Gid, content)
	a.NotifyToWechat(notify.GetGameType(), content, func(group *wechat_model.WechatGroup) bool {
		//必须存在
		if group == nil || group.OpenWechatGroup == nil {
			return false
		}
		//KOL的不接收机器人的推送
		if strings.Index(group.OpenWechatGroup.NickName, config.Get().ParamsConfig.KOLStartTag) == 0 {
			return false
		}
		return true
	})
}
