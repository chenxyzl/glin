package share

import (
	"github.com/chenxyzl/glin/slog"
	"runtime/debug"
)

func Recover(logger ...slog.Logger) {
	err := recover()
	if err != nil {
		stackTrace := debug.Stack()
		if len(logger) > 0 && logger[0] != nil {
			logger[0].Errorf("panic recover, err:%v|stackTrace:%v", err, string(stackTrace))
		} else {
			slog.Errorf("panic recover, err:%v|stackTrace:%v", err, string(stackTrace))
		}
	}
}

func RecoverInfo(info string, logger ...slog.Logger) {
	if info == "" {
		Recover()
	} else {
		err := recover()
		if err != nil {
			stackTrace := debug.Stack()
			if len(logger) > 0 && logger[0] != nil {
				logger[0].Errorf("panic recover, %v|err:%v|stackTrace:%v", info, err, string(stackTrace))
			} else {
				slog.Errorf("panic recover, %v|err:%v|stackTrace:%v", info, err, string(stackTrace))
			}
		}
	}
}

func RecoverFunc(pc func(err any), logger ...slog.Logger) {
	if pc == nil {
		Recover()
	} else {
		err := recover()
		if err != nil {
			stackTrace := debug.Stack()
			if len(logger) > 0 && logger[0] != nil {
				logger[0].Errorf("panic recover, err:%v|stackTrace:%v", err, string(stackTrace))
			} else {
				slog.Errorf("panic recover, err:%v|stackTrace:%v", err, string(stackTrace))
			}
			pc(err)
		}
	}
}
