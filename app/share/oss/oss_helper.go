package oss_helper

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"laiya/config"
	"sync"
)

var HeadOssBucket *oss.Bucket
var lockHeadOssBucketIns = sync.Mutex{}
var PicOssBucket *oss.Bucket
var lockPicOssBucketIns = sync.Mutex{}

func GetHeadBuket() (*oss.Bucket, error) {
	if HeadOssBucket == nil {
		lockHeadOssBucketIns.Lock()
		defer lockHeadOssBucketIns.Unlock()
		if HeadOssBucket != nil {
			return HeadOssBucket, nil
		}
		client, err := oss.New(config.Get().WebConfig.HeadOssBuket.OssEndpoint, config.Get().WebConfig.HeadOssBuket.OssKey, config.Get().WebConfig.HeadOssBuket.OssSecret)
		if err != nil {
			return nil, err
		}
		bucket, err := client.Bucket(config.Get().WebConfig.HeadOssBuket.OssBuketName)
		if err != nil {
			return nil, err
		}
		HeadOssBucket = bucket
	}
	return HeadOssBucket, nil
}

func GetPicBuket() (*oss.Bucket, error) {
	if PicOssBucket == nil {
		lockPicOssBucketIns.Lock()
		defer lockPicOssBucketIns.Unlock()
		if PicOssBucket != nil {
			return PicOssBucket, nil
		}
		client, err := oss.New(config.Get().WebConfig.PicOssBuket.OssEndpoint, config.Get().WebConfig.PicOssBuket.OssKey, config.Get().WebConfig.PicOssBuket.OssSecret)
		if err != nil {
			return nil, err
		}
		bucket, err := client.Bucket(config.Get().WebConfig.PicOssBuket.OssBuketName)
		if err != nil {
			return nil, err
		}
		PicOssBucket = bucket
	}
	return PicOssBucket, nil
}

func GetHeadUrl(fileName string) string {
	return "https://" + config.Get().WebConfig.HeadOssBuket.OssBuketName + "." + config.Get().WebConfig.HeadOssBuket.OssEndpoint + "/" + fileName
}

func GetPicUrl(fileName string) string {
	return "https://" + config.Get().WebConfig.PicOssBuket.OssBuketName + "." + config.Get().WebConfig.PicOssBuket.OssEndpoint + "/" + fileName
}
