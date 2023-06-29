package grain

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/chenxyzl/glin/share"
	"github.com/chenxyzl/glin/slog"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"time"
)

type IEntity interface {
	GetKindType() share.EntityKind
	GetComponents() []IComponent
	GetComponentByName(string) IComponent
	GetLocalHandler(msgName protoreflect.FullName) *HandlerDef
}

type IActorRef interface {
	//GetCtx ...
	GetCtx() actor.Context //get ctx
}

type IActor interface {
	IEntity
	IActorRef

	//AddOptions option
	AddOptions(...ActorOption) //options
	getOptions() actorOptions
	checkTerminateIdleTime() //option.terminateIdleTime

	//life
	beforeInit(IActor, actor.Context) error //before call Init
	afterInit() error                       //after call Init
	beforeTerminate()                       //before call Terminate
	afterTerminate()                        //after call Terminate
	componentTick()                         //after call Tick
	handleNextStep()                        //called by next
	handleDelay(id uint64)                  //called by delay
	handleRepeat(id uint64)                 //called by delay

	//method
	SetLogger(slog.Logger)  //rewrite default logger
	GetLogger() slog.Logger //get logger
	Init() error            //Start
	Terminate()             //Stop
	Tick()                  //1 second
	Next(func(params ...any), ...any)
	Delay(dt time.Duration, exec func(...interface{}), params ...interface{}) uint64                               //
	Repeat(initial time.Duration, interval time.Duration, exec func(...interface{}), params ...interface{}) uint64 //
	RemoveDelayByIds(ids []uint64)                                                                                 //
	CleanAllDelayIds()
	SetReceiver(receive actor.ReceiveFunc)
	ResetDefaultReceiver()
	defaultReceiveInitErr(msg actor.Context)
	DefaultReceiveAfterInit(msg actor.Context)
	ReceiveDefault(msg proto.Message) //receive msg

	//plugin
	addPlugin(pluginType share.EntityKind, plugin *actor.PID)
	getPlugin(pluginEntityKind share.EntityKind) *actor.PID
}
