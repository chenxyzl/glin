package getui_helper

import (
	"context"
	"fmt"
	"github.com/chenxyzl/glin/slog"
	"github.com/dacker-soul/getui/auth"
	"github.com/dacker-soul/getui/publics"
	"laiya/config"
	"sync"
)

var getTuiTokenLockIns = sync.Mutex{}
var geTuiToken string

func getConfig() publics.GeTuiConfig {
	return publics.GeTuiConfig{
		AppId:        config.Get().WebConfig.GeTuiAppId, // 个推提供的id，密码等
		AppSecret:    config.Get().WebConfig.GeTuiAppSecret,
		AppKey:       config.Get().WebConfig.GeTuiAppKey,
		MasterSecret: config.Get().WebConfig.GeTuiMasterSecret,
	}
}
func getSettings() *publics.Settings {
	return config.Get().WebConfig.GeTuiSettings
}
func getNotification(title, body string) *publics.Notification {
	return &publics.Notification{
		Title:       title,
		Body:        body,
		ClickType:   "startapp", // 打开应用首页
		BadgeAddNum: 1,
	}
}
func getIosChannel(title, body string) *publics.IosChannel {
	return &publics.IosChannel{
		Type: "",
		Aps: &publics.Aps{
			Alert: &publics.Alert{
				Title: title,
				Body:  body,
			},
			ContentAvailable: 0,
		},
		AutoBadge:      "+1",
		PayLoad:        "",
		Multimedia:     nil,
		ApnsCollapseId: "",
	}
}
func getPushMessage(title, body string) *publics.PushMessage {
	return &publics.PushMessage{
		Duration:     "", // 手机端通知展示时间段
		Notification: getNotification(title, body),
		Transmission: "",
		Revoke:       nil,
	}
}
func getPushChannel(title, body string) *publics.PushChannel {
	return &publics.PushChannel{
		Ios: getIosChannel(title, body),
		Android: &publics.AndroidChannel{Ups: &publics.Ups{
			Notification: getNotification(title, body),
			TransMission: "", // 透传消息内容，与notification 二选一
		}},
	}
}

func getToken(refresh ...bool) (string, error) {
	isRefresh := false
	if len(refresh) > 0 {
		isRefresh = refresh[0]
	}
	if geTuiToken == "" || isRefresh {
		getTuiTokenLockIns.Lock()
		defer getTuiTokenLockIns.Unlock()
		//强制刷新
		if isRefresh {
			geTuiToken = ""
		}
		//生成新的
		if geTuiToken == "" {
			result, err := auth.GetToken(context.Background(), getConfig())
			if err != nil {
				return "", err
			}
			if result.Code != 0 {
				return "", fmt.Errorf("get geTuiToken err, result.Code:%v|Msg:%v", result.Code, result.Msg)
			}
			geTuiToken = result.Data.Token
		}
		slog.Infof("gen new geTuiToken: %v", geTuiToken)
	}
	return geTuiToken, nil
}
