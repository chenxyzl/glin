package register

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/chenxyzl/glin/grain"
	"github.com/chenxyzl/glin/slog"
	"laiya/robot_local/animal_party"
	"laiya/robot_local/overcooked"
	"laiya/share/global"
)

func init() {
	grain.RegisterRobot(global.RobotUidOvercooked, func(params ...any) grain.IActor {
		return overcooked.NewOvercookedActor(params[0].(uint64), params[1].(*actor.PID))
	})
	grain.RegisterRobot(global.RobotUidAnimalParty, func(params ...any) grain.IActor {
		return animal_party.NewAnimalPartyActor(params[0].(uint64), params[1].(*actor.PID))
	})
	slog.Info("register robot success...")
}
