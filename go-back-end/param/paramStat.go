package param

type ParamGetCompanyStat struct {
	baseParam
	CompanyID int64 `json:"company_id" form:"company_id" binding:"required"`
}
