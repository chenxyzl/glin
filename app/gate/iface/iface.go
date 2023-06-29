package iface

import (
	"github.com/chenxyzl/glin/grain"
)

type ISessionActor interface {
	grain.IActor
	GetUid() uint64
	ResetHeartbeatCheck()
}
type ISessionComponent interface {
	grain.IComponent
}
