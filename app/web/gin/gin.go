package gin

import (
	"fmt"
	"github.com/chenxyzl/glin/slog"
	"github.com/gin-gonic/gin"
	"laiya/config"
)

func StartGin() {
	url := fmt.Sprintf(":%v", config.Get().AConfig.WebHttpParam.Port)
	slog.Infof("gin starting, url:%v", url)
	r := gin.Default()
	r.Use(crossMiddle, filterIcon)
	//重定向
	startRedirectRouter(r)
	//老版本的router
	startOldRouter(r)
	//新版本的router
	startApiRouterV1(r)
	//gm的router
	startGmRouterV1(r)

	go func() {
		err := r.Run(url)
		if err != nil {
			slog.Fatalf("start gin http server err,", err)
		}
	}()
}
