package outer

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

type tRequest struct {
	rpcId   int32
	rpcName string
	request protoreflect.MessageType
	reply   protoreflect.MessageType
	routers []string
}

type tPush struct {
	rpcId   int32
	push    protoreflect.MessageType
	routers []string
}

// 请求应答对应的请求应答消息
var requestRpcIdMap = map[int32]tRequest{}

var requestRpcNameMap = map[string]tRequest{}

// push消息对应的rpcId
var pushRpcNameMap = map[protoreflect.FullName]tPush{}

// push消息对应的rpcId
var pushRpcIdMap = map[int32]tPush{}

func init() {
	type Empty struct{}
	pkName := reflect.TypeOf(Empty{}).PkgPath()
	lastPkName := pkName[strings.LastIndex(pkName, "/")+1:]
	fillDes, err := protoregistry.GlobalFiles.FindFileByPath(fmt.Sprintf("%s/router.proto", lastPkName))
	if err != nil {
		panic(err)
	}
	router := proto.GetExtension(fillDes.Options(), E_Router).(string)
	routers := strings.Split(router, ",")
	protoregistry.GlobalTypes.RangeMessages(func(messageType protoreflect.MessageType) bool {
		msgPkName := reflect.TypeOf(messageType.New().Interface()).Elem().PkgPath()
		//只注册自己包内的
		if pkName != msgPkName {
			return true
		}
		fullName := messageType.Descriptor().FullName()
		//只有最外层的消息才需要注册
		if strings.Count(string(fullName), ".") != 1 {
			return true
		}
		filed := messageType.Descriptor().Fields().ByName("rpcId")
		if filed == nil {
			return true
		}
		//获取rpcId
		rpcId := filed.Number()
		if rpcId < 10000 { //rpcId从最小的心跳开始
			return true
		}
		//获取rpc名字
		rpcName := strings.Split(string(fullName), ".")[1]
		//解析
		if strings.Index(rpcName, "Push") >= 0 {
			parseScMessage(int32(rpcId), fullName, routers)
		} else {
			parseCsMessage(int32(rpcId), fullName, routers)
		}
		return true
	})
}

// GetRequestMsgByRpcId 获取请求的消息体
func GetRequestMsgByRpcId(rpcId int32) (proto.Message, bool) {
	msg, ok := requestRpcIdMap[rpcId]
	if !ok {
		return nil, false
	}
	return msg.request.New().Interface().(proto.Message), true
}

// GetRequestMsgByName 获取请求的消息体
func GetRequestMsgByName(rpcName string) (proto.Message, bool) {
	msg, ok := requestRpcNameMap[rpcName]
	if !ok {
		return nil, false
	}
	return msg.request.New().Interface().(proto.Message), true
}

// GetPushRpcIdByMsg 获取push的rpcId
func GetPushRpcIdByMsg(msg proto.Message) (int32, bool) {
	typeName := msg.ProtoReflect().Descriptor().FullName()
	rpcId, ok := pushRpcNameMap[typeName]
	return rpcId.rpcId, ok
}

// GetPushMsgByRpcId 获取push的rpcId
func GetPushMsgByRpcId(rpcId int32) (proto.Message, bool) {
	msg, ok := pushRpcIdMap[rpcId]
	if !ok {
		return nil, false
	}
	return msg.push.New().Interface().(proto.Message), ok
}

// GetRequestByUri 通过uri获取request
func GetRequestByUri(uri string) (proto.Message, bool) {
	//去除?
	if len(uri) > 0 {
		idx := strings.Index(uri, "?")
		if idx >= 0 {
			uri = uri[:idx] // 去掉URI中的`?`以及之后的内容
		}
	}
	//取最后一个"/"之后的字符串
	if len(uri) > 0 {
		idx := strings.LastIndex(uri, "/")
		if idx >= 0 {
			uri = uri[idx+1:]
		}
	}
	//首字母转大写--因为http的uri格式的原因
	if len(uri) > 0 {
		first := unicode.ToUpper(rune(uri[0]))
		uri = string(first) + uri[1:]
	}
	//
	return GetRequestMsgByName(uri)
}

func checkRepeated(rpcId int32, fullName protoreflect.FullName, isRequest bool) {
	if isRequest {
		//检查请求是否重复
		rpcName := strings.Split(string(fullName), ".")[1]
		if _, ok := requestRpcNameMap[rpcName]; ok {
			panic(fmt.Sprintf("rpc name define repeated,request rpcName:%v", rpcName))
		}
	} else {
		//检查push是否重复
		pushName := fullName + ".Push"
		if _, ok := pushRpcNameMap[pushName]; ok {
			panic(fmt.Sprintf("rpc id define repeated, push pushMsgName:%v", pushName))
		}
	}
	//id不能重复
	if _, ok := requestRpcIdMap[rpcId]; ok {
		panic(fmt.Sprintf("rpc id define repeated,request rpcId:%v", rpcId))
	}
	//id不能重复
	for _, def := range pushRpcNameMap {
		if def.rpcId == rpcId {
			panic(fmt.Sprintf("rpc id define repeated, push pushRpcId:%v", rpcId))
		}
	}
}

func parseCsMessage(rpcId int32, fullName protoreflect.FullName, routers []string) {
	//重复检查
	checkRepeated(rpcId, fullName, true)
	//
	rpcName := strings.Split(string(fullName), ".")[1]
	//获取request和reply
	//
	requestMsgName := fullName + ".Request"
	requestType, err := protoregistry.GlobalTypes.FindMessageByName(requestMsgName)
	if err != nil {
		panic(fmt.Sprintf("find request message err, err:%v", err.Error()))
	}
	replyMsgName := fullName + ".Reply"
	replyType, err := protoregistry.GlobalTypes.FindMessageByName(replyMsgName)
	if err != nil {
		panic(fmt.Sprintf("find reply message err, err:%v", err.Error()))
	}
	//save
	rpc := tRequest{rpcId: rpcId, rpcName: rpcName, request: requestType, reply: replyType, routers: routers}
	requestRpcIdMap[rpcId] = rpc
	requestRpcNameMap[rpcName] = rpc
	return
}
func parseScMessage(rpcId int32, fullName protoreflect.FullName, routers []string) {
	//重复检查
	checkRepeated(rpcId, fullName, false)
	//check msg.push
	pushName := fullName + ".Push"
	pushType, err := protoregistry.GlobalTypes.FindMessageByName(pushName)
	if err != nil {
		panic(fmt.Sprintf("find request message err, err:%v", err.Error()))
	}
	//save
	rpc := tPush{rpcId: rpcId, push: pushType, routers: routers}
	pushRpcNameMap[pushName] = rpc
	pushRpcIdMap[rpcId] = rpc
}
func String() string {
	str := "请求应答接口:"
	//请求应答接口
	{
		var keys []int
		for key := range requestRpcIdMap {
			keys = append(keys, int(key))
		}
		sort.Sort(sort.IntSlice(keys))
		for _, rpcId := range keys {
			val := requestRpcIdMap[int32(rpcId)]
			str += "\nrpcId:" + strconv.Itoa(rpcId)
			str += "  rpcName:" + val.rpcName
			str += "  request:" + reflect.TypeOf(val.request.New().Interface()).Elem().Name()
			str += "  reply:" + reflect.TypeOf(val.reply.New().Interface()).Elem().Name()
		}
	}
	//推送接口
	{
		var keys []int
		for _, val := range pushRpcNameMap {
			keys = append(keys, int(val.rpcId))
		}
		str += "\n推送接口"
		sort.Sort(sort.IntSlice(keys))
		for _, rpcId := range keys {
			for pushName, val := range pushRpcNameMap {
				if val.rpcId != int32(rpcId) {
					continue
				}
				str += "\nrpcId:" + strconv.Itoa(int(val.rpcId))
				str += "  push:" + string(pushName)
				break
			}
		}
	}
	return str
}

func IsMessageInPackage(msg proto.Message) bool {
	packageName := "outer"
	desc := msg.ProtoReflect().Descriptor()
	if desc == nil {
		return false
	}
	fullName := string(desc.FullName())
	return fullName == packageName || fullName[:len(packageName)+1] == packageName+"."
}
