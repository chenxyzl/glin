package grain

import (
	"fmt"
	"github.com/chenxyzl/glin/share"
	"google.golang.org/protobuf/reflect/protoreflect"
	"reflect"
)

// RegisterComponent 注册handler
// @param remote: 是否允许是远程actor,远程actor的rpc不允许重复
func RegisterComponent[T IEntity, Com IComponent, ICom IComponent](remote bool) {
	var a T
	var t = reflect.New(reflect.TypeOf(a).Elem()).Interface().(T)
	var com Com
	var typ = reflect.TypeOf(t)
	var comTyp = reflect.TypeOf(com)
	var iComTyp = reflect.TypeOf(new(ICom)).Elem()
	if !comTyp.Implements(iComTyp) {
		panic(fmt.Errorf("actor component not imp interface, actor:%v|component:%v|interface:%v", typ.Elem().Name(), comTyp.Elem().Name(), iComTyp.Name()))
	}
	//component
	if components[t.GetKindType()] == nil {
		components[t.GetKindType()] = make([]*ComponentDef, 0)
		componentsMapping[t.GetKindType()] = make(map[string]*ComponentDef)
	}
	for componentsMapping[t.GetKindType()] != nil && componentsMapping[t.GetKindType()][iComTyp.Name()] != nil {
		panic(fmt.Errorf("actor component interface repeated, actor:%v|component:%v|interface:%v", typ.Elem().Name(), comTyp.Elem().Name(), iComTyp.Name()))
	}
	component := &ComponentDef{
		IfName:  iComTyp.Name(),
		ComName: comTyp.Elem().Name(),
		Com:     comTyp.Elem(),
	}
	components[t.GetKindType()] = append(components[t.GetKindType()], component)
	componentsMapping[t.GetKindType()][component.IfName] = component
	//handler
	v := suitableHandlerMethods[T, Com, ICom]()
	for n, c := range v {
		c.entityType = t.GetKindType()
		if remote {
			if _, ok := remoteRpcHandlers[n]; ok {
				panic(fmt.Errorf("handler has registed, n:%v", n))
			}
			remoteRpcHandlers[n] = c
		}
		//
		if _, ok := localHandlersMapping[c.GetEntityKind()]; !ok {
			localHandlersMapping[c.GetEntityKind()] = map[protoreflect.FullName]*HandlerDef{}
		}
		localHandlersMapping[c.GetEntityKind()][n] = c
	}
}

// RegisterPlugin 插件注册
// todo 注册插件~插件会在父actor创建时候跟谁创建(也可能懒加载创建)
// todo 考虑实现插件监听父消息，这样游戏机器人相关的可能会比较好处理
func RegisterPlugin[T IEntity, P IEntity](f GenPluginFunc, remote bool) {
	var a T
	var t = reflect.New(reflect.TypeOf(a).Elem()).Interface().(T)
	var a1 P
	var p = reflect.New(reflect.TypeOf(a1).Elem()).Interface().(P)
	//
	if pluginMapping[t.GetKindType()] == nil {
		pluginMapping[t.GetKindType()] = map[share.EntityKind]GenPluginFunc{}
	}
	if _, ok := pluginMapping[t.GetKindType()][p.GetKindType()]; ok {
		panic(fmt.Errorf("reptead plugin kind, host:%v|plugin:%v", t.GetKindType(), p.GetKindType()))
	}
	pluginMapping[t.GetKindType()][p.GetKindType()] = f
	pluginHandlers := localHandlersMapping[p.GetKindType()]
	hostHandlers := localHandlersMapping[t.GetKindType()]
	if hostHandlers == nil {
		localHandlersMapping[t.GetKindType()] = map[protoreflect.FullName]*HandlerDef{}
	}
	//把插件的方法注册到host中
	for name, ph := range pluginHandlers {
		if _, ok := localHandlersMapping[t.GetKindType()][name]; ok {
			panic(fmt.Errorf("repeated plugin local method name, host:%v|plugin:%v|name:%v", t.GetKindType(), p.GetKindType(), name))
		}
		pluginMappingHandler := &HandlerDef{
			entityType:       t.GetKindType(),
			Entity:           ph.Entity,
			IfName:           ph.IfName,
			Method:           ph.Method,
			pluginEntityType: ph.GetEntityKind(),
		}
		localHandlersMapping[t.GetKindType()][name] = pluginMappingHandler
		//如果是远程方法,还需要注册到远程方法中
		if remote {
			if _, ok := remoteRpcHandlers[name]; ok {
				panic(fmt.Errorf("repeated plugin remote method name, host:%v|plugin:%v|name:%v", t.GetKindType(), p.GetKindType(), name))
			}
			remoteRpcHandlers[name] = pluginMappingHandler
		}
	}
}
