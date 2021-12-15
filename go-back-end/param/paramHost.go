package param

type ParamGetHostList struct {
	baseParam
	CompanyID int64 `json:"company_id" form:"company_id" binding:"required"`
	Page
}

type ParamGetHostDetail struct {
	baseParam
	IP string `json:"ip" form:"ip" binding:"required"`
}
