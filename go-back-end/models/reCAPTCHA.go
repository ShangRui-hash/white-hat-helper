package models

//GoogleResp 谷歌reCAPTHCA验证码平台返回的响应
type GoogleResp struct {
	Success    bool     `json:"success"`
	ErrorCodes []string `json:"error-codes"`
}
