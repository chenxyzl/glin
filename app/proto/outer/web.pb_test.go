package outer

import (
	"laiya/share/global"
	"testing"
)

func TestGetCaptcha(t *testing.T) {
	a := &GetCaptcha_Request{
		Ph: "13020025527",
	}
	writePb(t, a)
}

func TestCheckCaptcha(t *testing.T) {
	a := &CheckCaptcha_Request{
		Ph:      "13020025527",
		Captcha: "LaiyaAppGodCaptcha2023",
	}
	writePb(t, a)
}

func TestUploadIcon(t *testing.T) {
	a := &UploadIcon_Request{}
	writePb(t, a)
}

func TestDeleteAccount(t *testing.T) {
	a := &DeleteAccount_Request{
		Uid: 85445580617577472,
	}
	writePb(t, a)
}

func TestGetRobotInfo(t *testing.T) {
	a := &GetRobotInfo_Request{
		Uid: global.RobotUidOvercooked,
	}
	writePb(t, a)
}

func TestUpdateRobotInfo(t *testing.T) {
	a := &UpdateRobotInfo_Request{
		RobotInfo: &UpdateRobotInfo_RobotInfo{
			Uid:  global.RobotUidOvercooked,
			Name: "二鸭-厨房摇人",
			Icon: "https://laiyahead.oss-cn-shanghai.aliyuncs.com/head_robot_looking_circular_1.png",
			Des:  "一只给饭友拉皮条攒局的AI机器人。如有任何建议，在房间《来鸭用户反馈》中，联系管理员哈",
		},
	}
	writePb(t, a)
}

func TestUpdateRobotInfo_AnimalParty(t *testing.T) {
	a := &UpdateRobotInfo_Request{
		RobotInfo: &UpdateRobotInfo_RobotInfo{
			Uid:  global.RobotUidAnimalParty,
			Name: "二鸭-动物派对摇人",
			Icon: "https://laiyahead.oss-cn-shanghai.aliyuncs.com/head_robot_looking_circular_1.png",
			Des:  "一只给动物派对攒局的AI机器人。如有任何建议，在房间《来鸭用户反馈》中，联系管理员哈",
		},
	}
	writePb(t, a)
}

func TestCheckClientVersion(t *testing.T) {
	a := &CheckClientVersion_Request{
		ClientVersion: "0.0.1",
	}
	writePb(t, a)
}
