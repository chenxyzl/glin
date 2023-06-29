package grain

import (
	"github.com/chenxyzl/glin/share"
	"github.com/chenxyzl/glin/slog"
	"reflect"
)

var components = make(map[share.EntityKind][]*ComponentDef)
var componentsMapping = make(map[share.EntityKind]map[string]*ComponentDef)

type ComponentDef struct {
	IfName  string       //interface name
	ComName string       //component name
	Com     reflect.Type // receiver of method
}

func GetComponentDef(kind share.EntityKind, ifName string) *ComponentDef {
	if len(componentsMapping) == 0 {
		slog.Warning("components size is 0, may not import actor_register")
	}
	if kindComponents, ok := componentsMapping[kind]; ok {
		if component, ok := kindComponents[ifName]; ok {
			return component
		}
	}
	return nil
}

func GetEntityAllComponentDef(kind share.EntityKind) []*ComponentDef {
	return components[kind]
}
