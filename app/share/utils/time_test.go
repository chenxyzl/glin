package utils

import (
	"testing"
	"time"
)

func TestCheckInSameDiyDay(t *testing.T) {
	tim := time.Now()
	now := time.Now().UnixMilli()
	next := time.Now().Add(2 * time.Minute).UnixMilli()
	b := CheckInSameDiyDay(now, next, tim.Hour(), tim.Minute()+1, tim.Second())
	if b {
		t.Error("must not same")
	}

	now = time.Now().UnixMilli()
	next = time.Now().Add(time.Minute).UnixMilli()
	b = CheckInSameDiyDay(now, next, tim.Hour(), tim.Minute(), tim.Second()*1)
	if !b {
		t.Error("must same")
	}
}
