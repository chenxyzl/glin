package hall_model

import (
	"laiya/proto/inner"
)

func (mod *Hall) UpdateGroupInfo(updateInfo *inner.Gr2HaUpdateGroup_Request) *Group {
	//查找群
	old, ok := mod.Groups[updateInfo.Gid]
	//更新数据
	if ok {
		//更新前先删除老的跳表的值----------否者score更新会不生效，甚至重复数据的问题
		mod.SortGroupList.Delete(old.Gid)
	}
	group := &Group{
		Gid:                  updateInfo.GetGid(),
		Name:                 updateInfo.GetName(),
		VoiceRoomName:        updateInfo.GetVoiceRoomName(),
		GameType:             updateInfo.GetGameType(),
		SortScore:            updateInfo.GetSortScore(),
		TempVoiceRoomVersion: updateInfo.GetTempVoiceRoomVersion(),
	}
	mod.Groups[updateInfo.Gid] = group
	mod.SortGroupList.Insert(group)
	mod.MarkDirty()
	return group
}

func (mod *Hall) DeleteGroup(gid uint64) {
	if _, ok := mod.Groups[gid]; !ok {
		return
	}
	mod.SortGroupList.Delete(gid)
	delete(mod.Groups, gid)
	mod.MarkDirty()
}

func (mod *Hall) UpdateNextNotifyTime(nextNotifyTime int64) {
	mod.NextNotifyTime = nextNotifyTime
	mod.MarkDirty()
}
