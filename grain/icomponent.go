package grain

import (
	"reflect"
)

type IComponent interface {
	bindEntity(entity IEntity)
	GetEntity() IEntity
	BeforeInit() error //actor Init完成前
	AfterInit() error  //actor Init完成前
	BeforeTerminate()  //actor Terminate完成前
	AfterTerminate()   //actor Terminate完成前
	Tick()
}

var _ IComponent = new(BaseComponent)

type BaseComponent struct {
	IComponent
	entity IEntity
}

func (b *BaseComponent) bindEntity(entity IEntity) {
	b.entity = entity
}

func (b *BaseComponent) GetEntity() IEntity {
	return b.entity
}

func GetComponent[T IComponent](com IEntity) T {
	var iComTyp = reflect.TypeOf(new(T)).Elem() //todo 考虑用组件名字来获取组件,减少反射
	return com.GetComponentByName(iComTyp.Name()).(T)
}

func GetComponent2[T IComponent](com IComponent) T {
	var iComTyp = reflect.TypeOf(new(T)).Elem()
	return com.GetEntity().GetComponentByName(iComTyp.Name()).(T)
}
