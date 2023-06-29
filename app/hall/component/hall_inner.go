package component

import (
	"laiya/config"
	"laiya/model/hall_model"
	"laiya/proto/code"
	"laiya/proto/inner"
)

func (a *HallComponent) HandleGr2HaUpdateGroup(req *inner.Gr2HaUpdateGroup_Request) (*inner.Gr2HaUpdateGroup_Reply, code.Code) {
	if req.GetGid() == 0 {
		a.GetLogger().Warnf("update group warn, gid is 0")
		return nil, code.Code_Ok
	}
	group := a.GetModel().UpdateGroupInfo(req)
	a.GetLogger().Infof("update group success, gid:%v|onMicCount:%v", req.GetGid(), group.GetPlayerCount())
	return &inner.Gr2HaUpdateGroup_Reply{}, code.Code_Ok
}

func (a *HallComponent) HandleGr2HaDeleteGroup(req *inner.Gr2HaDeleteGroup_Request) (*inner.Gr2HaDeleteGroup_Reply, code.Code) {
	a.GetModel().DeleteGroup(req.GetGid())
	a.GetLogger().Infof("delete group success, gid:%v", req.GetGid())
	return &inner.Gr2HaDeleteGroup_Reply{}, code.Code_Ok
}

func (a *HallComponent) HandleCheckInviteBroadcast(req *inner.Ho2HaCheckInviteBroadcast_Request) (*inner.Ho2HaCheckInviteBroadcast_Reply, code.Code) {
	var groups []*inner.Ho2HaCheckInviteBroadcast_MatchGroup
	conf := config.Get().ParamsConfig
	gameConfig := config.Get().GameTypesConfig.GetParams(req.GetGameType())
	a.GetModel().SortGroupList.ForRange(func(v *hall_model.Group) bool {
		//够了不查 //避免过大
		if len(groups) >= conf.CheckInviteBroadcastLookingGroupCount {
			return false //终止
		}
		//找到空房间不用再查了
		if v.GetPlayerCount() == 0 {
			return false //终止
		}
		//游戏类型不匹配不查
		if v.GameType != req.GetGameType() {
			return true
		}
		//人数查询
		if v.GetPlayerCount() >= gameConfig.FullPlayerCount {
			return true
		}
		//可用
		groups = append(groups, &inner.Ho2HaCheckInviteBroadcast_MatchGroup{
			Gid:                  v.Gid,
			PlayerCount:          v.GetPlayerCount(),
			TempVoiceRoomVersion: v.TempVoiceRoomVersion,
		})
		//
		return true
	}, true)
	a.GetLogger().Infof("get invite broadcast groups success, gameType:%v|count:%v|groups:%v", req.GetGameType(), len(groups), groups)
	return &inner.Ho2HaCheckInviteBroadcast_Reply{MatchGroups: groups}, code.Code_Ok
}
