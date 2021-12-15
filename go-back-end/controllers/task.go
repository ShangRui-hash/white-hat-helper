package controllers

import (
	"strconv"
	"web_app/logic"
	"web_app/param"
	"web_app/pkg/validate"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AddTaskHandler(c *gin.Context) {
	//1.接收传参，后端校验
	var param param.ParamAddTask
	if msg, err := validate.JSONParam(c, &param); err != nil {
		zap.L().Error("add task handler param error", zap.Error(err))
		RespErrMsg(c, CodeInvalidParam, msg)
		return
	}
	//2.业务逻辑
	task, err := logic.AddTask(&param)
	if err != nil {
		zap.L().Error("add task handler error", zap.Error(err))
		RespErr(c, CodeServerBusy)
		return
	}
	RespSuc(c, task)
}

func GetTaskListHandler(c *gin.Context) {
	//1.接收传参，后端校验
	var param param.Page
	if msg, err := validate.QueryParam(c, &param); err != nil {
		zap.L().Error("get task list handler param error", zap.Error(err))
		RespErrMsg(c, CodeInvalidParam, msg)
		return
	}
	//2.业务逻辑
	tasks, err := logic.GetTaskList(&param)
	if err != nil {
		zap.L().Error("get task list handler error", zap.Error(err))
		RespErr(c, CodeServerBusy)
		return
	}
	RespSuc(c, tasks)
}

func DeleteTaskHandler(c *gin.Context) {
	//1.接收传参，后端校验
	var param param.MetaID
	if msg, err := validate.JSONParam(c, &param); err != nil {
		zap.L().Error("delete task handler param error", zap.Error(err))
		RespErrMsg(c, CodeInvalidParam, msg)
		return
	}
	//2.业务逻辑
	err := logic.DeleteTask(param.ID)
	if err != nil {
		zap.L().Error("delete task handler error", zap.Error(err))
		RespErr(c, CodeServerBusy)
		return
	}
	RespSuc(c, nil)
}

func StartTaskHandler(c *gin.Context) {
	//1.接收传参
	taskID := c.Param("id")
	id, err := strconv.ParseInt(taskID, 10, 64)
	if err != nil {
		zap.L().Error("start task handler param error", zap.Error(err))
		RespErr(c, CodeInvalidParam)
		return
	}
	//2.业务逻辑
	if err := logic.StartTask(id); err != nil {
		zap.L().Error("start task handler error", zap.Error(err))
		RespErr(c, CodeServerBusy)
		return
	}
	RespSuc(c, nil)
}
