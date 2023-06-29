package grain

import (
	"fmt"
)

type RobotFactory func(...any) IActor

var robots = map[uint64]RobotFactory{}

func RegisterRobot(uid uint64, robotFactory RobotFactory) {
	if _, ok := robots[uid]; ok {
		panic(fmt.Sprintf("robot uid repeated, uid:%v", uid))
	}
	robots[uid] = robotFactory
}

func GetRobotFactory(robotUid uint64) RobotFactory {
	return robots[robotUid]
}
