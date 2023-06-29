package config_impl

import (
	"github.com/BurntSushi/toml"
	"time"
)

type ConstConfig struct {
	UserDefaultIcon                       string
	UserNameLenMin                        int
	UserNameLenMax                        int
	UserDesLenMax                         int
	GroupNameLenMin                       int
	GroupNameLenMax                       int
	GroupDesLenMax                        int
	VoiceRoomNameLenMin                   int
	VoiceRoomNameLenMax                   int
	DefaultVoiceRoomName                  string
	MaxChatHistory                        uint32
	MaxGroupCount                         uint32
	MaxCreateGroupCount                   uint32
	GroupPlayerMaxCount                   uint32
	GroupOnlineLimit                      int
	ChatMaxLen                            uint32
	GodAcc                                []string
	GodCaptcha                            string
	HeartbeatCheck                        time.Duration
	PlayerTerminateIdleTime               time.Duration
	GroupTerminateIdleTime                time.Duration
	ChatCacheSize                         uint32
	GroupChatKeepTime                     time.Duration
	PrivateChatCacheSize                  uint32
	PrivateChatKeepTime                   time.Duration
	PrivateChatKeepSize                   int
	PrivateChatTitle                      string
	GroupPluginActivityMaxCount           int
	GroupPluginActivityTitleMinLen        int
	GroupPluginActivityTitleMaxLen        int
	GroupPluginActivityDesMaxLen          int
	GroupPluginActivitySignUpRemarkMaxLen int
}

func (conf *ConstConfig) Parse(data []byte) error {
	_, err := toml.Decode(string(data), conf)
	return err
}
