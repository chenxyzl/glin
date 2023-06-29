package web_model

import (
	"context"
	"errors"
	"github.com/chenxyzl/glin/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"laiya/share/mongo_helper"
)

var _ model.IModel = new(Captcha)

type Captcha struct {
	model.Dirty    `bson:"-"` //不入库
	Ph             string     `bson:"_id,omitempty"`
	LastReqTime    int64      `bson:"last_req_time"`
	Captcha        string     `bson:"captcha"`
	CaptchaExpTime int64      `bson:"captcha_exp_time"`
	LoginTimes     uint64     `bson:"first_login"`
}

func (mod *Captcha) CoName() string {
	return "captcha"
}

func (mod *Captcha) Load() error {
	filter := bson.D{{"_id", mod.Ph}}
	col := mongo_helper.GetCol(mod)
	err := col.FindOne(context.TODO(), filter).Decode(mod)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil
	}
	return err
}

func (mod *Captcha) Save() error {
	filter := bson.D{{"_id", mod.Ph}}
	col := mongo_helper.GetCol(mod)
	_, err := col.UpdateOne(context.Background(), filter, bson.M{"$set": mod}, options.Update().SetUpsert(true))
	if err != nil {
		//todo save到磁盘,有需要从这里做数据恢复
	} else {
		mod.CleanDirty()
	}
	return err
}

func (mod *Captcha) Delete() error {
	filter := bson.D{{"_id", mod.Ph}}
	col := mongo_helper.GetCol(mod)
	_, err := col.DeleteOne(context.TODO(), filter)
	return err
}
