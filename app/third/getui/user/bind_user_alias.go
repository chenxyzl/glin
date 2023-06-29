package user

import (
	"context"
	"encoding/json"
	"github.com/dacker-soul/getui/publics"
)

type BindAliasParam struct {
	Cid   string `json:"cid"`   //是否必须:是	默认值:无	描述:cid，用户标识
	Alias string `json:"alias"` //是否必须:是	默认值:无	描述:别名，有效的别名组成。字母（区分大小写）、数字、下划线、汉字;长度<40字; 一个别名最多允许绑定10个cid。
}

// BindAliasParams 绑定别名参数
type BindAliasParams struct {
	DataList []BindAliasParam `json:"data_list"` //是否必须:是	默认值:无	描述:数据列表，数组长度不大于1000
}

// BindAliasResult 绑定别名返回
type BindAliasResult struct {
	publics.PublicResult
}

// BindUserAlias 绑定别名
func BindUserAlias(ctx context.Context, config publics.GeTuiConfig, token string, param *BindAliasParams) (*BindAliasResult, error) {

	url := publics.ApiUrl + config.AppId + "/user/alias"
	bodyByte, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}

	result, err := publics.RestFulRequest(ctx, bodyByte, url, "POST", token)
	if err != nil {
		return nil, err
	}

	var push *BindAliasResult
	if err := json.Unmarshal([]byte(result), &push); err != nil {
		return nil, err
	}

	return push, err
}
