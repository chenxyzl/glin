package main

import (
	"fmt"
	"github.com/asynkron/protoactor-go/cluster"
	"github.com/chenxyzl/glin"
	"github.com/chenxyzl/glin/config_watcher"
	"github.com/chenxyzl/glin/grain"
	"github.com/chenxyzl/glin/slog"
	"laiya/config"
	"laiya/group"
	_ "laiya/register"
	share2 "laiya/share"
	"laiya/share/mongo_helper"
	"laiya/share/version"
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
			return group.NewGroupActor()
		})))
	if err != nil {
		slog.Fatal(err)
	}
	defer glin.Stop()
	//register actor
	//
	slog.Infof("service start success")
	//wait stop
	glin.WaitStopSignal()
	//
	slog.Infof("service quit success")
}
