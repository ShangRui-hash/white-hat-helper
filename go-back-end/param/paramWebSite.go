package param

type ParamGetWebSiteList struct {
	Page
	CompanyID int64 `json:"company_id" form:"company_id" binding:"required"`
}
