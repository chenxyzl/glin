package iface

import (
	"github.com/chenxyzl/glin/grain"
	"github.com/eatmoreapple/openwechat"
	"laiya/model/wechat_model"
)

type IWechatActor interface {
	grain.IActor
	GetModel() *wechat_model.Wechat
	WechatLoginSuccess(bot *openwechat.Bot)
	WechatExit(bot *openwechat.Bot, err error)
	LookingGroup(wechatCtx *openwechat.MessageContext, group *openwechat.Group)
	NotifyToWechat(gameType string, content string, filter func(group *wechat_model.WechatGroup) bool) //filter返回true才认为可选
	UpdateOpenWechatGroupInfo(group *openwechat.Group)
}

type IWechatComponent interface {
	grain.IComponent
}
