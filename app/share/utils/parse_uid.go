package utils

import (
	"fmt"
	"github.com/asynkron/protoactor-go/actor"
	"strconv"
	"strings"
)

func substringBetween(s, start, end string) (string, error) {
	startIndex := strings.Index(s, start)
	if startIndex == -1 {
		return "", fmt.Errorf("start string not found")
	}

	endIndex := strings.Index(s, end)
	if endIndex == -1 {
		return "", fmt.Errorf("end string not found")
	}

	if endIndex <= startIndex {
		return "", fmt.Errorf("end index is less than or equal to start index")
	}

	return s[startIndex+1 : endIndex], nil
}

func ParseUid(pid *actor.PID) (uint64, error) {
	id := pid.GetId()
	var err error
	id, err = substringBetween(id, "/", "$")
	if err != nil {
		return 0, err
	}
	uid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return 0, err
	}
	return uid, nil
}
