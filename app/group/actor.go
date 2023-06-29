package group

import (
	"errors"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/chenxyzl/glin/grain"
	"github.com/chenxyzl/glin/share"
	"github.com/chenxyzl/glin/slog"
	"go.mongodb.org/mongo-driver/mongo"
	"laiya/config"
	"laiya/group/iface"
	"laiya/model/group_model"
	"laiya/share/global"
	"laiya/share/utils"
	"sync"
)

var _ iface.IGroupActor = new(GroupActor)

type GroupActor struct {
	grain.BaseActor
	model     *group_model.Group
	robotMaps map[uint64]*actor.PID
	once      sync.Once
}

func (a *GroupActor) GetKindType() share.EntityKind {
	return global.GroupKind
}

func (a *GroupActor) GetGid() uint64 {
	if a.GetModel() == nil {
		return 0
	}
	return a.GetModel().Gid
}
func (a *GroupActor) GetOwnerUid() uint64 {
	if a.GetModel() == nil {
		return 0
	}
	return a.GetModel().OwnerUid
}
func NewGroupActor() *GroupActor {
	ac := &GroupActor{robotMaps: make(map[uint64]*actor.PID)}
	ac.AddOptions(grain.WithTerminateIdleTime(config.Get().ConstConfig.GroupTerminateIdleTime))
	ac.AddOptions(grain.WithFilterIdleCheck(func() bool { return ac.GetModel().GetOnlineCount() > 0 }))
	return ac
}

func (a *GroupActor) SetModel(mod *group_model.Group, logger slog.Logger) {
	//删除db
	if mod == nil {
		a.model = mod
		return
	}
	if mod.Room == nil {
		mod.Room = make(map[uint64]*group_model.Room)
	}
	if mod.Players == nil {
		mod.Players = make(map[uint64]*group_model.Player)
	}
	//版本升级检查
	mod.CheckVersion(logger)
	//
	mod.NewOnlineSkipList()
	a.model = mod
}

func (a *GroupActor) GetModel() *group_model.Group {
	return a.model
}

func (a *GroupActor) Init() error {
	//
	//a.BaseActor.Init()
	//解析id
	gid, err := utils.ParseUid(a.GetCtx().Self())
	if err != nil {
		a.GetLogger().Errorf("parse gid from ctx err or group is 0, ctx:%v|gid:%v|err:%v", a.GetCtx().Self(), gid, err)
		return err
	}
	//重写logger
	a.SetLogger(slog.NewWith("gid", gid, "kind", a.GetKindType()))
	//加载数据
	mod := &group_model.Group{Gid: gid}
	err = mod.Load()
	if errors.Is(err, mongo.ErrNoDocuments) {
		a.GetLogger().Infof("group not exist")
		a.SetModel(&group_model.Group{Gid: gid}, a.GetLogger())
		a.SetReceiver(a.HandleCreate)
		//load empty
		return nil
	} else if err != nil {
		a.GetLogger().Errorf("load group data err, group:%v|err:%v", gid, err)
		return err
	}
	//load success
	a.SetModel(mod, a.GetLogger())
	//mod disband will
	if mod.IsDisband {
		a.GetLogger().Infof("group has disband")
		a.SetReceiver(a.HandleGroupHasDisband)
		return nil
	}
	//load success
	a.SetReceiver(a.HandleAfterInit)
	//
	a.GetLogger().Infof("group active")
	return nil
}

func (a *GroupActor) Terminate() {
	err := a.GetModel().Save()
	if err != nil {
		a.GetLogger().Errorf("terminal: save group err, err:%v", err)
	}

	a.GetLogger().Infof("group terminal")
}

func (a *GroupActor) Tick() {
	a.BaseActor.Tick()
	if a.GetModel() == nil {
		return
	}
	//脏数据保存检查
	if a.GetModel().IsDirty() {
		err := a.GetModel().Save()
		if err != nil {
			a.GetLogger().Errorf("tick: save group err, err:%v", err)
		} else {
			a.GetLogger().Infof("tick: save group success")
		}
	}
}

func (a *GroupActor) onceAfterInit() {
	a.CheckAddRobot()
}

func (a *GroupActor) GetCurrentRobot() *actor.PID {
	if len(a.robotMaps) == 0 {
		return nil
	}
	for u, v := range a.robotMaps {
		if config.Get().RobotConfig.GetRobotGameType(u) == a.GetModel().GameType {
			return v
		}
	}
	return nil
}
