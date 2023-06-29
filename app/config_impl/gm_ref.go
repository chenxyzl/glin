package config_impl

import (
	"golang.org/x/exp/slices"
)

func (conf *GmConfig) IsOfficialGroup(gid uint64) bool {
	if len(conf.OfficialGroup) == 0 {
		return false
	}
	return slices.Contains(conf.OfficialGroup, gid)
}

func (conf *GmConfig) IsOpUid(uid uint64) bool {
	if len(conf.OpUids) == 0 {
		return false
	}
	return slices.Contains(conf.OpUids, uid)
}

func (conf *GmConfig) CheckGmAuth(uid uint64, gmApiName string) bool {
	if len(conf.GmUids) == 0 {
		return false
	}
	apis, ok := conf.GmUids[uid]
	if !ok {
		return false
	}
	return slices.Contains(apis, gmApiName)
}
