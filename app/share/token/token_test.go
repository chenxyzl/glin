package token

import (
	"testing"
	"time"
)

func TestBuildToken(t *testing.T) {
	var uid uint64 = 88859854976144384
	var exp = 90 * 24 * time.Hour
	//exp = time.Second * 10
	token, e := BuildToken(uid, exp, "laiyagate@1692449108")
	if e != nil {
		t.Error(e)
	}
	nid, e := ParseToken(token, "laiyagate@1692449108")
	if e != nil {
		t.Error(e)
	}
	if nid != uid {
		t.Error()
	}
}
