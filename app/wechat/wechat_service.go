package wechat

import (
	"github.com/chenxyzl/glin/share"
	"github.com/chenxyzl/glin/slog"
	"github.com/eatmoreapple/openwechat"
	"laiya/config"
	"laiya/share/datarangers"
	"laiya/wechat/iface"
	"math/rand"
	"strings"
	"time"
)

// QrcodeUrl 打印登录二维码
func QrcodeUrl(uuid string) {
	qrcodeUrl := openwechat.GetQrcodeUrl(uuid)
	slog.Infof("访问下面网址扫描二维码登录")
	slog.Infof(qrcodeUrl)
}

// 消息格式化
func formatContent(msg string) string {
	if self == nil {
		return msg
	}
	var atFlag = "@" + self.NickName
	if strings.Contains(msg, "\u2005") {
		atFlag += "\u2005"
	}

	idx := strings.Index(msg, atFlag)
	if idx == 0 { // >= 改为==0 //
		idx = idx + len(atFlag)
		if idx < len(msg) {
			msg = msg[idx:]
		}
	}
	//半空格替换为全空格替换
	msg = strings.ReplaceAll(msg, "\u2005", " ")
	//去掉左右空格
	msg = strings.TrimLeft(msg, " ")
	msg = strings.TrimRight(msg, " ")
	return msg
}

// ForceUpdateGroups 强制更新群信息
const ForceUpdateGroups openwechat.Selector = "ForceUpdateGroups"

// ForceUpdateFriends 强制更新好友信息
const ForceUpdateFriends openwechat.Selector = "ForceUpdateFriends"

// 检查更新联系人
func checkCodeUpdate(resp openwechat.SyncCheckResponse) {
	//检查变更
	if resp.Selector == openwechat.SelectorAddOrDelContact || resp.Selector == openwechat.SelectorModContact || resp.Selector == ForceUpdateFriends {
		_, err := self.Friends(true)
		if err != nil {
			slog.Errorf("update friend err, err:%v", err)
		} else {
			slog.Infof("update friend success")
		}
	}
	if resp.Selector == openwechat.SelectorModChatRoom || resp.Selector == ForceUpdateGroups {
		_, err := self.Groups(true)
		if err != nil {
			slog.Errorf("update wechat group err, err:%v", err)
		} else {
			slog.Infof("update wechat group success")
		}
	}
}

// 更新群信息
func doUpdateGroup(ctx *openwechat.MessageContext, entity iface.IWechatActor) {
	//
	defer share.Recover()

	//自己的不管
	if ctx.IsSendBySelf() {
		return
	}
	//只处理自己的
	if ctx.ToUserName != self.UserName {
		return
	}
	//引用的消息不处理
	if strings.Contains(ctx.Content, "- - - - - - - - - - - - - - -") {
		return
	}
	//sender
	sender, _ := ctx.Sender()
	if sender == nil {
		return
	}
	//群
	if !ctx.IsSendByGroup() { //群消息
		return
	}
	//获取发送者
	senderUser, _ := ctx.SenderInGroup() //只at且处理人发的
	if senderUser == nil {
		return
	}
	//转为group
	group, _ := sender.AsGroup()
	if group == nil {
		return
	}
	//更新群信息
	entity.Next(func(params ...any) {
		group := params[0].(*openwechat.Group)
		entity.UpdateOpenWechatGroupInfo(group)
	}, group)
}

// 文本业务处理
func doText(ctx *openwechat.MessageContext, entity iface.IWechatActor) {
	//
	defer share.Recover()

	//自己的不管
	if ctx.IsSendBySelf() {
		return
	}
	//只处理自己的
	if ctx.ToUserName != self.UserName {
		return
	}
	//引用的消息不处理
	if strings.Contains(ctx.Content, "- - - - - - - - - - - - - - -") {
		return
	}
	//sender
	sender, _ := ctx.Sender()
	if sender == nil {
		return
	}
	//群
	if ctx.IsSendByGroup() { //群消息
		//只处理@消息
		if !ctx.IsAt() {
			return
		}
		//获取发送者
		senderUser, _ := ctx.SenderInGroup() //只at且处理人发的
		if senderUser == nil {
			return
		}
		//转为group
		group, _ := sender.AsGroup()
		if group == nil {
			return
		}
		//命令检查
		fixMsg := formatContent(ctx.Content)
		slog.Infof("recv msg by group, group:%v|sender:%v|fixMsg:%v", sender.NickName, senderUser.NickName, fixMsg)
		command := GetWechatCommand(fixMsg)
		switch command {
		case WechatCommandLookingGroup:
			entity.Next(func(params ...any) {
				ctx := params[0].(*openwechat.MessageContext)
				entity.LookingGroup(ctx, group)
			}, ctx)
		case WechatCommandUpdateWechatGroups:
			checkCodeUpdate(openwechat.SyncCheckResponse{RetCode: "0", Selector: ForceUpdateGroups})
			delayReply(ctx, config.Get().GameCopywritingConfig.WechatUpdateInfoSuccess, group.NickName)
		case WechatCommandUpdateUpdateWechatFriends:
			checkCodeUpdate(openwechat.SyncCheckResponse{RetCode: "0", Selector: ForceUpdateFriends})
			delayReply(ctx, config.Get().GameCopywritingConfig.WechatUpdateInfoSuccess, group.NickName)
		default:
			delayReply(ctx, config.Get().GameCopywritingConfig.WechatHelper, group.NickName)
		}
	} else if ctx.IsSendByFriend() { //好友消息直接回复helper命令
		slog.Infof("recv msg by friend, sender:%v|msg:%v", sender.NickName, ctx.Content)
		delayReply(ctx, config.Get().GameCopywritingConfig.WechatPrivateReply, "")
		return
	}
}

// 聊天信息记录
func recordText(ctx *openwechat.MessageContext) {
	sender, err := ctx.Sender()
	if err != nil {
		slog.Error(err)
		return
	}
	groupName := ""
	senderName := sender.NickName
	if ctx.IsSendByGroup() {
		senderInGroup, err := ctx.SenderInGroup()
		if err != nil {
			slog.Error(err)
			return
		}
		//群消息,发送者为群名字,
		groupName = senderName
		senderName = senderInGroup.NickName
	}
	//groupName为空表示私聊消息
	slog.Infof("wechat recv msg, sender:%v|groupName:%v|msg:%v", senderName, groupName, ctx.Content)
}

// 邀请消息这个消息不是文本消息,所以需要这么处理
func checkInvited(ctx *openwechat.MessageContext) {
	//
	defer share.Recover()
	//群消息
	u, _ := ctx.Sender()
	if u == nil {
		return
	}
	group, _ := u.AsGroup()
	if group == nil {
		return
	}
	//如果是邀请-记录邀请信息
	if !ctx.IsJoinGroup() {
		return
	}
	//邀请方式1
	{
		tag1 := "通过扫描"
		tag2 := "分享的二维码加入群聊"
		idx1 := strings.Index(ctx.Content, tag1)
		idx2 := strings.Index(ctx.Content, tag2)
		if idx1 >= 0 && idx2 >= 0 && idx2 > idx1 {
			p := ctx.Content[:idx1]
			op := ctx.Content[idx1+len(tag1) : idx2]
			p = strings.TrimLeft(p, "\"")
			p = strings.TrimRight(p, "\"")
			op = strings.TrimLeft(op, "\"")
			op = strings.TrimRight(op, "\"")
			datarangers.SendEvent("wechat_invited", map[string]interface{}{"group_name": group.NickName, "op_name": op, "player_name": p}, nil)
		}
	}
	//邀请方式2
	{
		tag1 := "邀请"
		tag2 := "加入了群聊"
		idx1 := strings.Index(ctx.Content, tag1)
		idx2 := strings.Index(ctx.Content, tag2)
		if idx1 >= 0 && idx2 >= 0 && idx2 > idx1 {
			op := ctx.Content[:idx1]
			p := ctx.Content[idx1+len(tag1) : idx2]
			p = strings.TrimLeft(p, "\"")
			p = strings.TrimRight(p, "\"")
			op = strings.TrimLeft(op, "\"")
			op = strings.TrimRight(op, "\"")
			datarangers.SendEvent("wechat_invited", map[string]interface{}{"group_name": group.NickName, "op_name": op, "player_name": p}, nil)
		}
	}
	slog.Infof("wechat_invited, groupName:%v|content:%v", group.NickName, ctx.Content)
}

// 回复消息
func delayReply(ctx *openwechat.MessageContext, text string, groupName string) {
	go func() {
		//
		defer share.Recover()
		//已读标记检查~暂停0.5
		time.Sleep(time.Second/2 + time.Duration(rand.Int63())%time.Second)
		ctx.AsRead()
		//1s后回复
		time.Sleep(time.Second/2 + time.Duration(rand.Int63())%time.Second)
		//为了让每次回复的内容都不一样--避免被检查出来
		ctx.ReplyText(text)
		datarangers.SendEvent("wechat_robot_say", map[string]interface{}{
			"wechat_group_name": groupName,
			"reply_to":          1,
		}, nil)
		slog.Infof("wechat_reply, groupName:%v|content:%v", groupName, text)
	}()
}
