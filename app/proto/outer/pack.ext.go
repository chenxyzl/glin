package outer

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"laiya/proto/code"
	"laiya/share/global"
)

type IReqPack interface {
	proto.Message
	ProtoType() string
	GetSn() uint64
	GetRpcId() int32
}
type IRspPack interface {
	proto.Message
	ProtoType() string
	GetSn() uint64
	GetRpcId() int32
	GetCode() code.Code
	SetCode(code.Code)
}

func (x *RspPack) SetCode(code code.Code) {
	x.Msg = code.GetCodeMsg()
}
func (x *RspJsonPack) SetCode(code code.Code) {
	x.Msg = code.GetCodeMsg()
}

func (x *ReqPack) ProtoType() string {
	return global.ProtocolPb
}
func (x *ReqJsonPack) ProtoType() string {
	return global.ProtocolJson
}
func (x *RspPack) ProtoType() string {
	return global.ProtocolPb
}
func (x *RspJsonPack) ProtoType() string {
	return global.ProtocolJson
}

func (x *PushPack) ToRspPack(protoType string) (IRspPack, error) {
	var out IRspPack
	if protoType == global.ProtocolJson {
		msg, ok := GetPushMsgByRpcId(x.GetRpcId())
		if !ok {
			return nil, fmt.Errorf("convert to RspPack err, not fount push msg by rpcId, rpcId:%v", x.GetRpcId())
		}
		err := proto.Unmarshal(x.GetData(), msg)
		if err != nil {
			return nil, fmt.Errorf("convert to RspPack Unmarshal err, rpcId:%v|err:%v", x.GetRpcId(), err)
		}
		str, err := PbJson.Marshal(msg)
		if err != nil {
			return nil, fmt.Errorf("convert to RspPack Marshal err, rpcId:%v|err:%v", x.GetRpcId(), err)
		}
		out = &RspJsonPack{
			RpcId: x.GetRpcId(),
			//Code:  code.Code_Ok,//push消息不包含错误码
			Data: string(str),
		}
	} else {
		out = &RspPack{
			RpcId: x.GetRpcId(),
			//Code:  code.Code_Ok,//push消息不包含错误码
			Data: x.GetData(),
		}
	}
	return out, nil
}
