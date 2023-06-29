package share2

import (
	"fmt"
	"github.com/chenxyzl/glin/slog"
	"laiya/config"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func OpenProfile() {
	arg0, err := exec.LookPath(os.Args[0])
	if err != nil {
		slog.Panic(err)
	}
	//
	absExecFile, err := filepath.Abs(arg0)
	if err != nil {
		slog.Panic(err)
	}
	_, execFile := filepath.Split(absExecFile)
	execFile = strings.TrimSuffix(execFile, ".exe")
	//获取端口配置
	port, ok := config.Get().AConfig.PprofPort[execFile]
	if !ok {
		slog.Infof("pprof not start, please config at a.toml")
		return
	}
	//开启服务
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		slog.Fatal(err)
	}
	addr := ln.Addr().String()
	slog.Infof("pprof started on addr:%v/debug/pprof/", addr)
	//开启监听
	go func() {
		err = http.Serve(ln, nil)
		if err != nil {
			slog.Fatal(err)
		}
	}()
}

//由于云服务器只会开放特定的端口,所以没办法随便选择
//func OpenProfile() {
//	ln, err := net.Listen("tcp", ":0")
//	if err != nil {
//		slog.Fatal(err)
//	}
//	addr := ln.Addr().String()
//	slog.Infof("pprof started on addr:%v/debug/pprof/", addr)
//	go func() {
//		err = http.Serve(ln, nil)
//		if err != nil {
//			slog.Fatal(err)
//		}
//		slog.Infof("pprof closed on addr:%v/debug/pprof/", addr)
//	}()
//}
