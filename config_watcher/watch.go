package config_watcher

import (
	"context"
	"fmt"
	"github.com/chenxyzl/glin/slog"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"os"
	"reflect"
	"strings"
)

const Example = `type Config struct {
    TestKey       WebConfig     ` + "`" + `key:"/config/test.toml"` + "`" + `default:"../config/test.toml"` + "`" + `
    TestKeyPrefix *WebConfig    ` + "`" + `key:"/config/test.toml,prefix"` + "`" + `
    TestRef1      WebRefConfig  ` + "`" + `ref:"TestKey"` + "`" + `
    TestRef2      *WebRefConfig ` + "`" + `ref:"TestKeyPrefix"` + "`" + `
}
`

type Tag string

const (
	Key     Tag = "key"     //mean etcd path
	Prefix  Tag = "prefix"  //mean etcd prefix key
	Ref     Tag = "ref"     //mean depends on normal key
	Default Tag = "default" //mean local default file
)

type Parser interface {
	Parse([]byte) error
}

type RefParser interface {
	RefParse(interface{}) error
}

type IConfig interface {
	AfterLoadAll()
}

// Watcher 监控
type Watcher struct {
	client  *clientv3.Client
	keyRoot string
	nodes   []tNode
}

// Start 启动
// @config example: Example
// @keyRoot example:"/laiya_v1/dev/config"
func (m *Watcher) Start(Endpoints []string, keyRoot string, config IConfig) error {
	//
	cfg := clientv3.Config{
		Endpoints: Endpoints,
	}
	//
	etcdClient, err := clientv3.New(cfg)
	if err != nil {
		return fmt.Errorf("cannot connect to etcd, err:%v", err)
	}
	//
	m.client = etcdClient
	m.keyRoot = keyRoot
	//
	m.nodes, err = m.analyze(config)
	if err != nil {
		return err
	}
	//
	m.print()
	//
	err = m.firstLoad()
	if err != nil {
		return err
	}
	//
	config.AfterLoadAll()
	//
	m.start()
	//
	return nil
}

func (m *Watcher) start() {
	for _, node := range m.nodes {
		go m.watch(node)
	}
}

func (m *Watcher) firstLoad() error {
	for _, node := range m.nodes {
		var rsp *clientv3.GetResponse
		var err error
		//load数据
		if node.prefix {
			rsp, err = m.client.Get(context.Background(), node.key, clientv3.WithPrefix())
		} else {
			rsp, err = m.client.Get(context.Background(), node.key)
		}
		if err != nil {
			return err
		}
		//解析
		if len(rsp.Kvs) == 0 || len(rsp.Kvs[0].Value) == 0 {
			fileValue, err := os.ReadFile(node.defaultPath)
			if err != nil {
				return fmt.Errorf("load default config err, key:%v|default:%v|err:%v", node.key, node.defaultPath, err)
			}
			err = node.parse(fileValue)
			if err != nil {
				return err
			}
			slog.Infof("config first load with default, default:%v|value:%v", node.defaultPath, node.value.Interface())
		} else {
			for _, kv := range rsp.Kvs {
				err = node.parse(kv.Value)
				if err != nil {
					return err
				}
				slog.Infof("config first load with etcd, key:%v|value:%v", node.key, node.value.Interface())
			}
		}
	}
	return nil
}

func (m *Watcher) watch(node tNode) {
	var wch clientv3.WatchChan
	if node.prefix {
		wch = m.client.Watch(context.Background(), node.key, clientv3.WithPrefix(), clientv3.WithPrevKV())
	} else {
		wch = m.client.Watch(context.Background(), node.key, clientv3.WithPrevKV())
	}

	for v := range wch {
		for _, v1 := range v.Events {
			//get value
			var data []byte
			if v1.Type == mvccpb.PUT {
				data = v1.Kv.Value
				slog.Infof("watcher key changed, key:%v", string(v1.Kv.Key))
			} else {
				slog.Infof("watcher key has delete, key:%v", string(v1.Kv.Key))
			}
			//if data is empty, d`not parse
			if len(data) <= 0 {
				continue
			}
			//analyze
			if err := node.parse(data); err != nil {
				slog.Errorf("watcher nodes analyze err, key:%v|err:%v", string(v1.Kv.Key), err)
				continue
			}
			slog.Infof("config changed, key:%v|key:%v|value:%v", string(v1.Kv.Key), node.key, node.value.Interface())
		}
	}

}

func (m *Watcher) print() {
	for _, node := range m.nodes {
		slog.Infof("watch node, key:%v", node.key)
	}
}

func (m *Watcher) analyze(config IConfig) ([]tNode, error) {
	ic := reflect.ValueOf(config).Elem()
	var nodes []tNode
	for i := 0; i < ic.NumField(); i++ {
		fieldType := ic.Type().Field(i)
		key := fieldType.Tag.Get(string(Key))
		defaultPath := fieldType.Tag.Get(string(Default))
		ref := fieldType.Tag.Get(string(Ref))
		//only one
		if key != "" && ref != "" {
			return nil, fmt.Errorf("tag only has one, example:%v", Example)
		}

		//
		if key != "" {
			arr := strings.Split(key, ",")
			prefix := false
			if len(arr) > 2 {
				return nil, fmt.Errorf("tag format err, example: %v", Example)
			}
			if len(arr) == 2 {
				if arr[1] != string(Prefix) {
					return nil, fmt.Errorf("tag not support err, has:%v|example: %v", arr[1], Example)
				}
				prefix = true
			}

			key = arr[0]
			node := tNode{
				prefix:      prefix,
				key:         m.keyRoot + key,
				filedName:   fieldType.Name,
				value:       ic.Field(i),
				defaultPath: defaultPath,
			}
			nodes = append(nodes, node)
		} else if ref != "" {
			found := false
			for j := range nodes {
				if nodes[j].filedName == ref {
					nodes[j].ref = append(nodes[j].ref, ic.Field(i))
					found = true
					break
				}
			}
			if !found {
				return nil, fmt.Errorf("ref must after depend, file:%v|depend:%v", fieldType.Name, ref)
			}
		} else {
			return nil, fmt.Errorf("tag format err, example: %v", Example)
		}
	}
	return nodes, nil
}
