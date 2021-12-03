package controllers

type ResCode int64

const (
	CodeSuccess ResCode = 10000 + iota
	CodeInvalidParam
	CodeInvalidUserOrPassword
	CodeServerBusy
	CodeNeedAuth
	CodeInvalidToken
	CodeInvalidOffset
	CodeInvalidCount
	CodeLoginFailed
	CodeWrongOldPassword
	CodeAlreadyExist
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:               "success",
	CodeInvalidParam:          "请求参数错误",
	CodeInvalidUserOrPassword: "用户名或者密码错误,请稍后再试",
	CodeServerBusy:            "服务繁忙",
	CodeNeedAuth:              "需要登录",
	CodeInvalidToken:          "无效的token",
	CodeInvalidOffset:         "offset不正确",
	CodeInvalidCount:          "count不正确",
	CodeLoginFailed:           "用户名或者密码错误",
	CodeWrongOldPassword:      "原密码不正确",
	CodeAlreadyExist:          "已经存在",
}

//Msg 返回错误码对应的错误信息
func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}
