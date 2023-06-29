package outer

import (
	"laiya/proto/common"
	"laiya/share/global"
)

func GetUserType(uid uint64) common.UserDef_UserType {
	if global.IsRobot(uid) {
		return common.UserDef_Robot
	}
	return common.UserDef_User
}
