package outer

import (
	"encoding/base64"
	"fmt"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"os"
	"reflect"
	"strings"
	"testing"
)

func writePb(t *testing.T, message proto.Message) {
	data, err := proto.Marshal(message)
	if err != nil {
		t.Error(err)
	}
	parentName := strings.TrimSuffix(string(message.ProtoReflect().Descriptor().FullName()), ".Request")
	parent, err := protoregistry.GlobalTypes.FindMessageByName(protoreflect.FullName(parentName))
	if err != nil {
		t.Error(err)
	}
	filed := parent.Descriptor().Fields().ByName("rpcId")
	if filed == nil {
		t.Error()
	}
	fileName := fmt.Sprintf("../../bin/%d_%s.bin", filed.Number(), reflect.TypeOf(message).Elem().Name())
	err = os.WriteFile(fileName, data, os.FileMode(0644|os.O_TRUNC))
	if err != nil {
		t.Error()
	}
}

func writeJson(t *testing.T, message proto.Message) {
	data, err := PbJson.Marshal(message)
	if err != nil {
		t.Error(err)
	}
	parentName := strings.TrimSuffix(string(message.ProtoReflect().Descriptor().FullName()), ".Request")
	parent, err := protoregistry.GlobalTypes.FindMessageByName(protoreflect.FullName(parentName))
	if err != nil {
		t.Error(err)
	}
	filed := parent.Descriptor().Fields().ByName("rpcId")
	if filed == nil {
		t.Error()
	}
	fileName := fmt.Sprintf("../../bin/%d_%s.json", filed.Number(), reflect.TypeOf(message).Elem().Name())
	err = os.WriteFile(fileName, data, os.FileMode(0644|os.O_TRUNC))
	if err != nil {
		t.Error()
	}
}

func writePbPack(t *testing.T, message proto.Message) {
	parentName := strings.TrimSuffix(string(message.ProtoReflect().Descriptor().FullName()), ".Request")
	parent, err := protoregistry.GlobalTypes.FindMessageByName(protoreflect.FullName(parentName))
	if err != nil {
		t.Error(err)
	}
	filed := parent.Descriptor().Fields().ByName("rpcId")
	if filed == nil {
		t.Error()
	}
	pack := &ReqPack{
		Sn:    0,
		RpcId: int32(filed.Number()),
	}
	req, err := proto.Marshal(message)
	if err != nil {
		t.Error()
	}
	pack.Data = req
	//
	data, err := proto.Marshal(pack)
	if err != nil {
		t.Error()
	}
	base64Data := base64.StdEncoding.EncodeToString(data)

	fileName := fmt.Sprintf("../../bin/%d_%s.pack.bin", filed.Number(), reflect.TypeOf(message).Elem().Name())
	err = os.WriteFile(fileName, []byte(base64Data), os.FileMode(0644|os.O_TRUNC))
	if err != nil {
		t.Error()
	}
}

func writeJsonPack(t *testing.T, message proto.Message) {
	parentName := strings.TrimSuffix(string(message.ProtoReflect().Descriptor().FullName()), ".Request")
	parent, err := protoregistry.GlobalTypes.FindMessageByName(protoreflect.FullName(parentName))
	if err != nil {
		t.Error(err)
	}
	filed := parent.Descriptor().Fields().ByName("rpcId")
	if filed == nil {
		t.Error()
	}
	pack := &ReqPack{
		Sn:    0,
		RpcId: int32(filed.Number()),
	}
	req, err := PbJson.Marshal(message)
	if err != nil {
		t.Error()
	}
	pack.Data = req
	//
	data, err := PbJson.Marshal(pack)
	if err != nil {
		t.Error()
	}

	fileName := fmt.Sprintf("../../bin/%d_%s.pack.json", filed.Number(), reflect.TypeOf(message).Elem().Name())
	err = os.WriteFile(fileName, data, os.FileMode(0644|os.O_TRUNC))
	if err != nil {
		t.Error()
	}
}
