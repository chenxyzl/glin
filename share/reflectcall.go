package share

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

func SafeCall(method reflect.Method, receiver reflect.Value, args ...reflect.Value) (rets []reflect.Value, err error) {
	args = append([]reflect.Value{receiver}, args...)
	defer RecoverInfo(fmt.Sprintf("methodName:%s|args:%v", method.Name, args))
	r := method.Func.Call(args)
	return r, err
}

func CurrentPackageName() string {
	pc, _, _, _ := runtime.Caller(2)
	parts := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	pl := len(parts)
	packageName := ""
	funcName := parts[pl-1]

	if parts[pl-2][0] == '(' {
		funcName = parts[pl-2] + "." + funcName
		packageName = strings.Join(parts[0:pl-2], ".")
	} else {
		packageName = strings.Join(parts[0:pl-1], ".")
	}

	return packageName
}
