package web_old_model

type OldAccount struct {
	Account string `bson:"_id,omitempty"`
	Uid     uint64 `bson:"uid"`
}

type OldPlayer struct {
	Uid  uint64 `bson:"_id,omitempty"`
	Name string `bson:"name"`
	Icon string `bson:"icon"`
	Des  string `bson:"des"`
}
