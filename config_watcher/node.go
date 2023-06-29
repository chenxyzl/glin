package config_watcher

import (
	"fmt"
	"reflect"
)

// tNode is an Etcd key to load
type tNode struct {
	//
	prefix bool
	// key is the etcd key
	key string
	//
	filedName string
	//defaultPath is the local default config
	defaultPath string
	// node value
	value reflect.Value
	// ref value
	ref []reflect.Value
}

func (n *tNode) parse(data []byte) error {
	//
	t := n.value.Type()
	var cv reflect.Value
	if t.Kind() == reflect.Ptr {
		cv = reflect.New(t.Elem())
	} else {
		cv = reflect.New(t)
	}

	c := cv.Interface()

	//parse
	v, ok := c.(Parser)
	if !ok {
		return fmt.Errorf("cv can not convert to Parse,filed:%v|key:%v", t.Elem().Name(), n.key)
	}
	err := v.Parse(data)
	if err != nil {
		return err
	}

	//解析
	var refValues []reflect.Value
	for _, refValue := range n.ref {
		reft := refValue.Type()
		var refcv reflect.Value
		if reft.Kind() == reflect.Ptr {
			refcv = reflect.New(reft.Elem())
		} else {
			refcv = reflect.New(reft)
		}
		//
		refc := refcv.Interface()

		//parse ref
		refv, ok := refc.(RefParser)
		if !ok {
			return fmt.Errorf("ref cv can not convert to RefParse,filed:%v|key:%v", reft.Name(), n.key)
		}
		err := refv.RefParse(v)
		if err != nil {
			return err
		}

		//append
		if reft.Kind() != reflect.Ptr {
			refValues = append(refValues, refcv.Elem())
		} else {
			refValues = append(refValues, refcv)
		}
	}

	//赋值
	for i, ref := range n.ref {
		ref.Set(refValues[i])
	}

	if t.Kind() != reflect.Ptr {
		n.value.Set(cv.Elem())
	} else {
		n.value.Set(cv)
	}
	return nil
}
