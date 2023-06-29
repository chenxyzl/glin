package getui_helper

import (
	"context"
	"fmt"
	"github.com/chenxyzl/glin/slog"
	"github.com/dacker-soul/getui/publics"
	"github.com/dacker-soul/getui/push/single"
	"strconv"
	"time"
)

func PushMsg(uids []uint64, title, body string) {
	if len(uids) == 0 {
		slog.Warningf("push uids count 0, ignore")
		return
	}
	//if len(uids) == 1 {
	//	go PushToSingleUid(uids[0], title, body)
	//} else {
	go PushToMultiUids(uids, title, body)
	slog.Infof("getui push to, uids:%v|title:%v|body:%v", uids, title, body)
	//}
}

func PushToSingleUid(uid uint64, title, body string) {
	alias := fmt.Sprintf("%d", uid)
	//
	params := single.PushSingleAliasParam{
		RequestId:   strconv.FormatInt(time.Now().UnixNano(), 10), // 请求唯一标识号
		Audience:    &publics.Audience{Alias: []string{alias}},
		Settings:    getSettings(),
		PushMessage: getPushMessage(title, body),
		PushChannel: getPushChannel(title, body),
	}
	slog.Infof("getui push single params, uid:%v|params:%v", uid, params)
	token, err := getToken()
	if err != nil {
		slog.Errorf("get token err, err:%v|alias:%v", err, alias)
		return
	}
	//执行单推
	result, err := single.PushSingleByAlias(context.Background(), getConfig(), token, &params)
	if err != nil {
		slog.Errorf("first push single by alias err, err:%v|alias:%v", err, alias)
		return
	}
	//发送成功
	if result.Code == 0 {
		slog.Infof("first push single by alias success,alias:%v", alias)
		return
	}
	//结果检查
	if result.Code != 10001 {
		//其他错误
		slog.Errorf("first push single by alias err, result.Code:%v|Msg:%v|alias:%v", result.Code, result.Msg, alias)
		return
	}
	//只有token过期一种可能了--刷新token重试
	token, err = getToken(true)
	if err != nil {
		slog.Errorf("refresh token err, err:%v", err)
		return
	}
	//用新token继续执行单推
	result, err = single.PushSingleByAlias(context.Background(), getConfig(), token, &params)
	if err != nil {
		slog.Errorf("second push single by alias err, err:%v|alias:%v", err, alias)
		return
	}
	//结果检查
	if result.Code != 0 {
		//其他错误
		slog.Errorf("second push single by alias err, result.Code:%v|Msg:%v|alias:%v", result.Code, result.Msg, alias)
		return
	}
	slog.Infof("push single success, alias:%v", alias)
}

func PushToMultiUids(uids []uint64, title, body string) {
	//if len(uids) < 2 {
	//	slog.Warningf("push uids count little than 2, ignore, len:%v", len(uids))
	//	return
	//}
	var aliasList []string
	for _, uid := range uids {
		aliasList = append(aliasList, fmt.Sprintf("%d", uid))
	}
	settings := getSettings()
	pushMessage := getPushMessage(title, body)
	pushChannel := getPushChannel(title, body)
	request := time.Now().UnixNano()
	//
	for len(aliasList) > 0 {
		//接口单次允许最大批处理个数位200
		var sendAliasList []string
		if len(aliasList) > 200 {
			sendAliasList = aliasList[:200]
			aliasList = aliasList[200:]
		} else {
			sendAliasList = aliasList
			aliasList = nil
		}
		//构造发送的批次数据
		params := single.PushSingleBatchAliasParam{
			IsAsync: false,
			MsgList: make([]*single.PushSingleAliasParam, 0),
		}
		for _, alias := range sendAliasList {
			request++
			param := &single.PushSingleAliasParam{RequestId: strconv.FormatInt(request, 10), // 请求唯一标识号
				Audience: &publics.Audience{ // 目标用户
					Alias: []string{alias},
				},
				Settings:    settings,
				PushMessage: pushMessage,
				PushChannel: pushChannel,
			}
			params.MsgList = append(params.MsgList, param)
		}
		slog.Infof("getui push batch params, uid:%v|params:%v", sendAliasList, params)
		//获取token
		token, err := getToken()
		if err != nil {
			slog.Errorf("get token err, err:%v|sendAliasList:%v", err, sendAliasList)
			continue
		}
		//执行批量推送
		result, err := single.PushSingleByBatchAlias(context.Background(), getConfig(), token, &params)
		if err != nil {
			slog.Errorf("first push batch by alias err, err:%v|sendAliasList:%v", err, sendAliasList)
			continue
		}
		//发送成功
		if result.Code == 0 {
			slog.Infof("first push batch by alias success, sendAliasList:%v|result:%v", sendAliasList, result)
			continue
		}
		//结果检查
		if result.Code != 10001 {
			//其他错误
			slog.Errorf("first push batch by alias err, result.Code:%v|Msg:%v|sendAliasList:%v|", result.Code, result.Msg, sendAliasList)
			continue
		}
		//只有token过期一种可能了--刷新token重试
		token, err = getToken(true)
		if err != nil {
			slog.Errorf("refresh token err, err:%vsendAliasList:%v", err, sendAliasList)
			continue
		}
		//用新token继续执行第二次推送
		result, err = single.PushSingleByBatchAlias(context.Background(), getConfig(), token, &params)
		if err != nil {
			slog.Errorf("second push batch by alias err, err:%v|sendAliasList:%v", err, sendAliasList)
			continue
		}
		//结果检查
		if result.Code != 0 {
			//其他错误
			slog.Errorf("second push batch by alias err, result.Code:%v|Msg:%v|sendAliasList:%v|", result.Code, result.Msg, sendAliasList)
			continue
		}
		slog.Infof("push batch success, sendAliasList:%v|result:%v", sendAliasList, result)
	}
}

func TestGeTuiToken() error {
	_, err := getToken()
	return err
}
