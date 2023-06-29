package gin

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
	"laiya/proto/code"
	"laiya/proto/outer"
	"laiya/share/global"
	"net/http"
)

// 回复ok
func replyOk(c *gin.Context, reply proto.Message) {
	if c.MustGet(global.ProtoName) == global.ProtocolJson {
		c.Render(http.StatusOK, outer.BuildCusRenderOk(reply))
	} else {
		c.ProtoBuf(http.StatusOK, outer.BuildHttpOk(reply))
	}
}

// 回复错误
func replyErr(c *gin.Context, cod code.Code) {
	if c.MustGet(global.ProtoName) == global.ProtocolJson {
		c.Render(http.StatusOK, outer.BuildCusRenderErr(cod))
	} else {
		c.ProtoBuf(http.StatusOK, outer.BuildHttpErr(cod))
	}
}
