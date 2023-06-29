package register

import (
	"github.com/chenxyzl/glin/grain"
	"github.com/chenxyzl/glin/slog"
	"laiya/group"
	"laiya/plugin/activity"
)

func init() {
	grain.RegisterPlugin[*group.GroupActor, *activity.ActivityActor](activity.NewActivityActor, true)
	slog.Info("register plugin success...")
}
