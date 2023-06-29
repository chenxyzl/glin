package web_model

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/chenxyzl/glin/model"
	"github.com/chenxyzl/glin/slog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"laiya/share/mongo_helper"
)

var _ model.IModel = new(God)

type God struct {
	model.Dirty `bson:"-"` //不入库
	Uid         uint64     `bson:"_id,omitempty"`
	Level       uint64     `bson:"power_level"`
}

func (mod *God) CoName() string {
	return "god"
}

func (mod *God) Load() error {
	filter := bson.D{{"_id", mod.Uid}}
	col := mongo_helper.GetCol(mod)
	err := col.FindOne(context.TODO(), filter).Decode(mod)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil
	}
	return err
}
func (mod *God) Save() error {
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

func (mod *God) Delete() error {
	filter := bson.D{{"_id", mod.Uid}}
	col := mongo_helper.GetCol(mod)
	_, err := col.DeleteOne(context.TODO(), filter)
	return err
}
