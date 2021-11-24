package validator

import (
	"context"
	"encoding/json"
)

//必须传入指针
func Json(v interface{}, jsonRaw string) (err error) {
	if err = Validate.Var(jsonRaw, "required,json"); err != nil {
		return
	}

	if err = json.Unmarshal([]byte(jsonRaw), v); err != nil {
		return
	}

	err = Validate.Struct(v)

	return
}

//必须传入指针
func JsonCtx(ctx context.Context, v interface{}, jsonRaw string) (err error) {
	if err = Validate.VarCtx(ctx, jsonRaw, "required,json"); err != nil {
		return
	}

	if err = json.Unmarshal([]byte(jsonRaw), v); err != nil {
		return
	}

	err = Validate.StructCtx(ctx, v)

	return
}
