package config_impl

import (
	"github.com/BurntSushi/toml"
)

type GameParams struct {
	FullPlayerCount   uint32
	PlayerCountWeight map[string]uint64
}

type GameTypesConfig struct {
	GameTypes  []string
	GameParams map[string]GameParams
}

func (conf *GameTypesConfig) Parse(data []byte) error {
	var tempConfig GameTypesConfig
	_, err := toml.Decode(string(data), &tempConfig)
	if err != nil {
		return err
	}
	//
	*conf = tempConfig
	return nil
}

////////////////////////////////////////////////////////
