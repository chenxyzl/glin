package activity

import (
	"github.com/chenxyzl/glin/grain"
	"github.com/chenxyzl/glin/share"
	"github.com/chenxyzl/glin/slog"
	"laiya/model/plugin_model/group_plugin_activity_model"
	"laiya/plugin/activity/iface"
	"laiya/share/global"
)

var _ iface.IActivityActor = new(ActivityActor)

type ActivityActor struct {
	grain.BaseActor
	host  iface.GroupActorRef
	model *group_plugin_activity_model.PluginActivity
}

func (a *ActivityActor) SetModel(mod *group_plugin_activity_model.PluginActivity) {
	a.model = mod
}

func (a *ActivityActor) GetModel() *group_plugin_activity_model.PluginActivity {
	return a.model
}
func (a *ActivityActor) GetHost() iface.GroupActorRef {
	return a.host
}

func (a *ActivityActor) GetKindType() share.EntityKind {
	return global.GroupPluginActivityKind
}

func NewActivityActor(host grain.IActorRef) grain.IActor {
	groupActorRef := host.(iface.GroupActorRef)
	act := &ActivityActor{host: groupActorRef}
	act.SetLogger(slog.NewWith("gid", groupActorRef.GetGid(), "kind", act.GetKindType(), "plugin", "activity"))
	return act
}

func (a *ActivityActor) Init() error {
	//
	mod := &group_plugin_activity_model.PluginActivity{Gid: a.GetHost().GetGid()}
	err := mod.Load()
	if err != nil {
		a.GetLogger().Errorf("load group plugin activity data err, gid:%v|err:%v", a.GetHost().GetGid(), err)
		return err
	}
	a.SetModel(mod)
	//
	a.GetLogger().Infof("group plugin activity active")
	return nil
}
func (a *ActivityActor) Terminate() {
	//
	err := a.GetModel().Save()
	if err != nil {
		a.GetLogger().Errorf("save group plugin activity err, hall:%v|err:%v", a.GetModel().Gid, err)
	}
	//
	a.GetLogger().Infof("group plugin activity terminal success")
}

func (a *ActivityActor) Tick() {
	a.BaseActor.Tick()
	if a.GetModel() == nil {
		return
	}
	if a.GetModel().IsDirty() {
		err := a.GetModel().Save()
		if err != nil {
			a.GetLogger().Errorf("tick: save group plugin activity err, err:%v", err)
		} else {
			a.GetLogger().Infof("tick: save group plugin activity success")
		}
	}
}
