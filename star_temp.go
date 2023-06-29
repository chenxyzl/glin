package glin

import "sync"

var gs *Star
var lock sync.RWMutex

// setStar 设置go不支持struct的泛型函数, 暂时这么写
func setStar(s *Star) {
	lock.Lock()
	defer lock.Unlock()
	if gs != nil {
		panic("glin has set")
	}
	gs = s
}

// GetStar 设置go不支持struct的泛型函数, 暂时这么写
func GetStar() *Star {
	lock.RLock()
	defer lock.RUnlock()
	if gs == nil {
		panic("glin not set")
	}
	return gs
}
