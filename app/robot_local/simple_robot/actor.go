package simple_robot

import (
	"fmt"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/chenxyzl/glin/grain"
	"github.com/chenxyzl/glin/share"
	"laiya/model/robot_local_model"
	"laiya/robot_local/simple_robot/iface"
	"reflect"
)

var _ iface.ISimpleRobotActor[robot_local_model.ILocalSimpleRobotModel] = new(SimpleRobotActor[robot_local_model.ILocalSimpleRobotModel])

type SimpleRobotActor[T robot_local_model.ILocalSimpleRobotModel] struct {
	grain.BaseActor
	iface.ISimpleRobotOnly[T]
	Gid    uint64
	Parent *actor.PID
	model  T
}

func (a *SimpleRobotActor[T]) Init() error {
	a.GetLogger().Infof("robot local init")
	//加载数据
	var t T
	var typ = reflect.TypeOf(t)
	var entity = reflect.New(typ.Elem()).Interface()
	var mod = entity.(T)
	mod.InitUid(fmt.Sprintf("%d_%d", a.Gid, a.GetRobotId()))
	err := mod.Load()
	if err != nil {
		a.GetLogger().Errorf("load robot local data err, group:%v|err:%v", a.Gid, err)
		return err
	}
	//load success
	a.setModel(mod)
	return nil
}

func (a *SimpleRobotActor[T]) Terminate() {
	//
	err := a.GetModel().Save()
	if err != nil {
		a.GetLogger().Errorf("save robot local err, robotId:%v|err:%v", a.GetModel().GetRobotId(), err)
	}
	//
	a.GetLogger().Infof("robot local terminal")
}

func (a *SimpleRobotActor[T]) Tick() {
	a.BaseActor.Tick()
	if share.IsNil(a.GetModel()) {
		return
	}
	if !share.IsNil(a.GetModel()) && a.GetModel().IsDirty() {
		err := a.GetModel().Save()
		if err != nil {
			a.GetLogger().Errorf("tick: save robot local err, err:%v", err)
		} else {
			a.GetLogger().Infof("tick: save robot local success")
		}
	}
}

func (a *SimpleRobotActor[T]) setModel(mod T) {
	//
	a.model = mod
}
func (a *SimpleRobotActor[T]) GetHostGid() uint64 {
	return a.Gid
}

// GetHost 获取宿主
func (a *SimpleRobotActor[T]) GetHost() *actor.PID {
	return a.Parent
}

func (a *SimpleRobotActor[T]) GetModel() T {
	return a.model
}

func (a *SimpleRobotActor[T]) GetKindType() share.EntityKind {
	var t T
	return t.GetKindType()
}

func (a *SimpleRobotActor[T]) GetRobotId() uint64 {
	var t T
	return t.GetRobotUid()
}
