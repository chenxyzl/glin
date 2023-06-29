package sub_model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"laiya/share/mongo_helper"
)

var groupGameTypFindOps = options.Find().SetProjection(bson.D{{"_id", 1}, {"name", 1}, {"game_type", 1}})

// GroupGameTyp 群最大的消息id
type GroupGameTyp struct {
	Gid      uint64 `bson:"_id,omitempty"`
	Name     string `bson:"name"`      //不要手动修改
	GameType string `bson:"game_type"` //不要手动修改
}

// Load 加载玩家头像数据
func (group *GroupGameTyp) Load() error {
	filter := bson.D{{"_id", group.Gid}}
	col := mongo_helper.GetColByStr(groupIns.CoName())
	err := col.FindOne(context.TODO(), filter).Decode(group)
	return err
}
