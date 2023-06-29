package component

import (
	"errors"
	"github.com/chenxyzl/glin/slog"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"laiya/model/sub_model"
	"laiya/proto/code"
	"laiya/proto/outer"
	"laiya/share/global"
)

func (w *WebComponent) HandleGetRobotInfo(req *outer.GetRobotInfo_Request, ctx *gin.Context) (*outer.GetRobotInfo_Reply, code.Code) {
	robot := &sub_model.PlayerHead{Uid: req.GetUid()}
	err := robot.Load()
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, code.Code_RobotNotExist
	}
	return &outer.GetRobotInfo_Reply{
		RobotInfo: &outer.GetRobotInfo_RobotInfo{
			Uid:  robot.Uid,
			Name: robot.Name,
			Icon: robot.Icon,
			Des:  robot.Des,
		},
	}, code.Code_Ok
}

func (w *WebComponent) HandleUpdateRobotInfo(req *outer.UpdateRobotInfo_Request, ctx *gin.Context) (*outer.UpdateRobotInfo_Reply, code.Code) {
	robotInfo := req.GetRobotInfo()
	//check
	if !global.IsRobot(robotInfo.GetUid()) {
		slog.Errorf("not robot uid, uid:%v", robotInfo.GetUid())
		return nil, code.Code_RobotNotExist
	}
	//update
	robot := &sub_model.PlayerHead{Uid: robotInfo.GetUid()}
	if err := robot.UpdateNameIconDes(robotInfo.GetName(), robotInfo.GetIcon(), robotInfo.GetDes()); err != nil {
		slog.Errorf("update robot info err, robot:%v|err:%v", robotInfo.GetUid(), err)
		return nil, code.Code_Error
	}
	//
	slog.Infof("update robot info success, robot:%v", robot)
	return &outer.UpdateRobotInfo_Reply{}, code.Code_Ok
}
