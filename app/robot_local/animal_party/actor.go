package animal_party

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/chenxyzl/glin/slog"
	"laiya/model/robot_local_model"
	"laiya/robot_local/simple_robot"
	simplerobotcomponent "laiya/robot_local/simple_robot/component"
	"laiya/robot_local/simple_robot/iface"
)

type AnimalPartyActor = simple_robot.SimpleRobotActor[*robot_local_model.RobotLocalAnimalPartyModel]
type AnimalPartyComponent = simplerobotcomponent.SimpleRobotComponent[*robot_local_model.RobotLocalAnimalPartyModel]
type IAnimalPartyComponent = iface.ISimpleRobotComponent

func NewAnimalPartyActor(gid uint64, parent *actor.PID) *AnimalPartyActor {
	act := &AnimalPartyActor{Gid: gid, Parent: parent}
	act.SetLogger(slog.NewWith("robotId", act.GetRobotId(), "gid", gid, "kind", act.GetKindType()))
	return act
}
