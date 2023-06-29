package config_impl

import (
	"github.com/BurntSushi/toml"
)

type ClientDynamic map[string]string

func (conf *ClientDynamic) Parse(data []byte) error {
	tempConfig := &ClientDynamic{}
	if _, err := toml.Decode(string(data), tempConfig); err != nil {
		return err
	}
	if err := tempConfig.parseInviteBroadcastGameTypes(); err != nil {
		return err
	}
	*conf = *tempConfig
	return nil
}
