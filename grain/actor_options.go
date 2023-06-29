package grain

import "time"

type actorOptions struct {
	terminateIdleTime time.Duration //空闲多长时间后删除
	filterIdleCheck   func() bool
}
type ActorOption func(*actorOptions)

func WithTerminateIdleTime(timeout time.Duration) ActorOption {
	return func(o *actorOptions) { o.terminateIdleTime = timeout }
}

func WithFilterIdleCheck(check func() bool) ActorOption {
	return func(o *actorOptions) { o.filterIdleCheck = check }
}
