package config_impl

import (
	"fmt"
	"laiya/share/utils"
)

func (conf *WebConfig) GetShorUrl(groupId uint64) string {
	short := utils.Uint64ToBase62(groupId)
	return fmt.Sprintf("%v/%v/%v", conf.ShortUrl, conf.ShorUrlGroup, short)
}
