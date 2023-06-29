// Package ws todo 后续这里的代码全部改为代码生成器生成
package gate

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/chenxyzl/glin/grain"
	"google.golang.org/protobuf/proto"
	"laiya/proto/code"
	"laiya/proto/outer"
	"laiya/share/call"
	"reflect"
)

// HandleAfterInit 刚建立链接
func (a *SessionActor) HandleAfterInit(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *outer.PushPack:
		{
			a.sess.SafeWritePush(msg) //push类消息直接发送给客户端
			if msg.RpcId == 10003 {   //踢人的消息id为10003
				a.sess.Close(true)
				a.GetLogger().Infof("player kick out with close sess")
			}
		}
	case proto.Message:
		//
		fullName := msg.ProtoReflect().Descriptor().FullName()
		isLocal := grain.IsLocalMethod(a, fullName)
		var rsp proto.Message
		var cod code.Code
		if isLocal {
			rsp, cod = call.LocalRequestHandler(fullName, a, msg)
		} else {
			rsp, cod = call.RequestRemote[proto.Message](a.uid, msg)
		}
		if cod != code.Code_Ok {
			a.GetLogger().Errorf("call rpc err, rpc:%v|targetKind:%v|isLocal:%v|cod:%v", fullName, a.GetKindType(), isLocal, cod)
			ctx.Respond(outer.BuildError(cod))
			return
		}
		ctx.Respond(rsp)
	default:
		a.GetLogger().Errorf("msg unregister in connected,name:%v|msg:%v", reflect.TypeOf(msg).Elem().Name(), msg)
	}
}
