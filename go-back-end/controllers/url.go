package controllers

import (
	"web_app/logic"
	"web_app/param"
	"web_app/pkg/validate"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func StartURLDirScanHandler(c *gin.Context) {
	//1.接收传参
	var param param.ParamStartURLDirScan
	if msg, err := validate.QueryParam(c, &param); err != nil {
		zap.L().Error("validate param failed", zap.Error(err))
		RespErrMsg(c, CodeInvalidParam, msg)
		return
	}
	//2.业务逻辑层
	err := logic.StartURLDirScan(&param)
	if err != nil {
		zap.L().Error("url dir scan failed", zap.Error(err))
		RespErr(c, CodeServerBusy)
		return
	}
	RespSuc(c, nil)
}

func StopURLDirScanHandler(c *gin.Context) {
	//1.接收传参
	var param param.ParamStopURLDirScan
	if msg, err := validate.QueryParam(c, &param); err != nil {
		zap.L().Error("validate param failed", zap.Error(err))
		RespErrMsg(c, CodeInvalidParam, msg)
		return
	}
	//2.业务逻辑层
	err := logic.StopURLDirScan(&param)
	if err != nil {
		zap.L().Error("stop url dir scan failed", zap.Error(err))
		RespErr(c, CodeServerBusy)
		return
	}
	RespSuc(c, nil)
}

func DeleteURLSubDirHandler(c *gin.Context) {
	//1.接收传参
	var param param.ParamDeleteURLSubDir
	if msg, err := validate.JSONParam(c, &param); err != nil {
		zap.L().Error("validate param failed", zap.Error(err))
		RespErrMsg(c, CodeInvalidParam, msg)
		return
	}
	//2.业务逻辑层
	err := logic.DeleteURLSubDir(&param)
	if err != nil {
		zap.L().Error("delete url sub dir failed", zap.Error(err))
		RespErr(c, CodeServerBusy)
		return
	}
	RespSuc(c, nil)
}
