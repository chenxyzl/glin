package common

import (
	"github.com/chenxyzl/glin/grain"
	"laiya/proto/code"
	"laiya/proto/outer"
)

func DefaultRecvIfInitErr(a grain.IActor) {
	a.GetCtx().Respond(outer.BuildError(code.Code_InnerError))
	a.GetCtx().Poison(a.GetCtx().Self())
	a.GetLogger().Errorf("actor init failed,poison self,actor:%v", a.GetCtx().Self())
}
