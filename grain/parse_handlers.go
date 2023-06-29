package grain

import (
	"fmt"
	"github.com/chenxyzl/glin/share"
	"github.com/chenxyzl/glin/slog"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"reflect"
	"strings"
)

var remoteRpcHandlers = make(map[protoreflect.FullName]*HandlerDef) // all handler method
var localHandlersMapping = make(map[share.EntityKind]map[protoreflect.FullName]*HandlerDef)
var pluginMapping = make(map[share.EntityKind]map[share.EntityKind]GenPluginFunc)

var typeOfProtoMsg = reflect.TypeOf(new(proto.Message)).Elem()

type GenPluginFunc = func(host IActorRef) IActor

type HandlerDef struct {
	entityType       share.EntityKind
	pluginEntityType share.EntityKind //如果是插件类型者表示的是宿主的类型
	Entity           reflect.Value    // receiver of method
	IfName           string
	Method           reflect.Method // method stub
}

func (h *HandlerDef) GetPluginKind() share.EntityKind {
	if h == nil {
		return ""
	}
	return h.pluginEntityType
}

func (h *HandlerDef) GetEntityKind() share.EntityKind {
	return h.entityType
}

func GetRemoteHandlerDef(handlerName protoreflect.FullName) *HandlerDef {
	if len(handlerName) == 0 {
		slog.Warning("remote handlerName size is 0, may not import actor_register")
	}
	return remoteRpcHandlers[handlerName]
}

func GetLocalHandlerDef(kind share.EntityKind, handlerName protoreflect.FullName) *HandlerDef {
	if len(handlerName) == 0 {
		slog.Warning("local handlerName size is 0, may not import actor_register")
	}
	if localHandlersMapping[kind] == nil {
		return nil
	}
	return localHandlersMapping[kind][handlerName]
}

func IsLocalMethod(kind IEntity, fullName protoreflect.FullName) bool {
	hs, ok := localHandlersMapping[kind.GetKindType()]
	if !ok {
		return false
	}
	if _, o := hs[fullName]; !o {
		return false
	}
	return true
}

func suitableHandlerMethods[T IEntity, COM IComponent, ICom IComponent]() map[protoreflect.FullName]*HandlerDef {
	var a T
	var t = reflect.New(reflect.TypeOf(a).Elem()).Interface().(T)
	var com COM
	var typ = reflect.TypeOf(t)
	var comTyp = reflect.TypeOf(com)
	var entity = reflect.New(typ.Elem()).Interface()
	var iComTyp = reflect.TypeOf(new(ICom)).Elem()
	methods := make(map[protoreflect.FullName]*HandlerDef)
	for m := 0; m < comTyp.NumMethod(); m++ {
		method := comTyp.Method(m)
		if isHandlerMethod(method) {
			// rewrite handler name
			//mn := method.Name
			param1 := reflect.New(method.Type.In(1).Elem()).Interface()
			fullName := param1.(proto.Message).ProtoReflect().Descriptor().FullName()
			//mn := method.Type.In(1).Elem().Name()
			handler := &HandlerDef{
				entityType: t.GetKindType(),
				Entity:     reflect.ValueOf(entity),
				IfName:     iComTyp.Name(),
				Method:     method,
			}
			if _, ok := methods[fullName]; ok {
				err := fmt.Errorf("repeated handler, %v", fullName)
				panic(err)
			}
			methods[fullName] = handler
		}
	}
	return methods
}

//func isExported(name string) bool {
//	w, _ := utf8.DecodeRuneInString(name)
//	return unicode.IsUpper(w)
//}

// isHandlerMethod decide a method is suitable handler method
func isHandlerMethod(method reflect.Method) bool {
	mt := method.Type
	// Method must be exported.
	if method.PkgPath != "" {
		return false
	}
	// Method needs two or three ins: receiver, protoMessage, other ops params
	if mt.NumIn() < 2 {
		return false
	}

	if mt.NumOut() == 2 { // out == 2; means request/reply
		//匹配参数1的类型 必须是proto 且名字为{mn}Req
		if t1 := mt.In(1); t1.Kind() != reflect.Ptr || !t1.Implements(typeOfProtoMsg) || !strings.HasSuffix(t1.Elem().Name(), "_Request") {
			return false
		}
		//匹配返回值1的类型 必须是proto 且名字必须是{mn}Rsp
		if t1 := mt.Out(0); t1.Kind() != reflect.Ptr || !t1.Implements(typeOfProtoMsg) || !strings.HasSuffix(t1.Elem().Name(), "_Reply") {
			return false
		}
		//匹配返回值2的类型 必须是code
		if t1 := mt.Out(1); t1.Kind() != reflect.Int32 || t1.Name() != "Code" {
			return false
		}
		return true
	} else if mt.NumOut() == 0 { // out == 0; means notify
		//匹配参数1的类型 必须是proto 且名字为{mn}Req
		if t1 := mt.In(1); t1.Kind() != reflect.Ptr || !t1.Implements(typeOfProtoMsg) || !strings.HasSuffix(t1.Elem().Name(), "_Notify") {
			return false
		}
		return true
	}

	return false
}
