package wechat

import (
	"github.com/asynkron/protoactor-go/actor"
	"google.golang.org/protobuf/proto"
	"laiya/proto/code"
	"laiya/proto/inner"
	"laiya/proto/outer"
	"reflect"
	"strings"
)

// HandleWaitStarWechat 等待微信登录
func (a *WechatActor) HandleWaitStarWechat(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *inner.AwakeWechat_Request:
		a.DefaultReceiveAfterInit(ctx)
	case proto.Message:
		reqName := msg.ProtoReflect().Descriptor().FullName()
		if strings.HasSuffix(string(reqName), ".Request") {
			ctx.Respond(outer.BuildError(code.Code_WechatWaitingLogin))
		}
		a.GetLogger().Infof("wechat waiting login, please wait,actor:%v", ctx.Self())
	default:
		ctx.Respond(outer.BuildError(code.Code_Error))
		a.GetLogger().Errorf("msg unregister in connected,gid:%v|name:%v|msg:%v", a.GetModel().Uid, reflect.TypeOf(msg).Elem().Name(), msg)
	}
}
