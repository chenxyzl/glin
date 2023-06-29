package player_group

import (
	"github.com/chenxyzl/glin/grain"
	"laiya/home/iface"
)

type GroupComponent struct {
	grain.BaseComponent
	iface.IPlayerActor
	subscribeGids map[uint64]bool
}

func (a *GroupComponent) BeforeInit() error {
	a.IPlayerActor = a.GetEntity().(iface.IPlayerActor)
	a.subscribeGids = make(map[uint64]bool)

	return nil
}
func (a *GroupComponent) AfterInit() error {
	if a.GetModel() == nil {
		return nil
	}
	//a.GetModel().AddGroupDl = a.AddSbGroupActive
	//a.GetModel().RemoveGroupDl = a.RemoveSbGroupActive
	////下一帧去检查注册群监听
	//a.Next(func(params ...any) {
	//	a.initSbGroupActive()
	//})
	return nil
} //actor Init完成前

func (a *GroupComponent) BeforeTerminate() {
	if a.GetModel() == nil {
		return
	}
	////下一帧去移除
	//a.removeAllSbGroupActive()
}                                         //actor Terminate完成前
func (a *GroupComponent) AfterTerminate() {} //actor Terminate完成前
func (a *GroupComponent) Tick()           {}
