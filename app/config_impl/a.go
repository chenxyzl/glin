package config_impl

import (
	"github.com/BurntSushi/toml"
	"time"
)

type (
	Livekit struct {
		Host      string
		ApiKey    string
		ApiSecret string
	}
	WsParam struct {
		Host         string
		Port         string
		Path         string
		WriteTimeout time.Duration
	}
	HttpParam struct {
		Host string
		Port string
	}
	Mongo struct {
		Url    string
		DbName string
	}

	// AConfig 基础公共配置
	AConfig struct {
		AppKey          string
		TokenExpireTime time.Duration
		PprofPort       map[string]int
		Livekit         Livekit
		WebHttpParam    HttpParam
		WsParam         WsParam
		Mongo           Mongo
		OldMongo        Mongo
	}
)

func (conf *AConfig) Parse(data []byte) error {
	_, err := toml.Decode(string(data), conf)
	return err
}
