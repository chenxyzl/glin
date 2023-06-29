package robot_local_model

import (
	"github.com/chenxyzl/glin/model"
	"github.com/chenxyzl/glin/share"
	"laiya/share/global"
)

var _ model.IModel = new(RobotLocalOverCookedModel)

type robotLocalOverCookedParams struct{}

func (x robotLocalOverCookedParams) CoName() string {
	return "robot_local_overcooked"
}

func (x robotLocalOverCookedParams) GetKindType() share.EntityKind {
	return global.LocalOvercookedKind
}

func (x robotLocalOverCookedParams) GetRobotUid() uint64 {
	return global.RobotUidOvercooked
}

type RobotLocalOverCookedModel = LocalSimpleRobotModel[robotLocalOverCookedParams]
