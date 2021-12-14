package validator

import (
	"time"

	"github.com/haiyiyun/utils/help"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	Decoder.RegisterCustomTypeFunc(func(vals []string) (interface{}, error) {
		return time.Parse("2006-01-02", vals[0])
	}, help.Date{})

	Decoder.RegisterCustomTypeFunc(func(vals []string) (interface{}, error) {
		return time.Parse("2006-01-02 15:04:05", vals[0])
	}, help.DateTime{})

	Decoder.RegisterCustomTypeFunc(func(vals []string) (interface{}, error) {
		return primitive.ObjectIDFromHex(vals[0])
	}, primitive.ObjectID{})

}
