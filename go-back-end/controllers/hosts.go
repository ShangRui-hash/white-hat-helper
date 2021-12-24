package controllers

import (
	"web_app/logic"
	"web_app/param"
	"web_app/pkg/validate"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//GetHostListHandler 获取主机列表
func GetHostListHandler(c *gin.Context) {
	//1.接收传参
	var param param.ParamGetHostList
	if msg, err := validate.QueryParam(c, &param); err != nil {
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

func GetHostBaseInfoHandler(c *gin.Context) {
	//1.接收传参
	var param param.ParamGetHostDetail
	if msg, err := validate.QueryParam(c, &param); err != nil {
		zap.L().Error("validate param failed", zap.Error(err))
		RespErrMsg(c, CodeInvalidParam, msg)
		return
	}
	//2.业务逻辑层
	hostBaseInfo, err := logic.GetHostBaseInfo(param.IP)
	if err != nil {
		zap.L().Error("get host base info failed", zap.Error(err))
		RespErr(c, CodeServerBusy)
		return
	}
	RespSuc(c, hostBaseInfo)
}

func GetHostPortInfoHandler(c *gin.Context) {
	//1.接收传参
	var param param.ParamGetHostDetail
	if msg, err := validate.QueryParam(c, &param); err != nil {
		zap.L().Error("validate param failed", zap.Error(err))
		RespErrMsg(c, CodeInvalidParam, msg)
		return
	}
	//2.业务逻辑层
	ports, err := logic.GetHostPortInfo(param.IP)
	if err != nil {
		zap.L().Error("get host port info failed", zap.Error(err))
		RespErr(c, CodeServerBusy)
		return
	}
	RespSuc(c, ports)
}

//GetHostWebInfoHandler 分页获取主机对应的web服务列表
func GetHostWebInfoHandler(c *gin.Context) {
	//1.接收传参
	var param param.ParamGetWebInfo
	if msg, err := validate.QueryParam(c, &param); err != nil {
		zap.L().Error("validate param failed", zap.Error(err))
		RespErrMsg(c, CodeInvalidParam, msg)
		return
	}
	//2.业务逻辑层
	webs, err := logic.GetHostWebInfo(param.IP, param.Offset, param.Count)
	if err != nil {
		zap.L().Error("get host web info failed", zap.Error(err))
		RespErr(c, CodeServerBusy)
		return
	}
	RespSuc(c, webs)
}
