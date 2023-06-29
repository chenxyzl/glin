package robot_local_model

import (
	"github.com/chenxyzl/glin/model"
	"github.com/chenxyzl/glin/share"
	"laiya/share/global"
)

var _ model.IModel = new(RobotLocalAnimalPartyModel)

type robotLocalAnimalPartModParams struct{}

func (x robotLocalAnimalPartModParams) CoName() string {
	return "robot_local_animal_party"
}

func (x robotLocalAnimalPartModParams) GetKindType() share.EntityKind {
	return global.LocalAnimalPartyKind
}

func (x robotLocalAnimalPartModParams) GetRobotUid() uint64 {
	return global.RobotUidAnimalParty
}

type RobotLocalAnimalPartyModel = LocalSimpleRobotModel[robotLocalAnimalPartModParams]
