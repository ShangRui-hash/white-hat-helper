package models

type Validator interface {
	Validate() error
}
