package component

import (
	"fmt"
	"laiya/config"
	"laiya/model/sub_model"
	"laiya/proto/code"
	"laiya/proto/inner"
	"laiya/share/call"
	"laiya/share/global"
)

func (a *SimpleRobotComponent[T]) doCommandHelper(sender uint64, content string) {
	a.GetLogger().Infof("command helper")
	gameType := config.Get().RobotConfig.GetRobotGameType(a.ISimpleRobotActor.GetRobotId())
	copywritingConfig := config.Get().GameCopywritingConfig.Get(gameType)
	call.Send(a.GetCtx(), a.GetHost(), &inner.Ro2GrRobotSayChatMsg_Notify{
		RobotId:    a.ISimpleRobotActor.GetRobotId(),
		ReplyToUid: sender,
		Content:    copywritingConfig.GroupHelper,
	})
}

func (a *SimpleRobotComponent[T]) doCommandLookingGroup(sender uint64, content string) {
	mod := sub_model.GroupGameTyp{Gid: a.GetHostGid()}
	if err := mod.Load(); err != nil {
		a.GetLogger().Errorf("load group game type err, err:%v", mod.GameType)
		return
	}
	reply, cod := call.RequestRemote[*inner.We2HaLookingGroupAvailable_Reply](global.HallUID, &inner.We2HaLookingGroupAvailable_Request{GameType: mod.GameType})
	if cod != code.Code_Ok {
		a.GetLogger().Errorf("req hall looking available group  err, group:%v|cod:%v", a.GetHostGid(), cod)
		return
	}

	gameType := config.Get().RobotConfig.GetRobotGameType(a.ISimpleRobotActor.GetRobotId())
	copywritingConfig := config.Get().GameCopywritingConfig.Get(gameType)

	if len(reply.Groups) == 0 {
		call.Send(a.GetCtx(), a.GetHost(), &inner.Ro2GrRobotSayChatMsg_Notify{
			RobotId:    a.ISimpleRobotActor.GetRobotId(),
			ReplyToUid: sender,
			Content:    copywritingConfig.NotFoundGroup,
		})
		a.GetLogger().Infof("req hall looking available group, not found, group:%v", a.GetHostGid())
		return
	}
	var items string
	var itemT = copywritingConfig.LookingGroupItem
	for i, availableGroup := range reply.Groups {
		if i != 0 {
			items += "\n"
		}
		items += fmt.Sprintf(itemT, availableGroup.GroupName, availableGroup.GetVoiceRoomName(), config.Get().WebConfig.GetShorUrl(availableGroup.GetGid()), availableGroup.OnMicPlayerCount)
	}
	var total = fmt.Sprintf(copywritingConfig.LookingGroup, len(reply.Groups), items)
	call.Send(a.GetCtx(), a.GetHost(), &inner.Ro2GrRobotSayChatMsg_Notify{
		RobotId:    a.ISimpleRobotActor.GetRobotId(),
		ReplyToUid: sender,
		Content:    total,
	})

	a.GetLogger().Infof("command looking group, ret:%v", reply.GetGroups())
}
