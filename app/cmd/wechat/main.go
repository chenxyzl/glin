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
	"laiya/proto/code"
	"laiya/proto/inner"
	_ "laiya/register"
	share2 "laiya/share"
	"laiya/share/call"
	"laiya/share/global"
	"laiya/share/mongo_helper"
	"laiya/share/version"
	"laiya/wechat"
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
			return wechat.NewWechatActor()
		})))
	if err != nil {
		slog.Fatal(err)
	}
	defer glin.Stop()
	//register actor
	//
	slog.Infof("service start success")

	//延时唤醒
	go func() {
		time.Sleep(time.Second / 2)
		_, cod := call.RequestRemote[proto.Message](global.WechatUID, &inner.AwakeWechat_Request{})
		if cod != code.Code_Ok {
			slog.Fatalf("start wechat err, cod:%v", cod)
		}
	}()

	//wait stop
	glin.WaitStopSignal()
	//
	slog.Infof("service quit success")
}
