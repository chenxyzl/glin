package config_impl

import (
	"github.com/BurntSushi/toml"
)

type Banner struct {
	Name string
	Des  string
	Icon string
	Link string
}

type HallBannersConfig struct {
	Banners []Banner
}

func (conf *HallBannersConfig) Parse(data []byte) error {
	var tempConfig HallBannersConfig
	_, err := toml.Decode(string(data), &tempConfig)
	if err != nil {
		return err
	}
	*conf = tempConfig
	return nil
}

////////////////////////////////////////////////////////
