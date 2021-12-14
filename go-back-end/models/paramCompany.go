package models

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
