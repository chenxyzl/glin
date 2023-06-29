package group_model

import (
	"fmt"
	"laiya/share/algorithm/skiplist"
	"strings"
	"testing"
)

func TestSkip(t *testing.T) {

	a := skiplist.New[*OnlineSortT](&cmpOnline{})
	a.Insert(&OnlineSortT{
		uid: 1,
	})
	a.Insert(&OnlineSortT{
		uid: 2,
	})
	a.Insert(&OnlineSortT{
		uid: 3,
	})
	fmt.Println("range:")
	a.ForRange(func(sortT *OnlineSortT) bool {
		fmt.Println(sortT)
		return true
	})
	fmt.Println("range:reverse")
	a.ForRange(func(sortT *OnlineSortT) bool {
		fmt.Println(sortT)
		return true
	}, true)
	fmt.Println("xxxxx")
	b := a.GetRange(1, 3, true)
	for _, v := range b {
		fmt.Println(v)
	}
}
func TestOnlineSortT_Score(t *testing.T) {
	a := (uint64(1) << 33) | 1<<32 | uint64(0)
	b := (uint64(1) << 33) | 0<<32 | uint64(0)
	c := (uint64(0) << 33) | 1<<32 | uint64(0)
	fmt.Println(a, b, c)
}

func TestClean(t *testing.T) {
	var num uint64 = 1 // 假设这是您的64位无符号整型数
	// 使用位掩码保留右边的33位
	mask := uint64(1<<33 - 1)
	result := num & mask

	fmt.Printf("原始值: %b\n", num)
	fmt.Printf("处理后的结果: %b\n", result)
}

func TestStrContain(t *testing.T) {
	a := ""
	b := "aaa"
	c := strings.Contains(b, a)
	println(c)
}
