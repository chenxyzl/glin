package gin

import (
	"github.com/chenxyzl/glin/slog"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
	"laiya/proto/code"
	"laiya/share/call"
	"laiya/share/global"
	"strconv"
	"time"
)

// hook
func groupPluginHook(c *gin.Context) {
	c.Set(global.ParamGid, uint64(0)) //占位 避免后面类型强转失败
	c.Set(global.ParamPlugin, "")     //占位 避免后面类型强转失败
	gidStr := c.Param(global.ParamGid)
	pluginName := c.Param(global.ParamPlugin)
	gid, err := strconv.ParseUint(gidStr, 10, 64)
	if err != nil {
		replyErr(c, code.Code_UrlFormatErr)
		c.Abort()
		slog.Errorf("parse gid error; uri:%s|gid:%s|err:%v", c.Request.URL, gidStr, err)
		return
	}
	if pluginName == "" {
		replyErr(c, code.Code_UrlFormatErr)
		c.Abort()
		slog.Errorf("pluginName not found; uri:%s|gid:%s|err:%v", c.Request.URL, gidStr, err)
		return
	}
	pluginType := string(global.ParamPlugin) + "_" + string(global.GroupKind) + "_" + pluginName
	//
	c.Set(global.ParamGid, gid)           //
	c.Set(global.ParamPlugin, pluginType) //
	c.Next()
}

// 实际处理业务
func v1GroupPluginMsg(c *gin.Context) {
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
			slog.Errorf("deal group plugin msg error, uri:%v|cost:%vms|req:%v|cod:%v|err:%v", uri, cost.Milliseconds(), req, cod, err)
		} else {
			slog.Infof("deal group plugin msg success, uri:%v|cost:%vms|req:%v|rsp:%v", uri, cost.Milliseconds(), req, rsp)
		}
	}()

	rsp, cod = call.RequestRemote[proto.Message](uid, req)
	//
	if cod != code.Code_Ok {
		replyErr(c, cod)
	} else {
		replyOk(c, rsp)
	}
}
