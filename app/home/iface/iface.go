package iface

import (
	"github.com/chenxyzl/glin/grain"
	"laiya/model/home_model"
)

type IPlayerActor interface {
	grain.IActor
	GetModel() *home_model.Player
}
type IGroupComponent interface {
	grain.IComponent
	AddSbGroupActive(gid uint64)
	RemoveSbGroupActive(gid uint64)
}
type IPlayerComponent interface {
	grain.IComponent
	TellGroupOnline(gid uint64)
}
