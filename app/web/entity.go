package web

import (
	"github.com/chenxyzl/glin/grain"
	"github.com/chenxyzl/glin/share"
	"google.golang.org/protobuf/reflect/protoreflect"
	"laiya/share/global"
	"laiya/web/component"
	"laiya/web/iface"
	"reflect"
)

var WebEntityIns = NewWebEntity()

type WebEntity struct {
	components        []grain.IComponent
	componentsMapping map[string]grain.IComponent
}

func NewWebEntity() *WebEntity {
	components := make([]grain.IComponent, 0)
	componentsMapping := make(map[string]grain.IComponent)
	components = append(components, &component.WebComponent{})
	componentsMapping[reflect.TypeOf(new(iface.IWebComponent)).Elem().Name()] = &component.WebComponent{}
	return &WebEntity{components: components, componentsMapping: componentsMapping}
}

func (w *WebEntity) GetKindType() share.EntityKind {
	return global.WebKind
}

func (w *WebEntity) GetComponents() []grain.IComponent {
	return w.components
}

func (w *WebEntity) GetComponentByName(s string) grain.IComponent {
	return w.componentsMapping[s]
}
func (a *WebEntity) GetLocalHandler(msgName protoreflect.FullName) *grain.HandlerDef {
	return grain.GetLocalHandlerDef(a.GetKindType(), msgName)
}
