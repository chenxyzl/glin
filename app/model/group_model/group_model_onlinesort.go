package group_model

import (
	"github.com/asynkron/protoactor-go/actor"
	"laiya/share/algorithm/skiplist"
	"time"
)

func (mod *Group) NewOnlineSkipList() {
	mod.OnlineSort = skiplist.New[*OnlineSortT](&cmpOnline{})
}

func (mod *Group) Online(uid uint64, session *actor.PID) {
	//不在群的不处理--前面接口已经拦截过了
	if mod.GetPlayer(uid) == nil {
		return
	}
	newState := &OnlineSortT{uid: uid, sess: session}
	//先删除老的--并记录state
	onlineState, ok := mod.OnlineSort.Get(uid)
	if ok {
		mod.OnlineSort.Delete(uid)
		newState.SetInGroup(onlineState.IsNowInGroup())
	}
	//插入新的
	mod.OnlineSort.Insert(newState)
}

func (mod *Group) Offline(uid uint64, sess *actor.PID) {
	mod.OnlineSort.Delete(uid)
}
func (mod *Group) Enter(uid uint64, session *actor.PID) {
	//群用户何临时用户都重新设置在群中
	newState := &OnlineSortT{uid: uid, sess: session}
	//先删除老的--并记录state
	_, ok := mod.OnlineSort.Get(uid)
	if ok {
		mod.OnlineSort.Delete(uid)
	}
	newState.SetInGroup(true)
	//插入新的
	mod.OnlineSort.Insert(newState)
	//更新活跃时间
	mod.LastActivityTime = time.Now().Unix()
}

func (mod *Group) Exit(uid uint64) {
	//群用户和临时用户分开处理
	if mod.Players[uid] != nil { //群内用户只需要设置不在当前群
		state, ok := mod.OnlineSort.Get(uid)
		if ok {
			state.SetInGroup(false)
		}
	} else { //临时用户需要从在线列表中移除
		mod.OnlineSort.Delete(uid)
	}
	//更新活跃时间
	mod.LastActivityTime = time.Now().Unix()
}

type cmpOnline struct {
}

func (a *cmpOnline) CmpScore(v1 *OnlineSortT, v2 *OnlineSortT) int {
	s1 := v1.Score()
	s2 := v2.Score()
	switch {
	case s1 < s2:
		return -1
	case s1 == s2:
		return 0
	default:
		return 1
	}
}

func (a *cmpOnline) CmpKey(v1 *OnlineSortT, v2 *OnlineSortT) int {
	s1 := v1.Key()
	s2 := v2.Key()
	switch {
	case s1 < s2:
		return -1
	case s1 == s2:
		return 0
	default:
		return 1
	}
}

type OnlineSortT struct {
	uid                  uint64
	nowInGroup           bool
	sess                 *actor.PID
	lastHeartbeatTimeSec int64 //最近一次心跳时间
}

func (a *OnlineSortT) Key() uint64 {
	return a.uid
}

func (a *OnlineSortT) Score() uint64 {
	return a.uid
}

func (a *OnlineSortT) SetInGroup(status bool) {
	a.nowInGroup = status
}

func (a *OnlineSortT) IsNowInGroup() bool {
	return a.nowInGroup
}

func (a *OnlineSortT) GetSess() *actor.PID {
	if a == nil {
		return nil
	}
	return a.sess
}
