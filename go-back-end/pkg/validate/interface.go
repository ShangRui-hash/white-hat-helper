package validate

import "github.com/gin-gonic/gin/binding"

type Validator interface {
	Validate() error
}

type Receiver interface {
	ShouldBindWith(param interface{}, contentType binding.Binding) error
}
