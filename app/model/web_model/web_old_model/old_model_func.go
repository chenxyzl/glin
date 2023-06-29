package web_old_model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"laiya/model/web_model/web_old_model/old_mongo"
)

func (old *OldAccount) loadOldAccount() error {
	filter := bson.D{{"_id", old.Account}}
	col := old_mongo.GetCol("account")
	err := col.FindOne(context.TODO(), filter).Decode(old)
	//账号不存在
	//if err != nil {
	//	if err != mongo.ErrNoDocuments {
	//		return nil, fmt.Errorf("load old account err, account:%v|err:%v", account, err)
	//	}
	//	return nil, err
	//}
	//return player, nil
	return err
}

func (oldPlayer *OldPlayer) loadOldPlayer() error {
	filter := bson.D{{"_id", oldPlayer.Uid}}
	col := old_mongo.GetCol("player")
	err := col.FindOne(context.TODO(), filter).Decode(oldPlayer)
	//账号不存在
	//if err != nil {
	//	if err != mongo.ErrNoDocuments {
	//		return nil, fmt.Errorf("load old player err, uid:%v|err:%v", uid, err)
	//	}
	//	return nil, nil
	//}
	//return player, nil
	return err
}

func (oldPlayer *OldPlayer) LoadOldPlayerByAccount(account string) error {
	oldAccount := &OldAccount{Account: account}
	err := oldAccount.loadOldAccount()
	if err != nil {
		return err
	}
	oldPlayer.Uid = oldAccount.Uid
	return oldPlayer.loadOldPlayer()
}
