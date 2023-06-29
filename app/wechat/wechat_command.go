package wechat

import (
	"laiya/config"
)

type TWechatCommand string

const (
	WechatCommandUnknown                   TWechatCommand = ""
	WechatCommandLookingGroup              TWechatCommand = "LookingGroup"
	WechatCommandUpdateWechatGroups        TWechatCommand = "UpdateWechatGroups"
	WechatCommandUpdateUpdateWechatFriends TWechatCommand = "UpdateWechatFriends"
)

func GetWechatCommand(content string) TWechatCommand {
	fixMsg := formatContent(content)
	for _, s := range config.Get().RobotConfig.Command.LookingGroup {
		if fixMsg == s {
			return WechatCommandLookingGroup
		}
	}
	for _, s := range config.Get().RobotConfig.Command.UpdateWechatGroups {
		if fixMsg == s {
			return WechatCommandUpdateWechatGroups
		}
	}
	for _, s := range config.Get().RobotConfig.Command.UpdateWechatFriends {
		if fixMsg == s {
			return WechatCommandUpdateUpdateWechatFriends
		}
	}
	return WechatCommandUnknown
}
