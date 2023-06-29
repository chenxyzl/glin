package component

import (
	"github.com/chenxyzl/glin/uuid"
	"golang.org/x/exp/slices"
	"laiya/model/plugin_model/group_plugin_activity_model"
	"laiya/model/sub_model"
	"laiya/proto/code"
	"laiya/proto/inner"
	"laiya/proto/outer"
	"laiya/share/call"
)

func (a *ActivityComponent) HandleGetGroupAllActivity(req *outer.GetGroupAllActivity_Request) (*outer.GetGroupAllActivity_Reply, code.Code) {
	var arr []*outer.GroupActivitySummary
	for _, activity := range a.GetModel().Activities {
		if req.GetOnlySelf() && //如果要获取自己的则需要判断是否是(自己或已加入)
			(req.GetUid() != activity.OwnerUid || !slices.ContainsFunc(activity.Players, func(player *group_plugin_activity_model.Player) bool { return req.GetUid() == player.Uid })) {
			continue
		}
		arr = append(arr, activity.ToSummary())
	}
	a.GetLogger().Infof("get all activity success")
	return &outer.GetGroupAllActivity_Reply{Activities: arr}, code.Code_Ok
}
func (a *ActivityComponent) HandleGetGroupActivityDetail(req *outer.GetGroupActivityDetail_Request) (*outer.GetGroupActivityDetail_Reply, code.Code) {
	idx := slices.IndexFunc(a.GetModel().Activities, func(item *group_plugin_activity_model.ActivityItem) bool { return req.GetActivityId() == item.Aid })
	if idx < 0 {
		a.GetLogger().Errorf("not found activity by aid, aid:%v", req.GetActivityId())
		return nil, code.Code_GroupActivityNotExist
	}
	activity := a.GetModel().Activities[idx]
	ret := &outer.GetGroupActivityDetail_GroupActivityDetail{
		ActivityId:    activity.Aid,
		ActivityTitle: activity.Title,
		ActivityDes:   activity.Des,
		BeginTime:     activity.BeginTime,
		EndTime:       activity.EndTime,
		OwnerUid:      activity.OwnerUid,
		Players:       nil,
	}
	for _, player := range activity.Players {
		ret.Players = append(ret.Players, &outer.GetGroupActivityDetail_GroupActivityDetail_Player{
			Uid:    player.Uid,
			Remark: player.Remark,
		})
	}
	a.GetLogger().Infof("get activity success, aid:%v", req.GetActivityId())
	return &outer.GetGroupActivityDetail_Reply{Activity: ret}, code.Code_Ok
}
func (a *ActivityComponent) HandleCreateGroupActivity(req *outer.CreateGroupActivity_Request) (*outer.CreateGroupActivity_Reply, code.Code) {
	if req.GetBeginTime() >= req.GetEndTime() {
		a.GetLogger().Errorf("activity time illegal, begin:%v|end:%v", req.GetBeginTime(), req.GetEndTime())
		return nil, code.Code_GroupActivityTimeIllegal
	}
	if cod := checkActivityTitle(req.GetActivityTitle()); cod != code.Code_Ok {
		a.GetLogger().Errorf("check title err, title:%v|cod:%v", req.GetActivityTitle(), cod)
		return nil, cod
	}
	if cod := checkActivityDes(req.GetActivityDes()); cod != code.Code_Ok {
		a.GetLogger().Errorf("check des err, des:%v|cod:%v", req.GetActivityDes(), cod)
		return nil, cod
	}
	groupOwner := a.GetHost().GetOwnerUid()
	if req.GetUid() != groupOwner {
		a.GetLogger().Errorf("not groupOwner, no access to create activity, req:%v|groupOwner:%v", req.GetUid(), groupOwner)
		return nil, code.Code_NoAccess
	}
	activity := &group_plugin_activity_model.ActivityItem{
		Aid:       uuid.Generate(),
		Players:   nil,
		Title:     req.GetActivityTitle(),
		Des:       req.GetActivityDes(),
		BeginTime: req.GetBeginTime(),
		EndTime:   req.GetEndTime(),
		OwnerUid:  groupOwner,
	}
	cod := a.GetModel().AddActivity(activity)
	if cod != code.Code_Ok {
		a.GetLogger().Errorf("add activity err, aid:%v|cod:%v", activity.Aid, cod)
		return nil, cod
	}
	a.GetLogger().Infof("create activity success, aid:%v", activity.Aid)
	return &outer.CreateGroupActivity_Reply{Activity: activity.ToSummary()}, code.Code_Ok
}
func (a *ActivityComponent) HandleUpdateGroupActivity(req *outer.UpdateGroupActivity_Request) (*outer.UpdateGroupActivity_Reply, code.Code) {
	if req.GetBeginTime() >= req.GetEndTime() {
		a.GetLogger().Errorf("activity time illegal, begin:%v|end:%v", req.GetBeginTime(), req.GetEndTime())
		return nil, code.Code_GroupActivityTimeIllegal
	}
	if cod := checkActivityTitle(req.GetActivityTitle()); cod != code.Code_Ok {
		a.GetLogger().Errorf("check title err, title:%v|cod:%v", req.GetActivityTitle(), cod)
		return nil, cod
	}
	if cod := checkActivityDes(req.GetActivityDes()); cod != code.Code_Ok {
		a.GetLogger().Errorf("check des err, des:%v|cod:%v", req.GetActivityDes(), cod)
		return nil, cod
	}
	activity, cod := a.GetModel().UpdateActivity(req)
	if cod != code.Code_Ok {
		a.GetLogger().Errorf("add activity err, aid:%v|cod:%v", req.GetActivityId(), cod)
		return nil, cod
	}
	if req.GetUid() != activity.OwnerUid {
		a.GetLogger().Errorf("not owner, no access to update activity, req:%v|owner:%v", req.GetUid(), activity.OwnerUid)
		return nil, code.Code_NoAccess
	}
	a.GetLogger().Infof("update activity success, aid:%v", req.GetActivityId())
	return &outer.UpdateGroupActivity_Reply{Activity: activity.ToSummary()}, code.Code_Ok
}
func (a *ActivityComponent) HandleDeleteGroupActivity(req *outer.DeleteGroupActivity_Request) (*outer.DeleteGroupActivity_Reply, code.Code) {
	activity, _ := a.GetModel().GetActivity(req.GetActivityId())
	if activity != nil {
		groupOwner := a.GetHost().GetOwnerUid()
		if req.GetUid() != groupOwner && activity.OwnerUid != req.GetGid() {
			a.GetLogger().Errorf("not groupOwner and owner, no access to delete activity, req:%v|groupOwner:%v|owner:%v", req.GetUid(), groupOwner, activity.OwnerUid)
			return nil, code.Code_NoAccess
		}
		a.GetModel().DeleteActivity(req.GetActivityId())
	}
	a.GetLogger().Infof("delete activity success, aid:%v", req.GetActivityId())
	return &outer.DeleteGroupActivity_Reply{}, code.Code_Ok
}
func (a *ActivityComponent) HandleSignUpGroupActivity(req *outer.SignUpGroupActivity_Request) (*outer.SignUpGroupActivity_Reply, code.Code) {
	activity, cod := a.GetModel().GetActivity(req.GetActivityId())
	if cod != code.Code_Ok {
		a.GetLogger().Errorf("sign up activity err, aid:%v|uid:%v|cod:%v", req.GetActivityId(), req.GetUid(), cod)
		return nil, cod
	}
	if activity.IsExpire() {
		a.GetLogger().Errorf("sign up activity err, activity expired, aid:%v|uid:%v|begin:%v|end:%v", req.GetActivityId(), req.GetUid(), activity.BeginTime, activity.EndTime)
		return nil, code.Code_GroupActivityExpired
	}
	if cod := checkActivitySignUpRemark(req.GetRemark()); cod != code.Code_Ok {
		a.GetLogger().Errorf("check remark err, remark:%v|cod:%v", req.GetRemark(), cod)
		return nil, cod
	}
	player := activity.SignUp(req.GetUid(), req.GetRemark())
	a.GetLogger().Infof("sign in activity success, aid:%v|uid:%v", req.GetActivityId(), req.GetUid())
	return &outer.SignUpGroupActivity_Reply{Player: &outer.SignUpGroupActivity_Player{Uid: player.Uid, Remark: player.Remark}}, code.Code_Ok
}
func (a *ActivityComponent) HandleSignOutGroupActivity(req *outer.SignOutGroupActivity_Request) (*outer.SignOutGroupActivity_Reply, code.Code) {
	activity, _ := a.GetModel().GetActivity(req.GetActivityId())
	if activity != nil {
		activity.SignOut(req.GetUid())
	}
	if activity.IsExpire() {
		a.GetLogger().Errorf("sign out activity err, activity expired, aid:%v|uid:%v|begin:%v|end:%v", req.GetActivityId(), req.GetUid(), activity.BeginTime, activity.EndTime)
		return nil, code.Code_GroupActivityExpired
	}
	a.GetLogger().Infof("sign out activity success, aid:%v|uid:%v", req.GetActivityId(), req.GetUid())
	return &outer.SignOutGroupActivity_Reply{}, code.Code_Ok
}
func (a *ActivityComponent) HandleBroadcastGroupActivityInfo(req *outer.BroadcastGroupActivityInfo_Request) (*outer.BroadcastGroupActivityInfo_Reply, code.Code) {
	//产品说取消所有者校验
	//owner := a.GetHost().GetOwnerUid()
	//if req.GetUid() != owner {
	//	a.GetLogger().Errorf("not owner, no access to broad activity info msg, req:%v|owner:%v", req.GetUid(), owner)
	//	return nil, code.Code_NoAccess
	//}
	activity, cod := a.GetModel().GetActivity(req.GetActivityId())
	if cod != code.Code_Ok {
		a.GetLogger().Errorf("not found activity:%v", req.GetActivityId())
		return nil, cod
	}
	player, err := sub_model.LoadPlayerHead(req.GetUid())
	if err != nil {
		a.GetLogger().Errorf("load player err, err:%v", err)
		return nil, code.Code_InnerError
	}
	chatMsg := &outer.ChatMsg_GroupActivityMsg{
		SenderName:    player.Name,
		SenderIcon:    player.Icon,
		ActivityId:    activity.Aid,
		ActivityTitle: activity.Title,
		BeginTime:     activity.BeginTime,
		EndTime:       activity.EndTime,
	}
	call.Send(a.GetCtx(), a.GetHost().GetCtx().Self(), &inner.NotifyGroupActivityInfo_Notify{ActivityInfo: chatMsg, SenderUid: req.GetUid()})
	a.GetLogger().Infof("broadcast activity info success, aid:%v", req.GetActivityId())
	return &outer.BroadcastGroupActivityInfo_Reply{}, code.Code_Ok
}
func (a *ActivityComponent) HandleBroadcastGroupActivityOnMic(req *outer.BroadcastGroupActivityOnMic_Request) (*outer.BroadcastGroupActivityOnMic_Reply, code.Code) {
	activity, cod := a.GetModel().GetActivity(req.GetActivityId())
	if cod != code.Code_Ok {
		a.GetLogger().Errorf("not found activity:%v", req.GetActivityId())
		return nil, cod
	}
	if req.GetUid() != activity.OwnerUid {
		a.GetLogger().Errorf("not owner, no access to send activity onMic msg, req:%v|owner:%v", req.GetUid(), activity.OwnerUid)
		return nil, code.Code_NoAccess
	}
	var uids []uint64
	for _, player := range activity.Players {
		uids = append(uids, player.Uid)
	}
	call.Send(a.GetCtx(), a.GetHost().GetCtx().Self(), &inner.NotifyGroupInviteOnMic_Notify{SenderUid: req.GetUid(), ActivityTitle: activity.Title, Uids: uids})
	a.GetLogger().Infof("broadcast activity onMic msg success, aid:%v", req.GetActivityId())
	return &outer.BroadcastGroupActivityOnMic_Reply{}, code.Code_Ok
}
