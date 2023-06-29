package wechat_model

import (
	"github.com/eatmoreapple/openwechat"
	"laiya/config"
	"laiya/share/global"
	"regexp"
	"strings"
)

// gameTypeReg 游戏类型正则获取
var gameTypeReg = regexp.MustCompile(`【(.*?)】`)

func (mod *WechatGroup) GetGameType() string {
	return GetGameType(mod.OpenWechatGroup)
}

func GetGameType(wechatGroup *openwechat.Group) string {
	if wechatGroup == nil {
		return global.DefaultGameType
	}
	//昵称和重命名的任意一个包含都可,优先获取重命名中的
	name := wechatGroup.RemarkName + wechatGroup.NickName + wechatGroup.DisplayName
	if len(name) == 0 {
		return global.DefaultGameType
	}
	result := gameTypeReg.FindStringSubmatch(name)
	if len(result) <= 1 {
		return global.DefaultGameType
	}
	str := result[1]
	gameTypes := config.Get().GameTypesConfig.GameTypes
	for _, gameType := range gameTypes {
		if strings.Contains(gameType, str) {
			return gameType
		}
	}
	return global.DefaultGameType
}
