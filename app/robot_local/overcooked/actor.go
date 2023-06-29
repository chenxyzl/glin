package overcooked

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/chenxyzl/glin/slog"
	"laiya/model/robot_local_model"
	"laiya/robot_local/simple_robot"
	simplerobotcomponent "laiya/robot_local/simple_robot/component"
	"laiya/robot_local/simple_robot/iface"
)

type OvercookedActor = simple_robot.SimpleRobotActor[*robot_local_model.RobotLocalOverCookedModel]
type OvercookedComponent = simplerobotcomponent.SimpleRobotComponent[*robot_local_model.RobotLocalOverCookedModel]
type IOvercookedComponent = iface.ISimpleRobotComponent

func NewOvercookedActor(gid uint64, parent *actor.PID) *OvercookedActor {
	act := &OvercookedActor{Gid: gid, Parent: parent}
	act.SetLogger(slog.NewWith("robotId", act.GetRobotId(), "gid", gid, "kind", act.GetKindType()))
	return act
}
