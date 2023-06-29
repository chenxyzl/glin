package component

import (
	"laiya/config"
	"laiya/model/hall_model"
	"laiya/proto/inner"
	"laiya/share/call"
	"laiya/share/global"
	"math/rand"
	"strings"
	"time"
)

func (a *HallComponent) checkTimerNotify(now time.Time) {
	hour := int64(now.Hour())
	if hour == 0 {
		hour = 24
	}
	//时间初始化
	if a.GetModel().NextNotifyTime == 0 {
		a.GetModel().UpdateNextNotifyTime(now.Unix())
	}
	//修正时间~重新生成~避免时间调整后不生效
	maxNextTime := config.Get().ParamsConfig.GetHallMaxNextNotifyTime(now)
	if a.GetModel().NextNotifyTime > maxNextTime.Unix() {
		a.GetModel().UpdateNextNotifyTime(maxNextTime.Unix())
		a.GetLogger().Infof("checkTimerNotify ,fix  next notify time, next:%v", maxNextTime)
		return
	}
	//时间检查
	if a.GetModel().NextNotifyTime > now.Unix() {
		//a.GetLogger().Infof("checkTimerNotify ,not now, now:%v|next:%v", now.Unix(), a.GetModel().NextNotifyTime)
		return
	}

	//设置下一次广播时间-----------要在时间段之前检查才能保证不出现集中推送
	next := config.Get().ParamsConfig.GenHallNextNotifyTime(now)
	a.GetModel().UpdateNextNotifyTime(next.Unix())

	//时间段检查
	inTime := hour >= config.Get().ParamsConfig.AllNotifyBeginHour && hour <= config.Get().ParamsConfig.AllNotifyEndHour
	if !inTime {
		a.GetLogger().Infof("checkTimerNotify ,now in sleep time, hour:%v|begin:%v|end:%v", hour, config.Get().ParamsConfig.AllNotifyBeginHour, config.Get().ParamsConfig.AllNotifyEndHour)
		return
	}
	//开始广播
	a.timerNotify()
}

func (a *HallComponent) timerNotify() {
	gameTypeGroups := make(map[string][]*hall_model.Group)
	a.GetModel().SortGroupList.ForRange(func(v *hall_model.Group) bool {
		//只广播有人的房间
		if v.GetPlayerCount() <= 0 {
			return false
		}
		//todo 临时特殊处理 不广播私有房间
		if strings.Contains(v.Name, "私房") {
			return true
		}
		//只广播有人且不满的
		if v.GetPlayerCount() >= config.Get().GameTypesConfig.GetParams(v.GameType).FullPlayerCount {
			return true
		}
		//
		if _, ok := gameTypeGroups[v.GameType]; !ok {
			gameTypeGroups[v.GameType] = make([]*hall_model.Group, 0)
		}
		gameTypeGroups[v.GameType] = append(gameTypeGroups[v.GameType], v)
		return true
	}, true)
	if len(gameTypeGroups) == 0 {
		a.GetLogger().Infof("timerNotify ,no groups")
		return
	}
	for gameType, groups := range gameTypeGroups {
		if len(groups) == 0 {
			continue
		}
		if len(groups) > config.Get().ParamsConfig.LookingGroupCount {
			//打乱顺序
			rand.Shuffle(len(groups), func(i, j int) {
				groups[i], groups[j] = groups[j], groups[i]
			})
			//获取前面n个
			groups = groups[:config.Get().ParamsConfig.LookingGroupCount]
		}
		a.timerNotifyContent(gameType, groups)
	}
	a.GetLogger().Infof("timerNotify success")
}

func (a *HallComponent) timerNotifyContent(gameType string, groups []*hall_model.Group) {
	if len(groups) == 0 {
		a.GetLogger().Infof("timerNotifyContent ,no groups, gameType:%v", gameType)
		return
	}
	var v []*inner.Ha2WechatLookingMsg_Group
	for _, group := range groups {
		v = append(v, &inner.Ha2WechatLookingMsg_Group{
			Gid:           group.Gid,
			GroupName:     group.Name,
			VoiceRoomName: group.VoiceRoomName,
			PlayerCount:   group.GetPlayerCount(),
		})
	}
	call.NotifyRemote(global.WechatUID, &inner.Ha2WechatLookingMsg_Notify{
		GameType: gameType,
		Groups:   v,
	})

	a.GetLogger().Infof("timerNotifyContent ,gameType:%v|groups:%v", gameType, len(groups))
}
