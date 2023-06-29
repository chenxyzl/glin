package config_impl

import (
	"github.com/BurntSushi/toml"
	"time"
)

type ParamsConfig struct {
	//PlatformTypes                         []string
	ContinuityNotifyTimes                 int64
	LookingGroupCount                     int
	NotifyGroupCount                      int
	WechatInviteBroadcastInterval         time.Duration
	CheckInviteBroadcastTimeInterval      time.Duration
	CheckInviteBroadcastLookingGroupCount int
	CheckInviteBroadcastHeartbeatTimeout  time.Duration
	InviteBroadcastImCopywritingInterval  time.Duration
	VoiceSyncCheckInterval                uint64
	HallNotifyTimeInterval                time.Duration
	HallNotifyTimeIntervalRange           time.Duration
	AllNotifyBeginHour                    int64
	AllNotifyEndHour                      int64
	PlayerDefaultGameType                 string
	LookingNotifyWaitingTime              time.Duration
	KOLStartTag                           string
}

func (conf *ParamsConfig) Parse(data []byte) error {
	var tempConfig ParamsConfig
	_, err := toml.Decode(string(data), &tempConfig)
	if err != nil {
		return err
	}
	//
	*conf = tempConfig
	return nil
}

////////////////////////////////////////////////////////
