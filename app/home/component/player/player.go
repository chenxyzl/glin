package player

import (
	"github.com/chenxyzl/glin/grain"
	"laiya/home/iface"
)

type PlayerComponent struct {
	grain.BaseComponent
	iface.IPlayerActor
}

func (a *PlayerComponent) BeforeInit() error {
	a.IPlayerActor = a.GetEntity().(iface.IPlayerActor)
	return nil
}
func (a *PlayerComponent) AfterInit() error {
	return nil
}                                           //actor Init完成前
func (a *PlayerComponent) BeforeTerminate() {} //actor Terminate完成前
func (a *PlayerComponent) AfterTerminate()  {} //actor Terminate完成前
func (a *PlayerComponent) Tick() {
	if a.GetModel() == nil {
		return
	}
	a.checkInviteBroadcast()
}
