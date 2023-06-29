package gate

import (
	"context"
	"errors"
	"fmt"
	"github.com/chenxyzl/glin"
	"github.com/chenxyzl/glin/grain"
	"github.com/chenxyzl/glin/share"
	"github.com/chenxyzl/glin/slog"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"io"
	"laiya/config"
	"laiya/proto/code"
	"laiya/proto/common"
	"laiya/share/global"
	"laiya/share/token"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var serverList []*http.Server
var sessionMap = sync.Map{}

func helpersHighLevelHandler(w http.ResponseWriter, r *http.Request) {
	//recover
	defer share.Recover()
	//升级为websocket
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		slog.Errorf("upgrade error: %s", err)
		return
	}
	defer conn.Close()
	//解析uri
	uri := r.RequestURI
	rq, err := url.Parse(uri)
	if err != nil {
		slog.Errorf("session parse uri error; uri:%s|err:%v", uri, err)
		return
	}
	//解析query
	m, err := url.ParseQuery(rq.RawQuery)
	if err != nil {
		slog.Errorf("session parse query error; uri:%s|rq:%s|err:%v", uri, rq.RawQuery, err)
		return
	}
	//获取token
	t := m.Get(global.ParamToken)
	if t == "" {
		slog.Error(fmt.Errorf("session get token is nil; uri:%s|rq:%s|m:%v", uri, rq.RawQuery, m))
		return
	}
	//解析token
	uid, err := token.ParseToken(t, config.Get().AConfig.AppKey)
	if err != nil {
		slog.Error(fmt.Errorf("session parse token error; uri:%s|token:%s|cod:%v", uri, t, err))
		return
	}
	platform := m.Get("platform")
	platformType := common.LoginPlatformType_unknownPlatform
	if platform != "" {
		v, ok := common.LoginPlatformType_value[platform]
		if !ok {
			slog.Error(fmt.Errorf("session parse login platform error; uri:%s|uid:%s|platform:%v|cod:%v", uri, uid, platform, err))
			return
		}
		platformType = common.LoginPlatformType(v)
	}
	//记录sess
	sess := NewSess(conn, uid, platformType, nil, m.Get(global.ProtoName))
	if oldSess, ok := sessionMap.LoadOrStore(sess.Key(), sess); ok {
		sess.logger.Infof("sess connected, but find old sess, will close old sess, token:%v", t)
		(oldSess.(*Session)).Close(true)
	}
	sess.logger.Infof("sess connected,token:%v|remoteIp:%v", t, conn.RemoteAddr())
	//退出时候删除sess
	defer func() {
		sess.Close()
		sessionMap.Delete(sess.Key())
		sess.logger.Infof("sess closed")
	}()
	//创建actor
	sess.ac = glin.GetSystem().NewLocalActor(func() grain.IActor {
		return NewSessionActor(sess, uid, platformType)
	})
	//退出时候 actor下线
	defer func() {
		glin.GetSystem().GetCluster().ActorSystem.Root.Stop(sess.ac)
		sess.logger.Infof("player offline, cod:%v", code.Code_Ok)
	}()
	//登录成功
	sess.logger.Infof("login success")
	//开启消息监听
	for {
		bts, _, err := wsutil.ReadClientData(sess.conn)
		if err == io.EOF {
			sess.logger.Infof("sess normal close1, %v", err)
			return
		}
		var netError net.Error
		if errors.As(err, &netError) && netError.Timeout() {
			sess.logger.Infof("sess normal close2, %v", netError)
			return
		}
		var closeError wsutil.ClosedError
		if errors.As(err, &closeError) &&
			(closeError.Code == ws.StatusNoStatusRcvd ||
				closeError.Code == ws.StatusNormalClosure ||
				closeError.Code == ws.StatusGoingAway) {
			sess.logger.Infof("sess normal close3, %v", closeError)
			return
		}
		if err != nil {
			sess.logger.Errorf("read message error: %v", err)
			return
		}
		//处理消息,有错误不需要退出
		err = sess.OnMessage(bts)
		if err != nil {
			sess.logger.Errorf("on message err err:%v", err)
		}
	}
}

func StateWs() {
	slog.Infof("ws start...")
	http.HandleFunc(config.Get().AConfig.WsParam.Path, helpersHighLevelHandler)
	addr := net.JoinHostPort(config.Get().AConfig.WsParam.Host, config.Get().AConfig.WsParam.Port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		slog.Panicf("listen %q error: %v", addr, err)
	}

	s := new(http.Server)
	go func() {
		err := s.Serve(ln)
		if err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				slog.Panic(err)
			} else {
				slog.Infof("ws server closed")
			}
		}
	}()
	slog.Infof("ws start success, %s (%v%v)", ln.Addr(), addr, config.Get().AConfig.WsParam.Path)
	serverList = append(serverList, s)
}

func CloseWs() {
	for _, server := range serverList {
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		if err := server.Shutdown(ctx); err != nil {
			slog.Error(err)
		}
	}
}

var tickTimes = 0

func PrintOnlineCount() {
	//
	tickTimes++
	//
	if tickTimes%60 != 0 {
		return
	}
	slog.Infof("online session num, num:%d", lenSyncMap(&sessionMap))
}

func lenSyncMap(m *sync.Map) int {
	var i int
	m.Range(func(k, v interface{}) bool {
		i++
		return true
	})
	return i
}
