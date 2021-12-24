package param

type ParamGetHostList struct {
	Page
	CompanyID int64 `json:"company_id" form:"company_id" binding:"required"`
}

type ParamGetHostDetail struct {
	baseParam
	IP string `json:"ip" form:"ip" binding:"required"`
}

type ParamGetWebInfo struct {
	Page
	IP string `json:"ip" form:"ip" binding:"required"`
}
