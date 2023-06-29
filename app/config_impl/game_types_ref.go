package config_impl

import "laiya/share/global"

func (conf *GameTypesConfig) GetParams(gameType string) GameParams {
	config, ok := conf.GameParams[gameType]
	if ok {
		return config
	}
	//找不到文案配置就走默认的
	return conf.GameParams[global.DefaultGameType]
}
