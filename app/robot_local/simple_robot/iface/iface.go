package iface

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/chenxyzl/glin/grain"
	"laiya/model/robot_local_model"
)

type ISimpleRobotOnly[T robot_local_model.ILocalSimpleRobotModel] interface {
	GetRobotId() uint64
	GetHost() *actor.PID
	GetHostGid() uint64
	GetModel() T
}

type ISimpleRobotActor[T robot_local_model.ILocalSimpleRobotModel] interface {
	grain.IActor
	ISimpleRobotOnly[T]
}
type ISimpleRobotComponent interface {
	grain.IComponent
}
