package config_impl

import (
	"github.com/BurntSushi/toml"
)

type WechatConfig struct {
	KeepAliveNickName string
}

func (conf *WechatConfig) Parse(data []byte) error {
	_, err := toml.Decode(string(data), conf)
	return err
}
