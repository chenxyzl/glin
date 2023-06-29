package component

import (
	"laiya/config"
	"laiya/model/hall_model"
	"laiya/proto/code"
	"laiya/proto/inner"
	"laiya/proto/outer"
	"laiya/share/global"
	"strconv"
	"strings"
)

func (a *HallComponent) HandleGetRecommendGroupList(req *outer.GetRecommendGroupList_Request) (*outer.GetRecommendGroupList_Reply, code.Code) {
	//a.GetModel().SortGroupList.
	ret := &outer.GetRecommendGroupList_Reply{}
	page := req.GetPage()
	//todo 临时代码-产品不需要随机
	if page == 0 {
		page = 1
	}
	if page == 0 {
		//随机获取
		for _, group := range a.GetModel().Groups {
			//随机获取
			ret.GroupIds = append(ret.GroupIds, group.Gid)
			if len(ret.GroupIds) >= int(global.PageGroupCount) {
				break
			}
		}
	} else {
		//获取页数
		pageCount := global.PageGroupCount
		begin := (page-1)*pageCount + 1
		end := page * pageCount
		list := a.GetModel().SortGroupList.GetRange(begin, end, true)
		for _, group := range list {
			ret.GroupIds = append(ret.GroupIds, group.Gid)
			//a.GetLogger().Infof("in group user count, group:%v|count:%v|times:%v|score:%v", group.Gid, group.Score()>>32, group.Score()&0x00000000FFFFFFFF, group.Score())
		}
	}
	a.GetLogger().Infof("get recommend list, page:%v|rsp:%v", req.GetPage(), ret)
	return ret, code.Code_Ok
}

func (a *HallComponent) HandleSearchGroup(req *outer.SearchGroup_Request) (*outer.SearchGroup_Reply, code.Code) {
	ret := &outer.SearchGroup_Reply{}
	//获取名字 描述 游戏 类型等关键字
	page := req.GetPage()
	if page > 0 { //默认从第一页开始
		page -= 1
	}
	//首页查询加上id搜索
	if page == 0 {
		//先按照id搜索
		uid, err := strconv.ParseUint(req.GetStrVal(), 10, 64)
		if err == nil && uid > 0 {
			//获取指定id
			group, ok := a.GetModel().Groups[uid]
			if ok {
				ret.GroupIds = append(ret.GroupIds, group.Gid)
			}
		}
	}
	//名字+游戏类型+描述
	needSkip := page * global.PageGroupCount
	nowSkip := uint32(0)
	//todo 低效的实现方式,暂时如果功能优先
	a.GetModel().SortGroupList.ForRange(func(v *hall_model.Group) bool {
		//不满足查询条件跳过
		if !(strings.Contains(v.GameType, req.GetStrVal()) ||
			strings.Contains(v.Name, req.GetStrVal()) ||
			strings.Contains(v.VoiceRoomName, req.GetStrVal())) {
			return true
		}
		//满足，先计算页数
		if needSkip > nowSkip {
			nowSkip++
			return true
		}
		ret.GroupIds = append(ret.GroupIds, v.Gid)
		//最多获取个数
		return len(ret.GroupIds) < int(global.PageGroupCount)
	}, true)
	a.GetLogger().Infof("get search list, page:%v|search:%v|rsp:%v", req.GetPage(), req.GetStrVal(), ret)
	return ret, code.Code_Ok
}

func (a *HallComponent) HandleLookingGroupAvailable(req *inner.We2HaLookingGroupAvailable_Request) (*inner.We2HaLookingGroupAvailable_Reply, code.Code) {
	var groups []*inner.We2HaLookingGroupAvailable_Group
	conf := config.Get().ParamsConfig
	gameConfig := config.Get().GameTypesConfig.GetParams(req.GetGameType())
	a.GetModel().SortGroupList.ForRange(func(v *hall_model.Group) bool {
		//够了不查
		if len(groups) >= conf.LookingGroupCount {
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
		groups = append(groups, &inner.We2HaLookingGroupAvailable_Group{
			Gid:              v.Gid,
			GroupName:        v.Name,
			VoiceRoomName:    v.VoiceRoomName,
			OnMicPlayerCount: int32(v.GetPlayerCount()),
		})
		//
		return true
	}, true)
	a.GetLogger().Infof("get available groups success, gameType:%v|count:%v|groups:%v", req.GetGameType(), len(groups), groups)
	return &inner.We2HaLookingGroupAvailable_Reply{Groups: groups}, code.Code_Ok
}
