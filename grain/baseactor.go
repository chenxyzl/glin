package grain

import (
	"fmt"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/scheduler"
	"github.com/chenxyzl/glin/share"
	"github.com/chenxyzl/glin/slog"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"reflect"
	"strings"
	"time"
)

var _ IActor = new(BaseActor)

type execFunc struct {
	f          func(...interface{})
	params     []interface{}
	cancelFunc scheduler.CancelFunc
}

type BaseActor struct {
	IActor
	ctx        actor.Context
	logger     slog.Logger
	behavior   actor.ReceiveFunc
	options    actorOptions
	plugins    map[share.EntityKind]*actor.PID
	execFuncId uint64
	//
	lastActiveTimeSec int64 //最后一次活跃时间
	//
	scheduler *scheduler.TimerScheduler
	//next
	nextFunc       []execFunc
	nextCancelFunc scheduler.CancelFunc
	//tick
	tickCancelFunc scheduler.CancelFunc
	//
	delayFunc map[uint64]execFunc
	//
	repeatFunc map[uint64]execFunc
	//
	components   []IComponent
	componentMap map[string]IComponent
}

func (a *BaseActor) beforeInit(ia IActor, context actor.Context) error {
	a.lastActiveTimeSec = time.Now().Unix()
	a.IActor = ia
	a.ctx = context
	if a.logger == nil {
		a.logger = slog.DefaultLogger
	}
	a.plugins = make(map[share.EntityKind]*actor.PID)
	a.components = make([]IComponent, 0)
	a.componentMap = make(map[string]IComponent)
	a.delayFunc = make(map[uint64]execFunc)
	a.repeatFunc = make(map[uint64]execFunc)
	a.scheduler = scheduler.NewTimerScheduler(a.ctx)
	a.tickCancelFunc = a.scheduler.SendRepeatedly(0, time.Second, a.ctx.Self(), &Tick{})
	//new
	componentTypes := GetEntityAllComponentDef(a.GetKindType())
	for _, componentType := range componentTypes {
		component := reflect.New(componentType.Com).Interface().(IComponent)
		component.bindEntity(a.IActor)
		a.components = append(a.components, component)
		a.componentMap[componentType.IfName] = component
	}
	//init
	for _, component := range a.components {
		if err := component.BeforeInit(); err != nil {
			return err
		}
	}
	return nil
}
func (a *BaseActor) afterInit() error {
	//init
	for _, component := range a.components {
		if err := component.AfterInit(); err != nil {
			return err
		}
	}
	return nil
}
func (a *BaseActor) beforeTerminate() {
	//清除延时
	a.CleanAllDelayIds()
	a.CleanAllRepeatIds()
	//
	if a.tickCancelFunc != nil {
		a.tickCancelFunc()
		a.tickCancelFunc = nil
	}
	//terminate
	for _, component := range a.components {
		component.BeforeTerminate()
	}
}
func (a *BaseActor) afterTerminate() {
	//next事件允许延后处理
	if a.nextCancelFunc != nil {
		a.nextCancelFunc()
		a.nextCancelFunc = nil
	}
	//terminate
	for _, component := range a.components {
		component.AfterTerminate()
	}
}
func (a *BaseActor) componentTick() {
	//terminate
	for _, component := range a.components {
		component.Tick()
	}
}

// func (a *BaseActor) GetKindType() share.EntityKind         { return share.EntityKindUnknown }
func (a *BaseActor) GetComponents() []IComponent               { return a.components }
func (a *BaseActor) GetComponentByName(name string) IComponent { return a.componentMap[name] }
func (a *BaseActor) GetLocalHandler(msgName protoreflect.FullName) *HandlerDef {
	return GetLocalHandlerDef(a.GetKindType(), msgName)
}

func (a *BaseActor) GetCtx() actor.Context        { return a.ctx }
func (a *BaseActor) SetLogger(logger slog.Logger) { a.logger = logger }
func (a *BaseActor) GetLogger() slog.Logger {
	if a == nil {
		return slog.DefaultLogger
	}
	return a.logger
}

// func (a *BaseActor) Init()                                 {}
func (a *BaseActor) Terminate()                            {}
func (a *BaseActor) Tick()                                 { a.checkTerminateIdleTime() }
func (a *BaseActor) SetReceiver(receive actor.ReceiveFunc) { a.behavior = receive }
func (a *BaseActor) ResetDefaultReceiver()                 { a.behavior = nil }
func (a *BaseActor) defaultReceiveInitErr(msg actor.Context) {
	a.GetCtx().Respond(&Error{Code: int32(23)}) //ActorInitError--循环引用问题
	a.GetCtx().Poison(a.GetCtx().Self())
	a.GetLogger().Errorf("actor init failed,poison self,actor:%v", a.GetCtx().Self())
}
func (a *BaseActor) DefaultReceiveAfterInit(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case proto.Message:
		reqName := msg.ProtoReflect().Descriptor().FullName()
		if strings.HasSuffix(string(reqName), ".Notify") {
			LocalNotifyHandler(reqName, a, msg)
		} else {
			rsp, cod := LocalRequestHandler(reqName, a, msg)
			if cod != 1 { //code.Code_Ok
				a.GetLogger().Errorf("call rpc err, self:%v rpc:%v|targetKind:%v|isLocal:%v|cod:%v", a.GetCtx().Self(), reqName, a.GetKindType(), true, cod)
				ctx.Respond(&Error{Code: cod})
				return
			}
			ctx.Respond(rsp)
		}
	default:
		ctx.Respond(&Error{Code: int32(2)}) //code.Code_Error
		a.GetLogger().Errorf("msg unregister in connected, self:%v|name:%v|msg:%v", a.GetCtx().Self(), reflect.TypeOf(msg).Elem().Name(), msg)
	}
}
func (a *BaseActor) ReceiveDefault(msg proto.Message) {
	defer share.RecoverInfo(fmt.Sprintf("actor deal msg err, msgType:%v", msg.ProtoReflect().Descriptor().FullName()), a.GetLogger())
	//更新上次活跃时间
	a.lastActiveTimeSec = time.Now().Unix()
	//真正开始执行
	if a.behavior == nil {
		//没有设置则走默认的
		a.DefaultReceiveAfterInit(a.GetCtx())
	} else {
		//走自定义的
		a.behavior.Receive(a.GetCtx())
	}
}

func (a *BaseActor) handleNextStep() {
	//单次最多执行10个
	var batchCount = 10
	var list []execFunc
	if len(a.nextFunc) > 0 {
		if len(a.nextFunc) < batchCount {
			batchCount = len(a.nextFunc)
		}
		list = a.nextFunc[:batchCount]
		a.nextFunc = a.nextFunc[batchCount:]
	}
	//执行
	for _, f := range list {
		f.f(f.params...)
	}
	//有剩余则继续
	if len(a.nextFunc) > 0 {
		a.nextCancelFunc = a.scheduler.SendOnce(0, a.ctx.Self(), &NextStep{})
	}
}

func (a *BaseActor) handleDelay(id uint64) {
	f, ok := a.delayFunc[id]
	if !ok {
		//may be canceled
		return
	}
	//delete
	f.cancelFunc = nil
	delete(a.delayFunc, id)
	f.f(f.params...)
}

func (a *BaseActor) handleRepeat(id uint64) { //called by delay
	f, ok := a.repeatFunc[id]
	if !ok {
		//may be canceled
		return
	}
	f.f(f.params...)
}

func (a *BaseActor) Next(exec func(...interface{}), params ...interface{}) {
	a.nextFunc = append(a.nextFunc, execFunc{f: exec, params: params})
	a.nextCancelFunc = a.scheduler.SendOnce(0, a.ctx.Self(), &NextStep{})
}

func (a *BaseActor) Delay(dt time.Duration, exec func(...interface{}), params ...interface{}) uint64 {
	if len(a.delayFunc) > 100 {
		a.GetLogger().Warnf("delay func size too bigger, please check, size:%v", len(a.delayFunc))
	}
	//
	a.execFuncId++
	//
	a.delayFunc[a.execFuncId] = execFunc{
		f:          exec,
		params:     params,
		cancelFunc: a.scheduler.SendOnce(dt, a.GetCtx().Self(), &Delay{Id: a.execFuncId}),
	}
	return a.execFuncId
}

func (a *BaseActor) RemoveDelayByIds(ids []uint64) {
	for _, id := range ids {
		f, ok := a.delayFunc[id]
		if !ok {
			a.GetLogger().Warnf("remove delay func error, func not exist, id:%v", id)
			continue
		}
		if f.cancelFunc != nil {
			f.cancelFunc()
			f.cancelFunc = nil
		}
		delete(a.delayFunc, id)
	}
}

func (a *BaseActor) CleanAllDelayIds() {
	for id, f := range a.delayFunc {
		if f.cancelFunc != nil {
			f.cancelFunc()
			f.cancelFunc = nil
		}
		delete(a.delayFunc, id)
	}
}

func (a *BaseActor) Repeat(initial time.Duration, interval time.Duration, exec func(...interface{}), params ...interface{}) uint64 { //
	if len(a.repeatFunc) > 100 {
		a.GetLogger().Warnf("repeated func size too bigger, please check, size:%v", len(a.repeatFunc))
	}
	//
	a.execFuncId++
	//
	a.repeatFunc[a.execFuncId] = execFunc{
		f:          exec,
		params:     params,
		cancelFunc: a.scheduler.SendRepeatedly(initial, interval, a.GetCtx().Self(), &Repeat{Id: a.execFuncId}),
	}
	return a.execFuncId
}

func (a *BaseActor) RemoveRepeatByIds(ids []uint64) {
	for _, id := range ids {
		f, ok := a.repeatFunc[id]
		if !ok {
			a.GetLogger().Warnf("remove repeat func error, func not exist, id:%v", id)
			continue
		}
		if f.cancelFunc != nil {
			f.cancelFunc()
			f.cancelFunc = nil
		}
		delete(a.repeatFunc, id)
	}
}

func (a *BaseActor) CleanAllRepeatIds() {
	for id, f := range a.repeatFunc {
		if f.cancelFunc != nil {
			f.cancelFunc()
			f.cancelFunc = nil
		}
		delete(a.repeatFunc, id)
	}
}

func (a *BaseActor) getOptions() actorOptions {
	return a.options
}

// AddOptions add options
func (a *BaseActor) AddOptions(opts ...ActorOption) {
	for _, option := range opts {
		option(&a.options)
	}
}

func (a *BaseActor) checkTerminateIdleTime() {
	if a.options.filterIdleCheck != nil && a.options.filterIdleCheck() {
		return
	}
	if a.options.terminateIdleTime <= 0 {
		return
	}
	if time.Now().Unix() < a.lastActiveTimeSec+int64(a.options.terminateIdleTime/time.Second) {
		return
	}
	a.GetCtx().Poison(a.GetCtx().Self())
	a.GetLogger().Infof("actor idle too long, will stop self, last:%v|config:%v", a.lastActiveTimeSec, a.options.terminateIdleTime)
}
func (a *BaseActor) addPlugin(pluginType share.EntityKind, plugin *actor.PID) {
	a.plugins[pluginType] = plugin
}
func (a *BaseActor) getPlugin(pluginEntityKind share.EntityKind) *actor.PID {
	return a.plugins[pluginEntityKind]
}
