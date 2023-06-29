package user

import (
	"context"
	"encoding/json"
	"github.com/dacker-soul/getui/publics"
)

type DisBindAliasParam struct {
	Cid   string `json:"cid"`   //是否必须:是	默认值:无	描述:用户标识
	Alias string `json:"alias"` //是否必须:是	默认值:无	描述:别名
}

// DisBindAliasParams 解绑别名参数
type DisBindAliasParams struct {
	DataList []DisBindAliasParam `json:"data_list"` //是否必须:是	默认值:无	描述:数据列表，数组长度不大于1000
}

// DisBindAliasResult 解绑别名返回
type DisBindAliasResult struct {
	publics.PublicResult
}

// DisBindUserAlias 解绑别名
func DisBindUserAlias(ctx context.Context, config publics.GeTuiConfig, token string, param *DisBindAliasParams) (*DisBindAliasResult, error) {

	url := publics.ApiUrl + config.AppId + "/user/alias"
	bodyByte, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}

	result, err := publics.RestFulRequest(ctx, bodyByte, url, "DELETE", token)
	if err != nil {
		return nil, err
	}

	var push *DisBindAliasResult
	if err := json.Unmarshal([]byte(result), &push); err != nil {
		return nil, err
	}

	return push, err
}
