package gin

import (
	"github.com/chenxyzl/glin/slog"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"io"
	"laiya/proto/code"
)

// 读取body
func readPbBody(c *gin.Context, req proto.Message) code.Code {
	//用类型来反序列化data
	var buf []byte
	var err error
	buf, err = io.ReadAll(c.Request.Body)
	uri := c.Request.URL.RequestURI()
	if err != nil {
		slog.Errorf("read body err, uri:%v|err:%v", uri, err)
		return code.Code_DataReadErr
	}

	if len(buf) > 0 {
		if err = proto.Unmarshal(buf, req); err != nil {
			slog.Errorf("proto data unmarshal err, uri:%v|data:%v|err:%v", uri, string(buf), err)
			return code.Code_MsgUnmarshalErr
		}
	}
	return code.Code_Ok
}

// 读取body
func readJsonBody(c *gin.Context, req proto.Message) code.Code {
	//用类型来反序列化data
	var buf []byte
	var err error
	buf, err = io.ReadAll(c.Request.Body)
	uri := c.Request.URL.RequestURI()
	if err != nil {
		slog.Errorf("read body data err, uri:%v|err:%v", uri, err)
		return code.Code_DataReadErr
	}

	if len(buf) > 0 {
		if err = protojson.Unmarshal(buf, req); err != nil {
			slog.Errorf("json data unmarshal err, uri:%v|data:%v|err:%v", uri, string(buf), err)
			return code.Code_MsgUnmarshalErr
		}
	}
	return code.Code_Ok
}
