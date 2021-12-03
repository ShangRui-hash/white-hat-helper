package validate

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"web_app/models"

	"go.uber.org/zap"
)

//VerifyReCAPTCHAToken 效验谷歌验证码token的真伪
func VerifyReCAPTCHAToken(token, publicKey string) bool {
	postParams := make(url.Values)
	postParams.Set("secret", publicKey)
	postParams.Set("response", token)
	resp, err := http.PostForm("https://www.recaptcha.net/recaptcha/api/siteverify", postParams)
	if err != nil {
		zap.L().Error("http.PostForm failed", zap.Error(err))
		return false
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		zap.L().Error("ioutil.ReadAll failed", zap.Error(err))
		return false
	}
	//解析响应报文
	var googleResp models.GoogleResp
	err = json.Unmarshal(body, &googleResp)
	if err != nil {
		zap.L().Error("json.Unmarshal failed", zap.Error(err))
		return false
	}
	if !googleResp.Success {
		err = errors.New(googleResp.ErrorCodes[0])
		zap.L().Error("google.Resp.Success is false,hackers on the move!", zap.Error(err))
		return false
	}
	return true
}
