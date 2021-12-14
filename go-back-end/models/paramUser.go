package models

//ParamsRegister 注册需要参数
type ParamsRegister struct {
	baseParam
	Password  string `json:"password" binding:"required"`
	ReCAPTCHA string `form:"g_recaptcha_response" json:"g_recaptcha_response" binding:"required"`
}

//ParamLogout 退出登录
type ParamLogout struct {
	baseParam
	Token string `json:"token" binding:"required"`
}

//ParamLogin 用户名，密码登录需要的参数
type ParamLogin struct {
	baseParam
	Username  string `form:"username" json:"username" binding:"required"`
	Password  string `form:"password" json:"password" binding:"required"`
	ReCAPTCHA string `form:"g_recaptcha_response" json:"g_recaptcha_response" binding:"required"`
}
