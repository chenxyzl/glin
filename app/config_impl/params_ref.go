package config_impl

import (
	"math/rand"
	"time"
)

// GetHallMaxNextNotifyTime 最大发布间隔
func (conf *ParamsConfig) GetHallMaxNextNotifyTime(now time.Time) time.Time {
	return now.Add(conf.HallNotifyTimeInterval).Add(conf.HallNotifyTimeIntervalRange)
}

// GenHallNextNotifyTime 生成下一个发布时间
func (conf *ParamsConfig) GenHallNextNotifyTime(now time.Time) time.Time {
	now = now.Add(conf.HallNotifyTimeInterval)
	if conf.HallNotifyTimeIntervalRange > 0 {
		now = now.Add(time.Duration(rand.Int63()) % conf.HallNotifyTimeIntervalRange)
	}
	return now
}
