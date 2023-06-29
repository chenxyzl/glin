package component

import (
	"laiya/config"
	"laiya/proto/code"
)

func checkActivityTitle(str string) code.Code {
	l := len([]rune(str))
	conf := config.Get().ConstConfig
	if l < conf.GroupPluginActivityTitleMinLen || l > conf.GroupPluginActivityTitleMaxLen {
		return code.Code_GroupActivityTitleLenIllegal
	}
	return code.Code_Ok
}

func checkActivityDes(str string) code.Code {
	l := len([]rune(str))
	conf := config.Get().ConstConfig
	if l > conf.GroupPluginActivityDesMaxLen {
		return code.Code_GroupActivityDesLenIllegal
	}
	return code.Code_Ok
}

func checkActivitySignUpRemark(str string) code.Code {
	l := len([]rune(str))
	conf := config.Get().ConstConfig
	if l > conf.GroupPluginActivitySignUpRemarkMaxLen {
		return code.Code_GroupActivitySignUpRewardLenIllegal
	}
	return code.Code_Ok
}
