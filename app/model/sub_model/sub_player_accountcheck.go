package sub_model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"laiya/share/mongo_helper"
)

var playerAccountProjectionCheck = bson.D{{"_id", 1}, {"account", 1}}
var playerAccountFindOneOpsCheck = options.FindOne().SetProjection(playerAccountProjectionCheck)

// PlayerAccountCheck 用户的名字和icon
type PlayerAccountCheck struct {
	Uid     uint64 `bson:"_id,omitempty"`
	Account string `bson:"account"`
}

// Load 加载玩家头像数据
func (player *PlayerAccountCheck) Load() error {
	filter := bson.D{{"_id", player.Uid}}
	col := mongo_helper.GetColByStr(playerIns.CoName())
	err := col.FindOne(context.TODO(), filter, playerAccountFindOneOpsCheck).Decode(player)
	return err
}
