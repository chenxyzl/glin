package config_impl

import (
	"github.com/BurntSushi/toml"
	"github.com/dacker-soul/getui/publics"
	"time"
)

type OssBuket struct {
	OssBuketName string
	OssKey       string
	OssSecret    string
	OssEndpoint  string
}

type WebConfig struct {
	//验证码
	CaptchaExpireTime      time.Duration
	CaptchaReqIntervalTime time.Duration
	//短信
	SmsAccessKeyId     string
	SmsAccessKeySecret string
	SmsEndpoint        string
	SmsSignName        string
	SmsTemplateCode    string
	//阿里云的oss--头像
	HeadOssBuket OssBuket
	//阿里云的oss--图片(聊天)
	PicOssBuket OssBuket
	//长短链
	ShorUrlGroup string
	ShortUrl     string
	LongUrl      string
	//个推
	GeTuiAppId        string
	GeTuiAppSecret    string
	GeTuiAppKey       string
	GeTuiMasterSecret string
	GeTuiSettings     *publics.Settings
}

func (conf *WebConfig) Parse(data []byte) error {
	_, err := toml.Decode(string(data), conf)
	if err != nil {
		return err
	}
	return nil
}
