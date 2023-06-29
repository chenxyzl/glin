package component

import (
	"github.com/chenxyzl/glin/grain"
	"laiya/gate/iface"
	"laiya/proto/code"
	"laiya/proto/outer"
)

type SessionComponent struct {
	grain.BaseComponent
	iface.ISessionActor
}

func (a *SessionComponent) BeforeInit() error {
	a.ISessionActor = a.GetEntity().(iface.ISessionActor)
	return nil
}
func (a *SessionComponent) AfterInit() error {
	a.ResetHeartbeatCheck()
	return nil
}                                            //actor Init完成前
func (a *SessionComponent) BeforeTerminate() {} //actor Terminate完成前
func (a *SessionComponent) AfterTerminate()  {} //actor Terminate完成前
func (a *SessionComponent) HandleHeartbeat(req *outer.Heartbeat_Request) (*outer.Heartbeat_Reply, code.Code) {
	a.ISessionActor.GetLogger().Infof("session heartbeat")
	//重置心跳检测
	a.ResetHeartbeatCheck()
	return nil, code.Code_Ok
}
