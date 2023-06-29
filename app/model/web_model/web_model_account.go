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

var _ model.IModel = new(Account)

type Account struct {
	model.Dirty `bson:"-"` //不入库
	Account     string     `bson:"_id,omitempty"`
	Uid         uint64     `bson:"uid"`
}

func (mod *Account) CoName() string {
	return "account"
}

func (mod *Account) Load() error {
	filter := bson.D{{"_id", mod.Account}}
	col := mongo_helper.GetCol(mod)
	err := col.FindOne(context.TODO(), filter).Decode(mod)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil
	}
	return err
}
func (mod *Account) Save() error {
	filter := bson.D{{"_id", mod.Account}}
	col := mongo_helper.GetCol(mod)
	_, err := col.UpdateOne(context.Background(), filter, bson.M{"$set": mod}, options.Update().SetUpsert(true))
	if err != nil {
		tmp, _ := json.Marshal(mod)
		slog.Errorf("save error, save to log file, _id:%v|data:%s", mod.Account, tmp)
	} else {
		mod.CleanDirty()
	}
	return err
}

func (mod *Account) Delete() error {
	filter := bson.D{{"_id", mod.Account}}
	col := mongo_helper.GetCol(mod)
	_, err := col.DeleteOne(context.TODO(), filter)
	return err
}
