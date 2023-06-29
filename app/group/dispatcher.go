package group

import (
	"github.com/asynkron/protoactor-go/actor"
	"google.golang.org/protobuf/proto"
	"laiya/proto/code"
	"laiya/proto/inner"
	"laiya/proto/outer"
	"laiya/share/call"
	"laiya/share/global"
	"reflect"
)

func (a *GroupActor) HandleCreate(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *inner.Ho2GrCreateGroup_Request:
		a.DefaultReceiveAfterInit(ctx)
	default:
		ctx.Respond(outer.BuildError(code.Code_GroupNotExist))
		a.GetLogger().Errorf("msg unregister in connected, mean group not exist, gid:%v|name:%v|msg:%v", a.GetModel().Gid, reflect.TypeOf(msg).Elem().Name(), msg)
		a.SetModel(nil, a.GetLogger())
		ctx.Poison(ctx.Self())
	}
}
func (a *GroupActor) HandleAfterInit(ctx actor.Context) {
	a.once.Do(a.onceAfterInit)
	a.DefaultReceiveAfterInit(ctx)
}
func (a *GroupActor) HandleGroupHasDisband(ctx actor.Context) {
	msg := ctx.Message()
	//return
	ctx.Respond(outer.BuildError(code.Code_GroupNotExist))

	//todo 临时修复部分脏数据
	//从大厅移除--
	_, cod := call.RequestRemote[proto.Message](global.HallUID, &inner.Gr2HaDeleteGroup_Request{Gid: a.GetModel().Gid})
	if cod != code.Code_Ok {
		a.GetLogger().Errorf("disband group recv, hall delete group err, cod:%v", cod)
	}

	//stop self
	a.GetCtx().Poison(a.GetCtx().Self())
	//log
	a.GetLogger().Errorf("group is disband, mean group not exist, gid:%v|name:%v|msg:%v", a.GetModel().Gid, reflect.TypeOf(msg).Elem().Name(), msg)
}
