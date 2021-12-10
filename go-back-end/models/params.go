package models

//ParamsRegister 注册需要参数
type ParamsRegister struct {
	Password  string `json:"password" binding:"required"`
	ReCAPTCHA string `form:"g_recaptcha_response" json:"g_recaptcha_response" binding:"required"`
}

//ParamLogout 退出登录
type ParamLogout struct {
	Token string `json:"token" binding:"required"`
}

//ParamLogin 用户名，密码登录需要的参数
type ParamLogin struct {
	Username  string `form:"username" json:"username" binding:"required"`
	Password  string `form:"password" json:"password" binding:"required"`
	ReCAPTCHA string `form:"g_recaptcha_response" json:"g_recaptcha_response" binding:"required"`
}

//ParamAddCompany 添加公司需要的参数
type ParamAddCompany struct {
	Name string `json:"name" binding:"required"`
}

//Page 分页参数
type Page struct {
	Offset int `json:"offset" form:"offset"`
	Count  int `json:"count" form:"count" binding:"required"`
}

//ParamGetCompanyList 获取公司需要的参数
type ParamGetCompanyList struct {
	Page
}

type ParamDeleteCompany struct {
	ID int `json:"id" binding:"required"`
}

type ParamUpdateCompany struct {
	ID   int    `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type ParamAddTask struct {
	CompanyID int    `json:"company_id" binding:"required"`
	Name      string `json:"name" binding:"required"`
	ScanArea  string `json:"scan_area" binding:"required"`
}

type ParamGetHostList struct {
	CompanyID int64 `json:"company_id" form:"company_id" binding:"required"`
	Page
}

type ParamGetHostDetail struct {
	IP string `json:"ip" form:"ip" binding:"required"`
}
