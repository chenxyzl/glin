package component

import (
	"github.com/chenxyzl/glin/grain"
	"laiya/plugin/activity/iface"
)

var _ iface.IActivityComponent = new(ActivityComponent)

type ActivityComponent struct {
	grain.BaseComponent
	iface.IActivityActor
}

func (a *ActivityComponent) BeforeInit() error {
	a.IActivityActor = a.GetEntity().(iface.IActivityActor)
	return nil
}
func (a *ActivityComponent) AfterInit() error { return nil } //actor Init完成前
func (a *ActivityComponent) BeforeTerminate() {}             //actor Terminate完成前
func (a *ActivityComponent) AfterTerminate()  {}             //actor Terminate完成前
func (a *ActivityComponent) Tick()            {}
