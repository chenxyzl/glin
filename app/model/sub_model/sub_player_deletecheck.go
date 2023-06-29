package sub_model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"laiya/share/mongo_helper"
)

var playerDeleteProjectionCheck = bson.D{{"_id", 1}, {"is_deleted", 1}}
var playerDeleteFindOneOpsCheck = options.FindOne().SetProjection(playerDeleteProjectionCheck)

// PlayerDeleteCheck 用户的名字和icon
type PlayerDeleteCheck struct {
	Uid       uint64 `bson:"_id,omitempty"`
	IsDeleted bool   `bson:"is_deleted"`
}

// Load 加载玩家头像数据
func (player *PlayerDeleteCheck) Load() error {
	filter := bson.D{{"_id", player.Uid}}
	col := mongo_helper.GetColByStr(playerIns.CoName())
	err := col.FindOne(context.TODO(), filter, playerDeleteFindOneOpsCheck).Decode(player)
	return err
}
