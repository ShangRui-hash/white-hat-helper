package logic

import (
	"web_app/dao/mysql"
	"web_app/models"

	"go.uber.org/zap"
)

func AddCompany(param *models.ParamAddCompany) (company *models.Company, err error) {
	id, err := mysql.AddCompany(param)
	if err != nil {
		zap.L().Error("mysql.AddCompany failed", zap.Error(err))
		return nil, err
	}
	return &models.Company{
		ID:   id,
		Name: param.Name,
	}, nil
}

func IsCompanyExist(name string) (exist bool, err error) {
	company, err := mysql.GetCompanyByName(name)
	if err != nil {
		zap.L().Error("get company by name failed", zap.Error(err))
		return true, err
	}
	return company != nil, nil
}

func GetCompanyList(param *models.ParamGetCompanyList) ([]*models.Company, error) {
	return mysql.GetCompanyList(param)
}
