package config

import (
	"github.com/chenxyzl/glin/slog"
	"laiya/config_impl"
)

type Config struct {
	AConfig                   config_impl.AConfig                   `key:"a.toml" default:"../config_impl/a.toml"`
	ClientVersion             config_impl.ClientVersion             `key:"client_version.toml" default:"../config_impl/client_version.toml"`
	ClientDynamic             config_impl.ClientDynamic             `key:"client_dynamic.toml" default:"../config_impl/client_dynamic.toml"`
	WebConfig                 config_impl.WebConfig                 `key:"web.toml" default:"../config_impl/web.toml"`
	ConstConfig               config_impl.ConstConfig               `key:"const.toml" default:"../config_impl/const.toml"`
	ErrorConfig               config_impl.ErrorConfig               `key:"error.toml" default:"../config_impl/error.toml"`
	GmConfig                  config_impl.GmConfig                  `key:"gm.toml" default:"../config_impl/gm.toml"`
	ParamsConfig              config_impl.ParamsConfig              `key:"params.toml" default:"../config_impl/params.toml"`
	GameTypesConfig           config_impl.GameTypesConfig           `key:"game_types.toml" default:"../config_impl/game_types.toml"`
	GamePlatformsConfig       config_impl.GamePlatformsConfig       `key:"game_platforms.toml" default:"../config_impl/game_platforms.toml"`
	GroupPermissionsConfig    config_impl.GroupPermissionsConfig    `key:"group_permissions.toml" default:"../config_impl/group_permissions.toml"`
	HallBannersConfig         config_impl.HallBannersConfig         `key:"hall_banners.toml" default:"../config_impl/hall_banners.toml"`
	RobotConfig               config_impl.RobotConfig               `key:"robot.toml" default:"../config_impl/robot.toml"`
	GameCopywritingConfig     config_impl.GameCopywritingConfig     `key:"game_copywriting.toml" default:"../config_impl/game_copywriting.toml"`
	WechatConfig              config_impl.WechatConfig              `key:"wechat.toml" default:"../config_impl/wechat.toml"`
	GroupPluginActivityConfig config_impl.GroupPluginActivityConfig `key:"group_plugin_activity.toml" default:"../config_impl/group_plugin_activity.toml"`
}

var conf = &Config{}

func Get() *Config {
	return conf
}

func (conf *Config) AfterLoadAll() {
	slog.Infof("load config success")
}
