package grain

import (
	"fmt"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/cluster"
	"github.com/chenxyzl/glin/share"
	"github.com/chenxyzl/glin/slog"
	"google.golang.org/protobuf/proto"
)

type Grain struct {
	ctx          actor.Context
	inner        IActor
	isInitFailed bool
	isCluster    bool
}

func (a *Grain) GetLogger() slog.Logger {
	if a == nil || a.inner == nil {
		return slog.DefaultLogger
	}
	return a.inner.GetLogger()
}

// Receive ensures the lifecycle of the actor for the received message
func (a *Grain) Receive(ctx actor.Context) {
	msgName := share.GetTypeName(ctx.Message())
	defer share.RecoverInfo(fmt.Sprintf("handler panic, msg:%v", msgName), a.GetLogger())
	//
	a.ctx = ctx
	switch msg := ctx.Message().(type) {
	case *cluster.ClusterInit:
		//cluster actor will run
		if !a.isCluster {
			break
		}
		a.init(ctx)
	case *actor.Started: // pass
		//local actor will run
		if a.isCluster {
			break
		}
		a.init(ctx)
	case *actor.ReceiveTimeout:
		ctx.Poison(ctx.Self())
	case *actor.PoisonPill:
		ctx.Stop(ctx.Self())
	case *actor.Stopping:
		a.inner.beforeTerminate()
	case *actor.Stopped:
		a.inner.Terminate()
		a.inner.afterTerminate()
		a.GetLogger().Infof("Grain stopped:%v ", ctx.Self())
	case *actor.Restarting:
		a.GetLogger().Infof("Grain restarting:%v", a.ctx.Self())
	case actor.AutoReceiveMessage: // --不要换位置(因为上面的消息很多消息都属于这个,避免被拦截处理)
	case actor.SystemMessage: // --不要换位置(因为上面的消息很多消息都属于这个,避免被拦截处理)
	case *cluster.GrainRequest: //--不要换位置
	case *NextStep:
		a.inner.handleNextStep()
	case *Delay:
		a.inner.handleDelay(msg.GetId())
	case *Repeat:
		a.inner.handleRepeat(msg.GetId())
	case *Tick:
		a.inner.Tick()
		a.inner.componentTick()
	case proto.Message:
		protoMsgName := proto.MessageName(msg)
		msgName = string(protoMsgName)
		//检查是插件方法还是本地方法,如果是插件方法转发到插件中
		handler := a.inner.GetLocalHandler(protoMsgName)
		if handler.GetPluginKind() != "" {
			plugin := a.inner.getPlugin(handler.GetPluginKind())
			if plugin == nil {
				selfKind, ok := pluginMapping[a.inner.GetKindType()]
				if !ok {
					a.GetLogger().Errorf("plugin not found 1, plugin:%v|msg:%v", plugin, protoMsgName)
					break
				}
				pluginType, ok := selfKind[handler.GetPluginKind()]
				if !ok {
					a.GetLogger().Errorf("plugin not found 2, plugin:%v|msg:%v", plugin, protoMsgName)
					break
				}
				//创建插件
				plugin = a.inner.GetCtx().Spawn(actor.PropsFromProducer(func() actor.Actor {
					return &Grain{inner: pluginType(a.inner), isCluster: false}
				}))
				//添加插件
				a.inner.addPlugin(handler.GetPluginKind(), plugin)
			}
			a.inner.GetCtx().Forward(plugin)
		} else {
			a.inner.ReceiveDefault(msg)
		}
	default:
		a.GetLogger().Infof("recv unknown msg, msg:%v", msgName)
	}
}

func (a *Grain) init(ctx actor.Context) {
	identity := "local"
	if a.isCluster {
		identity = cluster.GetClusterIdentity(ctx.(actor.ExtensionContext)).GetIdentity()
	}

	for {
		a.isInitFailed = true
		if err := a.inner.beforeInit(a.inner, a.ctx); err != nil {
			a.GetLogger().Errorf("Grain start before init self:%v|identity:%v|err:%v", ctx.Self(), identity, err)
			break
		}
		if err := a.inner.Init(); err != nil {
			a.GetLogger().Errorf("Grain start init self:%v|identity:%v|err:%v", ctx.Self(), identity, err)
			break
		}
		if err := a.inner.afterInit(); err != nil {
			a.GetLogger().Errorf("Grain start after self:%v|identity:%v|err:%v", ctx.Self(), identity, err)
			break
		}
		a.isInitFailed = false
		break //
	}
	if a.isInitFailed {
		a.inner.SetReceiver(a.inner.defaultReceiveInitErr)
	}
	a.GetLogger().Infof("Grain start success:%v|identity:%v", a.ctx.Self(), identity)
}
