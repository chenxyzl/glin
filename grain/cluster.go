package grain

import (
	"fmt"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/cluster"
	"github.com/asynkron/protoactor-go/cluster/clusterproviders/etcd"
	"github.com/asynkron/protoactor-go/cluster/identitylookup/disthash"
	"github.com/asynkron/protoactor-go/remote"
	"github.com/chenxyzl/glin/slog"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"os"
	"os/exec"
	"path/filepath"
	"sync/atomic"
	"time"
)

// System 集群
type System struct {
	cluster *cluster.Cluster
	pid     string
	running int32
}

// Start 启动集群
func (s *System) Start(clusterName string, etcdUrl []string, options ...cluster.ConfigOption) error {
	//only one system running
	err := s.cas(0, 1)
	if err != nil {
		return err
	}
	//
	system := actor.NewActorSystem()
	provider, err := etcd.NewWithConfig("/glin", clientv3.Config{
		Endpoints:   etcdUrl,
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		return err
	}
	config := remote.Configure("localhost",
		0,
		remote.WithDialOptions(grpc.WithKeepaliveParams(keepalive.ClientParameters{PermitWithoutStream: true}),
			grpc.WithBlock(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		),
	)
	lookup := disthash.New()
	clusterConfig := cluster.Configure(clusterName, provider, lookup, config, options...)
	s.cluster = cluster.New(system, clusterConfig)
	s.cluster.StartMember()
	//
	s.createPid()
	//
	slog.Info("system start success")
	//
	return nil
}

// Stop 关闭集群
func (s *System) Stop() error {
	err := s.cas(1, 0)
	if err != nil {
		return err
	}
	slog.Info("system stopped begin")
	s.cluster.Shutdown(true)
	s.removePid()
	slog.Info("system stopped success")
	return nil
}

// NewLocalActor 创建本地actor
func (s *System) NewLocalActor(factory func() IActor, opts ...actor.PropsOption) *actor.PID {
	if s.getRunning() != 1 {
		slog.Panicf("system not running")
	}
	props := actor.PropsFromProducer(func() actor.Actor {
		a := &Grain{inner: factory(), isCluster: false}
		return a
	}, opts...)
	return s.cluster.ActorSystem.Root.Spawn(props)
}

// RegisterClusterActor 注册远程actor
func RegisterClusterActor(factory func() IActor, opts ...actor.PropsOption) *cluster.Kind {
	test := factory()
	props := actor.PropsFromProducer(func() actor.Actor {
		return &Grain{inner: factory(), isCluster: true}
	}, opts...)
	kind := cluster.NewKind(string(test.GetKindType()), props)
	return kind
}

// GetCluster 获取集群
func (s *System) GetCluster() *cluster.Cluster {
	return s.cluster
}

// cas 检查重复运行
func (s *System) cas(old, new int32) error {
	//检查是否在运行
	if !atomic.CompareAndSwapInt32(&s.running, old, new) {
		return fmt.Errorf("state not in state, state:%v", old)
	}
	return nil
}

func (s *System) getRunning() int32 {
	return atomic.LoadInt32(&s.running)
}

// createPid 创建pid文件
func (s *System) createPid() {
	// create pid file
	arg0, err := exec.LookPath(os.Args[0])
	if err != nil {
		slog.Panic(err)
	}
	absExecFile, err := filepath.Abs(arg0)
	if err != nil {
		slog.Panic(err)
	}
	execDir, execFile := filepath.Split(absExecFile)
	pid := execDir + execFile + ".pid"
	//
	err = os.WriteFile(pid, []byte(fmt.Sprintf("%d", os.Getpid())), 0644)
	if err != nil {
		slog.Fatal(err)
	}
	//
	s.pid = pid
	slog.Infof("create pid, pid:%v", s.pid)
}

// removePid 删除pid文件
func (s *System) removePid() {
	if s.pid != "" {
		err := os.Remove(s.pid)
		if err != nil {
			slog.Error(err)
		}
		slog.Infof("remove pid, pid:%v", s.pid)
	}
}
