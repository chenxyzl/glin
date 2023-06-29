package main

import (
	"fmt"
	"github.com/asynkron/protoactor-go/cluster"
	"github.com/chenxyzl/glin"
	"github.com/chenxyzl/glin/config_watcher"
	"github.com/chenxyzl/glin/grain"
	"github.com/chenxyzl/glin/slog"
	"google.golang.org/protobuf/proto"
	"laiya/config"
	"laiya/hall"
	"laiya/proto/code"
	"laiya/proto/inner"
	_ "laiya/register"
	share2 "laiya/share"
	"laiya/share/call"
	"laiya/share/global"
	"laiya/share/mongo_helper"
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
	mongo_helper.Connect()
	defer mongo_helper.Close()
	//frame
	err = glin.Start(
		version.Get().AppName,
		version.Get().Etcd,
		rootKey+"uuid/",
		version.String(),
		cluster.WithKinds(grain.RegisterClusterActor(func() grain.IActor {
			return hall.NewHallActor()
		})))
	if err != nil {
		slog.Fatal(err)
	}
	defer glin.Stop()
	//延时唤醒
	go func() {
		time.Sleep(time.Second / 2)
		_, cod := call.RequestRemote[proto.Message](global.HallUID, &inner.AwakeHall_Request{})
		if cod != code.Code_Ok {
			slog.Fatalf("start hall err, cod:%v", cod)
		}
	}()
	//
	slog.Infof("service start success")
	//wait stop
	glin.WaitStopSignal()
	//
	slog.Infof("service quit success")
}
