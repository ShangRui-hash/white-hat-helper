package param

import (
	"errors"
	"web_app/dao/mysql"

	"go.uber.org/zap"
)

type ParamDeleteCompany struct {
	baseParam
	ID int `json:"id" binding:"required"`
}

type ParamUpdateCompany struct {
	baseParam
	ID   int    `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}

//ParamGetCompanyList 获取公司需要的参数
type ParamGetCompanyList struct {
	baseParam
	Page
}

//ParamAddCompany 添加公司需要的参数
type ParamAddCompany struct {
	baseParam
	Name string `json:"name" binding:"required"`
}

func (param *ParamAddCompany) Validate() error {
	//查重
	exist, err := param.isCompanyExist()
	if err != nil {
		zap.L().Error("add company handler check company exist error", zap.Error(err))
		return err
	}
	if exist {
		return errors.New("公司已存在")
	}
	return nil
}

func (param *ParamAddCompany) isCompanyExist() (exist bool, err error) {
	company, err := mysql.GetCompanyByName(param.Name)
	if err != nil {
		zap.L().Error("get company by name failed", zap.Error(err))
		return true, err
	}
	return company != nil, nil
}
