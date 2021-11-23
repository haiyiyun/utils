package validator

import (
	"net/url"
)

func Form(v interface{}, values url.Values) error {
	err := Decoder.Decode(v, values)
	if err != nil {
		return err
	}

	return Validate.Struct(v)
}
