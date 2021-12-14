package models

type Company struct {
	MetaID
	MetaTime
	AssetCount int64  `json:"asset_count"`
	TaskCount  int64  `json:"task_count"`
	Name       string `json:"name" db:"name"`
}
