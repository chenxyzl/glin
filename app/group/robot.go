package group

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/chenxyzl/glin"
	"github.com/chenxyzl/glin/grain"
	"laiya/config"
	"laiya/group/iface"
	"laiya/model/group_model"
	"laiya/model/sub_model"
	"laiya/proto/inner"
	"laiya/proto/outer"
	"laiya/share/call"
)

func (a *GroupActor) GetRobotByUid(robotId uint64) *actor.PID {
	return a.robotMaps[robotId]
}

func (a *GroupActor) CheckAddRobot() {
	//初始化各种机器人
	for robotId, gameType := range config.Get().RobotConfig.Robots {
		//判断游戏类型是否对的上
		if gameType != a.GetModel().GameType {
			continue
		}
		//builder
		fun := grain.GetRobotFactory(robotId)
		if fun == nil {
			a.GetLogger().Errorf("RobotFactory not found, robotId:%v", robotId)
			continue
		}
		//check add
		firstAdd := !a.GetModel().IsPlayerInGroup(robotId)
		if firstAdd {
			robotInfo := sub_model.PlayerHead{Uid: robotId}
			err := robotInfo.Load()
			if err != nil {
				a.GetLogger().Errorf("load robot info err, robotId:%v|err:%v", robotId, err)
				continue
			}
			a.GetModel().AddPlayer(&group_model.Player{Uid: robotId})
			a.GetLogger().Infof("add robotId to group, robotId:%v", robotId)

			//添加机器人后的的通知
			grain.GetComponent[iface.IGroupComponent](a).NotifyChatMsg(
				&outer.ChatMsg_AddRobotMsg{RobotId: robotId, RobotName: robotInfo.Name, RobotIcon: robotInfo.Icon},
				a.GetModel().Gid)
		}
		//
		robot := glin.GetSystem().NewLocalActor(func() grain.IActor {
			return fun(a.GetModel().Gid, a.GetCtx().Self())
		})
		//
		a.robotMaps[robotId] = robot
		//
		if firstAdd {
			call.Send(a.GetCtx(), robot, &inner.Gr2RoAddRobot_Notify{})
		}
		//
		a.GetLogger().Infof("wakeup robot success, robotId:%v|robot:%v", robotId, robot)
	}
}
