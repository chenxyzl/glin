package crypt

import (
	"testing"
)

const cryptKey = "AES256Key-32Characters1234567890"

func TestA(t *testing.T) {
	pwd := "fakdjfhlakhfpoie"
	p, err := Encrypt(pwd, cryptKey)
	if err != nil {
		t.Error(err)
	}
	p1, err := Encrypt(pwd, cryptKey)
	if err != nil {
		t.Error(err)
	}
	pwd1, err := Decrypt(p, cryptKey)
	if err != nil {
		t.Error(err)
	}
	if p != p1 {
		//t.Error()
	}
	if pwd != pwd1 {
		t.Error()
	}
}
