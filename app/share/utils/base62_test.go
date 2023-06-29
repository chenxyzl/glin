package utils

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func BenchmarkUint64ToHashedBase26String(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Uint64ToHashedBase52(uint64(i))
	}
}

func TestA(t *testing.T) {
	//check hash
	maxUint64 := uint64(1660118221883179265)
	base52 := Uint64ToHashedBase52(maxUint64)
	fmt.Printf("The max uint64:%v value in hash is: %s\n", maxUint64, base52)
	fmt.Printf("The value:%v in hash is: %s\n", 42, Uint64ToHashedBase52(42))
	fmt.Printf("The value:%v in hash is: %s\n", 43, Uint64ToHashedBase52(43))
	//check base62
	base62 := Uint64ToBase62(maxUint64)
	u1, err := Base62ToUint64(base62)
	if err != nil {
		t.Error(err)
	}
	if u1 != maxUint64 {
		t.Error()
	}
	fmt.Println(base62, u1)
	//uid parse
	url := "https://www.laiyaplaygame.com/laiya/groups/1660118221883179264/1660118221883179265"
	tag := "groups/"
	idx := strings.Index(url, tag)
	if idx < 0 {
		t.Error()
	}
	sub := url[idx+len(tag):]
	arr := strings.Split(sub, "/")
	if len(arr) != 2 {
		t.Error()
	}
	roomId, err := strconv.ParseUint(arr[1], 10, 64)
	if err != nil {
		t.Error(err)
	}
	if roomId != 1660118221883179265 {
		t.Error()
	}
}
