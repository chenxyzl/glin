package version

import (
	"fmt"
	"github.com/chenxyzl/glin/grain"
	"laiya/proto/outer"
	"os"
	"runtime"
	"strings"
	"time"
)

// 编译时注入
var (
	BuildTime  = ""
	CommitID   = ""
	CommitTime = ""
	Branch     = "dev"
	GoVer      = runtime.Version()
)

func String() string {
	str := "版本信息:" +
		"\nBuild Time:" + BuildTime +
		"\nBranch:" + Branch +
		"\nCommit ID:" + CommitID +
		"\nCommit Time:" + CommitTime +
		"\nAppName:" + Get().AppName +
		"\nEtcd:" + strings.Join(Get().Etcd, ";") +
		"\nGoVer:" + GoVer
	if options.V {
		fmt.Println(str)
		fmt.Println(outer.String())
		fmt.Println(grain.String()) //todo 这个因为依赖顺序问题,打印不出来,暂时不好解决
		time.Sleep(time.Second)
		os.Exit(0)
	}
	return str
}
