package logic

import (
	"web_app/dao/mysql"
	"web_app/dao/redis"
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
		MetaID: models.MetaID{
			ID: id,
		},
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
	//1.查询公司基本信息
	companyList, err := mysql.GetCompanyList(param)
	if err != nil {
		return nil, err
	}

	for i := range companyList {
		//2.查询公司资产数
		count, err := redis.GetAssetCount(companyList[i].ID)
		if err != nil {
			zap.L().Error("redis.GetHostCount failed", zap.Error(err))
			return nil, err
		}
		companyList[i].AssetCount = count
		//3.查询任务数
		taskCount, err := mysql.GetTaskCount(companyList[i].ID)
		if err != nil {
			zap.L().Error("mysql.GetTaskCount failed", zap.Error(err))
			return nil, err
		}
		companyList[i].TaskCount = taskCount
	}
	return companyList, nil

}

//DeleteCompany 删除公司
func DeleteCompany(param *models.ParamDeleteCompany) error {
	return mysql.DeleteCompany(param.ID)
}

//UpdateCompany 更新公司
func UpdateCompany(param *models.ParamUpdateCompany) error {
	return mysql.UpdateCompany(param)
}
