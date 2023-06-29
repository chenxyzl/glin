package utils

import (
	"fmt"
	"github.com/twmb/murmur3"
	"strings"
)

const hashBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const hashBytesLen = len(hashBytes)

func Uint64ToHashedBase52(n uint64) string {
	hasher := murmur3.New64()
	hasher.Write([]byte(fmt.Sprint(n)))
	u := hasher.Sum64()
	if u == 0 {
		return string(hashBytes[0])
	}

	result := ""
	for u > 0 {
		remainder := u % uint64(hashBytesLen)
		result = string(hashBytes[remainder]) + result
		u = u / uint64(hashBytesLen)
	}
	return result
}

const Base62 = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const Base62Len = len(Base62)

func Uint64ToBase62(num uint64) string {
	var result string
	for num > 0 {
		result = string(Base62[num%uint64(Base62Len)]) + result
		num = num / uint64(Base62Len)
	}
	return result
}

func Base62ToUint64(str string) (uint64, error) {
	var result uint64
	for i := 0; i < len(str); i++ {
		idx := strings.IndexByte(Base62, str[i])
		if idx == -1 {
			return 0, fmt.Errorf("invalid character in input string: %c", str[i])
		}
		result = result*62 + uint64(idx)
	}
	return result, nil
}
