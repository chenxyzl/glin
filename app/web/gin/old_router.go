package gin

import (
	"github.com/chenxyzl/glin/grain"
	"github.com/chenxyzl/glin/slog"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
	"laiya/proto/code"
	"laiya/proto/outer"
	"laiya/share/call"
	"laiya/share/global"
	"laiya/web"
	"net/http"
	"time"
)

const prefixWeb = "/web"
const prefixShare = "/share"

func startOldRouter(r *gin.Engine) {
	//pb接口
	group := r.Group(prefixWeb)
	group.Use(firstHook)
	//无需授权的
	group.Group("/getCaptcha").Any("", readPbBodyNoUid)
	group.Group("/checkCaptcha").Any("", readPbBodyNoUid)
	group.Group("/getGameTypes").Any("", readPbBodyNoUid)
	group.Group("/getThirdPlatformConfig").Any("", readPbBodyNoUid)
	group.Group("/getHallBannersConfig").Any("", readPbBodyNoUid)
	//group.Group("/getClientDynamicConfig").Any("", readPbBodyNoUid)

	//上传接口-不需要授权
	//r.Group(prefixWeb+"/uploadWechatQrcode").Any("", ignoreBodyRouterNoUid)
	//上传-需要授权
	group.Group("/uploadIcon").Use(authMiddle).Any("", ignoreBodyRouterWithUid)
	//需要授权的
	group.Use(authMiddle).Any("/:path", readPbBodyWithUid)

	//json接口
	jsonGroup := r.Group(prefixShare)
	jsonGroup.Use(firstHook)
	jsonGroup.Group("/getShareGroupCard").Any("", readJsonBodyNoUid)
}

// 不需要读取body的请求~需要uid的
func ignoreBodyRouterWithUid(c *gin.Context) {
	uid := c.MustGet(global.ParamUid)
	req := c.MustGet(global.ParamRequest).(proto.Message)
	doPbMsg(c, req, uid.(uint64))
}

// 需要读取body的请求~不需要uid的
func readPbBodyNoUid(c *gin.Context) {
	req := c.MustGet(global.ParamRequest).(proto.Message)
	cod := readPbBody(c, req)
	if cod != code.Code_Ok {
		c.ProtoBuf(http.StatusOK, outer.BuildHttpErr(cod))
		return
	}
	doPbMsg(c, req, 0)
}

// 需要读取body的请求~需要uid的
func readPbBodyWithUid(c *gin.Context) {
	uid := c.MustGet(global.ParamUid)
	req := c.MustGet(global.ParamRequest).(proto.Message)
	cod := readPbBody(c, req)
	if cod != code.Code_Ok {
		c.ProtoBuf(http.StatusOK, outer.BuildHttpErr(cod))
		return
	}
	doPbMsg(c, req, uid.(uint64))
}

// 需要读取body的请求~不需要uid的
func readJsonBodyNoUid(c *gin.Context) {
	req := c.MustGet(global.ParamRequest).(proto.Message)
	cod := readJsonBody(c, req)
	if cod != code.Code_Ok {
		c.ProtoBuf(http.StatusOK, outer.BuildHttpErr(cod))
		return
	}
	doJsonMsg(c, req, 0)
}

// 实际处理业务
func doPbMsg(c *gin.Context, req proto.Message, uid uint64) {
	var rsp proto.Message
	var cod code.Code
	var err error
	var uri = c.Request.URL.RequestURI()
	startTime := time.Now()
	defer func() {
		cost := time.Now().Sub(startTime)
		if cod != code.Code_Ok {
			slog.Errorf("deal msg error, uri:%v|cost:%vms|req:%v|cod:%v|err:%v", uri, cost.Milliseconds(), req, cod, err)
		} else {
			slog.Infof("deal msg success, uri:%v|cost:%vms|req:%v|rsp:%v", uri, cost.Milliseconds(), req, rsp)
		}
	}()

	reqName := req.ProtoReflect().Descriptor().FullName()
	if grain.IsLocalMethod(web.WebEntityIns, reqName) {
		rsp, cod = call.LocalRequestHandler(reqName, web.WebEntityIns, req, c)
	} else {
		rsp, cod = call.RequestRemote[proto.Message](uid, req)
	}
	if cod != code.Code_Ok {
		c.ProtoBuf(http.StatusOK, outer.BuildHttpErr(cod))
		return
	}
	c.ProtoBuf(http.StatusOK, outer.BuildHttpOk(rsp))
}

// 实际处理业务
func doJsonMsg(c *gin.Context, req proto.Message, uid uint64) {
	var rsp proto.Message
	var cod code.Code
	var err error
	var uri = c.Request.URL.RequestURI()
	startTime := time.Now()
	defer func() {
		cost := time.Now().Sub(startTime)
		if cod != code.Code_Ok {
			slog.Errorf("deal msg error, uri:%v|cost:%vms|req:%v|cod:%v|err:%v", uri, cost.Milliseconds(), req, cod, err)
		} else {
			slog.Infof("deal msg success, uri:%v|cost:%vms|req:%v|rsp:%v", uri, cost.Milliseconds(), req, rsp)
		}
	}()

	reqName := req.ProtoReflect().Descriptor().FullName()
	if grain.IsLocalMethod(web.WebEntityIns, reqName) {
		rsp, cod = call.LocalRequestHandler(reqName, web.WebEntityIns, req, c)
	} else {
		rsp, cod = call.RequestRemote[proto.Message](uid, req)
	}
	if cod != code.Code_Ok {
		c.Render(http.StatusOK, outer.BuildCusRenderErr(cod))
		return
	}
	c.Render(http.StatusOK, outer.BuildCusRenderOk(rsp))
}
