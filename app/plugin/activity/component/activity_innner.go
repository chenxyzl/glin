package component

import (
	"laiya/proto/inner"
)

// HandlePlayerUnfavoriteGroup todo 玩家离开群的时候应该需要从报名列表中删除,目前暂时不调用这个逻辑
func (a *ActivityComponent) HandlePlayerUnfavoriteGroup(req *inner.Gr2AcPlayerUnfavoriteGroup_Notify) {
	for _, activity := range a.GetModel().Activities {
		if activity.IsExpire() {
			continue
		}
		activity.SignOut(req.GetUid())
	}
}
