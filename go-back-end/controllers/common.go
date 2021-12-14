package controllers

import (
	"web_app/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

//ValidateJSONParam 效验JSON传参，如果参数不符合效验规则并返回响应
func ValidateJSONParam(c *gin.Context, param models.Validator) (msg interface{}, err error) {
	return ValidateParam(c, param, binding.JSON)
}

//ValidateQueryParam 效验Query 传参
func ValidateQueryParam(c *gin.Context, param models.Validator) (msg interface{}, err error) {
	return ValidateParam(c, param, binding.Query)
}

//ValidateParam 接收并效验参数
func ValidateParam(c *gin.Context, param models.Validator, contentType binding.Binding) (msg interface{}, err error) {
	//1.接收参数
	err = c.ShouldBindWith(param, contentType)
	if nil == err {
		return nil, nil
	}
	//2.基本校验
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		//非validator.ValidationErrors类型错误
		return codeMsgMap[CodeInvalidParam], err
	}
	if errs != nil {
		//validator.ValidationErrors类型错误则进行翻译
		return removeTopStruct(errs.Translate(trans)), errs
	}
	//3.自定义校验
	if err := param.Validate(); err != nil {
		return err.Error(), err
	}
	return "", nil
}
