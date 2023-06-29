package online_model

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/chenxyzl/glin/model"
	"github.com/chenxyzl/glin/slog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"laiya/share/algorithm/skiplist"
	"laiya/share/mongo_helper"
)

var _ model.IModel = new(Online)

type Player struct {
	Uid          uint64 `bson:"_id,omitempty"`
	LastTickTime int64  `bson:"last_tick_time"`
}

type Online struct {
	model.Dirty    `bson:"-"`                  //不入库
	OnlineId       uint64                      `bson:"online_id"`
	SortPlayerList *skiplist.SkipList[*Player] `bson:"-"` //不入库-玩家在线时间排序的跳表
}

func (mod *Online) CoName() string {
	return "online"
}

func (mod *Online) Load() error {
	filter := bson.D{{"_id", mod.OnlineId}}
	col := mongo_helper.GetCol(mod)
	err := col.FindOne(context.TODO(), filter).Decode(mod)
	//账号不存在
	if errors.Is(err, mongo.ErrNoDocuments) {
		_, err = col.InsertOne(context.TODO(), mod)
		if err != nil {
			return fmt.Errorf("create %s err, onlineId:%v|err:%v", mod.CoName(), mod.OnlineId, err)
		}
	} else if err != nil {
		return fmt.Errorf("load %s err, onlineId:%v|err:%v", mod.CoName(), mod.OnlineId, err)
	}
	return err
}

func (mod *Online) Save() error {
	filter := bson.D{{"_id", mod.OnlineId}}
	col := mongo_helper.GetCol(mod)
	_, err := col.UpdateOne(context.Background(), filter, bson.M{"$set": mod}, options.Update().SetUpsert(true))
	if err != nil {
		tmp, _ := json.Marshal(mod)
		slog.Errorf("save error, save to log file, _id:%v|data:%s", mod.OnlineId, tmp)
	} else {
		mod.CleanDirty()
	}
	return err
}

func (mod *Online) Delete() error {
	return nil
}
