package call

import (
	"github.com/chenxyzl/glin/grain"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"laiya/proto/code"
)

// LocalRequestHandler 调用本地方法--需要传入接收器
func LocalRequestHandler(rpcName protoreflect.FullName, entity grain.IEntity, params ...any) (proto.Message, code.Code) {
	v1, v2 := grain.LocalRequestHandler(rpcName, entity, params...)
	return v1, code.Code(v2)
}
