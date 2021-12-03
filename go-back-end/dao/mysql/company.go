package mysql

import (
	"web_app/models"
)

//AddCompany 添加公司
func AddCompany(company *models.ParamAddCompany) (int64, error) {
	sql := `insert into company(name) values(?)`
	ret, err := db.Exec(sql, company.Name)
	if err != nil {
		return 0, err
	}
	return ret.LastInsertId()
}

//GetCompanyByName 根据公司名称获取公司信息
func GetCompanyByName(name string) (*models.Company, error) {
	sql := `select id, name from company where name = ?`
	rows, err := db.Query(sql, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		company := &models.Company{}
		err := rows.Scan(&company.ID, &company.Name)
		if err != nil {
			return nil, err
		}
		return company, nil
	}
	return nil, nil
}

func GetCompanyList(param *models.ParamGetCompanyList) (companyList []*models.Company, err error) {
	sql := `select * from company limit ?,?`
	err = db.Select(&companyList, sql, param.Offset, param.Count)
	return companyList, err
}
