package outer

import "google.golang.org/protobuf/proto"

type GroupMsg interface {
	proto.Message
	GetGid() uint64
}

type UidMsg interface {
	proto.Message
	GetUid() uint64
}

type RedirectUidMsg interface {
	proto.Message
	GetRedirectUid() uint64
}

type RoomMsg interface {
	proto.Message
	GetRoomId() uint64
}
