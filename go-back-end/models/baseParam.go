package models

import "time"

type baseParam struct{}

func (b *baseParam) Validate() error {
	return nil
}

//Page 分页参数
type Page struct {
	baseParam
	Offset int `json:"offset" form:"offset"`
	Count  int `json:"count" form:"count" binding:"required"`
}

type MetaID struct {
	baseParam
	ID int64 `json:"id" db:"id"`
}

type MetaTime struct {
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
