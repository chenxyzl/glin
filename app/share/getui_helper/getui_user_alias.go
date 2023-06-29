package getui_helper

import (
	"context"
	"fmt"
	"github.com/chenxyzl/glin/slog"
	"github.com/dacker-soul/getui/user"
)

func BindUidCid(uid uint64, cid string) error {
	//
	params := user.BindAliasParams{
		DataList: []user.BindAliasParam{{
			Cid:   cid,
			Alias: fmt.Sprintf("%d", uid),
		}},
	}
	token, err := getToken()
	if err != nil {
		return fmt.Errorf("get token err, uid:%v|cid:%v|err:%v", uid, cid, err)
	}
	//
	result, err := user.BindUserAlias(context.Background(), getConfig(), token, &params)
	if err != nil {
		return fmt.Errorf("first bind alias by cid err, uid:%v|cid:%v|err:%v", uid, cid, err)
	}
	//发送成功
	if result.Code == 0 {
		slog.Infof("first bind alias by cid success, uid:%v|cid:%v", uid, cid)
		return nil
	}
	//结果检查
	if result.Code != 10001 {
		//其他错误
		return fmt.Errorf("first bind alias by cid err, uid:%v|cid:%v|result.Code:%v|Msg:%v", uid, cid, result.Code, result.Msg)

	}
	//只有token过期一种可能了--刷新token重试
	token, err = getToken(true)
	if err != nil {
		return fmt.Errorf("refresh token err, uid:%v|cid:%v|err:%v", uid, cid, err)
	}
	//用新token继续执行
	result, err = user.BindUserAlias(context.Background(), getConfig(), token, &params)
	if err != nil {
		return fmt.Errorf("second bind alias by cid err, uid:%v|cid:%v|err:%v", uid, cid, err)
	}
	//结果检查
	if result.Code != 0 {
		//其他错误
		return fmt.Errorf("second bind alias by cid err, uid:%v|cid:%v|result.Code:%v|Msg:%v", uid, cid, result.Code, result.Msg)
	}
	slog.Infof("bind alias success, uid:%v|cid:%v", uid, cid)
	return nil
}

func DisBindUidCid(uid uint64, cid string) error {
	//
	params := user.DisBindAliasParams{
		DataList: []user.DisBindAliasParam{{
			Cid:   cid,
			Alias: fmt.Sprintf("%d", uid),
		}},
	}
	token, err := getToken()
	if err != nil {
		return fmt.Errorf("get token err, err:%v", err)
	}
	//
	result, err := user.DisBindUserAlias(context.Background(), getConfig(), token, &params)
	if err != nil {
		return fmt.Errorf("first disbind alias by cid err, uid:%v|cid:%v|err:%v", uid, cid, err)
	}
	//发送成功
	if result.Code == 0 {
		slog.Infof("first disbind alias by cid success, uid:%v|cid:%v", uid, cid)
		return nil
	}
	//结果检查
	if result.Code != 10001 {
		//其他错误
		return fmt.Errorf("first disbind alias by cid err, uid:%v|cid:%v|result.Code:%v|Msg:%v", uid, cid, result.Code, result.Msg)

	}
	//只有token过期一种可能了--刷新token重试
	token, err = getToken(true)
	if err != nil {
		return fmt.Errorf("refresh token err, uid:%v|cid:%v|err:%v", uid, cid, err)
	}
	//用新token继续执行
	result, err = user.DisBindUserAlias(context.Background(), getConfig(), token, &params)
	if err != nil {
		return fmt.Errorf("second disbind alias by cid err, uid:%v|cid:%v|err:%v", uid, cid, err)
	}
	//结果检查
	if result.Code != 0 {
		//其他错误
		return fmt.Errorf("second disbind alias by cid err, uid:%v|cid:%v|result.Code:%v|Msg:%v", uid, cid, result.Code, result.Msg)
	}
	slog.Infof("disbind alias success, uid:%v|cid:%v", uid, cid)
	return nil
}
