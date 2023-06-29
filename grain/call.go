package grain

import (
	"github.com/chenxyzl/glin/share"
	"github.com/chenxyzl/glin/slog"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"reflect"
)

// LocalRequestHandler 调用本地方法--需要传入接收器
func LocalRequestHandler(rpcName protoreflect.FullName, entity IEntity, params ...any) (proto.Message, int32) {
	handler := entity.GetLocalHandler(rpcName)
	if handler == nil {
		slog.Errorf("handler not fond, entity:%v|rpcName:%v", entity.GetKindType(), rpcName)
		return nil, 12 //code.Code_NotFoundMethod
	}
	//获取接收者
	receiver := entity.GetComponentByName(handler.IfName)
	if receiver == nil {
		slog.Errorf("component not fond, entity:%v|rpcName:%v|component:%v", entity.GetKindType(), rpcName, handler.IfName)
		return nil, 13 //code.Code_NotFoundComponent
	}
	//参数转换
	var paramsValue []reflect.Value
	for _, param := range params {
		paramsValue = append(paramsValue, reflect.ValueOf(param))
	}
	//
	rets, err := share.SafeCall(handler.Method, reflect.ValueOf(receiver), paramsValue...)
	if err != nil {
		slog.Errorf("call local function err, rpcName:%v|err:%v", rpcName, err)
		return nil, 2 //code.Code_Error
	}
	//
	if len(rets) != 2 {
		slog.Errorf("rets must len, now:%v", len(rets))
		return nil, 2 //code.Code_Error
	}
	ret0 := rets[0].Interface()
	msg, ok := ret0.(proto.Message)
	if !ok {
		slog.Errorf("type convert err, want proto.message, ret0:%v", ret0)
		return nil, 2 //code.Code_Error
	}
	ret1 := rets[1].Interface()
	cod, ok := ret1.(protoreflect.Enum)
	if !ok {
		slog.Errorf("type convert err, want protoreflect.Enum.Code, ret1:%v", ret1)
		return nil, 2 //code.Code_Error
	}
	return msg, int32(cod.Number())
}

// LocalNotifyHandler 调用本地方法--需要传入接收器
func LocalNotifyHandler(rpcName protoreflect.FullName, entity IEntity, params ...any) {
	handler := entity.GetLocalHandler(rpcName)
	if handler == nil {
		slog.Errorf("handler not fond, entity:%v|rpcName:%v", entity.GetKindType(), rpcName)
		return
	}
	//获取接收者
	receiver := entity.GetComponentByName(handler.IfName)
	if receiver == nil {
		slog.Errorf("component not fond, entity:%v|rpcName:%v|component:%v", entity.GetKindType(), rpcName, handler.IfName)
		return
	}
	//参数转换
	var paramsValue []reflect.Value
	for _, param := range params {
		paramsValue = append(paramsValue, reflect.ValueOf(param))
	}
	//
	_, err := share.SafeCall(handler.Method, reflect.ValueOf(receiver), paramsValue...)
	if err != nil {
		slog.Errorf("call local function err, rpcName:%v|err:%v", rpcName, err)
		return
	}
}
