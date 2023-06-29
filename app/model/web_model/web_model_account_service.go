package web_model

import (
	"errors"
	"github.com/chenxyzl/glin/uuid"
	"laiya/share/mongo_helper"
)

func GetAccountData(accountId string) (*Account, bool, error) {
	isNew := false
	data, err := mongo_helper.Transaction(func() (interface{}, error) {
		account := &Account{
			Account: accountId,
		}
		err := account.Load()
		if err != nil {
			return nil, err
		}
		if account.Uid == 0 {
			account.Uid = uuid.Generate()
			isNew = true
		}
		err = account.Save()
		if err != nil {
			return nil, err
		}
		return account, nil
	})
	if err != nil {
		return nil, false, err
	}
	account := data.(*Account)
	if account == nil {
		return nil, false, errors.New("data cannot convert to Account")
	}
	return account, isNew, nil
}
