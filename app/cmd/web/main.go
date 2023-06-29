package main

import (
	"fmt"
	"github.com/chenxyzl/glin"
	"github.com/chenxyzl/glin/config_watcher"
	"github.com/chenxyzl/glin/slog"
	"laiya/config"
	"laiya/model/web_model/web_old_model/old_mongo"
	_ "laiya/register"
	share2 "laiya/share"
	"laiya/share/getui_helper"
	"laiya/share/mongo_helper"
	"laiya/share/version"
	"laiya/web/gin"
)

func main() {
	slog.Infof("service starting:\n%v", version.String())
	defer slog.Sync()
	//
	rootKey := fmt.Sprintf("/%v/%v/", version.Get().AppName, version.Branch)
	//config
	err := (&config_watcher.Watcher{}).Start(version.Get().Etcd, rootKey+"config/", config.Get())
	if err != nil {
		slog.Fatal(err)
	}
	//开启pprof--需要在config加载后
	share2.OpenProfile()
	//mongo
	mongo_helper.Connect()
	defer mongo_helper.Close()
	//老的mongo 临时保存
	old_mongo.Connect()
	defer old_mongo.Close()
	//frame
	err = glin.Start(version.Get().AppName, version.Get().Etcd, rootKey+"uuid/", version.String())
	if err != nil {
		slog.Fatal(err)
	}
	defer glin.Stop()
	//
	check()
	//启动gin
	gin.StartGin()
	//
	slog.Infof("service start success")
	//wait stop
	glin.WaitStopSignal()
	//
	slog.Infof("service quit success")
}

func check() {
	//个推
	err := getui_helper.TestGeTuiToken()
	if err != nil {
		slog.Fatal(err)
	}
}
