package outer

import (
	"testing"
)

func TestGetGroupCard(t *testing.T) {
	a := &GetGroupCard_Request{
		Gid: 117231500950241280,
	}
	writePb(t, a)
}

func TestGetShareGroupCard(t *testing.T) {
	a := &GetShareGroupCard_Request{
		Gid: 117231500950241280,
	}
	writePb(t, a)
}

func TestCreateGroup(t *testing.T) {
	a := &CreateGroup_Request{
		Name:     "AAA",
		Icon:     "BBB",
		Des:      "CCC",
		GameType: "DDD",
	}
	writePb(t, a)
}

func TestFavoriteGroup(t *testing.T) {
	a := &FavoriteGroup_Request{
		Gid: 117231500950241280,
	}
	writePb(t, a)
}

func TestUnfavoriteGroup(t *testing.T) {
	a := &UnfavoriteGroup_Request{
		Gid: 1691366974028554240,
	}
	writePb(t, a)
}

func TestDisbandGroup(t *testing.T) {
	a := &DisbandGroup_Request{
		Gid: 1691366974028554240,
	}
	writePb(t, a)
}

func TestEnterGroup(t *testing.T) {
	a := &EnterGroup_Request{
		Gid: 117231500950241280,
	}
	writePb(t, a)
}

func TestExitGroup(t *testing.T) {
	a := &ExitGroup_Request{
		Gid: 117231500950241280,
	}
	writePb(t, a)
}

func TestUpdateGroupInfo(t *testing.T) {
	a := &UpdateGroupInfo_Request{
		Gid:   117231500950241280,
		Items: []*UpdateGroupInfo_Item{{Type: UpdateGroupInfo_Des, StrVar: "CCCC"}},
	}
	writePb(t, a)
}

func TestSendGroupChatMsg(t *testing.T) {
	a := &SendGroupChatMsg_Request{
		Uid:     105714784397045760,
		Gid:     83663319443876864,
		Content: "嘿嘿嘿",
	}
	writePb(t, a)
}

func TestGetGroupChatMsg(t *testing.T) {
	a := &GetGroupChatMsg_Request{
		Gid:     1691392512571547648,
		BeginId: 0,
	}
	writePb(t, a)
}
