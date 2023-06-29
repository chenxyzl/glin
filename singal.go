package glin

import (
	"github.com/chenxyzl/glin/share"
	"github.com/chenxyzl/glin/slog"
	"os"
	"os/signal"
	"syscall"
)

func callFuncSlice(fs []func()) {
	defer share.Recover()
	for _, f := range fs {
		if f != nil {
			f()
		}
	}
}

func WaitStopSignal(beforeFunc ...func()) {
	// signal.Notify的ch信道是阻塞的(signal.Notify不会阻塞发送信号), 需要设置缓冲
	signals := make(chan os.Signal, 1)
	// It is not possible to block SIGKILL or syscall.SIGSTOP
	signal.Notify(signals, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
Quit:
	for {
		sig := <-signals
		slog.Infof("get signal %s", sig.String())
		//
		switch sig {
		case syscall.SIGHUP:
			slog.Infof("get signal do nothings")
		default:
			slog.Infof("get signal will quit")
			break Quit
		}
	}
	//
	callFuncSlice(beforeFunc)
	//
	slog.Infof("app exit success")
}
