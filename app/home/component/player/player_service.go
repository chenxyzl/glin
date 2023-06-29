package player

import (
	"laiya/proto/code"
	"laiya/proto/inner"
	"laiya/proto/outer"
	"laiya/share/call"
)

func (a *PlayerComponent) TellGroupOnline(gid uint64) {
	msg := &inner.Ho2GrPlayerOnline_Request{
		Uid:     a.GetModel().Uid,
		Session: a.GetCtx().Self()}
	_, cod := call.RequestRemote[*inner.Ho2GrPlayerOnline_Reply](gid, msg)
	if cod == code.Code_GroupNotExist {
		a.GetModel().RemoveGroup(gid)
		a.GetLogger().Infof("group not found, may be disband, gid:%v", gid)
		//通知客户端
		call.Send(a.GetCtx(), a.GetCtx().Self(), outer.BuildPushPack(&outer.PushGroupDisband_Push{Gid: gid}))
	} else if cod == code.Code_NotInGroup {
		a.GetModel().RemoveGroup(gid)
		a.GetLogger().Infof("not in group, may be kicked, gid:%v", gid)
	} else if cod != code.Code_Ok {
		a.GetLogger().Errorf("request gid online err, gid:%v|cod:%v", gid, cod)
	}
}
