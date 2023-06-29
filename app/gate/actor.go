package gate

import (
	"fmt"
	"github.com/chenxyzl/glin/grain"
	"github.com/chenxyzl/glin/share"
	"github.com/chenxyzl/glin/slog"
	"google.golang.org/protobuf/proto"
	"laiya/config"
	"laiya/gate/iface"
	"laiya/proto/code"
	"laiya/proto/common"
	"laiya/proto/inner"
	"laiya/share/call"
	"laiya/share/datarangers"
	"laiya/share/global"
)

var _ iface.ISessionActor = new(SessionActor)

type SessionActor struct {
	grain.BaseActor
	sess         *Session                 //session
	uid          uint64                   //session对应的uid
	platformType common.LoginPlatformType //登录平台
	timeoutId    uint64                   //超时id
}

func (a *SessionActor) GetKindType() share.EntityKind {
	return global.SessionKind
}

func (a *SessionActor) GetUid() uint64 {
	return a.uid
}

func NewSessionActor(sess *Session, uid uint64, platformType common.LoginPlatformType) *SessionActor {
	act := &SessionActor{sess: sess, uid: uid, platformType: platformType}
	act.SetLogger(slog.NewWith("uid", act.uid, "kind", act.GetKindType(), "platformType", act.platformType))
	return act
}

func (a *SessionActor) Init() error {
	//
	_, cod := call.RequestRemote[proto.Message](a.uid, &inner.C2HoSessionOnline_Request{
		LoginPlatformType: a.platformType,
		Session:           a.GetCtx().Self(),
	})
	if cod != code.Code_Ok {
		a.GetLogger().Errorf("session offline err, cod:%v", cod)
		return fmt.Errorf("sess init err, cod:%v", cod)
	}
	//
	a.SetReceiver(a.HandleAfterInit)
	//
	datarangers.SendEvent("player_ws_connect", map[string]interface{}{
		"uid": a.uid,
	}, nil)
	//
	a.GetLogger().Infof("sess online success")
	return nil
}
func (a *SessionActor) Terminate() {
	//关闭session
	a.sess.Close(true)
	//通知home断开链接
	_, cod := call.RequestRemote[proto.Message](a.uid, &inner.C2HoSessionOffline_Request{
		LoginPlatformType: a.platformType,
		Session:           a.GetCtx().Self()})
	if cod != code.Code_Ok {
		a.GetLogger().Errorf("session offline err, cod:%v", cod)
		return
	}
	//
	a.GetLogger().Infof("sess offline success")
}
func (a *SessionActor) ResetHeartbeatCheck() {
	//删除老的
	if a.timeoutId > 0 {
		a.RemoveDelayByIds([]uint64{a.timeoutId})
	}
	//新的
	a.timeoutId = a.Delay(config.Get().ConstConfig.HeartbeatCheck, func(i ...interface{}) {
		a.GetLogger().Warnf("session timeout, stop self")
		a.GetCtx().Poison(a.GetCtx().Self())
	})
}
