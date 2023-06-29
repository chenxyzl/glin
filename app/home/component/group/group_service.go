package player_group

import (
	"github.com/chenxyzl/glin"
	"laiya/share/global"
)

// 监听群激活的消息
func (a *GroupComponent) initSbGroupActive() {
	for _, group := range a.GetModel().Groups {
		a.AddSbGroupActive(group.Gid)
	}
}

// 监听群激活的消息
func (a *GroupComponent) removeAllSbGroupActive() {
	for _, group := range a.GetModel().Groups {
		a.RemoveSbGroupActive(group.Gid)
	}
}
func (a *GroupComponent) AddSbGroupActive(gid uint64) {
	//取消监听
	_, err := glin.GetSystem().GetCluster().SubscribeByPid(global.GetGroupActiveTopic(gid), a.GetCtx().Self())
	if err != nil {
		a.GetLogger().Errorf("add group active subscribe err, gid:%v|err:%v", gid, err)
		return
	}
	a.GetLogger().Infof("add group active subscribe success, gid:%v", gid)
}
func (a *GroupComponent) RemoveSbGroupActive(gid uint64) {
	//注册监听
	_, err := glin.GetSystem().GetCluster().UnsubscribeByPid(global.GetGroupActiveTopic(gid), a.GetCtx().Self())
	if err != nil {
		a.GetLogger().Errorf("remove group active subscribe err, gid:%v|err:%v", gid, err)
		return
	}
	a.GetLogger().Infof("remove group active subscribe success, gid:%v", gid)
}
