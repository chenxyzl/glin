package wechat

import (
	"github.com/chenxyzl/glin/share"
	"github.com/chenxyzl/glin/slog"
	"github.com/eatmoreapple/openwechat"
	"laiya/wechat/iface"
	"sync"
)

var bot *openwechat.Bot
var self *openwechat.Self
var lock sync.RWMutex

func StartWechat(entity iface.IWechatActor) {
	//check 避免多个节点重复登录
	lock.Lock()
	defer lock.Unlock()
	if bot != nil {
		bot.Exit()
		bot = nil
	}
	//
	var wechatErr error
	defer func() {
		//异常退出要告诉actor
		entity.Next(func(params ...any) {
			bot := params[0].(*openwechat.Bot)
			wechatErr := params[1].(error)
			entity.WechatExit(bot, wechatErr)
		}, bot, wechatErr)
		//函数退出检查
		if wechatErr != nil {
			slog.Errorf("wechat quit with wechatErr, wechatErr:%v", wechatErr)
		} else {
			slog.Infof("wechat close success...")
		}
	}()
	slog.Infof("wechat start...")

	//设置默认状态
	bot = openwechat.DefaultBot(openwechat.Desktop) // 桌面模式

	// 注册登陆二维码回调
	bot.UUIDCallback = QrcodeUrl

	// 创建热存储容器对象
	reloadStorage := openwechat.NewFileHotReloadStorage("wechat_hot_login.json")
	defer reloadStorage.Close()

	// 执行热登录
	wechatErr = bot.HotLogin(reloadStorage)
	if wechatErr != nil { // 热登录失败则执行普通登录
		if wechatErr = bot.Login(); wechatErr != nil {
			//不用做其他错误, defer中会处理退出相关事情
			slog.Errorf("wechat login wechatErr, wechatErr:%v", wechatErr)
			return
		}
	}

	// 获取登陆的用户
	self, wechatErr = bot.GetCurrentUser()
	if wechatErr != nil {
		//不用做其他错误, defer中会处理退出相关事情
		slog.Errorf("wechat get current user wechatErr, wechatErr:%v", wechatErr)
		return
	}
	slog.Infof("wechat get self:%v", self)
	// 获取所有的好友
	var friends openwechat.Friends
	friends, wechatErr = self.Friends()
	if wechatErr != nil {
		//不用做其他错误, defer中会处理退出相关事情
		slog.Errorf("wechat get friends wechatErr, wechatErr:%v", wechatErr)
		return
	}
	slog.Infof("wechat get friends:%v", friends)
	//更新群信息
	var groups openwechat.Groups
	groups, wechatErr = self.Groups()
	if wechatErr != nil {
		//不用做其他错误, defer中会处理退出相关事情
		slog.Errorf("wechat get groups wechatErr, wechatErr:%v", wechatErr)
		return
	}
	for _, group := range groups {
		entity.Next(func(params ...any) {
			group := params[0].(*openwechat.Group)
			entity.UpdateOpenWechatGroupInfo(group)
		}, group)
	}
	slog.Infof("wechat get groups:%v", groups)
	//有错误暂时不退出,给个高级别的警告
	errorTimes := 0
	bot.MessageErrorHandler = func(err error) error {
		slog.Errorf("wechat handler err, err:%v", err)
		errorTimes++
		if errorTimes <= 3 {
			return nil
		}
		//返回错误会退出block状态,并且返回错误,所以这里不用做其他处理
		return err
	}

	//状态同步
	bot.SyncCheckCallback = func(resp openwechat.SyncCheckResponse) {
		defer share.RecoverInfo("wechat sync panic")
		//
		if resp.Success() {
			if !resp.NorMal() {
				checkCodeUpdate(resp)
				slog.Infof("wechat sync code, code:%v|select:%v", resp.RetCode, resp.Selector)
			}
			return
		}
		//有错误不用处理--会进入MessageErrorHandler
		slog.Errorf("wechat sync code err, resp:%v", resp)
	}
	//消息处理
	dispatcher := openwechat.NewMessageMatchDispatcher()
	//dispatcher.SetAsync(true) //异步处理
	bot.MessageHandler = dispatcher.AsMessageHandler()
	//记录聊天
	dispatcher.OnText(recordText)
	//更新群信息
	dispatcher.OnGroup(func(ctx *openwechat.MessageContext) { doUpdateGroup(ctx, entity) })
	//处理其他文本消息
	dispatcher.OnText(func(ctx *openwechat.MessageContext) { errorTimes = 0 })
	//实际的业务处理
	dispatcher.OnText(func(ctx *openwechat.MessageContext) { doText(ctx, entity) })
	//处理记录邀请
	dispatcher.OnGroup(checkInvited)
	//
	entity.Next(func(params ...any) {
		bot := params[0].(*openwechat.Bot)
		entity.WechatLoginSuccess(bot)
	}, bot)
	//
	slog.Infof("chat start success...")
	//
	wechatErr = bot.Block()
}
