package call

import (
	"fmt"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/cluster"
	"github.com/asynkron/protoactor-go/remote"
	"github.com/chenxyzl/glin"
	"github.com/chenxyzl/glin/grain"
	"github.com/chenxyzl/glin/share"
	"github.com/chenxyzl/glin/slog"
	"google.golang.org/protobuf/proto"
	"laiya/proto/code"
	"laiya/proto/outer"
	"laiya/share/global"
	"reflect"
)

func Send(ctx actor.Context, pid *actor.PID, req any) {
	if ctx != nil && pid != nil {
		ctx.Send(pid, req)
	}
}

func Request[T any](pid *actor.PID, req any) (T, code.Code) {
	null := share.GetZero[T]()
	rsp, err := glin.GetSystem().GetCluster().ActorSystem.Root.RequestFuture(pid, req, global.RpcRequestTimeout).Result()
	if err != nil {
		slog.Error(err)
		return null, code.Code_InnerRpcError
	}
	switch t := rsp.(type) {
	case *grain.Error:
		return null, code.Code(t.Code)
	default:
		return t.(T), code.Code_Ok
	}
}

func checkAndSetUid(uid uint64, req proto.Message) error {
	if uid == 0 {
		return nil
	}
	//设置uid
	if !outer.IsMessageInPackage(req) {
		return nil
	}
	msg, ok := req.(outer.UidMsg)
	if !ok { //uid未设置则设置为自己的值
		return nil
	}
	//客户端已设置就不用处理
	if msg.GetUid() == uid {
		return nil
	}
	//客户端设置了且uid和目标的不同 则错误
	if msg.GetUid() != 0 {
		return fmt.Errorf("uid not eq, uid:%v|reqUid:%v", uid, msg.GetUid())
	}
	//强制设置uid
	value := reflect.ValueOf(req).Elem()                                            // 获取结构体的 Value
	field := value.FieldByName("Uid")                                               // 通过字段名获取 Value
	if field.IsValid() && field.CanSet() && field.Type().Kind() == reflect.Uint64 { // 确认该字段存在且可设置
		field.SetUint(uid) // 给字段赋值
	}
	return nil
}

func RequestRemote[T proto.Message](uid uint64, req proto.Message, opts ...cluster.GrainCallOption) (T, code.Code) {
	null := share.GetZero[T]()
	//每个请求都有个反射,性能可以差了,先暂时这么实现--因为用起来简单
	msgName := req.ProtoReflect().Descriptor().FullName()
	//rpcName := strings.TrimSuffix(msgName, "_Request")
	//rpcName = rpcName[strings.LastIndex(rpcName, "_")+1:]
	//handlerName := "Handle" + rpcName
	handler := grain.GetRemoteHandlerDef(msgName)
	if handler == nil {
		slog.Errorf("request not found rpc handler, msgName:%v", msgName)
		return null, code.Code_InnerError
	}
	//uid检查
	if err := checkAndSetUid(uid, req); err != nil {
		slog.Errorf("request uid check err, msgName:%v|err:%v", msgName, err)
		return null, code.Code_Error
	}
	//发送到目标id
	targetId := uid
	entityType := handler.GetEntityKind()
	switch entityType {
	case global.PlayerKind:
		if outer.IsMessageInPackage(req) {
			msg, ok := req.(outer.RedirectUidMsg)
			if ok { //此类消息转发到其他用户的对应接口处理
				targetId = msg.GetRedirectUid()
			}
		}
	case global.GroupKind:
		if outer.IsMessageInPackage(req) {
			msg, ok := req.(outer.GroupMsg)
			if !ok {
				slog.Errorf("request msg can not convert to outer.GroupMsg, kind:%v|msgName:%v", entityType, msgName)
				return null, code.Code_MsgDefineErr
			}
			targetId = msg.GetGid()
		}
	case global.HallKind:
		targetId = global.HallUID
	case global.WechatKind:
		targetId = global.WechatUID
	case global.SessionKind, global.WebKind, global.LocalOvercookedKind, global.LocalAnimalPartyKind:
		slog.Errorf("request not remote kind, kind:%v|msgName:%v", entityType, msgName)
		return null, code.Code_InnerError
	default:
		slog.Errorf("request not switch kind, kind:%v|msgName:%v", entityType, msgName)
		return null, code.Code_InnerError
	}
	//check
	if targetId == 0 {
		slog.Errorf("request targetUid is 0, uid:%v|kind:%v|rpc:%v|req:%v", uid, entityType, msgName, req)
		return null, code.Code_NotFoundTargetId
	}

	//
	opts = append(opts, cluster.WithTimeout(global.RpcRequestTimeout))
	//
	id := fmt.Sprintf("%d", targetId)
	rsp, err := glin.GetSystem().GetCluster().Call(id, string(entityType), req, opts...) //默认3s超时
	if err != nil {
		slog.Errorf("request call rpc err, uid:%v|kind:%v|rpc:%v|req:%v|err:%v", uid, entityType, msgName, req, err)
		return null, code.Code_InnerRpcError
	}
	//
	switch t := rsp.(type) {
	case *grain.Error:
		return null, code.Code(t.Code)
	default:
		return t.(T), code.Code_Ok
	}
}

func NotifyRemote(uid uint64, req proto.Message) {
	//每个请求都有个反射,性能可以差了,先暂时这么实现--因为用起来简单
	msgName := req.ProtoReflect().Descriptor().FullName()
	//rpcName := strings.TrimSuffix(msgName, "_Request")
	//rpcName = rpcName[strings.LastIndex(rpcName, "_")+1:]
	//handlerName := "Handle" + rpcName
	handler := grain.GetRemoteHandlerDef(msgName)
	if handler == nil {
		slog.Errorf("notify not found rpc handler, msgName:%v", msgName)
		return
	}
	//uid检查
	if err := checkAndSetUid(uid, req); err != nil {
		slog.Errorf("notify uid check err, msgName:%v|err:%v", msgName, err)
		return
	}
	//发送到目标id
	targetId := uid
	entityType := handler.GetEntityKind()
	switch entityType {
	case global.PlayerKind:
		if outer.IsMessageInPackage(req) {
			msg, ok := req.(outer.RedirectUidMsg)
			if ok { //此类消息转发到其他用户的对应接口处理
				targetId = msg.GetRedirectUid()
			}
		}
	case global.GroupKind:
		if outer.IsMessageInPackage(req) {
			msg, ok := req.(outer.GroupMsg)
			if !ok {
				slog.Errorf("notify msg can not convert to outer.GroupMsg, kind:%v|msgName:%v", entityType, msgName)
				return
			}
			targetId = msg.GetGid()
		}
	case global.HallKind:
		targetId = global.HallUID
	case global.WechatKind:
		targetId = global.WechatUID
	case global.SessionKind, global.WebKind, global.LocalOvercookedKind, global.LocalAnimalPartyKind:
		slog.Errorf("notify not remote kind, kind:%v|msgName:%v", entityType, msgName)
		return
	default:
		slog.Errorf("notify not switch kind, kind:%v|msgName:%v", entityType, msgName)
		return
	}
	//check
	if targetId == 0 {
		slog.Errorf("notify targetUid is 0, uid:%v|kind:%v|rpc:%v|req:%v", uid, entityType, msgName, req)
		return
	}

	//
	id := fmt.Sprintf("%d", targetId)
	//
	pid := glin.GetSystem().GetCluster().Get(id, string(entityType))
	if pid == nil {
		slog.Errorf("notify get pid from cluster err, uid:%v|kind:%v|rpc:%v|req:%v|err:%v", uid, entityType, msgName, req, remote.ErrUnknownError)
		return
	}
	//Send(glin.GetSystem().GetCluster().ActorSystem.Root,pid,req,opts)
	glin.GetSystem().GetCluster().ActorSystem.Root.Send(pid, req)
}
