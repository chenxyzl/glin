package outer

import "testing"

func TestGetGroupAllActivity(t *testing.T) {
	a := &GetGroupAllActivity_Request{
		Gid:      83663319443876864,
		Uid:      105714784397045760,
		OnlySelf: false,
	}
	writePb(t, a)
}

func TestCreateGroupActivity(t *testing.T) {
	a := &CreateGroupActivity_Request{
		Gid:           83663319443876864,
		Uid:           105714784397045760,
		ActivityTitle: "aaa",
		ActivityDes:   "xxx",
		BeginTime:     1,
		EndTime:       2,
	}
	writeJson(t, a)
}

func TestUpdateGroupActivity(t *testing.T) {
	a := &UpdateGroupActivity_Request{
		Gid:           83663319443876864,
		Uid:           105714784397045760,
		ActivityId:    117725305176876032,
		ActivityTitle: "bbb",
		ActivityDes:   "yyy",
		BeginTime:     1,
		EndTime:       2,
	}
	writeJson(t, a)
}
