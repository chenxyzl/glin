package sub_model

import (
	"context"
	"fmt"
	"github.com/chenxyzl/glin/slog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"laiya/share/mongo_helper"
)

var playerHeadProjection = bson.D{{"_id", 1}, {"name", 1}, {"icon", 1}, {"des", 1}}
var playerHeadFindOneOps = options.FindOne().SetProjection(playerHeadProjection)
var playerHeadFindOps = options.Find().SetProjection(playerHeadProjection)

// PlayerHead 用户的名字和icon
type PlayerHead struct {
	Uid  uint64 `bson:"_id,omitempty"`
	Name string `bson:"name"`
	Icon string `bson:"icon"`
	Des  string `bson:"des"`
}

// Load 加载玩家头像数据
func (player *PlayerHead) Load() error {
	filter := bson.D{{"_id", player.Uid}}
	col := mongo_helper.GetColByStr(playerIns.CoName())
	err := col.FindOne(context.TODO(), filter, playerHeadFindOneOps).Decode(player)
	return err
}

func LoadPlayerHead(uid uint64) (*PlayerHead, error) {
	player := &PlayerHead{Uid: uid}
	err := player.Load()
	return player, err
}

func BatchLoadPlayerHead(uids []uint64, logger slog.Logger) ([]*PlayerHead, error) {
	if len(uids) == 0 {
		return []*PlayerHead{}, nil
	}
	const max = 50
	if len(uids) > max {
		uids = uids[:max]
		logger.Warnf("data count too big, force fix to max, count:%v|max:%v", len(uids), max)
	}
	filter := bson.M{"_id": bson.M{"$in": uids}}
	// 执行查询
	cursor, err := mongo_helper.GetColByStr(playerIns.CoName()).Find(context.Background(), filter, playerHeadFindOps)
	if err != nil {
		return nil, fmt.Errorf("batch load player from db err, uids:%v|err:%v", uids, err)
	}
	defer cursor.Close(context.Background())
	//
	var players []*PlayerHead
	// 遍历结果
	for cursor.Next(context.Background()) {
		var player = &PlayerHead{}
		if err := cursor.Decode(&player); err != nil {
			return nil, fmt.Errorf("unmarshal player from db err, err:%v", err)
		}
		//
		players = append(players, player)
	}
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("mongo cursor err, uids:%v|err:%v", uids, err)
	}
	return players, nil
}

func (player *PlayerHead) UpdateNameIconDes(name string, icon string, des string) error {
	filter := bson.D{{"_id", player.Uid}}
	update := bson.D{{"$set", bson.D{
		{"name", name},
		{"icon", icon},
		{"des", des},
	}}}
	// 执行更新操作
	_, err := mongo_helper.GetColByStr(playerIns.CoName()).UpdateOne(context.TODO(), filter, update, options.Update().SetUpsert(true))
	if err == nil {
		player.Name = name
		player.Icon = icon
		player.Des = des
	}
	return err
}
