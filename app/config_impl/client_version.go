package config_impl

import "github.com/BurntSushi/toml"

type LastVersion struct {
	Version string
	Title   string
	Detail  string
}

type ClientVersion struct {
	EnableCheck      bool
	LastVersion      LastVersion
	VersionAvailable map[string]bool
}

func (conf *ClientVersion) Parse(data []byte) error {
	_, err := toml.Decode(string(data), conf)
	return err
}
