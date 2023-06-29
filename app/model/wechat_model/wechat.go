package wechat_model

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/chenxyzl/glin/model"
	"github.com/chenxyzl/glin/slog"
	"github.com/eatmoreapple/openwechat"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"laiya/share/mongo_helper"
)

var _ model.IModel = new(Wechat)

type WechatGroup struct {
	RecentNotifyTimes int64             `bson:"recent_notify_times"`
	OpenWechatGroup   *openwechat.Group `bson:"-"` //不入库
}

type Wechat struct {
	model.Dirty `bson:"-"`              //不入库
	Uid         uint64                  `bson:"_id,omitempty"`
	Groups      map[string]*WechatGroup `bson:"-"` //groups 暂时不入库--因为解决不了改名之后的脏数据问题
	Bot         *openwechat.Bot         `bson:"-"` //不入库
}

func (mod *Wechat) CoName() string {
	return "wechat"
}

func (mod *Wechat) Load() error {
	filter := bson.D{{"_id", mod.Uid}}
	col := mongo_helper.GetCol(mod)
	err := col.FindOne(context.TODO(), filter).Decode(mod)
	//账号不存在
	if errors.Is(err, mongo.ErrNoDocuments) {
		_, err = col.InsertOne(context.TODO(), mod)
		if err != nil {
			return fmt.Errorf("create %s err, uid:%v|err:%v", mod.CoName(), mod.Uid, err)
		}
	} else if err != nil {
		return fmt.Errorf("load %s err, uid:%v|err:%v", mod.CoName(), mod.Uid, err)
	}
	return err
}

func (mod *Wechat) Save() error {
	filter := bson.D{{"_id", mod.Uid}}
	col := mongo_helper.GetCol(mod)
	_, err := col.UpdateOne(context.Background(), filter, bson.M{"$set": mod}, options.Update().SetUpsert(true))
	if err != nil {
		tmp, _ := json.Marshal(mod)
		slog.Errorf("save error, save to log file, _id:%v|data:%s", mod.Uid, tmp)
	} else {
		mod.CleanDirty()
	}
	return err
}

func (mod *Wechat) Delete() error {
	return nil
}
