package component

import (
	"github.com/chenxyzl/glin/grain"
	"laiya/hall/iface"
	"laiya/proto/code"
	"laiya/proto/inner"
	"time"
)

type HallComponent struct {
	grain.BaseComponent
	iface.IHallActor
}

func (a *HallComponent) BeforeInit() error {
	a.IHallActor = a.GetEntity().(iface.IHallActor)
	return nil
}
func (a *HallComponent) AfterInit() error {
	return nil
}                                         //actor Init完成前
func (a *HallComponent) BeforeTerminate() {} //actor Terminate完成前
func (a *HallComponent) AfterTerminate()  {} //actor Terminate完成前
func (a *HallComponent) Tick() {
	if a.GetModel() == nil {
		return
	}
	a.checkTimerNotify(time.Now())
}

func (a *HallComponent) HandleAwakeHall(notify *inner.AwakeHall_Request) (*inner.AwakeHall_Reply, code.Code) {
	a.GetLogger().Infof("hall awake success")
	return &inner.AwakeHall_Reply{}, code.Code_Ok
}
