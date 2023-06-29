package glin

import (
	"context"
	"fmt"
	"github.com/asynkron/protoactor-go/cluster"
	"github.com/chenxyzl/glin/grain"
	"github.com/chenxyzl/glin/slog"
	"github.com/chenxyzl/glin/uuid"
	"math/rand"
	"sync/atomic"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	ttlTime = 10 //ttl的单位都是秒
)

type STATUS uint32

const (
	SET STATUS = iota
	TTL STATUS = iota
)

type Star struct {
	id      uint64
	client  *clientv3.Client
	lease   clientv3.Lease
	leaseId clientv3.LeaseID
	status  STATUS //状态机
	path    string
	val     interface{}
	system  *grain.System
	running int32
}

func (s *Star) key() string {
	return s.path + fmt.Sprintf("%d", s.id)
}

func (s *Star) nextId() uint64 {
	return rand.Uint64()%uuid.GetMaxNodeId() + 1
}

func (s *Star) set() bool {
	//设置key
	key := s.key()
	tx := s.client.Txn(context.Background())
	//key no exist
	leaseResp, err := s.lease.Grant(context.Background(), ttlTime)
	if err != nil {
		return false
	}
	s.leaseId = leaseResp.ID
	tx.If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).
		Then(clientv3.OpPut(key, fmt.Sprintf("%v", s.val), clientv3.WithLease(s.leaseId))).
		Else()
	txnRes, err := tx.Commit()
	if err != nil || !txnRes.Succeeded { //抢锁失败
		s.id = s.nextId()
		return false
	}

	s.status = TTL

	//更新uuid生成器
	err = uuid.Init(s.id)
	if err != nil {
		slog.Fatalf("init uuid err, s.id:%v|err:%v", s.id, err)
	}
	slog.Infof("service register success, service:\n%vuuid:%v", s.val, s.id)
	return true
}

func (s *Star) ttl() {
	//保持ttl
	_, err := s.lease.KeepAliveOnce(context.Background(), s.leaseId)
	if err != nil {
		s.status = SET
		slog.Errorf("service ttl failed, will register:\n%vuuid:%v", s.val, s.id)
	} else {
		time.Sleep((ttlTime * 2 / 3) * time.Second)
	}
}

func (s *Star) run() {
	for s.getRunning() == 1 {
		switch s.status {
		case SET:
			s.set()
			time.Sleep(time.Second * 1) //间隔1秒一次避免把etcd弄炸了
		case TTL:
			s.ttl()
		}
	}
}

func (s *Star) start() {
	go s.run()
}

// cas 检查重复运行
func (s *Star) cas(old, new int32) error {
	//检查是否在运行
	if !atomic.CompareAndSwapInt32(&s.running, old, new) {
		return fmt.Errorf("state not in state, state:%v", old)
	}
	return nil
}

func (s *Star) getRunning() int32 {
	return atomic.LoadInt32(&s.running)
}

func GetSystem() *grain.System {
	return GetStar().system
}

// Start 注册服务
func Start(appName string, etcdUrl []string, path string, val interface{}, options ...cluster.ConfigOption) error {
	s := &Star{}
	slog.Infof("glin starting...")
	if err := s.cas(0, 1); err != nil {
		return fmt.Errorf("glin is in running, appName:%v", appName)
	}
	//
	etcdClient, err := clientv3.New(clientv3.Config{Endpoints: etcdUrl})
	if err != nil {
		return fmt.Errorf("cannot connect to etcd:%v|err:%v", etcdUrl, err)
	}
	lease := clientv3.NewLease(etcdClient)
	s.id = s.nextId()
	s.client = etcdClient
	s.lease = lease
	s.path = path
	s.val = val
	s.system = &grain.System{}
	err = s.system.Start(appName, etcdUrl, options...)
	if err != nil {
		return err
	}
	setStar(s)
	//register and get uuid
	for !s.set() {
		time.Sleep(time.Second*1 + time.Duration(rand.Uint64())%time.Second) //间隔随机间隔1s+避免把etcd弄炸了
	}
	//
	s.start()
	//
	slog.Infof("glin start success")
	return nil
}

func Stop() {
	s := GetStar()
	if s == nil {
		return
	}
	if err := s.cas(1, 0); err != nil {
		slog.Error(err)
	}
	if s.client != nil {
		_, err := s.lease.Revoke(context.Background(), s.leaseId)
		if err != nil {
			slog.Error(err)
		}
	}
	if s.system != nil {
		err := s.system.Stop()
		if err != nil {
			slog.Error(err)
		}
	}
	return
}
