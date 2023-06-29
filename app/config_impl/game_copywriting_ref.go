package config_impl

import "laiya/share/global"

func (conf *GameCopywritingConfig) Get(gameType string) GameCopywriting {
	config, ok := conf.Copywriting[gameType]
	if ok {
		return config
	}
	//找不到文案配置就走默认的
	return conf.Copywriting[global.DefaultGameType]
}
