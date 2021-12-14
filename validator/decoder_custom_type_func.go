package validator

import (
	"time"

	"github.com/haiyiyun/utils/help"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	Decoder.RegisterCustomTypeFunc(func(vals []string) (interface{}, error) {
		d := help.Date{}
		t, err := time.Parse("2006-01-02", vals[0])
		d.Time = t
		return d, err
	}, help.Date{})

	Decoder.RegisterCustomTypeFunc(func(vals []string) (interface{}, error) {
		dt := help.DateTime{}
		t, err := time.Parse("2006-01-02 15:04:05", vals[0])
		dt.Time = t
		return dt, err
	}, help.DateTime{})

	Decoder.RegisterCustomTypeFunc(func(vals []string) (interface{}, error) {
		return primitive.ObjectIDFromHex(vals[0])
	}, primitive.ObjectID{})

}
