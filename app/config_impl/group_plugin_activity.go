package config_impl

import (
	"github.com/BurntSushi/toml"
)

type GroupPluginActivityConfig struct {
	NotOnMicCopywriting  string
	CallOnMicCopywriting string
}

func (conf *GroupPluginActivityConfig) Parse(data []byte) error {
	var tempConfig GroupPluginActivityConfig
	_, err := toml.Decode(string(data), &tempConfig)
	if err != nil {
		return err
	}
	//
	*conf = tempConfig
	return nil
}

////////////////////////////////////////////////////////
