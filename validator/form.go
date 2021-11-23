package validator

import (
	"net/url"
)

func FormStruct(v interface{}, values url.Values) (err error) {
	if err = Decoder.Decode(&v, values); err != nil {
		return
	}

	err = Validate.Struct(v)

	return
}
