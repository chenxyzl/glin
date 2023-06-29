package online_model

import (
	"laiya/share/global"
	"time"
)

func (p *Player) Key() uint64 {
	return p.Uid
}

func (p *Player) Score() uint64 {
	return uint64(p.LastTickTime)
}

func (mod *Online) UpdatePlayerOnlineTick(uid uint64) *Player {
	//删除老数据
	mod.SortPlayerList.Delete(uid)
	//增加排名数据
	player := &Player{
		Uid:          uid,
		LastTickTime: time.Now().UnixNano(),
	}
	mod.SortPlayerList.Insert(player)
	return player
}

func (mod *Online) CheckOnline(dt time.Duration) {
	times := 0
	now := time.Now().UnixNano()
	for {
		//单次最多循环maxLoop次
		if times >= global.MaxLoop {
			break
		}
		if mod.SortPlayerList.Length() == 0 {
			break
		}
		first := mod.SortPlayerList.First().Value()
		if first == nil {
			break
		}
		if now-first.LastTickTime < int64(dt) {
			break
		}
		//删除过期的
		mod.SortPlayerList.Delete(first.Key())
		//次数++
		times++
	}
}
