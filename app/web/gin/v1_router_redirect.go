package gin

import (
	"fmt"
	"github.com/chenxyzl/glin/slog"
	"github.com/gin-gonic/gin"
	"laiya/config"
	"laiya/proto/code"
	"laiya/share/utils"
	"net/http"
	"strings"
)

// 重定向
func startRedirectRouter(r *gin.Engine) {
	group := r.Group(config.Get().WebConfig.ShorUrlGroup)
	group.Any("/:path", v1RedirectRouter)
}

// 重定向
func v1RedirectRouter(c *gin.Context) {
	uri := c.Request.URL.RequestURI()
	uri = strings.TrimLeft(uri, config.Get().WebConfig.ShorUrlGroup+"/")

	idx := strings.Index(uri, "?")
	params := ""
	if idx >= 0 {
		uri = uri[:idx] // 去掉URI中的`?`以及之后的内容
		params = uri[idx:]
	}
	gid, err := utils.Base62ToUint64(uri)
	if err != nil {
		replyErr(c, code.Code_UrlFormatErr)
		slog.Infof("url format err, url:%v", c.Request.URL.RequestURI())
		return
	}
	url := fmt.Sprintf("%v%v%v", config.Get().WebConfig.LongUrl, gid, params)
	slog.Infof("short link redirect,shortUrl:%v|longUrl:%v", c.Request.URL, url)
	c.Redirect(http.StatusFound, url)
}
