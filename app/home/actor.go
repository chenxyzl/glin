package home

import (
	"errors"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/chenxyzl/glin/grain"
	"github.com/chenxyzl/glin/share"
	"github.com/chenxyzl/glin/slog"
	"go.mongodb.org/mongo-driver/mongo"
	"laiya/config"
	"laiya/home/iface"
	"laiya/model/home_model"
	"laiya/proto/common"
	"laiya/proto/outer"
	"laiya/share/global"
	"laiya/share/utils"
	"sync"
)

var _ iface.IPlayerActor = new(PlayerActor)

type PlayerActor struct {
	grain.BaseActor
	model *home_model.Player
	once  sync.Once
}

func (a *PlayerActor) GetKindType() share.EntityKind {
	return global.PlayerKind
}

func NewPlayerActor() *PlayerActor {
	act := &PlayerActor{}
	act.AddOptions(grain.WithTerminateIdleTime(config.Get().ConstConfig.PlayerTerminateIdleTime))
	act.AddOptions(grain.WithFilterIdleCheck(func() bool { return act.GetModel().IsOnline() }))
	return act
}

func (a *PlayerActor) SetModel(mod *home_model.Player) {
	if mod.Groups == nil {
		mod.Groups = make(map[uint64]*home_model.Group)
	}
	if mod.Friends == nil {
		mod.Friends = make(map[uint64]*home_model.Friend)
	}
	if mod.Sess == nil {
		mod.Sess = make(map[common.LoginPlatformType]*actor.PID)
	}
	if mod.ThirdPlatformAccount == nil {
		mod.ThirdPlatformAccount = make(map[int32]string)
	}
	if mod.LastInviteBroadcastGroup == nil {
		mod.LastInviteBroadcastGroup = make([]home_model.LastInviteBroadcastGroup, 0)
	}
	if mod.LastChatIdMap == nil {
		mod.LastChatIdMap = make(map[uint64]uint64)
	}
	if mod.ChatMsgMap == nil {
		mod.ChatMsgMap = make(map[uint64][]*outer.ChatMsg)
	}
	a.model = mod
}

func (a *PlayerActor) GetModel() *home_model.Player {
	return a.model
}

func (a *PlayerActor) Init() error {
	uid, err := utils.ParseUid(a.GetCtx().Self())
	if err != nil {
		a.GetLogger().Errorf("parse uid from ctx err, ctx:%v|err:%v", a.GetCtx().Self(), err)
		return err
	}
	a.SetLogger(slog.NewWith("uid", uid, "kind", a.GetKindType()))
	player := &home_model.Player{Uid: uid}
	err = player.Load()
	if errors.Is(err, mongo.ErrNoDocuments) {
		a.GetLogger().Errorf("user not exist ,uid:%v|err:%v", uid, err)
		a.SetReceiver(a.HandleNoUser)
		return nil
	} else if err != nil {
		a.GetLogger().Errorf("load player data err,uid:%v|err:%v", uid, err)
		return err
	}
	//
	a.SetModel(player)
	//
	if player.IsDeleted {
		a.GetLogger().Errorf("user has deleted ,uid:%v|err:%v", uid, err)
		a.SetReceiver(a.HandleUserDeleted)
	} else {
		a.SetReceiver(a.HandleAfterInit)
	}
	a.GetLogger().Infof("player active")
	return nil
}

func (a *PlayerActor) Terminate() {
	err := a.GetModel().Save()
	if err != nil {
		a.GetLogger().Errorf("save player err, uid:%v|err:%v", a.GetModel().Uid, err)
	}
	//_, err = glin.GetSystem().GetCluster().UnsubscribeByPid(global.TopicInviteBroadcast, a.GetCtx().Self())
	//if err != nil {
	//	a.GetLogger().Errorf("unsubscribe TopicInviteBroadcast err, err:%v", err)
	//}
	a.GetLogger().Infof("player terminal")
}
func (a *PlayerActor) Tick() {
	a.BaseActor.Tick()
	if a.GetModel() == nil {
		return
	}
	if a.GetModel().IsDirty() {
		err := a.GetModel().Save()
		if err != nil {
			a.GetLogger().Errorf("tick: save player err, err:%v", err)
		} else {
			a.GetLogger().Infof("tick: save player success")
		}
	}
}
