package token

import (
	"fmt"
	"github.com/form3tech-oss/jwt-go"
	"time"
)

type MyCustomClaims struct {
	UID uint64 `json:"uid"`
	jwt.StandardClaims
}

func BuildToken(uid uint64, exp time.Duration, appKey string) (string, error) {
	// 一个24小时有效的token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, MyCustomClaims{
		UID: uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(exp).Unix(),
		},
	})
	tokenString, err := token.SignedString([]byte(appKey))
	if err != nil {

		return "", fmt.Errorf("signed token error; err:%v", err)
	}
	return tokenString, nil
}

func ParseToken(token string, appKey string) (uint64, error) {
	a := MyCustomClaims{}
	_, err := jwt.ParseWithClaims(token, &a, func(token *jwt.Token) (interface{}, error) {
		return []byte(appKey), nil
	})
	if err != nil {
		return 0, fmt.Errorf("parse token error; err:%v", err)
	}

	if a.UID == 0 {
		return 0, fmt.Errorf("not found uid or uid is 0, map:%v", a)
	}
	return a.UID, nil
}
