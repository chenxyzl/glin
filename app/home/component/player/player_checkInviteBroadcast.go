package player

import (
	"fmt"
	"laiya/config"
	"laiya/proto/code"
	"laiya/proto/inner"
	"laiya/proto/outer"
	"laiya/share/algorithm/weight_random"
	"laiya/share/call"
	"laiya/share/global"
	"math/rand"
	"time"
)

func (a *PlayerComponent) checkInviteBroadcast() {
	//session全部断开则认为已掉线
	if len(a.GetModel().Sess) <= 0 {
		return
	}
	now := time.Now()
	gameConfig := config.Get().ParamsConfig
	//上次状态同步时间小于当前则认为当前不在可以接收的状态 跳过检查
	if a.GetModel().LastSyncUserStateTime+int64(gameConfig.CheckInviteBroadcastHeartbeatTimeout) < now.UnixNano() {
		return
	}
	//发送间隔过小则跳过
	if a.GetModel().LastCheckInviteBroadcastTime+int64(gameConfig.CheckInviteBroadcastTimeInterval) > now.UnixNano() {
		return
	}
	//记录
	a.GetModel().LastCheckInviteBroadcastTime = now.UnixNano()
	gameType := gameConfig.PlayerDefaultGameType
	//
	ret, cod := call.RequestRemote[*inner.Ho2HaCheckInviteBroadcast_Reply](global.HallUID, &inner.Ho2HaCheckInviteBroadcast_Request{GameType: gameType})
	if cod != code.Code_Ok {
		a.GetLogger().Errorf("check invite broad cast err,cod:%v", cod)
		return
	}
	//没有合适的群
	if len(ret.MatchGroups) == 0 {
		a.GetLogger().Info("check invite broad cast success, but no group found")
		return
	}

	var sl weight_random.WeightSlice[*randomItem]
Label:
	for _, group := range ret.MatchGroups {
		//不推荐当前所在的群
		if a.GetModel().NowInGroup == group.GetGid() {
			continue Label
		}
		//不推荐进入过的群
		for _, broadcastGroup := range a.GetModel().LastInviteBroadcastGroup {
			if broadcastGroup.Gid == group.GetGid() && broadcastGroup.TempVoiceRoomVersion == group.GetTempVoiceRoomVersion() {
				continue Label
			}
		}
		sl = append(sl, &randomItem{gameType: gameType, gid: group.Gid, playerCount: group.GetPlayerCount(), tempVoiceRoomVersion: group.GetTempVoiceRoomVersion()})
	}
	item := sl.Rand(rand.Uint64())
	if item == nil {
		a.GetLogger().Infof("check invite broad cast success, but all group ignore, ret:%v|ignore:%v", ret.MatchGroups, a.GetModel().LastInviteBroadcastGroup)
		return
	}
	//广播客户端
	a.GetModel().AddLastInviteGroup(item.gid, item.tempVoiceRoomVersion)
	call.Send(a.GetCtx(), a.GetCtx().Self(), outer.BuildPushPack(&outer.PushInviteBroadcast_Push{Gid: item.gid}))
	a.GetLogger().Infof("check invite broad cast success, select:%v", ret)
}

type randomItem struct {
	gameType             string
	gid                  uint64
	playerCount          uint32
	tempVoiceRoomVersion uint64
}

func (r randomItem) Weight() uint64 {
	conf := config.Get().GameTypesConfig.GetParams(r.gameType)
	count := fmt.Sprintf("%d", r.playerCount)
	return conf.PlayerCountWeight[count]
}
