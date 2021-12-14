package validator

import (
	"regexp"

	"github.com/haiyiyun/validator"
)

func init() {
	Validate.RegisterValidation("chinamobile", func(fl validator.FieldLevel) bool {
		chinaMobileRegexString := `^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$`
		chinaMobileRegex := regexp.MustCompile(chinaMobileRegexString)

		return chinaMobileRegex.MatchString(fl.Field().String())
	})
}
