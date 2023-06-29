package group_model

import (
	"github.com/chenxyzl/glin/slog"
)

type upgradeVersionFunc func(*Group, slog.Logger)
type upgradeVersionT []upgradeVersionFunc

// 获取版本对应的升级函数
func (x *upgradeVersionT) getVersionFunc(version uint64) upgradeVersionFunc {
	return upgradeVersion[version-1]
}

// 获取当前版本
func (x *upgradeVersionT) getNowVersion() uint64 {
	return uint64(len(upgradeVersion))
}

// CheckVersion 版本升级检查
// @return 上个版本
// @return 是否有版本变化
func (mod *Group) CheckVersion(logger slog.Logger) {
	for version := mod.Version + 1; version <= upgradeVersion.getNowVersion(); version++ {
		slog.Infof("db need upgrade, oldVersion:%v|newVersion:%d", mod.Version, version)
		//升级
		upgradeVersion.getVersionFunc(version)(mod, logger)
		//设置当前版本
		mod.Version = version
		//脏标记标记
		mod.MarkDirty()
		//
		slog.Infof("db upgrade success, newVersion:%d", mod.Version)
	}
}
