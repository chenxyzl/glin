package datarangers

import (
	"github.com/chenxyzl/glin/slog"
	"gopkg.in/yaml.v2"
	"laiya/share/version"
	"os"
	"strconv"

	// 引入 sdk
	sdk "github.com/volcengine/datarangers-sdk-go"
	"time"
)

var confIns *sdk.SysConf

var (
	appId int64 = 0
	uuid        = ""
)

func init() {
	// 初始化, YAML_PATH 替换成真实的配置文件
	confIns = &sdk.SysConf{}
	yamlFile, err := os.ReadFile("../config_impl/datarangers.yaml")
	if err != nil {
		slog.Fatalf("init config fail->  : " + err.Error())
		return
	}
	if err = yaml.Unmarshal(yamlFile, confIns); err != nil {
		slog.Fatalf("init config fail->  : " + err.Error())
		return
	}
	if len(confIns.AppKeys) != 1 {
		slog.Fatalf("datarangers.yaml not set app keys")
		return
	}
	//init appId
	for k := range confIns.AppKeys {
		appId = k
		uuid = strconv.Itoa(int(appId))
	}
	err = sdk.InitBySysConf(confIns)
	//err = sdk.InitByFile("../config/datarangers.yaml")
	if err != nil {
		slog.Fatalf("init datarangers err, err:%v", err)
	}
}

func SendEvent(eventName string, eventPrams map[string]interface{}, commonPrams map[string]interface{}) {
	// 事件公共属性，如果不需要的化，可以传nil
	//commonParams := map[string]interface{}{
	//	"common_string1": "common_value1",
	//	"common_int1":    13,
	//}
	//eventName := "app_send_event_info"
	// 事件发生时间
	localTimeMs := time.Now().UnixMilli()
	if version.Branch == "dev" {
		eventName += "_debug"
	}
	eventV3 := &sdk.EventV3{
		Event:       eventName,
		LocalTimeMs: &localTimeMs,
		Params:      eventPrams,
	}
	sdk.SendEventInfo(sdk.APP, appId, uuid, eventV3, commonPrams)
}
