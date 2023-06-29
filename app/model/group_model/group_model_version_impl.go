package group_model

import (
	"github.com/chenxyzl/glin/slog"
	"laiya/proto/common"
)

var upgradeVersion upgradeVersionT = []upgradeVersionFunc{
	upgradeVersion1,
	upgradeVersion2,
}

// 增加群加入权限
func upgradeVersion1(mod *Group, logger slog.Logger) {
	mod.JoinAccess = true
	logger.Infof("db upgrade, add joinAccess")
}

// 增加普通用户权限
func upgradeVersion2(mod *Group, logger slog.Logger) {
	for _, player := range mod.Players {
		player.Permissions = int32(common.GroupPermissionsDef_User)
	}
	logger.Infof("db upgrade, add player permission")
}
