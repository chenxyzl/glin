package iface

import (
	"github.com/chenxyzl/glin/grain"
	"laiya/model/hall_model"
)

type IHallActor interface {
	grain.IActor
	GetModel() *hall_model.Hall
}

type IHallComponent interface {
	grain.IComponent
}
