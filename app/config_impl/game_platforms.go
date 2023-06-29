package config_impl

import (
	"github.com/BurntSushi/toml"
	"strconv"
)

type Platform struct {
	Name string
	Icon string
	Idx  int32
}

type GamePlatformsConfig struct {
	PlatformTypes    map[int32]Platform
	PlatformTypesTmp map[string]Platform
}

func (conf *GamePlatformsConfig) Parse(data []byte) error {
	var tempConfig GamePlatformsConfig
	_, err := toml.Decode(string(data), &tempConfig)
	if err != nil {
		return err
	}
	//CheckRobotConfigGameType
	tempConfig.PlatformTypes = make(map[int32]Platform)
	for s, s2 := range tempConfig.PlatformTypesTmp {
		i, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		tempConfig.PlatformTypes[int32(i)] = s2
	}
	//
	*conf = tempConfig
	return nil
}

////////////////////////////////////////////////////////
