package hall_model

import (
	"fmt"
	"laiya/share/algorithm/skiplist"
	"laiya/share/global"
)

func (mod *Hall) NewGroupSort() {
	mod.SortGroupList = skiplist.New[*Group](&cmpHallGroup{})
	for _, c := range mod.Groups {
		mod.SortGroupList.Insert(c)
	}
	mod.SortGroupList.ForRange(func(v *Group) bool {
		fmt.Println(v.Name, v.Score())
		return true
	})
}

// CleanOldPlayerCount //重启时候需要清除在线人数~避免之前的人数污染
func (mod *Hall) CleanOldPlayerCount() {
	for _, group := range mod.Groups {
		//重启时候需要清除在线人数~避免之前的人数污染
		group.CleanPlayerCount()
	}
}

func (c *Group) Key() uint64 {
	return c.Gid
}

func (c *Group) Score() uint64 {
	return c.SortScore
}

func (c *Group) CleanPlayerCount() {
	//
	mask := uint64(1<<global.PlayerCountOffset - 1)
	c.SortScore = c.SortScore & mask
}

func (c *Group) GetPlayerCount() uint32 {
	if c == nil {
		return 0
	}
	return uint32(c.Score() >> global.PlayerCountOffset)
}

// ChangeScore 跳表修改的hook
func (c *Group) ChangeScore(score uint64, before func(), after func()) {
	if before != nil {
		before()
	}
	c.SortScore = score
	if after != nil {
		after()
	}
}

type cmpHallGroup struct {
}

func (a *cmpHallGroup) CmpScore(v1 *Group, v2 *Group) int {
	s1 := v1.Score()
	s2 := v2.Score()
	switch {
	case s1 < s2:
		return -1
	case s1 == s2:
		return 0
	default:
		return 1
	}
}

func (a *cmpHallGroup) CmpKey(v1 *Group, v2 *Group) int {
	s1 := v1.Key()
	s2 := v2.Key()
	switch {
	case s1 < s2:
		return -1
	case s1 == s2:
		return 0
	default:
		return 1
	}
}
