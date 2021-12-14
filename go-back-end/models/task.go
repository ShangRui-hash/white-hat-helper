package models

type Task struct {
	MetaID
	CompanyID int64 `json:"company_id" db:"company_id"`
	MetaTime
	Name         string   `json:"name" db:"name"`
	Status       string   `json:"status"`
	ScanArea     string   ` db:"scan_area"`
	ScanAreaList []string `json:"scan_area"`
	CompanyName  string   `json:"company_name" db:"company_name"`
}
