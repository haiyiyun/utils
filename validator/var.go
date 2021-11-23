package validator

import (
	"github.com/haiyiyun/validator"
	"github.com/haiyiyun/validator/form"
)

var (
	Decoder  = form.NewDecoder()
	Validate = validator.New()
)
