package config_impl

import "github.com/BurntSushi/toml"

type GameCopywriting struct {
	GroupHelper           string
	NotFoundGroup         string
	LookingNotifySuccess  string
	LookingNotifyWechat   string
	TimerNotifyWechatItem string
	TimerNotifyWechat     string
	LookingGroupItem      string
	LookingGroup          string
	OnMicPlayer0Guide     string
	OnMicPlayer1Guide     string
}

type GameCopywritingConfig struct {
	WechatPrivateReply      string
	WechatHelper            string
	GroupNotConfigGameType  string
	RobotNotActive          string
	WaitingInviteBroadcast0 string
	WaitingInviteBroadcast1 string
	WaitingInviteBroadcast2 string
	WechatUpdateInfoSuccess string
	LookingNotifyInterval   string
	LookingNotifyWaiting    string
	Copywriting             map[string]GameCopywriting
}

func (conf *GameCopywritingConfig) Parse(data []byte) error {
	_, err := toml.Decode(string(data), conf)
	return err
}
