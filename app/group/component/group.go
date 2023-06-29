package component

import (
	"github.com/chenxyzl/glin/grain"
	"laiya/group/iface"
)

type GroupComponent struct {
	grain.BaseComponent
	iface.IGroupActor
}

func (a *GroupComponent) BeforeInit() error {
	a.IGroupActor = a.GetEntity().(iface.IGroupActor)
	return nil
}
func (a *GroupComponent) AfterInit() error {
	a.Next(func(params ...any) { a.publishActive() })
	return nil
}                                          //actor Init完成前
func (a *GroupComponent) BeforeTerminate() {} //actor Terminate完成前
func (a *GroupComponent) AfterTerminate()  {} //actor Terminate完成前
func (a *GroupComponent) Tick() {
	if a.GetModel() == nil {
		return
	}
	//检查麦上用户状态
	a.checkVoicePlayer()
	//如果状态同步
	a.checkSyncLivekitRoomState()
}
