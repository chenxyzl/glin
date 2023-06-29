package outer

import (
	"testing"
)

func TestHeartbeat(t *testing.T) {
	a := &Heartbeat_Request{}
	writePb(t, a)
}

func TestGetUserInfo(t *testing.T) {
	a := &GetUserInfo_Request{}
	writePb(t, a)
}

func TestGetUserInfo_json(t *testing.T) {
	a := &GetUserInfo_Request{}
	writeJson(t, a)
}
func TestGetUserInfo_websocket(t *testing.T) {
	a := &GetUserInfo_Request{}
	writePbPack(t, a)
}

func TestGetUserInfo_websocket_json(t *testing.T) {
	a := &GetUserInfo_Request{}
	writeJsonPack(t, a)
}

func TestSetUserInfo(t *testing.T) {
	a := &SetUserInfo_Request{
		Name: "aaa",
		Icon: "bbb",
		Des:  "ccc",
	}
	writePb(t, a)
}

func TestUpdateUserInfo(t *testing.T) {
	a := &UpdateUserInfo_Request{
		Type:   UpdateUserInfo_Des,
		StrVar: "ccc",
	}
	writePb(t, a)
}

func TestGetUserCard(t *testing.T) {
	a := &GetUserCard_Request{
		RedirectUid: 1690637782555557888,
	}
	writePb(t, a)
}

func TestGetUserCardDetail(t *testing.T) {
	a := &GetUserCardDetail_Request{
		RedirectUid: 1690637782555557888,
	}
	writePb(t, a)
}

func TestGetThirdPlatformConfig(t *testing.T) {
	a := &GetThirdPlatformConfig_Request{}
	writePb(t, a)
}
