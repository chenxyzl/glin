package outer

import (
	"fmt"
	"github.com/chenxyzl/glin/grain"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"laiya/proto/code"
	"laiya/share/global"
	"strings"
)

// BuildErrPack 构造错误消息
func BuildErrPack(pkg IReqPack, code code.Code) IRspPack {
	var iRsp IRspPack
	if pkg.ProtoType() == global.ProtocolJson {
		iRsp = &RspJsonPack{
			Sn:    pkg.GetSn(),
			RpcId: pkg.GetRpcId(),
		}
		iRsp.SetCode(code)
	} else {
		iRsp = &RspPack{
			Sn:    pkg.GetSn(),
			RpcId: pkg.GetRpcId(),
		}
		iRsp.SetCode(code)
	}
	return iRsp
}

// BuildReplyPack 构造Reply消息
func BuildReplyPack(pkg IReqPack, msg proto.Message) IRspPack {
	var v []byte
	var err error
	var out IRspPack
	if pkg.ProtoType() == global.ProtocolJson {
		v, err = PbJson.Marshal(msg)
		if err != nil {
			panic(err) //流程正确情况下不会触发,所以无所谓
		}
		out = &RspJsonPack{
			Sn:    pkg.GetSn(),
			RpcId: pkg.GetRpcId(),
			Code:  code.Code_Ok,
			Data:  string(v),
		}
	} else {
		v, err = proto.Marshal(msg)
		if err != nil {
			panic(err) //流程正确情况下不会触发,所以无所谓
		}
		out = &RspPack{
			Sn:    pkg.GetSn(),
			RpcId: pkg.GetRpcId(),
			Code:  code.Code_Ok,
			Data:  v,
		}
	}

	return out
}

// BuildPushPack 构造push消息
func BuildPushPack(msg proto.Message) *PushPack {
	//避免多层包装
	if v, ok := msg.(*PushPack); ok {
		return v
	}
	v, err := proto.Marshal(msg)
	if err != nil {
		panic(err) //流程正确情况下不会触发,所以无所谓
	}
	rpcId, found := GetPushRpcIdByMsg(msg)
	if !found {
		panic(fmt.Sprintf("get rpcId not found, typeName:%v", msg.ProtoReflect().Descriptor().Name())) //流程正确情况下不会触发,所以无所谓
	}
	out := &PushPack{
		RpcId: rpcId,
		Data:  v,
	}
	return out
}

// BuildError 构造rpc错误
func BuildError(code code.Code) *grain.Error {
	return &grain.Error{Code: int32(code)}
}

// BuildChatMsg 构造聊天消息
func BuildChatMsg(msg proto.Message) (*ChatMsg, error) {
	content, err := anypb.New(msg)
	if err != nil {
		return nil, err
	}
	msgType := strings.TrimSuffix(strings.TrimPrefix(content.TypeUrl, "type.googleapis.com/outer.ChatMsg."), "Msg")
	msgTypeId, ok := ChatMsg_MsgType_value[msgType]
	if !ok {
		return nil, fmt.Errorf("chat MsgType not have val, val:%v|typ:%v", content.TypeUrl, msgType)
	}
	ret := &ChatMsg{
		MsgType: ChatMsg_MsgType(msgTypeId),
		Content: content,
	}
	return ret, nil
}

// BuildHttpErr 构造http回复 Error
func BuildHttpErr(code code.Code) *HttpReply {
	return &HttpReply{Code: code, Msg: code.GetCodeMsg()}
}

// BuildHttpOk 构造http回复 OK
func BuildHttpOk(rsp proto.Message) *HttpReply {
	data, err := proto.Marshal(rsp)
	if err != nil {
		return nil
	}
	return &HttpReply{Code: code.Code_Ok, Data: data}
}

func BuildCusRenderErr(code code.Code) *CusRender {
	return &CusRender{Code: code}
}

func BuildCusRenderOk(data proto.Message) *CusRender {
	return &CusRender{Code: code.Code_Ok, Data: data}
}
