package gin

import (
	"github.com/chenxyzl/glin/slog"
	"github.com/gin-gonic/gin"
	"laiya/config"
	"laiya/proto/code"
	"laiya/share/global"
)

// gm
func startGmRouterV1(r *gin.Engine) {
	const gmPath = "/" + global.GmUrlBase + "/v1"
	group := r.Group(gmPath)
	group.Use(firstHook)
	group.Use(v1ParseBody).Use(v1AuthMiddle).Use(gmAuthMiddle).Any("/:path", v1DoMsg)
}

// hook
func gmAuthMiddle(c *gin.Context) {
	uid := c.MustGet(global.ParamUid).(uint64)
	uri := c.Request.RequestURI
	if !config.Get().GmConfig.CheckGmAuth(uid, uri) {
		replyErr(c, code.Code_GmAuthError)
		c.Abort()
		slog.Errorf("gm auth error; uri:%s|uid:%v", c.Request.RequestURI, uid)
		return
	}
	c.Next()
}
