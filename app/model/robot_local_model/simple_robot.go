package robot_local_model

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/chenxyzl/glin/model"
	"github.com/chenxyzl/glin/share"
	"github.com/chenxyzl/glin/slog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"laiya/share/mongo_helper"
)

var _ model.IModel = new(LocalSimpleRobotModel[ModParams])

var _ ILocalSimpleRobotModel = new(LocalSimpleRobotModel[ModParams])

type ILocalSimpleRobotModel interface {
	model.IModel
	InitUid(robotId string)
	GetRobotId() string
	GetLastWechatInviteBroadcastTime() int64
	SetLastWechatInviteBroadcastTime(time int64)
	GetKindType() share.EntityKind
	GetRobotUid() uint64
}

type ModParams interface {
	CoName() string
	GetKindType() share.EntityKind
	GetRobotUid() uint64
}
type LocalSimpleRobotModel[T ModParams] struct {
	model.Dirty                   `bson:"-"` //不入库
	RobotId                       string     `bson:"_id,omitempty"`
	LastWechatInviteBroadcastTime int64      `bson:"last_wechat_invite_broadcast_time"` //上次微信邀请广播时间
}

func (mod *LocalSimpleRobotModel[T]) CoName() string {
	var a T
	return a.CoName()
}

func (mod *LocalSimpleRobotModel[T]) Load() error {
	filter := bson.D{{"_id", mod.RobotId}}
	col := mongo_helper.GetCol(mod)
	err := col.FindOne(context.TODO(), filter).Decode(mod)
	//账号不存在
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("load %s err, robotId:%v|err:%v", mod.CoName(), mod.RobotId, err)
		}
		//
		_, err = col.InsertOne(context.TODO(), mod)
		if err != nil {
			return fmt.Errorf("create %s err, robotId:%v|err:%v", mod.CoName(), mod.RobotId, err)
		}
	}
	return err
}

func (mod *LocalSimpleRobotModel[T]) Save() error {
	filter := bson.D{{"_id", mod.RobotId}}
	col := mongo_helper.GetCol(mod)
	_, err := col.ReplaceOne(context.Background(), filter, bson.M{"$set": mod}, options.Replace().SetUpsert(true))
	if err != nil {
		tmp, _ := json.Marshal(mod)
		slog.Errorf("save error, save to log file, _id:%v|data:%s", mod.RobotId, tmp)
	} else {
		mod.CleanDirty()
	}
	return err
}

func (mod *LocalSimpleRobotModel[T]) Delete() error {
	filter := bson.D{{"_id", mod.RobotId}}
	col := mongo_helper.GetCol(mod)
	_, err := col.DeleteOne(context.TODO(), filter)
	return err
}

func (mod *LocalSimpleRobotModel[T]) GetRobotId() string {
	return mod.RobotId
}

func (mod *LocalSimpleRobotModel[T]) InitUid(robotId string) {
	mod.RobotId = robotId
}

func (mod *LocalSimpleRobotModel[T]) GetLastWechatInviteBroadcastTime() int64 {
	return mod.LastWechatInviteBroadcastTime
}

func (mod *LocalSimpleRobotModel[T]) SetLastWechatInviteBroadcastTime(time int64) {
	mod.LastWechatInviteBroadcastTime = time
}

func (mod *LocalSimpleRobotModel[T]) GetKindType() share.EntityKind {
	var t T
	return t.GetKindType()
}

func (mod *LocalSimpleRobotModel[T]) GetRobotUid() uint64 {
	var t T
	return t.GetRobotUid()
}
