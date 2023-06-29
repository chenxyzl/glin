package player_group

import (
	"github.com/chenxyzl/glin/grain"
	"github.com/chenxyzl/glin/uuid"
	"google.golang.org/protobuf/proto"
	"laiya/common_service"
	"laiya/config"
	"laiya/home/iface"
	"laiya/model/home_model"
	"laiya/proto/code"
	"laiya/proto/inner"
	"laiya/proto/outer"
	"laiya/share/call"
)

func (a *GroupComponent) HandleGr2HoGroupDisband(req *inner.Gr2HoGroupDisband_Request) (*inner.Gr2HoGroupDisband_Reply, code.Code) {
	_, ok := a.GetModel().Groups[req.GetGid()]
	if !ok { //不存在直接ok
		return &inner.Gr2HoGroupDisband_Reply{}, code.Code_Ok
	}
	a.GetModel().RemoveGroup(req.GetGid())
	a.GetLogger().Infof("group disband notify, gid:%v", req.GetGid())
	return &inner.Gr2HoGroupDisband_Reply{}, code.Code_Ok
}

func (a *GroupComponent) HandleCreateGroup(req *outer.CreateGroup_Request) (*outer.CreateGroup_Reply, code.Code) {
	selfGroupCount := a.GetModel().SelfGroupCount()
	//非运营只能创建n个群
	if selfGroupCount >= int(config.Get().ConstConfig.MaxCreateGroupCount) && !config.Get().GmConfig.IsOpUid(a.GetModel().Uid) {
		a.GetLogger().Errorf("create group count full, now:%v", selfGroupCount)
		return nil, code.Code_CreateGroupCountFull
	}
	if len(a.GetModel().Groups) >= int(config.Get().ConstConfig.MaxGroupCount) {
		a.GetLogger().Errorf("all group count full, now:%v", len(a.GetModel().Groups))
		return nil, code.Code_GroupCountFull
	}
	if cod := common_service.CheckGroupName(req.GetName()); cod != code.Code_Ok {
		a.GetLogger().Errorf("check group name err, name:%v", req.GetName())
		return nil, cod
	}
	if cod := common_service.CheckGroupDes(req.GetDes()); cod != code.Code_Ok {
		a.GetLogger().Errorf("check group des err, des:%v", req.GetDes())
		return nil, cod
	}
	//
	gid := uuid.GenerateInterval()
	_, cod := call.RequestRemote[*inner.Ho2GrCreateGroup_Reply](gid, &inner.Ho2GrCreateGroup_Request{
		Owner: &inner.Ho2GrCreateGroup_Player{
			Uid:     a.GetModel().Uid,
			Session: a.GetCtx().Self(),
		},
		GroupName: req.GetName(),
		GroupIcon: req.GetIcon(),
		GroupDes:  req.GetDes(),
		GameType:  req.GetGameType(),
	})
	if cod != code.Code_Ok {
		a.GetLogger().Errorf("req create group err, cod:%v", cod)
		return nil, cod
	}
	//
	a.GetModel().AddGroup(&home_model.Group{
		Gid:    gid,
		IsSelf: true,
	})
	a.GetLogger().Infof("create group success, gid:%v", gid)
	return &outer.CreateGroup_Reply{Gid: gid}, code.Code_Ok
}

func (a *GroupComponent) HandleFavoriteGroup(req *outer.FavoriteGroup_Request) (*outer.FavoriteGroup_Reply, code.Code) {
	if len(a.GetModel().Groups) >= int(config.Get().ConstConfig.MaxGroupCount) {
		a.GetLogger().Errorf("group count full, now:%v", len(a.GetModel().Groups))
		return nil, code.Code_GroupCountFull
	}
	ret, cod := call.RequestRemote[*inner.Ho2GrFavoriteGroup_Reply](req.GetGid(), &inner.Ho2GrFavoriteGroup_Request{Uid: a.GetModel().Uid})
	if cod != code.Code_Ok {
		a.GetLogger().Errorf("req favorite group err, cod:%v", cod)
		return nil, cod
	}
	a.GetModel().AddGroup(&home_model.Group{
		Gid:    req.GetGid(),
		IsSelf: ret.TransferOwner,
	})
	a.GetLogger().Infof("favorite group success, gid:%v|transferOwner:%v", req.GetGid(), ret.TransferOwner)
	return &outer.FavoriteGroup_Reply{}, code.Code_Ok
}

func (a *GroupComponent) HandleUnfavoriteGroup(req *outer.UnfavoriteGroup_Request) (*outer.UnfavoriteGroup_Reply, code.Code) {
	group, ok := a.GetModel().Groups[req.GetGid()]
	if !ok { //不在群中直接ok
		return &outer.UnfavoriteGroup_Reply{}, code.Code_Ok
	}
	if group.IsSelf {
		a.GetLogger().Errorf("self group can not unfavorite, cid:%v", req.GetGid())
		return nil, code.Code_SelfGroup
	}
	_, cod := call.RequestRemote[*inner.Ho2GrUnfavoriteGroup_Reply](req.GetGid(), &inner.Ho2GrUnfavoriteGroup_Request{
		Uid: a.GetModel().Uid,
	})
	if cod != code.Code_Ok && cod != code.Code_GroupNotExist {
		a.GetLogger().Errorf("req unfavorite group err, cod:%v", cod)
		return nil, cod
	}
	a.GetModel().RemoveGroup(req.GetGid())
	a.GetLogger().Infof("player unfavorite group, gid:%v", req.GetGid())
	return &outer.UnfavoriteGroup_Reply{}, code.Code_Ok
}

func (a *GroupComponent) HandleDisbandGroup(req *outer.DisbandGroup_Request) (*outer.DisbandGroup_Reply, code.Code) {
	group, ok := a.GetModel().Groups[req.GetGid()]
	if !ok { //不存在直接ok
		return &outer.DisbandGroup_Reply{}, code.Code_Ok
	}
	if !group.IsSelf {
		a.GetLogger().Errorf("only self group can disband, cid:%v", req.GetGid())
		return nil, code.Code_NotSelfGroup
	}
	//
	_, cod := call.RequestRemote[*inner.Ho2GrDisbandGroup_Reply](req.GetGid(), &inner.Ho2GrDisbandGroup_Request{
		Uid: a.GetModel().Uid,
	})
	if cod != code.Code_Ok && cod != code.Code_GroupNotExist {
		a.GetLogger().Errorf("request group disband err, cod:%v", cod)
		return nil, cod
	}
	//
	a.GetModel().RemoveGroup(req.GetGid())
	//
	a.GetLogger().Infof("player disband group, gid:%v", req.GetGid())
	return &outer.DisbandGroup_Reply{}, code.Code_Ok
}

func (a *GroupComponent) HandleEnterGroup(req *outer.EnterGroup_Request) (*outer.EnterGroup_Reply, code.Code) {
	//注: 可以不收藏 直接进入
	//if _, ok := a.GetModel().Groups[req.GetGid()]; !ok {
	//	a.GetLogger().Errorf("not in group, uid:%v", req.GetGid())
	//	return nil, code.Code_NotInGroup
	//}
	//注:可以先进入群,但不收藏
	//退出老的群
	if a.GetModel().NowInGroup > 0 {
		_, cod := call.RequestRemote[proto.Message](a.GetModel().NowInGroup, &inner.Ho2GrExitGroup_Request{Uid: a.GetModel().Uid})
		if cod == code.Code_GroupNotExist {
			//注: 不存在不处理
			a.GetModel().RemoveGroup(a.GetModel().NowInGroup)
			//return nil, cod
		} else if cod != code.Code_Ok {
			a.GetLogger().Errorf("exit group err, gid:%v|cod:%v", a.GetModel().NowInGroup, cod)
			return nil, cod
		}
		//log
		gid := a.GetModel().NowInGroup
		//reset
		a.GetModel().NowInGroup = 0
		a.GetLogger().Infof("exit old group, gid:%v", gid)
	}
	//
	_, cod := call.RequestRemote[*inner.Ho2GrEnterGroup_Reply](req.GetGid(), &inner.Ho2GrEnterGroup_Request{
		Session: a.GetCtx().Self(),
		Uid:     a.GetModel().Uid,
	})
	if cod == code.Code_GroupNotExist {
		a.GetModel().RemoveGroup(a.GetModel().NowInGroup)
		return nil, cod
	} else if cod != code.Code_Ok {
		a.GetLogger().Errorf("enter group err, gid:%v|cod:%v", req.GetGid(), cod)
		return nil, cod
	}
	//
	a.GetModel().NowInGroup = req.GetGid()
	a.GetLogger().Infof("enter group success, gid:%v", req.GetGid())
	return &outer.EnterGroup_Reply{}, code.Code_Ok
}

func (a *GroupComponent) HandleExitGroup(req *outer.ExitGroup_Request) (*outer.ExitGroup_Reply, code.Code) {
	if a.GetModel().NowInGroup <= 0 {
		return &outer.ExitGroup_Reply{}, code.Code_Ok
	}
	_, cod := call.RequestRemote[proto.Message](a.GetModel().NowInGroup, &inner.Ho2GrExitGroup_Request{Uid: a.GetModel().Uid})
	if cod == code.Code_GroupNotExist {
		a.GetModel().RemoveGroup(a.GetModel().NowInGroup)
		//return nil, cod
	} else if cod != code.Code_Ok {
		a.GetLogger().Errorf("exit group err, gid:%v|cod:%v", a.GetModel().NowInGroup, cod)
		return nil, cod
	}
	//log
	gid := a.GetModel().NowInGroup
	//reset
	a.GetModel().NowInGroup = 0
	a.GetLogger().Infof("exit group success, gid:%v", gid)
	return &outer.ExitGroup_Reply{}, code.Code_Ok
}
func (a *GroupComponent) HandleGroupActiveEvent(req *inner.GroupActiveEvent_Notify) {
	if a.GetModel().IsOnline() {
		grain.GetComponent[iface.IPlayerComponent](a).TellGroupOnline(req.GetGid())
		a.GetLogger().Infof("group active, tell group player online, gid:%v", req.GetGid())
	}
}
