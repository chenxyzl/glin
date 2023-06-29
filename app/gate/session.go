package gate

import (
	"fmt"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/chenxyzl/glin/share"
	"github.com/chenxyzl/glin/slog"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"laiya/config"
	"laiya/proto/code"
	"laiya/proto/common"
	"laiya/proto/outer"
	"laiya/share/call"
	"laiya/share/global"
	"net"
	"sync"
	"time"
)

type Session struct {
	conn         net.Conn
	uid          uint64
	platformType common.LoginPlatformType
	ac           *actor.PID
	logger       slog.Logger
	writeLock    sync.RWMutex
	isClosed     bool
	protoType    string
}

func NewSess(conn net.Conn, uid uint64, platformType common.LoginPlatformType, ac *actor.PID, protoType string) *Session {
	return &Session{
		conn:         conn,
		uid:          uid,
		platformType: platformType,
		ac:           ac,
		logger:       slog.NewWith("uid", uid, "platformType", platformType, "protoType", protoType),
		writeLock:    sync.RWMutex{},
		isClosed:     false,
		protoType:    protoType,
	}
}

func (sess *Session) Key() string {
	return fmt.Sprintf("session_%d_%v", sess.uid, sess.platformType)
}

func (sess *Session) OnMessage(data []byte) error {
	defer share.Recover()
	//
	var iPkg outer.IReqPack
	var request proto.Message
	var found bool
	if sess.protoType == global.ProtocolJson {
		pkg := &outer.ReqJsonPack{}
		iPkg = pkg
		err := protojson.Unmarshal(data, pkg)
		if err != nil {
			return fmt.Errorf("protojson.Unmarshal:%#v", err)
		}
		request, found = outer.GetRequestMsgByRpcId(pkg.GetRpcId())
		if !found {
			sess.SafeWritePkg(outer.BuildErrPack(pkg, code.Code_Error))
			return fmt.Errorf("rpc msg not found, rpcId:%v", pkg.GetRpcId())
		}
		if len(pkg.GetData()) > 0 {
			if err := protojson.Unmarshal([]byte(pkg.GetData()), request); err != nil {
				sess.SafeWritePkg(outer.BuildErrPack(pkg, code.Code_Error))
				return fmt.Errorf("data unmarshal err, rpcId:%v|,err:%v", pkg.RpcId, err)
			}
		}
	} else {
		pkg := &outer.ReqPack{}
		iPkg = pkg
		err := proto.Unmarshal(data, pkg)
		if err != nil {
			return fmt.Errorf("proto.Unmarshal:%#v", err)
		}
		request, found = outer.GetRequestMsgByRpcId(pkg.GetRpcId())
		if !found {
			sess.SafeWritePkg(outer.BuildErrPack(pkg, code.Code_Error))
			return fmt.Errorf("rpc msg not found, rpcId:%v", pkg.GetRpcId())
		}
		if len(pkg.GetData()) > 0 {
			if err = proto.Unmarshal(pkg.GetData(), request); err != nil {
				sess.SafeWritePkg(outer.BuildErrPack(pkg, code.Code_Error))
				return fmt.Errorf("data unmarshal err, rpcId:%v|,err:%v", pkg.RpcId, err)
			}
		}
	}
	//发送消息给目标actor
	reply, cod := call.Request[proto.Message](sess.ac, request)
	if cod != code.Code_Ok {
		sess.SafeWritePkg(outer.BuildErrPack(iPkg, cod))
	} else {
		sess.SafeWritePkg(outer.BuildReplyPack(iPkg, reply))
	}
	return nil
}

func (sess *Session) SafeWritePush(push *outer.PushPack) {
	rspPack, err := push.ToRspPack(sess.protoType)
	if err != nil {
		sess.logger.Errorf("convert push msg to RspPack err, err:%v", err)
		return
	}
	sess.SafeWritePkg(rspPack)
}

func (sess *Session) SafeWritePkg(rsp outer.IRspPack) {
	var data []byte
	var err error
	if sess.protoType == global.ProtocolJson {
		data, err = outer.PbJson.Marshal(rsp)
		if err != nil {
			sess.logger.Errorf("SafeWritePkg: protojson.Marshal, rsp:%v|err:%v", rsp, err)
			return
		}
	} else {
		data, err = proto.Marshal(rsp)
		if err != nil {
			sess.logger.Errorf("SafeWritePkg: proto.Marshal, rsp:%v|err:%v", rsp, err)
			return
		}
	}
	//lock
	sess.writeLock.Lock()
	defer sess.writeLock.Unlock()

	if sess.IsClosed() {
		sess.logger.Warnf("sess closed, not allow write")
		return
	}

	//timeout
	err = sess.conn.SetWriteDeadline(time.Now().Add(config.Get().AConfig.WsParam.WriteTimeout))
	if err != nil {
		_ = sess.conn.Close()
		sess.logger.Errorf("set write dead line err, err:%v", err)
		return
	}

	//write
	err = wsutil.WriteServerMessage(sess.conn, ws.OpBinary, data)
	if err != nil {
		//写入失败则关闭连接
		_ = sess.conn.Close()
		sess.logger.Errorf("send msg to client err, rpcId:%v|code:%v|error: %v", rsp.GetRpcId(), rsp.GetCode(), err)
		return
	}
	sess.logger.Infof("send msg to client success, rpcId:%v", rsp.GetRpcId())

}

func (sess *Session) IsClosed() bool {
	return sess.isClosed
}

func (sess *Session) Close(closeSession ...bool) {
	sess.writeLock.Lock()
	defer sess.writeLock.Unlock()
	if !sess.IsClosed() {
		sess.isClosed = true
		if len(closeSession) > 0 && closeSession[0] {
			_ = sess.conn.Close()
		}
	}
}
