package main

import (
	"fmt"
	"github.com/chenxyzl/glin"
	"github.com/chenxyzl/glin/config_watcher"
	"github.com/chenxyzl/glin/slog"
	"laiya/config"
	_ "laiya/register"
	"laiya/share/version"
)

func main() {
	//version
	slog.Infof("app version:\n%v", version.String())
	defer slog.Sync()

	//
	rootKey := fmt.Sprintf("/%v/%v/", version.Get().AppName, version.Branch)
	//config
	err := (&config_watcher.Watcher{}).Start(version.Get().Etcd, rootKey+"config/", config.Get())
	if err != nil {
		slog.Fatal(err)
	}
	fmt.Println(config.Get())
	//frame
	err = glin.Start(version.Get().AppName, version.Get().Etcd, rootKey+"uuid/", version.String())
	if err != nil {
		slog.Fatal(err)
	}
	defer glin.Stop()
	//wait stop
	glin.WaitStopSignal()
}
