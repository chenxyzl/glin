package config_impl

import (
	"golang.org/x/exp/slices"
	"laiya/proto/common"
)

func (conf *GroupPermissionsConfig) HasAccess(userPermissions common.GroupPermissionsDef_Type, access common.GroupPermissionsDef_Access) bool {
	accessWrapper, ok := conf.GroupPermissions[userPermissions.String()]
	if !ok {
		return false
	}
	return slices.Contains(accessWrapper.Access, int32(access.Number()))
}
