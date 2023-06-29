package sub_model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"laiya/share/mongo_helper"
)

var groupLastChatIdFindOps = options.Find().SetProjection(bson.D{{"_id", 1}, {"last_chat_id", 1}})

// GroupLastChatId 群最大的消息id
type GroupLastChatId struct {
	Gid        uint64 `bson:"_id,omitempty"`
	LastChatId uint64 `bson:"last_chat_id"` //不要手动修改-保存时计算
}

// Load 加载玩家头像数据
func (group *GroupLastChatId) Load() error {
	filter := bson.D{{"_id", group.Gid}}
	col := mongo_helper.GetColByStr(groupIns.CoName())
	err := col.FindOne(context.TODO(), filter).Decode(group)
	return err
}

func LoadMultiGroupLastChatId(gids []uint64) (map[uint64]uint64, error) {
	if len(gids) == 0 {
		return nil, nil
	}
	filter := bson.M{"_id": bson.M{"$in": gids}}
	// 执行查询
	cursor, err := mongo_helper.GetColByStr(groupIns.CoName()).Find(context.Background(), filter, groupLastChatIdFindOps)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	//
	ret := map[uint64]uint64{}
	// 遍历结果
	for cursor.Next(context.Background()) {
		var group GroupLastChatId
		if err := cursor.Decode(&group); err != nil {
			return nil, err
		}
		// 获取最后一个ChatMsg的MsgId
		ret[group.Gid] = group.LastChatId
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return ret, nil
}
