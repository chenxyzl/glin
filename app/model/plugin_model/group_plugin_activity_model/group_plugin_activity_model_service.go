package group_plugin_activity_model

import (
	"laiya/config"
	"laiya/proto/code"
	"laiya/proto/outer"
	"slices"
	"time"
)

func (activity *ActivityItem) ToSummary() *outer.GroupActivitySummary {
	return &outer.GroupActivitySummary{
		ActivityId:    activity.Aid,
		ActivityTitle: activity.Title,
		ActivityDes:   activity.Des,
		BeginTime:     activity.BeginTime,
		EndTime:       activity.EndTime,
		OwnerUid:      activity.OwnerUid,
	}
}

// AddActivity 增加活动
func (mod *PluginActivity) AddActivity(activity *ActivityItem) code.Code {
	if len(mod.Activities) >= config.Get().ConstConfig.GroupPluginActivityMaxCount {
		return code.Code_GroupActivityFull
	}
	mod.Activities = append(mod.Activities, activity)
	mod.MarkDirty()
	return code.Code_Ok
}

// GetActivity 获取活动
func (mod *PluginActivity) GetActivity(aid uint64) (*ActivityItem, code.Code) {
	//检查活动是否存在
	idx := slices.IndexFunc(mod.Activities, func(item *ActivityItem) bool { return aid == item.Aid })
	if idx < 0 {
		return nil, code.Code_GroupActivityNotExist
	}
	return mod.Activities[idx], code.Code_Ok
}

// UpdateActivity 修改活动
func (mod *PluginActivity) UpdateActivity(req *outer.UpdateGroupActivity_Request) (*ActivityItem, code.Code) {
	idx := slices.IndexFunc(mod.Activities, func(item *ActivityItem) bool { return req.GetActivityId() == item.Aid })
	if idx < 0 {
		return nil, code.Code_GroupActivityNotExist
	}
	activity := mod.Activities[idx]
	activity.Title = req.GetActivityTitle()
	activity.Des = req.GetActivityDes()
	activity.BeginTime = req.GetBeginTime()
	activity.EndTime = req.GetEndTime()
	mod.MarkDirty()
	return activity, code.Code_Ok
}

// DeleteActivity 删除活动
func (mod *PluginActivity) DeleteActivity(aid uint64) {
	mod.Activities = slices.DeleteFunc(mod.Activities, func(item *ActivityItem) bool { return item.Aid == aid })
	mod.MarkDirty()
}

// SignUp 报名活动
func (activity *ActivityItem) SignUp(uid uint64, remark string) *Player {
	idx := slices.IndexFunc(activity.Players, func(player *Player) bool { return player.Uid == uid })
	player := &Player{Uid: uid, Remark: remark}
	//新增还是修改
	if idx < 0 {
		//新增
		activity.Players = append(activity.Players, player)
	} else {
		//修改
		activity.Players[idx] = player
	}
	return player
}

// SignOut 取消报名
func (activity *ActivityItem) SignOut(uid uint64) {
	activity.Players = slices.DeleteFunc(activity.Players, func(player *Player) bool { return player.Uid == uid })
}

// IsExpire 活动是否过期
func (activity *ActivityItem) IsExpire() bool {
	now := time.Now().UnixMilli()
	return now >= activity.EndTime
}

// IsWillOpen 活动未开放
func (activity *ActivityItem) IsWillOpen() bool {
	now := time.Now().UnixMilli()
	return now <= activity.BeginTime
}

// IsAtOpen 活动在开放中
func (activity *ActivityItem) IsAtOpen() bool {
	now := time.Now().UnixMilli()
	return now > activity.BeginTime && now < activity.EndTime
}
