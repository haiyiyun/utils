package validator

import (
	"context"
	"net/url"
)

//必须传入指针
func FormStruct(v interface{}, values url.Values) (err error) {
	if err = Decoder.Decode(v, values); err != nil {
		return
	}

	err = Validate.Struct(v)

	return
}

//必须传入指针
func FormStructCtx(ctx context.Context, v interface{}, values url.Values) (err error) {
	if err = Decoder.Decode(v, values); err != nil {
		return
	}

	err = Validate.StructCtx(ctx, v)

	return
}
