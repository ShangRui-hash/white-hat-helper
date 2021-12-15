package mysql

import "web_app/models"

//AddCompany 添加公司
func AddCompany(name string) (int64, error) {
	sql := `insert into company(name) values(?)`
	ret, err := db.Exec(sql, name)
	if err != nil {
		return 0, err
	}
	return ret.LastInsertId()
}

func GetCompanyByID(id int64) (*models.Company, error) {
	var company models.Company
	sql := `select * from company where id = ?`
	err := db.Get(&company, sql, id)
	return &company, err
}

//GetCompanyByName 根据公司名称获取公司信息
func GetCompanyByName(name string) (*models.Company, error) {
	sql := `select * from company where name = ?`
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

//GetCompanyList 获取公司列表
func GetCompanyList(offset, count int) (companyList []*models.Company, err error) {
	sql := `select * from company limit ?,?`
	err = db.Select(&companyList, sql, offset, count)
	return companyList, err
}

func DeleteCompany(id int) error {
	sql := `delete from company where id = ?`
	_, err := db.Exec(sql, id)
	return err
}

func UpdateCompany(id int64, name string) error {
	sql := `update company set name = ? where id = ?`
	_, err := db.Exec(sql, name, id)
	return err
}
