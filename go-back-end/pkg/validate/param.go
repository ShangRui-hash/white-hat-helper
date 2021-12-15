package validate

import (
	"web_app/pkg/translator"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

//JSONParam 效验JSON传参，如果参数不符合效验规则并返回响应
func JSONParam(c Receiver, param Validator) (msg interface{}, err error) {
	return validateParam(c, param, binding.JSON)
}

//QueryParam 效验Query 传参
func QueryParam(c Receiver, param Validator) (msg interface{}, err error) {
	return validateParam(c, param, binding.Query)
}

//validateParam 接收并效验参数
func validateParam(c Receiver, param Validator, contentType binding.Binding) (msg interface{}, err error) {
	//1.接收参数
	err = c.ShouldBindWith(param, contentType)
	if nil == err {
		return nil, nil
	}
	//2.基本校验
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		//非validator.ValidationErrors类型错误
		return "invalid param", err
	}
	if errs != nil {
		//validator.ValidationErrors类型错误则进行翻译
		return translator.RemoveTopStruct(errs.Translate(translator.GetTranslator())), errs
	}
	//3.自定义校验
	if err := param.Validate(); err != nil {
		return err.Error(), err
	}
	return "", nil
}
