package main

import (
	"fmt"
	"github.com/chenxyzl/glin"
	"github.com/chenxyzl/glin/config_watcher"
	"github.com/chenxyzl/glin/slog"
	"github.com/robfig/cron"
	"laiya/config"
	"laiya/gate"
	_ "laiya/register"
	share2 "laiya/share"
	"laiya/share/version"
	"time"
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
	//mongo_helper.Connect()
	//defer mongo_helper.Close()
	//frame
	err = glin.Start(version.Get().AppName, version.Get().Etcd, rootKey+"uuid/", version.String())
	if err != nil {
		slog.Fatal(err)
	}
	defer glin.Stop()
	//启动websocket
	gate.StateWs()
	defer gate.CloseWs()
	//tick
	cron := startTick(tick)
	defer cron.Stop()
	//
	slog.Infof("service start success")
	//wait stop
	glin.WaitStopSignal()
	//
	slog.Infof("service quit success")
}

func tick(time int64) {
	gate.PrintOnlineCount()
}

func startTick(f func(timestamp int64)) *cron.Cron {
	tick := func() {
		now := time.Now().Unix()
		f(now)
	}
	cron2 := cron.New() //创建一个cron实例
	//执行定时任务（每5秒执行一次）
	err := cron2.AddFunc("*/1 * * * * *", tick)
	if err != nil {
		panic(err)
	}
	//启动/关闭
	cron2.Start()
	return cron2
}
