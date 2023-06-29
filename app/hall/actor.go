package hall

import (
	"github.com/chenxyzl/glin/grain"
	"github.com/chenxyzl/glin/share"
	"github.com/chenxyzl/glin/slog"
	"laiya/hall/iface"
	"laiya/model/hall_model"
	"laiya/share/global"
	"laiya/share/utils"
)

var _ iface.IHallActor = new(HallActor)

type HallActor struct {
	grain.BaseActor
	model *hall_model.Hall
}

func (a *HallActor) GetKindType() share.EntityKind {
	return global.HallKind
}

func NewHallActor() *HallActor {
	act := &HallActor{}
	return act
}

func (a *HallActor) SetModel(mod *hall_model.Hall) {
	if mod.Groups == nil {
		mod.Groups = make(map[uint64]*hall_model.Group)
	}
	mod.CleanOldPlayerCount()
	mod.NewGroupSort()
	a.model = mod
}

func (a *HallActor) GetModel() *hall_model.Hall {
	return a.model
}

func (a *HallActor) Init() error {
	hallId, err := utils.ParseUid(a.GetCtx().Self())
	if err != nil {
		a.GetLogger().Errorf("parse hallId from ctx err, ctx:%v|err:%v", a.GetCtx().Self(), err)
		return err
	}
	//
	a.SetLogger(slog.NewWith("hallId", hallId, "kind", a.GetKindType()))
	//
	mod := &hall_model.Hall{HallId: hallId}
	err = mod.Load()
	if err != nil {
		a.GetLogger().Errorf("load hall data err, hall:%v|err:%v", hallId, err)
		return err
	}
	a.SetModel(mod)
	//
	a.GetLogger().Infof("hall active")
	return nil
}

func (a *HallActor) Terminate() {
	//
	err := a.GetModel().Save()
	if err != nil {
		a.GetLogger().Errorf("save hall err, hall:%v|err:%v", a.GetModel().HallId, err)
	}
	//
	a.GetLogger().Infof("hall terminal")
}

func (a *HallActor) Tick() {
	a.BaseActor.Tick()
	if a.GetModel() == nil {
		return
	}
	if a.GetModel().IsDirty() {
		err := a.GetModel().Save()
		if err != nil {
			a.GetLogger().Errorf("tick: save hall err, err:%v", err)
		} else {
			a.GetLogger().Infof("tick: save hall success")
		}
	}
}
