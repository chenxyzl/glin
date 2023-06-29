package outer

import (
	"testing"
)

func TestGetRecommendGroupList(t *testing.T) {
	a := &GetRecommendGroupList_Request{
		Page: 0,
	}
	writePb(t, a)
}

func TestSearchGroup(t *testing.T) {
	a := &SearchGroup_Request{
		Page:   1,
		StrVal: "A",
	}
	writePb(t, a)
}
