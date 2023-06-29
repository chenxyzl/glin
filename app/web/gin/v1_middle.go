package gin

import (
	"github.com/chenxyzl/glin/grain"
	"github.com/chenxyzl/glin/share"
	"github.com/chenxyzl/glin/slog"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
	"laiya/config"
	"laiya/proto/code"
	"laiya/proto/outer"
	"laiya/share/call"
	"laiya/share/global"
	token2 "laiya/share/token"
	"laiya/web"
	"net/url"
	"strings"
	"time"
)

// hook
func firstHook(c *gin.Context) {
	c.Set(global.ParamUid, uint64(0)) //占位 避免后面类型强转失败
	c.Set(global.ProtoName, "")       //占位 避免后面类型强转失败
	//解析query
	params, err := url.ParseQuery(c.Request.URL.RawQuery)
	if err != nil {
		replyErr(c, code.Code_UrlFormatErr)
		c.Abort()
		slog.Errorf("parse query error; uri:%s|rq:%s|err:%v", c.Request.URL, c.Request.URL.RawQuery, err)
		return
	}
	//设置消息协议类型
	c.Set(global.ProtoName, params.Get(global.ProtoName))
	//获取uri对应的proto
	uri := c.Request.URL.RequestURI()
	req, ok := outer.GetRequestByUri(uri)
	if !ok {
		replyErr(c, code.Code_NotImpl)
		c.Abort()
		slog.Errorf("not found proto msg from uri, uri:%v", uri)
		return
	}
	//设置request
	c.Set(global.ParamRequest, req)
}

// 解析body
func v1ParseBody(c *gin.Context) {
	//获取uri对应的proto
	var req = c.MustGet(global.ParamRequest).(proto.Message)
	//解析body
	if c.MustGet(global.ProtoName) == global.ProtocolJson {
		//如果是json格式
		cod := readJsonBody(c, req)
		if cod != code.Code_Ok {
			replyErr(c, cod)
			c.Abort()
			return
		}
	} else {
		//否者认为是proto格式
		cod := readPbBody(c, req)
		if cod != code.Code_Ok {
			replyErr(c, cod)
			c.Abort()
			return
		}
	}
	//设置request
	c.Set(global.ParamRequest, req)
	//
	c.Next()
}

// 授权
func v1AuthMiddle(c *gin.Context) {
	//
	defer share.RecoverFunc(func(err any) {
		replyErr(c, code.Code_InnerError)
	})
	// 从请求头中获取授权令牌
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		replyErr(c, code.Code_TokenNotExist)
		c.Abort()
		return
	}
	// 检查授权令牌是否是Bearer类型
	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || (authHeaderParts[0] != "Bearer" && authHeaderParts[0] != "Bear") {
		replyErr(c, code.Code_TokenFormatErr)
		c.Abort()
		return
	}
	// 从请求头中获取 Token
	token := authHeaderParts[1]
	uid, err := token2.ParseToken(token, config.Get().AConfig.AppKey)
	if err != nil {
		slog.Errorf("token过期,uid:%v|token:%v", uid, token)
		replyErr(c, code.Code_TokenErr)
		c.Abort()
		return
	}
	c.Set(global.ParamToken, token)
	c.Set(global.ParamUid, uid)
	//
	c.Next()
}

// 实际处理业务
func v1DoMsg(c *gin.Context) {
	var req = c.MustGet(global.ParamRequest).(proto.Message)
	var uid = c.MustGet(global.ParamUid).(uint64)
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
	//
	if cod != code.Code_Ok {
		replyErr(c, cod)
	} else {
		replyOk(c, rsp)
	}
}
