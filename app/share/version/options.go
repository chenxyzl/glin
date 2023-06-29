package version

import (
	"errors"
	"github.com/chenxyzl/glin/slog"
	"github.com/jessevdk/go-flags"
	"os"
)

// Options 启动时候的命令行参数
type Options struct {
	Tag     string   `short:"t" long:"tag" description:"tag" default:"laiya"` //仅仅是为了进程过滤之类的
	AppName string   `short:"a" long:"app" description:"app name" default:"laiya_v2"`
	Etcd    []string `short:"e" long:"etcd" description:"etcd url" default:"127.0.0.1:2379"`
	V       bool     `short:"v" long:"version" description:"echo version and rpc"`
}

var (
	options Options
)

func Get() Options {
	return options
}

func init() {
	if _, err := flags.NewParser(&options, flags.Default).Parse(); err != nil {
		var flagsErr *flags.Error
		switch {
		case errors.As(err, &flagsErr):
			if flagsErr.Type == flags.ErrHelp {
				os.Exit(0)
			}
			slog.Fatalf("parse flags err, :%v", flagsErr)
		default:
			slog.Fatalf("parse flags err, :%v", flagsErr)
		}
	}
}
