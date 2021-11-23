package validator

import "encoding/json"

func FormJson(v interface{}, jsonRaw string) (err error) {
	if err = Validate.Var(jsonRaw, "required,json"); err != nil {
		return
	}

	if err = json.Unmarshal([]byte(jsonRaw), &v); err != nil {
		return
	}

	err = Validate.Struct(v)

	return
}
