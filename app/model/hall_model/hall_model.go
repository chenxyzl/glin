package hall_model

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

var _ model.IModel = new(Hall)

type Group struct {
	Gid                  uint64 `bson:"_id,omitempty"`
	Name                 string `bson:"name"`                    //用于搜索
	VoiceRoomName        string `bson:"temp_voice_room"`         //用于搜索
	GameType             string `bson:"game_type"`               //用于搜索
	SortScore            uint64 `bson:"sort_score"`              //用于排序
	TempVoiceRoomVersion uint64 `bson:"temp_voice_room_version"` //临时的语音房间版本,用于过滤重复的摇人--垃圾需求的临时实现
}

type Hall struct {
	model.Dirty    `bson:"-"`                 //不入库
	HallId         uint64                     `bson:"_id,omitempty"`
	Groups         map[uint64]*Group          `bson:"groups"`
	NextNotifyTime int64                      `bson:"next_notify_time"`
	SortGroupList  *skiplist.SkipList[*Group] `bson:"-"` //不入库-玩家在线时间排序的跳表
}

func (mod *Hall) CoName() string {
	return "hall"
}

func (mod *Hall) Load() error {
	filter := bson.D{{"_id", mod.HallId}}
	col := mongo_helper.GetCol(mod)
	err := col.FindOne(context.TODO(), filter).Decode(mod)
	//账号不存在
	if errors.Is(err, mongo.ErrNoDocuments) {
		_, err = col.InsertOne(context.TODO(), mod)
		if err != nil {
			return fmt.Errorf("create %s err, hallId:%v|err:%v", mod.CoName(), mod.HallId, err)
		}
	} else if err != nil {
		return fmt.Errorf("load %s err, hallId:%v|err:%v", mod.CoName(), mod.HallId, err)
	}
	return err
}

func (mod *Hall) Save() error {
	filter := bson.D{{"_id", mod.HallId}}
	col := mongo_helper.GetCol(mod)
	_, err := col.UpdateOne(context.Background(), filter, bson.M{"$set": mod}, options.Update().SetUpsert(true))
	if err != nil {
		tmp, _ := json.Marshal(mod)
		slog.Errorf("save error, save to log file, _id:%v|data:%s", mod.HallId, tmp)
	} else {
		mod.CleanDirty()
	}
	return err
}

func (mod *Hall) Delete() error {
	return nil
}
