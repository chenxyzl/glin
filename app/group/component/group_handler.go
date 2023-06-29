package component

import (
	"google.golang.org/protobuf/proto"
	"laiya/common_service"
	"laiya/config"
	"laiya/model/group_model"
	"laiya/model/sub_model"
	"laiya/proto/code"
	"laiya/proto/common"
	"laiya/proto/inner"
	"laiya/proto/outer"
	"laiya/share/call"
	"laiya/share/datarangers"
	"laiya/share/global"
	"slices"
	"time"
)

func (a *GroupComponent) HandleCS30000_GetGroupCard(req *outer.GetGroupCard_Request) (*outer.GetGroupCard_Reply, code.Code) {
	ret := &outer.GetGroupCard_Reply{
		GroupCard: &outer.GetGroupCard_GroupCard{
			Gid:               a.GetModel().Gid,
			OwnerUid:          a.GetModel().OwnerUid,
			Name:              a.GetModel().Name,
			Icon:              a.GetModel().Icon,
			Des:               a.GetModel().Des,
			GameType:          a.GetModel().GameType,
			OnMicPlayers:      a.GetModel().OnMicPlayers,
			ScreenSharingUids: a.GetModel().OnScreenSharingPlayers,
			IsOfficial:        config.Get().GmConfig.IsOfficialGroup(a.GetModel().Gid),
		},
	}
	a.GetLogger().Infof("get group card success")
	return ret, code.Code_Ok
}

func (a *GroupComponent) HandleCS30014_GetShareGroupCard(req *outer.GetShareGroupCard_Request) (*outer.GetShareGroupCard_Reply, code.Code) {
	ret := &outer.GetShareGroupCard_Reply{
		GroupCard: &outer.GetShareGroupCard_GroupCard{
			Gid:         a.GetModel().Gid,
			OwnerUid:    a.GetModel().OwnerUid,
			Name:        a.GetModel().Name,
			Icon:        a.GetModel().Icon,
			Des:         a.GetModel().Des,
			GameType:    a.GetModel().GameType,
			OnlineCount: uint32(a.GetModel().GetOnlineCount(true)),
			PlayerCount: uint32(a.GetModel().PlayerCount()),
			Players:     []*outer.GetShareGroupCard_GroupCard_Player{},
			IsOfficial:  config.Get().GmConfig.IsOfficialGroup(a.GetModel().Gid),
			JoinAccess:  a.GetModel().JoinAccess,
		},
	}
	var onlineUids []uint64
	const needCount = 4
	a.GetModel().OnlineSort.ForRange(func(v *group_model.OnlineSortT) bool {
		//最多获取n个
		if len(onlineUids) >= needCount {
			return false
		}
		//只返回在群内的
		if !a.GetModel().IsPlayerInGroup(v.Key()) {
			return true
		}
		onlineUids = append(onlineUids, v.Key())
		return true
	})
	//
	for _, uid := range a.GetModel().GetPlayers() {
		if len(onlineUids) >= needCount {
			break
		}
		//check contain
		if slices.Contains(onlineUids, uid) {
			continue
		}
		onlineUids = append(onlineUids, uid)
	}
	//
	players, err := sub_model.BatchLoadPlayerHead(onlineUids, a.GetLogger())
	if err != nil {
		a.GetLogger().Errorf("load players err, uids:%v", onlineUids)
		return nil, code.Code_InnerError
	}
	for _, player := range players {
		ret.GetGroupCard().Players = append(ret.GetGroupCard().Players, &outer.GetShareGroupCard_GroupCard_Player{
			Uid:  player.Uid,
			Name: player.Name,
			Icon: player.Icon,
			Des:  player.Des,
		})
	}
	a.GetLogger().Infof("get share group card success")
	return ret, code.Code_Ok
}

func (a *GroupComponent) HandleGetGroupDetail(req *outer.GetGroupDetail_Request) (*outer.GetGroupDetail_Reply, code.Code) {
	//
	var onlinePlayers []uint64
	var offlinePlayers []uint64
	var permissions = make(map[uint64]common.GroupPermissionsDef_Type)
	for _, v := range a.GetModel().Players {
		if a.GetModel().IsPlayerOnline(v.Uid) {
			onlinePlayers = append(onlinePlayers, v.Uid)
		} else {
			offlinePlayers = append(offlinePlayers, v.Uid)
		}
		//非普通用户返回权限返回
		p := a.GetModel().GetUserPermissions(v.Uid)
		if p != common.GroupPermissionsDef_UnknownType && p != common.GroupPermissionsDef_User {
			permissions[v.Uid] = p
		}
	}
	copywritingConfig := config.Get().GameCopywritingConfig.Get(a.GetModel().GameType)
	guideCopywriting := &outer.GetGroupDetail_GuideCopyWriting{
		OnMicPlayer0Guide: copywritingConfig.OnMicPlayer0Guide,
		OnMicPlayer1Guide: copywritingConfig.OnMicPlayer1Guide,
	}
	rsp := &outer.GetGroupDetail_Reply{Group: &outer.GetGroupDetail_Group{
		Gid:               a.GetModel().Gid,
		OwnerUid:          a.GetModel().OwnerUid,
		Name:              a.GetModel().Name,
		Icon:              a.GetModel().Icon,
		Des:               a.GetModel().Des,
		GameType:          a.GetModel().GameType,
		Players:           append(onlinePlayers, offlinePlayers...),
		OnlineCount:       uint32(len(onlinePlayers)),
		OnMicPlayers:      a.GetModel().OnMicPlayers,
		ScreenSharingUids: a.GetModel().OnScreenSharingPlayers,
		IsOfficial:        config.Get().GmConfig.IsOfficialGroup(a.GetModel().Gid),
		JoinAccess:        a.GetModel().JoinAccess,
		GuideCopyWriting:  guideCopywriting,
		VoiceRoomName:     a.GetModel().GetVoiceRoomName(),
		Permissions:       permissions,
	}}

	a.GetLogger().Infof("get group detail, rsp:%v", rsp)
	return rsp, code.Code_Ok

}

func (a *GroupComponent) HandleGetGroupUserInfo(req *outer.GetGroupUserInfo_Request) (*outer.GetGroupUserInfo_Reply, code.Code) {
	//
	a.GetLogger().Infof("get group user info, uids:%v", req.GetUids())
	//
	uids := req.GetUids()
	players, err := sub_model.BatchLoadPlayerHead(uids, a.GetLogger())
	if err != nil {
		a.GetLogger().Infof("batch load player from db err, uids:%v|err:%v", uids, err)
		return nil, code.Code_Error
	}
	//
	var pbPlayers []*outer.GetGroupUserInfo_Player
	for _, player := range players {
		pbPlayers = append(pbPlayers, &outer.GetGroupUserInfo_Player{
			Uid:      player.Uid,
			Name:     player.Name,
			Icon:     player.Icon,
			Des:      player.Des,
			IsOnline: a.GetModel().IsPlayerOnline(player.Uid),
			UserType: outer.GetUserType(player.Uid),
			//Permissions: a.GetModel().GetUserPermissions(player.Uid),
		})
	}
	return &outer.GetGroupUserInfo_Reply{Players: pbPlayers}, code.Code_Ok
}

func (a *GroupComponent) HandleHo2GrCreateGroup(req *inner.Ho2GrCreateGroup_Request) (*inner.Ho2GrCreateGroup_Reply, code.Code) {
	//
	ownerInfo := sub_model.PlayerHead{Uid: req.GetOwner().GetUid()}
	err := ownerInfo.Load()
	if err != nil {
		a.GetLogger().Errorf("load owner name and icon err, err:%v", err)
		return nil, code.Code_InnerError
	}
	group := &group_model.Group{
		Gid:        a.GetModel().Gid,
		CreateTime: time.Now().UnixMilli(),
		OwnerUid:   req.GetOwner().GetUid(),
		Name:       req.GetGroupName(),
		Icon:       req.GetGroupIcon(),
		Des:        req.GetGroupDes(),
		GameType:   req.GetGameType(),
		Room:       make(map[uint64]*group_model.Room),
		Players:    make(map[uint64]*group_model.Player),
	}
	//版本升级检查
	group.CheckVersion(a.GetLogger())
	//
	group.Players[req.GetOwner().GetUid()] = &group_model.Player{
		Uid: req.GetOwner().GetUid(),
	}
	room, cod := group.CreateRoom(config.Get().ConstConfig.DefaultVoiceRoomName, "")
	if cod != code.Code_Ok {
		a.GetLogger().Errorf("create room err, cod:%v", cod)
		return nil, cod
	}
	group.Room[room.RoomId] = room
	//立即保存
	err = group.Save()
	//
	if err != nil {
		a.GetLogger().Errorf("create group err, err:%v", err)
		return nil, code.Code_InnerError
	}
	//
	a.SetModel(group, a.GetLogger())
	//设置为在线--必须要再SetModel之后
	group.Online(req.GetOwner().GetUid(), req.GetOwner().GetSession())
	//
	a.UpdateToHall()
	//
	a.NotifyChatMsg(&outer.ChatMsg_CreateGroupMsg{UserName: ownerInfo.Name}, req.GetOwner().GetUid())
	//
	a.GetLogger().Infof("create group success")
	return &inner.Ho2GrCreateGroup_Reply{}, code.Code_Ok
}

func (a *GroupComponent) HandleUpdateGroupInfo(req *outer.UpdateGroupInfo_Request) (*outer.UpdateGroupInfo_Reply, code.Code) {
	if len(req.GetItems()) == 0 {
		a.GetLogger().Infof("change nothings")
		return &outer.UpdateGroupInfo_Reply{}, code.Code_Ok
	}
	for _, item := range req.GetItems() {
		if cod := common_service.CheckGroupInfo(item.GetType(), item.GetStrVar()); cod != code.Code_Ok {
			a.GetLogger().Errorf("check update group info err, type:%v|val:%v", item.GetType(), item.GetStrVar())
			return nil, cod
		}
	}

	for _, item := range req.GetItems() {
		switch item.GetType() {
		case outer.UpdateGroupInfo_Name:
			a.GetModel().SetGroupName(item.GetStrVar())
		case outer.UpdateGroupInfo_Icon:
			a.GetModel().SetGroupIcon(item.GetStrVar())
		case outer.UpdateGroupInfo_Des:
			a.GetModel().SetGroupDes(item.GetStrVar())
		case outer.UpdateGroupInfo_GameType:
			a.GetModel().SetGroupGameType(item.GetStrVar())
			//检查是否要添加对印的机器人
			a.CheckAddRobot()
		default:
			a.GetLogger().Errorf("not found enum impl, val:%v", item.GetType())
			return nil, code.Code_EnumNotFoundImpl
		}
	}
	//更新到大厅--因为搜索需要的关键字改了
	a.UpdateToHall()
	//
	a.GetLogger().Infof("update group info success, val:%v", req)
	return &outer.UpdateGroupInfo_Reply{}, code.Code_Ok
}

func (a *GroupComponent) HandleHo2GrFavoriteGroup(req *inner.Ho2GrFavoriteGroup_Request) (*inner.Ho2GrFavoriteGroup_Reply, code.Code) {
	if !a.GetModel().JoinAccess {
		a.GetLogger().Errorf("player join group access denied,uid:%v", req.GetUid())
		return nil, code.Code_JoinAccessDenied
	}
	//如果是机器人房间则转移拥有者转首个进入的
	if _, ok := a.GetModel().Players[req.GetUid()]; ok {
		return &inner.Ho2GrFavoriteGroup_Reply{TransferOwner: a.GetModel().OwnerUid == req.GetUid()}, code.Code_Ok
	}
	if len(a.GetModel().Players) >= int(config.Get().ConstConfig.GroupPlayerMaxCount) {
		a.GetLogger().Errorf("group player full,now:%v", len(a.GetModel().Players))
		return nil, code.Code_GroupPlayerFull
	}
	senderInfo := sub_model.PlayerHead{Uid: req.GetUid()}
	err := senderInfo.Load()
	if err != nil {
		a.GetLogger().Errorf("load sender name and icon err, err:%v", err)
		return nil, code.Code_InnerError
	}
	//
	beforeOwner := a.GetModel().OwnerUid
	transferOwner := global.IsRobot(a.GetModel().OwnerUid)
	//是否是机器人创建 交换所有者
	if transferOwner {
		a.GetModel().OwnerUid = req.GetUid()
	}
	//save
	player := &group_model.Player{Uid: req.GetUid()}
	a.GetModel().AddPlayer(player)
	//如果机器人创建 交换了所有者-更新到大厅
	if transferOwner {
		a.GetLogger().Infof("robot group change to player, robot:%v|uid:%v", beforeOwner, player.Uid)
		//此时下一帧更新到大厅
		a.UpdateToHall()
	}
	//通知有人收藏
	a.NotifyGroupMsg(&outer.PushGroupPlayerFavorite_Push{Uids: []uint64{req.GetUid()}}, true)
	//通知聊天消息
	a.NotifyChatMsg(&outer.ChatMsg_FavoriteGroupMsg{UserName: senderInfo.Name}, req.GetUid())
	//
	a.GetLogger().Infof("player favorite group, uid:%v", req.GetUid())
	return &inner.Ho2GrFavoriteGroup_Reply{TransferOwner: transferOwner}, code.Code_Ok
}

func (a *GroupComponent) HandleHo2GrUnfavoriteGroup(req *inner.Ho2GrUnfavoriteGroup_Request) (*inner.Ho2GrUnfavoriteGroup_Reply, code.Code) {
	uid := req.GetUid()
	if _, ok := a.GetModel().Players[uid]; !ok {
		return &inner.Ho2GrUnfavoriteGroup_Reply{}, code.Code_Ok
	}
	if a.GetModel().OwnerUid == uid {
		a.GetLogger().Errorf("self group, can not unfavorite")
		return nil, code.Code_SelfGroup
	}
	//再从群离开
	a.GetModel().RemovePlayer(uid)
	//通知有人取消收藏
	a.NotifyGroupMsg(&outer.PushGroupPlayerUnfavorite_Push{Uids: []uint64{req.GetUid()}}, true)
	//
	a.GetLogger().Infof("player unfavorite group, uid:%v", req.GetUid())
	//
	return &inner.Ho2GrUnfavoriteGroup_Reply{}, code.Code_Ok
}

func (a *GroupComponent) HandleHo2GrDisbandGroup(req *inner.Ho2GrDisbandGroup_Request) (*inner.Ho2GrDisbandGroup_Reply, code.Code) {
	//从大厅移除
	_, cod := call.RequestRemote[proto.Message](global.HallUID, &inner.Gr2HaDeleteGroup_Request{Gid: a.GetModel().Gid})
	if cod != code.Code_Ok {
		a.GetLogger().Errorf("disband group, hall delete group err, cod:%v", cod)
		return nil, cod
	}
	//广播解散
	a.NotifyGroupMsg(&outer.PushGroupDisband_Push{Gid: a.GetModel().Gid}, true)
	//通知解散
	for _, player := range a.GetModel().Players {
		a.Next(func(params ...any) {
			uid := params[0].(uint64)
			_, cod := call.RequestRemote[*inner.Gr2HoGroupDisband_Reply](uid, &inner.Gr2HoGroupDisband_Request{Gid: a.GetModel().Gid})
			if cod != code.Code_Ok {
				a.GetLogger().Errorf("notify group disband err, uid:%v|cod:%v", uid, cod)
			}
		}, player.Uid)
	}
	//标记已解散
	a.Next(func(...any) {
		a.GetLogger().Infof("group mark disband begin")
		//只需要设置已解散标记即可
		err := a.GetModel().Delete()
		if err != nil {
			a.GetLogger().Errorf("group mark disband save err, err:%v", err)
		}
		a.GetLogger().Infof("group mark disband end")
		//
		a.GetCtx().Poison(a.GetCtx().Self())
	})
	a.GetLogger().Infof("group disband, gid:%v", a.GetModel().Gid)

	return &inner.Ho2GrDisbandGroup_Reply{}, code.Code_Ok
}

func (a *GroupComponent) HandleHo2GrEnterGroup(req *inner.Ho2GrEnterGroup_Request) (*inner.Ho2GrEnterGroup_Reply, code.Code) {
	//注: 可以不收藏,直接进入
	//player, ok := a.GetModel().Players[req.GetPlayer().GetUid()]
	//if !ok {
	//	a.GetLogger().Errorf("not in group, uid:%v", req.GetPlayer().GetUid())
	//	return nil, code.Code_NotInGroup
	//}
	//
	if a.GetModel().OnlineSort.Length() > config.Get().ConstConfig.GroupOnlineLimit {
		return nil, code.Code_GroupOnlineFull
	}
	a.GetModel().Enter(req.GetUid(), req.GetSession())
	//通知有人进入房间
	a.NotifyGroupMsg(&outer.PushGroupPlayerEnter_Push{Uids: []uint64{req.GetUid()}}, true)
	//
	a.GetLogger().Infof("player enter group, uid:%v", req.GetUid())
	return &inner.Ho2GrEnterGroup_Reply{}, code.Code_Ok
}

func (a *GroupComponent) HandleHo2GrExitGroup(req *inner.Ho2GrExitGroup_Request) (*inner.Ho2GrExitGroup_Reply, code.Code) {
	//注: 可以不用收藏
	//player, ok := a.GetModel().Players[req.GetUid()]
	//if !ok {
	//	a.GetLogger().Errorf("not in group, uid:%v", req.GetUid())
	//	return nil, code.Code_NotInGroup
	//}
	//
	a.GetModel().Exit(req.GetUid())
	//通知有人离开房间
	a.NotifyGroupMsg(&outer.PushGroupPlayerExit_Push{Uids: []uint64{req.GetUid()}}, true)
	//
	a.GetLogger().Infof("player exit group, uid:%v", req.GetUid())
	return &inner.Ho2GrExitGroup_Reply{}, code.Code_Ok
}

func (a *GroupComponent) HandleSendGroupChatMsg(req *outer.SendGroupChatMsg_Request) (*outer.SendGroupChatMsg_Reply, code.Code) {
	//在群检查
	_, ok := a.GetModel().Players[req.GetUid()]
	if !ok {
		a.GetLogger().Errorf("not in group, uid:%v", req.GetUid())
		return nil, code.Code_NotInGroup
	}
	var l int
	if req.GetMsgType() == outer.ChatMsg_Pic {
		l = len([]rune(req.GetUrl()))
	} else {
		l = len([]rune(req.GetContent()))
	}
	//长度检查
	if l <= 0 {
		a.GetLogger().Errorf("not send empty msg", req.GetUid())
		return nil, code.Code_Ok
	}
	limit := int(config.Get().ConstConfig.ChatMaxLen)
	if l > limit {
		a.GetLogger().Errorf("chat len limit, now:%v|limit:%v", l, limit)
		return nil, code.Code_ChatLenLimit
	}
	//
	senderInfo := sub_model.PlayerHead{Uid: req.GetUid()}
	err := senderInfo.Load()
	if err != nil {
		a.GetLogger().Errorf("load sender name and icon err, err:%v", err)
		return nil, code.Code_InnerError
	}
	datarangers.SendEvent("group_players_im", map[string]interface{}{
		"group_id":   a.GetModel().Gid,
		"group_name": a.GetModel().Name,
		"sender":     req.GetUid(),
	}, nil)
	//广播群消息
	if req.GetMsgType() == outer.ChatMsg_Pic {
		a.NotifyChatMsg(&outer.ChatMsg_PicMsg{Url: req.GetUrl(), SenderName: senderInfo.Name, SenderIcon: senderInfo.Icon}, req.GetUid())
	} else {
		a.NotifyChatMsg(&outer.ChatMsg_NormalMsg{Text: req.GetContent(), SenderName: senderInfo.Name, SenderIcon: senderInfo.Icon}, req.GetUid())
	}
	//
	a.GetLogger().Infof("send chat, chat:%v", req)
	return &outer.SendGroupChatMsg_Reply{}, code.Code_Ok
}

func (a *GroupComponent) HandleGetGroupChatMsg(req *outer.GetGroupChatMsg_Request) (*outer.GetGroupChatMsg_Reply, code.Code) {
	//获取
	chatMsgList, lastId := a.GetModel().GetChatMsgByBeginId(req.GetBeginId(), a.GetLogger())
	//
	a.GetLogger().Infof("get group chat msg, range, beginId:%v|lastId:%v|rsp:%v", req.GetBeginId(), lastId, chatMsgList)
	return &outer.GetGroupChatMsg_Reply{MsgList: chatMsgList}, code.Code_Ok
}

func (a *GroupComponent) HandleUpdateGroupJoinAccess(req *outer.UpdateGroupJoinAccess_Request) (*outer.UpdateGroupJoinAccess_Reply, code.Code) {
	//有权限变更
	if req.GetJoinAccess() == a.GetModel().JoinAccess {
		a.GetLogger().Infof("update group join access equal,req:%v", req)
		return &outer.UpdateGroupJoinAccess_Reply{}, code.Code_Ok
	}
	//检查权限种类
	if req.GetJoinAccess() { //放开权限
		//更新到大厅
		a.UpdateToHall()
	} else { //禁止权限
		//从大厅移除
		_, cod := call.RequestRemote[proto.Message](global.HallUID, &inner.Gr2HaDeleteGroup_Request{Gid: a.GetModel().Gid})
		if cod != code.Code_Ok {
			a.GetLogger().Errorf("join access disable, hall delete group err, cod:%v", cod)
			return nil, cod
		}
	}
	//
	a.GetModel().SetJoinAccess(req.GetJoinAccess())
	//
	a.GetLogger().Infof("update group join access,req:%v", req)
	return &outer.UpdateGroupJoinAccess_Reply{}, code.Code_Ok
}

func (a *GroupComponent) UpdateToHall() {
	//下一帧再处理
	a.GetModel().NeedUpdateToHall = true
	a.Next(func(...any) {
		//检查是否能更新到大厅
		if global.IsRobot(a.GetModel().OwnerUid) || //机器人群不更新到大厅
			!a.GetModel().JoinAccess || //没有加入权限不能更新到大厅
			a.GetModel().IsDisband || //已解散不能更新到大厅
			!a.GetModel().NeedUpdateToHall { //已更新过
			return
		}
		a.GetModel().NeedUpdateToHall = false
		sortScore := a.GetModel().GetHallSortScore()
		_, cod := call.RequestRemote[proto.Message](global.HallUID, &inner.Gr2HaUpdateGroup_Request{
			Gid:                  a.GetModel().Gid,
			Name:                 a.GetModel().Name,
			VoiceRoomName:        a.GetModel().GetVoiceRoomName(),
			GameType:             a.GetModel().GameType,
			SortScore:            sortScore,
			TempVoiceRoomVersion: a.GetModel().TempVoiceRoomVersion,
		})
		if cod != code.Code_Ok {
			a.GetLogger().Errorf("notify hall update score err,hall:%v|sortScore:%v|playerCount:%v|cod:%v", global.HallUID, sortScore, a.GetModel().GetOnMicPlayerCount(), cod)
		} else {
			a.GetLogger().Infof("notify hall update score success,hall:%v|sortScore:%v|playerCount:%v", global.HallUID, sortScore, a.GetModel().GetOnMicPlayerCount())
		}
	})
}

func (a *GroupComponent) HandleSetGroupVoiceName(req *outer.SetGroupVoiceName_Request) (*outer.SetGroupVoiceName_Reply, code.Code) {
	if req.GetVoiceRoomName() == "" {
		a.GetLogger().Infof("req set group name empty, uid:%v|name:%v", req.GetUid(), req.GetVoiceRoomName())
		return &outer.SetGroupVoiceName_Reply{}, code.Code_Ok
	}
	if req.GetVoiceRoomName() == a.GetModel().GetVoiceRoomName() {
		a.GetLogger().Infof("req set group name equal, uid:%v|name:%v", req.GetUid(), req.GetVoiceRoomName())
		return &outer.SetGroupVoiceName_Reply{}, code.Code_Ok
	}
	cod := common_service.CheckVoiceRoomName(req.GetVoiceRoomName())
	if cod != code.Code_Ok {
		a.GetLogger().Errorf("check voice room name err, cod:%v,name:%v", cod, req.GetVoiceRoomName())
		return nil, cod
	}
	//
	a.GetModel().SetVoiceRoomName(req.GetVoiceRoomName(), true)
	//更新到大厅
	a.UpdateToHall()
	//
	a.NotifyGroupMsg(&outer.PushGroupVoiceNameChange_Push{
		Gid:           a.GetModel().Gid,
		VoiceRoomName: req.GetVoiceRoomName(),
	})
	//通知机器人
	robot := a.GetCurrentRobot()
	if robot == nil {
		a.GetLogger().Warnf("trigger invite broadcast, but not found active robot, gameType:%v", a.GetModel().GameType)
	} else {
		call.Send(a.GetCtx(), robot, &inner.Gr2RoNeedInviteBroadcast_Notify{
			Uid:            req.GetUid(),
			GroupName:      a.GetModel().Name,
			GroupVoiceName: a.GetModel().GetVoiceRoomName(),
			GameType:       a.GetModel().GameType,
		})
	}
	//
	a.GetLogger().Infof("req set group name equal, uid:%v|name:%v", req.GetUid(), req.GetVoiceRoomName())
	return &outer.SetGroupVoiceName_Reply{}, code.Code_Ok
}

func (a *GroupComponent) HandleChangeGroupUserPermissions(req *outer.ChangeGroupUserPermissions_Request) (*outer.ChangeGroupUserPermissions_Reply, code.Code) {
	//权限检查
	if !a.GetModel().CanPermissionsConfig(req.GetUid()) {
		a.GetLogger().Errorf("change user permissions err, no access, targetUid:%v|permissions:%v", req.GetTargetUid(), req.GetPermissions())
		return nil, code.Code_NoAccess
	}
	//改变
	cod := a.GetModel().SetUserPermissions(req.GetTargetUid(), req.GetPermissions())
	if cod != code.Code_Ok {
		a.GetLogger().Errorf("change user permissions cod err, targetUid:%v|permissions:%v|cod:%v", req.GetTargetUid(), req.GetPermissions(), cod)
		return nil, code.Code_NoAccess
	}
	//广播
	a.NotifyGroupMsg(
		&outer.PushGroupUserPermissionsChange_Push{
			Gid:         a.GetModel().Gid,
			OpUid:       req.GetUid(),
			TargetUid:   req.GetTargetUid(),
			Permissions: req.GetPermissions(),
		},
		true)
	//
	a.GetLogger().Infof("change user permissions success, targetUid:%v|permissions:%v", req.GetTargetUid(), req.GetPermissions())
	return &outer.ChangeGroupUserPermissions_Reply{}, code.Code_Ok
}
