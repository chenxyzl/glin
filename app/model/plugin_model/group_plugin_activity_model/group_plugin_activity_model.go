package group_plugin_activity_model

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
	"laiya/share/mongo_helper"
)

var _ model.IModel = new(PluginActivity)

type Player struct {
	Uid    uint64 `bson:"uid"`    //参与者的uid
	Remark string `bson:"remark"` //备注
}

type ActivityItem struct {
	Aid       uint64    `bson:"aid"`        //活动id
	Players   []*Player `bson:"players"`    //参与的用户
	Title     string    `bson:"title"`      //活动标题
	Des       string    `bson:"des"`        //活动描述
	BeginTime int64     `bson:"begin_time"` //开始时间
	EndTime   int64     `bson:"end_time"`   //结束时间
	OwnerUid  uint64    `bson:"owner_uid"`  //创建者
}

type PluginActivity struct {
	model.Dirty `bson:"-"`      //不入库
	Gid         uint64          `bson:"_id,omitempty"` //所属的群id
	Activities  []*ActivityItem `bson:"activities"`    //活动列表
}

func (mod *PluginActivity) CoName() string {
	return "plugin_activity"
}

func (mod *PluginActivity) Load() error {
	filter := bson.D{{"_id", mod.Gid}}
	col := mongo_helper.GetCol(mod)
	err := col.FindOne(context.TODO(), filter).Decode(mod)
	if errors.Is(err, mongo.ErrNoDocuments) {
		_, err = col.InsertOne(context.TODO(), mod)
		if err != nil {
			return fmt.Errorf("create %s err, gid:%v|err:%v", mod.CoName(), mod.Gid, err)
		}
	} else if err != nil {
		return fmt.Errorf("load %s err, gid:%v|err:%v", mod.CoName(), mod.Gid, err)
	}
	return err
}
func (mod *PluginActivity) Save() error {
	filter := bson.D{{"_id", mod.Gid}}
	col := mongo_helper.GetCol(mod)
	_, err := col.UpdateOne(context.Background(), filter, bson.M{"$set": mod}, options.Update().SetUpsert(true))
	if err != nil {
		tmp, _ := json.Marshal(mod)
		slog.Errorf("save error, save to log file, _id:%v|data:%s", mod.Gid, tmp)
	} else {
		mod.CleanDirty()
	}
	return err
}

func (mod *PluginActivity) Delete() error {
	filter := bson.D{{"_id", mod.Gid}}
	col := mongo_helper.GetCol(mod)
	_, err := col.DeleteOne(context.TODO(), filter)
	return err
}
