package wechat_model

import (
	"github.com/eatmoreapple/openwechat"
	"laiya/config"
	"testing"
)

func TestName(t *testing.T) {
	config.Get().GameTypesConfig.GameTypes = []string{"xy", "aa", "y"}
	wechatGroup := &openwechat.Group{User: &openwechat.User{NickName: "c【a】b【d】..", RemarkName: "x【y】z【w】.."}}
	name := GetGameType(wechatGroup)
	if name != "xy" {
		t.Error()
	}
}
