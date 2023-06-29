package config_impl

// CheckClientVersion 检查客户端版本号是否可用
func (conf *ClientVersion) CheckClientVersion(clientVersion string) bool {
	//配置未生效之前,认为所有版本都可用
	if !conf.EnableCheck {
		return true
	}
	//
	return conf.VersionAvailable[clientVersion]
}
