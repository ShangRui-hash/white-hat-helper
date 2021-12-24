package controllers

import (
	"web_app/logic"
	"web_app/param"
	"web_app/pkg/validate"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetWebSiteListHandler(c *gin.Context) {
	//1.接收传参
	var param param.ParamGetWebSiteList
	if msg, err := validate.QueryParam(c, &param); err != nil {
		zap.L().Error("validate param failed", zap.Error(err))
		RespErrMsg(c, CodeInvalidParam, msg)
		return
	}
	//2.业务逻辑层
	websites, err := logic.GetWebSiteList(param.CompanyID, param.Offset, param.Count)
	if err != nil {
		zap.L().Error("get website list failed", zap.Error(err))
		RespErr(c, CodeServerBusy)
		return
	}
	RespSuc(c, websites)
}
