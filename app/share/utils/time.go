package utils

import "time"

// GetBeginningOfDay1 获取零点时间
// @t 时间
func GetBeginningOfDay1(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

// GetBeginningOfDay 获取零点时间
// @millTime 时间戳 毫秒
func GetBeginningOfDay(millTime int64) time.Time {
	t := time.UnixMilli(millTime)
	return GetBeginningOfDay1(t)
}

// CheckInSameDiyDay 是否同一天
// @millTime1 时间戳 毫秒
// @millTime2 时间戳 毫秒
// @splitHour:@splitMin：@splitSec 分割时间为每天的 时:分:秒
func CheckInSameDiyDay(millTime1 int64, millTime2 int64, splitHour int, splitMin int, splitSec int) bool {
	splitHourDt := time.Duration(splitHour)
	splitMinDt := time.Duration(splitMin)
	splitSecDt := time.Duration(splitSec)
	t1 := time.UnixMilli(millTime1)
	t2 := time.UnixMilli(millTime2)
	t1 = t1.Add(-splitHourDt * time.Hour).Add(-splitMinDt * time.Minute).Add(-splitSecDt * time.Second)
	t2 = t2.Add(-splitHourDt * time.Hour).Add(-splitMinDt * time.Minute).Add(-splitSecDt * time.Second)
	bt1 := GetBeginningOfDay1(t1)
	bt2 := GetBeginningOfDay1(t2)
	return bt1.UnixMilli() == bt2.UnixMilli()
}

// CheckInSameDay 是否同一天
// @millTime1 时间戳 毫秒
// @millTime2 时间戳 毫秒
func CheckInSameDay(millTime1 int64, millTime2 int64) bool {
	return CheckInSameDiyDay(millTime1, millTime2, 0, 0, 0)
}
