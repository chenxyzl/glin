package wechat

import (
	"github.com/chenxyzl/glin/grain"
	"github.com/chenxyzl/glin/share"
	"github.com/chenxyzl/glin/slog"
	"github.com/eatmoreapple/openwechat"
	"google.golang.org/protobuf/proto"
	"laiya/model/wechat_model"
	"laiya/proto/code"
	"laiya/proto/inner"
	"laiya/share/call"
	"laiya/share/global"
	"laiya/share/utils"
	"time"
)

//var _ iface.IWechatActor = new(WechatActor)

type WechatActor struct {
	grain.BaseActor
	model *wechat_model.Wechat
}

func (a *WechatActor) GetKindType() share.EntityKind {
	return global.WechatKind
}

func NewWechatActor() *WechatActor {
	act := &WechatActor{}
	return act
}

func (a *WechatActor) SetModel(mod *wechat_model.Wechat) {
	if mod.Groups == nil {
		mod.Groups = make(map[string]*wechat_model.WechatGroup)
	}
	a.model = mod
}

func (a *WechatActor) GetModel() *wechat_model.Wechat {
	return a.model
}

func (a *WechatActor) Init() error {
	uid, err := utils.ParseUid(a.GetCtx().Self())
	if err != nil {
		a.GetLogger().Errorf("parse wechat from ctx err, ctx:%v|err:%v", a.GetCtx().Self(), err)
		return err
	}
	//
	a.SetLogger(slog.NewWith("wechat", uid, "kind", a.GetKindType()))
	//
	mod := &wechat_model.Wechat{Uid: uid}
	err = mod.Load()
	if err != nil {
		a.GetLogger().Errorf("load wechat data err, uid:%v|err:%v", uid, err)
		return err
	}
	a.SetModel(mod)
	a.SetReceiver(a.HandleWaitStarWechat)

	go StartWechat(a)

	//
	a.GetLogger().Infof("wechat active")
	return nil
}

func (a *WechatActor) Terminate() {
	//
	err := a.GetModel().Save()
	if err != nil {
		a.GetLogger().Errorf("save wechat err, _id:%v|err:%v", a.GetModel().Uid, err)
	}
	if a.GetModel().Bot != nil {
		a.GetLogger().Warnf("wechat terminal, will exit wechat")
		a.GetModel().Bot.Exit()
		a.GetModel().Bot = nil
	}
	//
	a.GetLogger().Infof("wechat terminal")
}

func (a *WechatActor) Tick() {
	a.BaseActor.Tick()
	if a.GetModel() == nil {
		return
	}
	if a.GetModel().IsDirty() {
		err := a.GetModel().Save()
		if err != nil {
			a.GetLogger().Errorf("tick: save wechat err, err:%v", err)
		} else {
			a.GetLogger().Infof("tick: save wechat success")
		}
	}
}

func (a *WechatActor) WechatLoginSuccess(bot *openwechat.Bot) {
	if a.GetModel().Bot != nil {
		a.GetLogger().Warnf("have a old wechat, will exit old wechat")
		a.GetModel().Bot.Exit()
		a.GetModel().Bot = nil
		a.GetLogger()
	}
	a.GetModel().Bot = bot
	//切换到默认的接收器
	a.ResetDefaultReceiver()
	//
	a.GetLogger().Infof("wechat login success")
}

const waitRestart = time.Second * 3

func (a *WechatActor) WechatExit(bot *openwechat.Bot, err error) {
	//避免退出了不是自己的
	if a.GetModel().Bot == bot {
		a.GetModel().Bot = nil
	}
	//log
	if err != nil {
		a.GetLogger().Errorf("wechat quit with err, err:%v", err)
	} else {
		a.GetLogger().Warnf("wechat quit")
	}
	//异常退出后需要重启
	if err != nil {
		//销毁WechatActor,等待重新登录
		a.GetCtx().Poison(a.GetCtx().Self())
		a.GetLogger().Warnf("wechat actor poison, will restart after %v", waitRestart)
		//
		go func() {
			//销毁后延时重新唤醒
			time.Sleep(waitRestart)
			_, cod := call.RequestRemote[proto.Message](global.WechatUID, &inner.AwakeWechat_Request{})
			if cod != code.Code_Ok {
				a.GetLogger().Error("start wechat err, cod:%v", cod)
			} else {
				a.GetLogger().Warnf("wechat actor restart success")
			}
		}()
	}
}
