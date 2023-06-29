package home

import (
	"github.com/asynkron/protoactor-go/actor"
	"google.golang.org/protobuf/proto"
	"laiya/proto/code"
	"laiya/proto/inner"
	"laiya/proto/outer"
	"laiya/share/call"
)

func (a *PlayerActor) HandleNoUser(ctx actor.Context) {
	ctx.Respond(outer.BuildError(code.Code_UserNotExist))
	ctx.Poison(ctx.Self())
	a.GetLogger().Errorf("user not exist,poison self,actor:%v", ctx.Self())
}

func (a *PlayerActor) HandleUserDeleted(ctx actor.Context) {
	_, ok := ctx.Message().(*inner.Web2HoDeletePlayer_Request)
	if !ok {
		ctx.Respond(outer.BuildError(code.Code_UserNotExist))
	} else {
		ctx.Respond(&inner.Web2HoDeletePlayer_Reply{Account: a.GetModel().Account})
	}
	ctx.Poison(ctx.Self())
	a.GetLogger().Errorf("user has be deleted,poison self,actor:%v", ctx.Self())
}

func (a *PlayerActor) HandleAfterInit(ctx actor.Context) {
	a.once.Do(a.onceAfterInit)
	switch msg := ctx.Message().(type) {
	case *outer.PushPack:
		{
			for _, pid := range a.GetModel().Sess {
				a.Next(func(params ...any) {
					call.Send(a.GetCtx(), params[0].(*actor.PID), params[1])
				}, pid, msg)
			}
		}
	case proto.Message:
		a.DefaultReceiveAfterInit(ctx)
	}
}
func (a *PlayerActor) onceAfterInit() {
	//_, err := glin.GetSystem().GetCluster().SubscribeByPid(global.TopicInviteBroadcast, a.GetCtx().Self())
	//if err != nil {
	//	a.GetLogger().Errorf("subscribe TopicInviteBroadcast err, err:%v", err)
	//}
}
