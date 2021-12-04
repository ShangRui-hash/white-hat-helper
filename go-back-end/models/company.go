package models

type Company struct {
	MetaID
	MetaTime
	Name string `json:"name" db:"name"`
}
