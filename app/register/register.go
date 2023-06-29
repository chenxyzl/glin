package register

import (
	"github.com/chenxyzl/glin/grain"
	"github.com/chenxyzl/glin/slog"
	"laiya/gate"
	gatecomponent "laiya/gate/component"
	gateiface "laiya/gate/iface"
	"laiya/group"
	groupcompoent "laiya/group/component"
	groupiface "laiya/group/iface"
	"laiya/hall"
	hallcomponent "laiya/hall/component"
	halliface "laiya/hall/iface"
	"laiya/home"
	homegroupcomponent "laiya/home/component/group"
	homeplayercomponent "laiya/home/component/player"
	homeiface "laiya/home/iface"
	"laiya/plugin/activity"
	"laiya/robot_local/animal_party"
	"laiya/robot_local/overcooked"
	"laiya/web"
	webcomponent "laiya/web/component"
	webiface "laiya/web/iface"
	"laiya/wechat"
	wechatcomponent "laiya/wechat/component"
	wechatiface "laiya/wechat/iface"

	activitycomponent "laiya/plugin/activity/component"
	activityiface "laiya/plugin/activity/iface"
)

func init() {
	slog.Info("register component begin...")
	//web 组件注册
	grain.RegisterComponent[*web.WebEntity, *webcomponent.WebComponent, webiface.IWebComponent](false)

	//gate 组件注册
	grain.RegisterComponent[*gate.SessionActor, *gatecomponent.SessionComponent, gateiface.ISessionComponent](false)

	//home 组件注册
	grain.RegisterComponent[*home.PlayerActor, *homeplayercomponent.PlayerComponent, homeiface.IPlayerComponent](true)
	grain.RegisterComponent[*home.PlayerActor, *homegroupcomponent.GroupComponent, homeiface.IGroupComponent](true)

	//hall 组件注册
	grain.RegisterComponent[*hall.HallActor, *hallcomponent.HallComponent, halliface.IHallComponent](true)

	//group 组件注册
	grain.RegisterComponent[*group.GroupActor, *groupcompoent.GroupComponent, groupiface.IGroupComponent](true)

	//wechat 组件注册
	grain.RegisterComponent[*wechat.WechatActor, *wechatcomponent.WechatComponent, wechatiface.IWechatComponent](true)

	//group_胡闹厨房机器人 组件注册
	grain.RegisterComponent[*overcooked.OvercookedActor, *overcooked.OvercookedComponent, overcooked.IOvercookedComponent](false)
	//group_动物派对机器人 组件注册
	grain.RegisterComponent[*animal_party.AnimalPartyActor, *animal_party.AnimalPartyComponent, animal_party.IAnimalPartyComponent](false)

	//group_插件_活动
	grain.RegisterComponent[*activity.ActivityActor, *activitycomponent.ActivityComponent, activityiface.IActivityComponent](false)

	//
	slog.Info("register component success...")
}
