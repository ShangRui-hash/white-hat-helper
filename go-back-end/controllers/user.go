package controllers

import (
	"os"
	"web_app/logic"
	"web_app/models"
	"web_app/pkg/validate"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//UserLoginHandler 用户登录
func UserLoginHandler(c *gin.Context) {
	//1.接收参数
	var params models.ParamLogin
	if msg, err := ValidateJSONParam(c, &params); err != nil {
		zap.L().Error("user register with invalid param", zap.Error(err))
		RespErrMsg(c, CodeInvalidParam, msg)
		return
	}
	//2.效验google验证码
	publicKey := os.Getenv("reCAPTCHA_public_key")
	if len(publicKey) == 0 {
		zap.L().Error("环境变量中未设置谷歌验证码的公钥")
		RespErr(c, CodeServerBusy)
		return
	}
	if !validate.VerifyReCAPTCHAToken(params.ReCAPTCHA, publicKey) {
		zap.L().Error("validate.VerifyReCAPTCHAToken(params.ReCAPTCHA, publicKey) failed")
		RespErr(c, CodeServerBusy)
		return
	}
	//3.业务逻辑
	token, err := logic.UserLogin(params)
	if err != nil {
		zap.L().Error("logic.UserLogin failed", zap.Error(err))
		RespErr(c, CodeInvalidUserOrPassword)
		return
	}
	RespSuc(c, token)
}

//LogoutHandler 退出登录
func LogoutHandler(c *gin.Context) {
	//1.接收传参,业务逻辑
	var params models.ParamLogout
	if msg, err := ValidateJSONParam(c, &params); err != nil {
		zap.L().Error("logout with invalid param", zap.Error(err))
		RespErrMsg(c, CodeInvalidParam, msg)
		return
	}
	//2.业务逻辑
	if err := logic.Logout(params.Token); err != nil {
		zap.L().Error("logic.Logout failed", zap.Error(err))
		RespErr(c, CodeServerBusy)
		return
	}
	RespSuc(c, nil)
}

//获取用户信息
func GetUserInfoHandler(c *gin.Context) {
	userid, exist := c.Get(CtxUserIDKey)
	if !exist {
		zap.L().Error("userid don't exist")
		RespErr(c, CodeServerBusy)
		return
	}
	//2.业务逻辑
	userinfo, err := logic.GetUserInfo(userid.(int64))
	if err != nil {
		zap.L().Error("logic.GetUserInfo failed", zap.Error(err))
		RespErr(c, CodeInvalidUserOrPassword)
		return
	}
	RespSuc(c, userinfo)
}
