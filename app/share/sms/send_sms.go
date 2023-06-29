package sms

import (
	"encoding/json"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/chenxyzl/glin/slog"
	"laiya/config"
	"laiya/proto/code"
)

func createClient() (_result *dysmsapi20170525.Client, _err error) {
	conf := &openapi.Config{
		// 必填，您的 AccessKey ID
		AccessKeyId: tea.String(config.Get().WebConfig.SmsAccessKeyId),
		// 必填，您的 AccessKey Secret
		AccessKeySecret: tea.String(config.Get().WebConfig.SmsAccessKeySecret),
		//
		Endpoint: tea.String(config.Get().WebConfig.SmsEndpoint),
	}
	_result = &dysmsapi20170525.Client{}
	_result, _err = dysmsapi20170525.NewClient(conf)
	return _result, _err
}

func SendCode(ph string, captcha string) code.Code {
	//发送验证码
	client, err := createClient()
	if err != nil {
		slog.Errorf("create dysms err, err:%v", err)
		return code.Code_InnerError
	}
	templateParams, err := json.Marshal(struct {
		Code string `json:"code"`
	}{Code: captcha})

	if err != nil {
		slog.Errorf("code marshal err, err:%v", err)
		return code.Code_InnerError
	}

	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers:  tea.String(ph),
		SignName:      tea.String(config.Get().WebConfig.SmsSignName),
		TemplateCode:  tea.String(config.Get().WebConfig.SmsTemplateCode),
		TemplateParam: tea.String(string(templateParams)),
	}

	rsp, err := client.SendSmsWithOptions(sendSmsRequest, &util.RuntimeOptions{})
	if err != nil {
		slog.Errorf("send dysms err, err:%v", err)
		return code.Code_InnerError
	}
	slog.Infof("send code ok, ph:%v|cod:%v|rsp:%v", ph, captcha, rsp)
	return code.Code_Ok
}
