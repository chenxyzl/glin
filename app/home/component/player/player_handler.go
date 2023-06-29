package player

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/chenxyzl/glin/uuid"
	"golang.org/x/exp/slices"
	"google.golang.org/protobuf/proto"
	"laiya/common_service"
	"laiya/config"
	"laiya/model/sub_model"
	"laiya/model/web_model"
	"laiya/proto/code"
	"laiya/proto/inner"
	"laiya/proto/outer"
	"laiya/share/call"
	"laiya/share/getui_helper"
	"time"
)

func (a *PlayerComponent) HandleC2HoSessionOnline(req *inner.C2HoSessionOnline_Request) (*inner.C2HoSessionOnline_Reply, code.Code) {
	a.GetLogger().Infof("player session online")
	//online
	firstSessOnline := a.GetModel().OnlineUpdateSess(req.GetLoginPlatformType(), req.GetSession(), func(oldSess *actor.PID) {
		call.Send(a.GetCtx(), oldSess, outer.BuildPushPack(&outer.PushKick_Push{Reason: outer.PushKick_LoginElsewhere}))
		a.GetLogger().Infof("player login else where, will kick out old")
	})
	_ = firstSessOnline
	//send to group
	for k := range a.GetModel().Groups {
		//延时发送 ,用next处理是因为担心加入的群过多,同步处理耗时太长,所以分批处理
		a.Next(func(params ...any) {
			gid := params[0].(uint64)
			a.TellGroupOnline(gid)
		}, k)
	}
	a.GetLogger().Infof("player session online success")
	return nil, code.Code_Ok
}

func (a *PlayerComponent) HandleC2HoSessionOffline(req *inner.C2HoSessionOffline_Request) (*inner.C2HoSessionOffline_Reply, code.Code) {
	a.GetLogger().Infof("player session offline")
	//offline
	lastSessOffline := a.GetModel().OfflineCleanSess(req.GetLoginPlatformType(), req.GetSession())
	//是否所有链接都离线~全都客户端链接离线才认为是离线
	if lastSessOffline {
		//send to gid
		msg := &inner.Ho2GrPlayerOffline_Request{Uid: a.GetModel().Uid, Session: a.GetCtx().Self()}
		for k := range a.GetModel().Groups {
			a.Next(func(params ...any) {
				gid := params[0].(uint64)
				_, cod := call.RequestRemote[proto.Message](gid, msg)
				if cod != code.Code_Ok {
					a.GetLogger().Errorf("request gid offline err, gid:%v|cod:%v", gid, cod)
				}
			}, k)
		}
	}
	////如果session全部断开则actor销毁
	//if len(a.GetModel().Sess) == 0 {
	//	a.Next(func(params ...any) {
	//		a.GetLogger().Info("player will stop, because all session close")
	//		a.GetCtx().Poison(a.GetCtx().Self())
	//	})
	//}
	a.GetLogger().Info("player session offline success")
	return nil, code.Code_Ok
}

func (a *PlayerComponent) HandleGetUserInfo(req *outer.GetUserInfo_Request) (*outer.GetUserInfo_Reply, code.Code) {
	var groupIds []uint64
	var selfGroupIds []uint64
	for gid, group := range a.GetModel().Groups {
		groupIds = append(groupIds, gid)
		if group.IsSelf {
			selfGroupIds = append(selfGroupIds, gid)
		}
	}
	rsp := &outer.GetUserInfo_Reply{
		PlayerInfo: &outer.GetUserInfo_PlayerInfo{
			Uid:          a.GetModel().Uid,
			Name:         a.GetModel().Name,
			Icon:         a.GetModel().Icon,
			Des:          a.GetModel().Des,
			GroupIds:     groupIds,
			SelfGroupIds: selfGroupIds,
		},
	}
	a.GetLogger().Infof("get user info, rsp:%v", rsp)
	return rsp, code.Code_Ok
}

func (a *PlayerComponent) HandleSetUserInfo(req *outer.SetUserInfo_Request) (*outer.SetUserInfo_Reply, code.Code) {
	if cod := common_service.CheckUserName(req.GetName()); cod != code.Code_Ok {
		a.GetLogger().Errorf("check user name err, name:%v", req.GetName())
		return nil, cod
	}
	if cod := common_service.CheckUserDes(req.GetDes()); cod != code.Code_Ok {
		a.GetLogger().Errorf("check user des err, des:%v", req.GetDes())
		return nil, cod
	}
	a.GetModel().SetUserName(req.GetName())
	a.GetModel().SetUserIcon(req.GetIcon())
	a.GetModel().SetUserDes(req.GetDes())
	//save
	err := a.GetModel().Save()
	if err != nil {
		a.GetLogger().Error(err)
	}
	a.GetLogger().Infof("set user info, model:%v", a.GetModel())
	return &outer.SetUserInfo_Reply{}, code.Code_Ok
}

func (a *PlayerComponent) HandleUpdateUserInfo(req *outer.UpdateUserInfo_Request) (*outer.UpdateUserInfo_Reply, code.Code) {
	if cod := common_service.CheckUserInfo(req.GetType(), req.GetStrVar()); cod != code.Code_Ok {
		a.GetLogger().Errorf("check update user info err, type:%v|val:%v", req.GetType(), req.GetStrVar())
		return nil, cod
	}
	switch req.GetType() {
	case outer.UpdateUserInfo_Name:
		a.GetModel().SetUserName(req.GetStrVar())
	case outer.UpdateUserInfo_Icon:
		a.GetModel().SetUserIcon(req.GetStrVar())
	case outer.UpdateUserInfo_Des:
		a.GetModel().SetUserDes(req.GetStrVar())
	}
	a.GetLogger().Infof("update user info success, req:%v|model:%v", req, a.GetModel())
	return &outer.UpdateUserInfo_Reply{}, code.Code_Ok
}

func (a *PlayerComponent) HandleGetUserCard(req *outer.GetUserCard_Request) (*outer.GetUserCard_Reply, code.Code) {
	rsp := &outer.GetUserCard_Reply{
		Player: &outer.GetUserCard_Player{
			Uid:      a.GetModel().Uid,
			Name:     a.GetModel().Name,
			Icon:     a.GetModel().Icon,
			Des:      a.GetModel().Des,
			IsOnline: a.GetModel().IsOnline(),
			UserType: outer.GetUserType(a.GetModel().Uid),
		},
	}
	a.GetLogger().Infof("get user card, rsp:%v", rsp)
	return rsp, code.Code_Ok
}

func (a *PlayerComponent) HandleCheckSync(req *outer.CheckSync_Request) (*outer.CheckSync_Reply, code.Code) {
	var gids []uint64
	for _, group := range a.GetModel().Groups {
		gids = append(gids, group.Gid)
		//todo fix 临时解决下
		delete(a.GetModel().LastChatIdMap, group.Gid)
	}
	ret, err := sub_model.LoadMultiGroupLastChatId(gids)
	if err != nil {
		a.GetLogger().Errorf("load group lastChatId err, gids:%v|err:%v", gids, err)
		return nil, code.Code_InnerError
	}
	rsp := &outer.CheckSync_Reply{
		GroupsMaxChatId:  ret,
		PrivateMaxChatId: a.GetModel().GetLastPrivateChatId(),
	}
	a.GetLogger().Infof("check sync, rsp:%v", rsp)
	return rsp, code.Code_Ok
}

func (a *PlayerComponent) HandleWeb2HoDeletePlayer(req *inner.Web2HoDeletePlayer_Request) (*inner.Web2HoDeletePlayer_Reply, code.Code) {
	//先退出群
	if len(a.GetModel().Groups) > 0 {
		a.GetLogger().Errorf("delete player account err, please quit all group, uid:%v|account:%v|groups:%v", a.GetModel().Uid, a.GetModel().Account, a.GetModel().Groups)
		return nil, code.Code_MustQuitAllGroupBefore
	}

	//先删除账户关联
	needDeleteAccount := &web_model.Account{Account: a.GetModel().Account}
	if err := needDeleteAccount.Delete(); err != nil {
		a.GetLogger().Errorf("delete player account err, uid:%v|account:%v|err:%v", a.GetModel().Uid, a.GetModel().Account, err)
	}
	//再删除用户数据
	a.GetModel().DeletePlayerData()
	//
	a.GetLogger().Infof("delete player success, uid:%v|account:%v", a.GetModel().Uid, a.GetModel().Account)
	return &inner.Web2HoDeletePlayer_Reply{Account: a.GetModel().Account}, code.Code_Ok
}

func (a *PlayerComponent) HandleUpdateThirdPlatformInfo(req *outer.UpdateThirdPlatformInfo_Request) (*outer.UpdateThirdPlatformInfo_Reply, code.Code) {
	if _, ok := config.Get().GamePlatformsConfig.PlatformTypes[req.GetId()]; !ok {
		return nil, code.Code_PlatformTypeNotExist
	}
	a.GetModel().SetThirdPlatformAccount(req.GetId(), req.GetCode())
	a.GetLogger().Infof("update third account success, req:%v", req)
	return &outer.UpdateThirdPlatformInfo_Reply{}, code.Code_Ok
}

func (a *PlayerComponent) HandleGetUserCardDetail(req *outer.GetUserCardDetail_Request) (*outer.GetUserCardDetail_Reply, code.Code) {
	rsp := &outer.GetUserCardDetail_Reply{
		Player: &outer.GetUserCardDetail_Player{
			Uid:      a.GetModel().Uid,
			Name:     a.GetModel().Name,
			Icon:     a.GetModel().Icon,
			Des:      a.GetModel().Des,
			IsOnline: a.GetModel().IsOnline(),
			UserType: outer.GetUserType(a.GetModel().Uid),
		},
	}
	a.GetLogger().Infof("get user card detail, rsp:%v", rsp)
	return rsp, code.Code_Ok
}

func (a *PlayerComponent) HandleSyncUserState(req *outer.SyncUserState_Request) (*outer.SyncUserState_Reply, code.Code) {
	rsp := &outer.SyncUserState_Reply{}
	if slices.Contains(req.GetActivityStateId(), int32(outer.SyncUserState_CanRecvInviteBroadcast)) {
		a.GetModel().LastSyncUserStateTime = time.Now().UnixNano()
	}
	a.GetLogger().Infof("sync user state, req:%v|rsp:%v", req, rsp)
	return rsp, code.Code_Ok
}

func (a *PlayerComponent) HandleGetThirdPlatformInfo(req *outer.GetThirdPlatformInfo_Request) (*outer.GetThirdPlatformInfo_Reply, code.Code) {
	rsp := &outer.GetThirdPlatformInfo_Reply{
		ThirdPlatformAccount: a.GetModel().ThirdPlatformAccount,
	}
	a.GetLogger().Infof("get user third platform, rsp:%v", rsp)
	return rsp, code.Code_Ok
}

func (a *PlayerComponent) HandleOutSendPrivateChatMsg(req *outer.SendPrivateChatMsg_Request) (*outer.SendPrivateChatMsg_Reply, code.Code) {
	var l int
	if req.GetMsgType() == outer.ChatMsg_Pic {
		l = len([]rune(req.GetUrl()))
	} else {
		l = len([]rune(req.GetContent()))
	}
	//长度检查
	if l <= 0 {
		a.GetLogger().Errorf("not send empty msg", a.GetModel().Uid)
		return nil, code.Code_Ok
	}
	limit := int(config.Get().ConstConfig.ChatMaxLen)
	if l > limit {
		a.GetLogger().Errorf("chat len limit, now:%v|limit:%v", l, limit)
		return nil, code.Code_ChatLenLimit
	}
	//允许发送给自己
	uuid := uuid.Generate()
	if req.GetTargetUid() != a.GetModel().Uid {
		_, cod := call.RequestRemote[proto.Message](req.GetTargetUid(), &inner.SendPrivateChatMsg_Request{
			SenderUid:  a.GetModel().Uid,
			SenderName: a.GetModel().Name,
			Uuid:       uuid,
			Content:    req.GetContent(),
			Url:        req.GetUrl(),
			MsgType:    req.GetMsgType(),
		})
		if cod != code.Code_Ok {
			a.GetLogger().Errorf("out send chat to target fail, target:%v|content:%v|cod:%v", req.GetTargetUid(), req.GetContent(), cod)
			return nil, cod
		}
	}
	//发送成功
	var chatMsg *outer.ChatMsg
	var err error
	if req.GetMsgType() == outer.ChatMsg_Pic {
		chatMsg, err = a.GetModel().AddPrivatePicChat(req.GetUrl(), a.GetModel().Uid, req.TargetUid, uuid, a.GetLogger())
	} else {
		chatMsg, err = a.GetModel().AddPrivateTextChat(req.GetContent(), a.GetModel().Uid, req.TargetUid, uuid, a.GetLogger())
	}
	if err != nil {
		a.GetLogger().Errorf("out add chat msg fail, target:%v|content:%v|err:%v", req.GetTargetUid(), req.GetContent(), err)
		return nil, code.Code_InnerError
	}
	//通知客户端
	call.Send(a.GetCtx(), a.GetCtx().Self(), outer.BuildPushPack(&outer.PushChatMsg_Push{MsgList: []*outer.ChatMsg{chatMsg}}))
	a.GetLogger().Infof("out send chat to target success, sender:%v|chatMsg:%v", req.GetTargetUid(), chatMsg)
	return &outer.SendPrivateChatMsg_Reply{}, code.Code_Ok
}
func (a *PlayerComponent) HandleInSendPrivateChatMsg(req *inner.SendPrivateChatMsg_Request) (*inner.SendPrivateChatMsg_Reply, code.Code) {
	var chatMsg *outer.ChatMsg
	var err error
	if req.GetMsgType() == outer.ChatMsg_Pic {
		chatMsg, err = a.GetModel().AddPrivatePicChat(req.GetUrl(), req.GetSenderUid(), a.GetModel().Uid, req.GetUuid(), a.GetLogger())
	} else {
		chatMsg, err = a.GetModel().AddPrivateTextChat(req.GetContent(), req.GetSenderUid(), a.GetModel().Uid, req.GetUuid(), a.GetLogger())
	}

	if err != nil {
		a.GetLogger().Errorf("in add chat msg fail, sender:%v|content:%v|err:%v", req.GetSenderUid(), req.GetContent(), err)
		return nil, code.Code_InnerError
	}
	//发送给client
	call.Send(a.GetCtx(), a.GetCtx().Self(), outer.BuildPushPack(&outer.PushChatMsg_Push{MsgList: []*outer.ChatMsg{chatMsg}}))
	//如果离线则增加推送
	if !a.GetModel().IsOnline() {
		withAtMsg := common_service.ParseChatMsgFormat(req.GetContent())
		//注意拼接的的中文的：
		pushContent := req.GetSenderName() + "：" + withAtMsg
		getui_helper.PushMsg([]uint64{a.GetModel().Uid}, config.Get().ConstConfig.PrivateChatTitle, pushContent)
	}
	//
	a.GetLogger().Infof("in send chat to target success, sender:%v|chatMsg:%v", req.GetSenderUid(), chatMsg)
	return &inner.SendPrivateChatMsg_Reply{}, code.Code_Ok
}

func (a *PlayerComponent) HandleGetPrivateChatMsg(req *outer.GetPrivateChatMsg_Request) (*outer.GetPrivateChatMsg_Reply, code.Code) {
	//获取
	chatMsgList, lastId := a.GetModel().GetChatMsgByRidBeginId(req.GetTargetUid(), req.GetBeginId(), a.GetLogger())
	//
	a.GetLogger().Infof("get private chat msg, range, rid:%v|beginId:%v|lastId:%v|rsp:%v", req.GetTargetUid(), req.GetBeginId(), lastId, chatMsgList)
	return &outer.GetPrivateChatMsg_Reply{MsgList: chatMsgList}, code.Code_Ok
}
