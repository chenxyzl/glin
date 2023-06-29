package config_impl

import (
	"github.com/BurntSushi/toml"
	"strconv"
)

type Command struct {
	LookingGroup        []string
	UpdateWechatGroups  []string
	UpdateWechatFriends []string
}

type RobotConfig struct {
	Command          Command
	Robots           map[uint64]string
	RobotGameTypeTmp map[string]string
}

func (conf *RobotConfig) Parse(data []byte) error {
	var tempConfig RobotConfig
	_, err := toml.Decode(string(data), &tempConfig)
	if err != nil {
		return err
	}
	//CheckRobotConfigGameType
	tempConfig.Robots = make(map[uint64]string)
	for s, s2 := range tempConfig.RobotGameTypeTmp {
		i, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		tempConfig.Robots[uint64(i)] = s2
	}
	//
	*conf = tempConfig
	return nil
}

////////////////////////////////////////////////////////
