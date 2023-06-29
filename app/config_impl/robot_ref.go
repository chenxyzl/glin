package config_impl

func (conf *RobotConfig) GetRobotGameType(uid uint64) string {
	if len(conf.Robots) == 0 {
		return ""
	}
	return conf.Robots[uid]
}
