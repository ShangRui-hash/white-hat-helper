package controllers

import (
	"strings"
	"web_app/logic"
	"web_app/param"
	"web_app/pkg/validate"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//AddCompanyHandler 添加公司
func AddCompanyHandler(c *gin.Context) {
	//1.接收传参
	var param param.ParamAddCompany
	if msg, err := validate.JSONParam(c, &param); err != nil {
		zap.L().Error("add company handler param error", zap.Error(err))
		RespErrMsg(c, CodeInvalidParam, msg)
		return
	}
	//2.业务逻辑层
	company, err := logic.AddCompany(&param)
	if err != nil {
		zap.L().Error("add company handler error", zap.Error(err))
		RespErr(c, CodeServerBusy)
		return
	}
	RespSuc(c, company)
}

//GetCompanyHandler 获取公司列表
func GetCompanyListHandler(c *gin.Context) {
	//1.接收传参
	var param param.ParamGetCompanyList
	if msg, err := validate.QueryParam(c, &param); err != nil {
		zap.L().Error("get company list handler param error", zap.Error(err))
		RespErrMsg(c, CodeInvalidParam, msg)
		return
	}
	//2.业务逻辑层
	companies, err := logic.GetCompanyList(&param)
	if err != nil {
		zap.L().Error("get company list handler error", zap.Error(err))
		RespErr(c, CodeServerBusy)
		return
	}
	RespSuc(c, companies)
}

//DeleteCompanyHandler 删除公司
func DeleteCompanyHandler(c *gin.Context) {
	//1.接收参数
	var param param.ParamDeleteCompany
	if msg, err := validate.JSONParam(c, &param); err != nil {
		zap.L().Error("delete company handler param error", zap.Error(err))
		RespErrMsg(c, CodeInvalidParam, msg)
		return
	}
	//2.业务逻辑层
	err := logic.DeleteCompany(&param)
	if err != nil {
		zap.L().Error("delete company handler error", zap.Error(err))
		RespErr(c, CodeServerBusy)
		return
	}
	RespSuc(c, nil)
}

func UpdateCompanyHandler(c *gin.Context) {
	//1.接收参数
	var param param.ParamUpdateCompany
	if msg, err := validate.JSONParam(c, &param); err != nil {
		zap.L().Error("update company handler param error", zap.Error(err))
		RespErrMsg(c, CodeInvalidParam, msg)
		return
	}
	//2.后端校验
	param.Name = strings.TrimSpace(param.Name)
	//3.业务逻辑层
	err := logic.UpdateCompany(&param)
	if err != nil {
		zap.L().Error("update company handler error", zap.Error(err))
		RespErr(c, CodeServerBusy)
		return
	}
	RespSuc(c, nil)
}
