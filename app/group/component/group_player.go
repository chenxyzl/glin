package component

import (
	"golang.org/x/exp/slices"
	"laiya/proto/code"
	"laiya/proto/inner"
	"laiya/proto/outer"
	"laiya/share/livekit_helper"
)

func (a *GroupComponent) HandleHo2GrPlayerOnline(req *inner.Ho2GrPlayerOnline_Request) (*inner.Ho2GrPlayerOnline_Reply, code.Code) {
	player, ok := a.GetModel().Players[req.GetUid()]
	if !ok {
		a.GetLogger().Errorf("not in group, uid:%v", req.GetUid())
		return nil, code.Code_NotInGroup
	}
	a.GetModel().Online(player.Uid, req.GetSession())
	//通知有人在线
	a.NotifyGroupMsg(&outer.PushGroupPlayerOnline_Push{Uids: []uint64{req.GetUid()}}, true)
	a.GetLogger().Infof("user online, uid:%v", req.GetUid())
	return &inner.Ho2GrPlayerOnline_Reply{LastChatId: a.GetModel().LastChatId}, code.Code_Ok
}

func (a *GroupComponent) HandleHo2GrPlayerOffline(req *inner.Ho2GrPlayerOffline_Request) (*inner.Ho2GrPlayerOffline_Reply, code.Code) {
	player, ok := a.GetModel().Players[req.GetUid()]
	if !ok {
		a.GetLogger().Errorf("not in group, uid:%v", req.GetUid())
		return nil, code.Code_NotInGroup
	}
	//如果在房间,需要更新语音
	if slices.Contains(a.GetModel().OnMicPlayers, req.GetUid()) {
		//强制离开
		_, err := livekit_helper.RemoveVoicePlayer(req.GetUid(), a.GetModel().Gid)
		if err != nil {
			a.GetLogger().Errorf("force live player from voice err, uid:%v|err:%v", req.GetUid(), err)
		}
		//下一秒刷新
		a.GetModel().MarkVoiceDirty()
		a.GetLogger().Infof("user on mic offline, uid:%v", req.GetUid())
	}
	//session
	a.GetModel().Offline(player.Uid, req.GetSession())
	//通知有人离线
	a.NotifyGroupMsg(&outer.PushGroupPlayerOffline_Push{Uids: []uint64{req.GetUid()}}, true)
	a.GetLogger().Infof("user offline, uid:%v", req.GetUid())
	return &inner.Ho2GrPlayerOffline_Reply{}, code.Code_Ok
}
