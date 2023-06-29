package iface

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/chenxyzl/glin/grain"
	"github.com/chenxyzl/glin/slog"
	"google.golang.org/protobuf/proto"
	"laiya/model/group_model"
)

type IGroupActor interface {
	grain.IActor
	GetGid() uint64
	GetOwnerUid() uint64
	GetModel() *group_model.Group
	SetModel(*group_model.Group, slog.Logger)
	GetCurrentRobot() *actor.PID
	GetRobotByUid(uid uint64) *actor.PID
	CheckAddRobot()
	//AtRobotChatMsg(toRobotId uint64, sendUid uint64, content string)
}

// IGroupComponent ..
type IGroupComponent interface {
	grain.IComponent
	UpdateToHall()
	NotifyChatMsg(msg proto.Message, senderId uint64)
	//NotifyRobotSay()
}
