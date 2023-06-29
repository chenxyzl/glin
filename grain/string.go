package grain

import (
	"github.com/chenxyzl/glin/share"
	"sort"
)

func String() string {
	str := "actor远程接口:"
	//for rpcName, handler := range remoteRpcHandlers {
	//	str += "\nactorName:" + handler.EntityType
	//	str += "  rpcName:" + rpcName
	//}
	m := map[share.EntityKind][]*HandlerDef{}
	for _, handler := range remoteRpcHandlers {
		m[handler.GetEntityKind()] = append(m[handler.GetEntityKind()], handler)
	}

	var keys []string
	for key := range m {
		keys = append(keys, string(key))
	}
	sort.Sort(sort.StringSlice(keys))
	for _, key := range keys {
		vs := m[share.EntityKind(key)]
		sort.Slice(vs, func(i, j int) bool {
			return vs[i].Method.Name < vs[j].Method.Name
		})
		for _, v := range vs {
			str += "\nactorName:" + key
			str += "  rpcName:" + v.Method.Name
			if v.Method.Type.NumOut() == 0 {
				str += "  notify:" + v.Method.Type.In(1).Elem().Name()
			} else {
				str += "  request:" + v.Method.Type.In(1).Elem().Name()
				str += "  reply:" + v.Method.Type.Out(0).Elem().Name()
			}
		}
	}
	return str
}
