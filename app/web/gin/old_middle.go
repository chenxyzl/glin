package gin

import (
	"github.com/chenxyzl/glin/share"
	"github.com/chenxyzl/glin/slog"
	"github.com/gin-gonic/gin"
	"laiya/config"
	"laiya/proto/code"
	"laiya/proto/outer"
	"laiya/share/global"
	token2 "laiya/share/token"
	"net/http"
	"strings"
)

func crossMiddle(c *gin.Context) {
	//
	defer share.RecoverFunc(func(err any) {
		c.ProtoBuf(http.StatusOK, outer.BuildHttpErr(code.Code_InnerError))
	})
	//
	//method := c.Request.Method
	origin := c.GetHeader("Origin") //请求头部
	if origin != "" {
		//c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Token")
		c.Header("Access-Control-Expose-Headers", "Access-Control-Allow-Headers, Token")
		c.Header("Access-Control-Max-Age", "172800")
		c.Header("Access-Control-Allow-Credentials", "true")
	}
	//if method == "OPTIONS" {
	//	c.AbortWithStatus(http.StatusNoContent)
	//}
	//
	slog.Infof("middle origin:%v", origin)
	//
	c.Next()
}

func filterIcon(c *gin.Context) {
	//
	defer share.RecoverFunc(func(err any) {
		c.ProtoBuf(http.StatusOK, outer.BuildHttpErr(code.Code_InnerError))
	})

	if strings.Contains(c.Request.URL.Path, "favicon.ico") {
		c.AbortWithStatus(404)
		return
	}
	c.Next()
}

func authMiddle(c *gin.Context) {
	//
	defer share.RecoverFunc(func(err any) {
		c.ProtoBuf(http.StatusOK, outer.BuildHttpErr(code.Code_InnerError))
	})

	// 从请求头中获取授权令牌
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.ProtoBuf(http.StatusOK, outer.BuildHttpErr(code.Code_TokenNotExist))
		c.Abort()
		return
	}
	// 检查授权令牌是否是Bearer类型
	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || (authHeaderParts[0] != "Bearer" && authHeaderParts[0] != "Bear") {
		c.ProtoBuf(http.StatusOK, outer.BuildHttpErr(code.Code_TokenFormatErr))
		c.Abort()
		return
	}
	// 从请求头中获取 Token
	token := authHeaderParts[1]
	uid, err := token2.ParseToken(token, config.Get().AConfig.AppKey)
	if err != nil {
		slog.Errorf("token过期,uid:%v|token:%v", uid, token)
		c.ProtoBuf(http.StatusOK, outer.BuildHttpErr(code.Code_TokenErr))
		c.Abort()
		return
	}
	//todo uid对应的用户数据必须要先存在~token生成只有我们内部有知道,这里就不校验uid对应的数据是否合法了
	c.Set(global.ParamToken, token)
	c.Set(global.ParamUid, uid)
	//
	c.Next()
}
