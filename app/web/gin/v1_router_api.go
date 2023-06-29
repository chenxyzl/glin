package gin

import (
	"github.com/gin-gonic/gin"
	"laiya/share/global"
)

// api
func startApiRouterV1(r *gin.Engine) {
	const path = "/" + global.WebUrlBase + "/v1"
	//创建group
	group := r.Group(path)
	group.Use(firstHook)
	//无需授权的
	group.Group("/getCaptcha").Use(v1ParseBody).Any("", v1DoMsg)
	group.Group("/checkCaptcha").Use(v1ParseBody).Any("", v1DoMsg)
	group.Group("/getGameTypes").Use(v1ParseBody).Any("", v1DoMsg)
	group.Group("/getThirdPlatformConfig").Use(v1ParseBody).Any("", v1DoMsg)
	group.Group("/getHallBannersConfig").Use(v1ParseBody).Any("", v1DoMsg)
	group.Group("/getClientDynamicConfig").Use(v1ParseBody).Any("", v1DoMsg)
	group.Group("/getShareGroupCard").Use(v1ParseBody).Any("", v1DoMsg)
	group.Group("/checkClientVersion").Use(v1ParseBody).Any("", v1DoMsg)
	//上传-需要授权-不需要解析body
	group.Group("/uploadIcon").Use(v1AuthMiddle).Any("", v1DoMsg)
	//群插件
	//group.Group("/group/:gid/plugin/:plugin").Use(groupPluginHook).Use(v1ParseBody).Use(v1AuthMiddle).Any("/:path", v1GroupPluginMsg)
	//其他-需要授权的也许接口
	//先解析body,后检查授权,因为reply的格式依赖于body解析的uri
	group.Use(v1ParseBody).Use(v1AuthMiddle).Any("/:path", v1DoMsg)
}
