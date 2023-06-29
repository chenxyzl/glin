package config_impl

import (
	"github.com/BurntSushi/toml"
)

type AccessWrapper struct {
	Access []int32
}

type GroupPermissionsConfig struct {
	GroupPermissions map[string]AccessWrapper
}

func (conf *GroupPermissionsConfig) Parse(data []byte) error {
	var tempConfig GroupPermissionsConfig
	_, err := toml.Decode(string(data), &tempConfig)
	if err != nil {
		return err
	}
	//
	*conf = tempConfig
	return nil
}

////////////////////////////////////////////////////////
