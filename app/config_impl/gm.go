package config_impl

import (
	"github.com/BurntSushi/toml"
	"strconv"
)

type GmConfig struct {
	//管理员账户
	Gods          []string
	OfficialGroup []uint64
	OpUids        []uint64
	GmUidsStr     map[string][]string
	GmUids        map[uint64][]string
}

func (conf *GmConfig) Parse(data []byte) error {
	var tempConfig GmConfig
	_, err := toml.Decode(string(data), &tempConfig)
	if err != nil {
		return err
	}
	//GmUids: string 2 uint64
	tempConfig.GmUids = make(map[uint64][]string)
	for s, s2 := range tempConfig.GmUidsStr {
		i, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		tempConfig.GmUids[uint64(i)] = s2
	}
	*conf = tempConfig
	return err
}
