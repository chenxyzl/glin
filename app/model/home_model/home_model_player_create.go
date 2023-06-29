package home_model

import (
	"context"
	"laiya/share/mongo_helper"
)

// CreatePlayerTyp 用户的名字和icon
type CreatePlayerTyp struct {
	Uid     uint64 `bson:"_id,omitempty"`
	Account string `bson:"account"`
	Name    string `bson:"name"`
	Icon    string `bson:"icon"`
	Des     string `bson:"des"`
}

func CreatePlayer(uid uint64, account string, name string, icon string, des string) (*CreatePlayerTyp, error) {
	player := &CreatePlayerTyp{
		Uid:     uid,
		Account: account,
		Name:    name,
		Icon:    icon,
		Des:     des,
	}
	col := mongo_helper.GetColByStr((&Player{}).CoName())
	_, err := col.InsertOne(context.Background(), player)
	return player, err
}
