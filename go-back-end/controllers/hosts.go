package controllers

import (
	"web_app/logic"
	"web_app/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//GetHostListHandler 获取主机列表
func GetHostListHandler(c *gin.Context) {
	//1.接收传参
	var param models.ParamGetHostList
	if msg, err := ValidateQueryParam(c, &param); err != nil {
		zap.L().Error("validate param failed", zap.Error(err))
		RespErrMsg(c, CodeInvalidParam, msg)
		return
	}
	//2.业务逻辑层
	hostList, err := logic.GetHostList(&param)
	if err != nil {
		zap.L().Error("get host list failed", zap.Error(err))
		RespErr(c, CodeServerBusy)
		return
	}
	RespSuc(c, hostList)
}

//GetHostDetailHandler 获取主机详情
func GetHostDetailHandler(c *gin.Context) {
	//1.接收传参
	var param models.ParamGetHostDetail
	if msg, err := ValidateQueryParam(c, &param); err != nil {
		zap.L().Error("validate param failed", zap.Error(err))
		RespErrMsg(c, CodeInvalidParam, msg)
		return
	}
	//2.业务逻辑层
	hostDetail, err := logic.GetHostDetail(param.IP)
	if err != nil {
		zap.L().Error("get host detail failed", zap.Error(err))
		RespErr(c, CodeServerBusy)
		return
	}
	RespSuc(c, hostDetail)
}
