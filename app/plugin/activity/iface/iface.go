package iface

import (
	"github.com/chenxyzl/glin/grain"
	"laiya/model/plugin_model/group_plugin_activity_model"
)

type GroupActorRef interface {
	grain.IActorRef
	GetGid() uint64      //注意线程安全问题
	GetOwnerUid() uint64 //注意线程安全问题
}
type IActivityActor interface {
	grain.IActor
	GetHost() GroupActorRef
	GetModel() *group_plugin_activity_model.PluginActivity
}
type IActivityComponent interface {
	grain.IComponent
}
