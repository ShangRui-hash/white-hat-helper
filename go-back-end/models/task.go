package models

import "time"

type MetaID struct {
	ID int64 `json:"id" db:"id"`
}

type MetaTime struct {
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

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
