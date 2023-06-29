package uuid

import (
	"fmt"
	"time"
)

/*
+-----------------------------------------------------------+
| 42 Bit Timestamp | 12 Bit WorkID | 3 Bit Time Sequence ID | 7 Bit Inc Sequence ID |
+-----------------------------------------------------------+
*/
const (
	//stepBits = subTagOffset + parentBits
	subTagOffset uint8  = 7                         // 子id自增位数-必须小于stepBits
	subIdMax     uint64 = -1 ^ (-1 << subTagOffset) // 子id序列号的最大值，用于检测溢出
	parentBits   uint8  = stepBits - subTagOffset   //父id自增位数
	parentMax    uint64 = -1 ^ (-1 << parentBits)   // 父id序列号的最大值，用于检测溢出
)

var intervalUUID *UUID

// initIntervalUUID 构造器
func initIntervalUUID(node uint64) error {
	// 如果超出节点的最大范围，产生一个 error
	if node > nodeMax {
		return fmt.Errorf("node number must be between 0 and %d", nodeMax)
	}
	// 生成并返回节点实例的指针
	intervalUUID = &UUID{
		timestamp: 0,
		node:      node,
		step:      0,
	}
	return nil
}

// GenerateInterval 生成间隔的uuid
// 注意GenerateInterval和Generate生成的id可能会一样,所以不要在全局uuid中混用
func GenerateInterval() uint64 {
	//
	intervalUUID.mu.Lock() // 保证并发安全, 加锁
	// 获取当前时间的时间戳 (毫秒数显示)
	now := time.Now().UnixMilli()
	//
	if intervalUUID.timestamp == now {
		//
		intervalUUID.step = (intervalUUID.step + 1) & parentMax
		// 当前 step 用完
		if intervalUUID.step == 0 {
			// 等待本毫秒结束
			for now <= intervalUUID.timestamp {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		// 本毫秒内 step 用完
		intervalUUID.step = 0
	}
	//
	intervalUUID.timestamp = now
	// 移位运算，生产最终 ID
	result := (uint64(now)-epoch)<<timeShift | (intervalUUID.node << nodeShift) | (intervalUUID.step << subTagOffset)
	//
	intervalUUID.mu.Unlock() // 方法运行完毕后解锁
	//
	return result
}

var AllUsedErr = fmt.Errorf("sub id all used")
var FormatErr = fmt.Errorf("uuid format err")

// GenerateIntervalSubId 根据uuid生成子id
// @uuid 父id
// @check 检查函数 返回true表示已用,false表示没用过。false会退出并返回这个值
func GenerateIntervalSubId(uuid uint64, check func(uint64) bool) (uint64, error) {
	v := uuid >> subTagOffset
	v = v << subTagOffset
	if v != uuid {
		return 0, FormatErr
	}
	for i := uint64(1); i <= subIdMax; i++ {
		n := uuid | i
		if !check(n) {
			return n, nil
		}
	}
	return 0, AllUsedErr
}

func ToIntervalBaseId(uuid uint64) uint64 {
	//前面64-RoomTagOffset是gid
	v := uuid >> subTagOffset
	v = v << subTagOffset
	return v
}

func ToIntervalIncId(uuid uint64) uint64 {
	return uuid - ToIntervalBaseId(uuid)
}
