package models

import "time"

type MetaID struct {
	ID int64 `json:"id" db:"id"`
}

type MetaTime struct {
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
