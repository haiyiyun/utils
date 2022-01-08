package validator

import (
	"time"

	"github.com/haiyiyun/utils/help"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	Decoder.RegisterCustomTypeFunc(func(vals []string) (interface{}, error) {
		d := help.Date{}
		zoneLayout := "Z07:00"
		localZone := time.Now().Format(zoneLayout)
		timeStr := vals[0] + "T00:00:00" + localZone
		t, err := time.Parse(time.RFC3339, timeStr)
		d.Time = t
		return d, err
	}, help.Date{})

	Decoder.RegisterCustomTypeFunc(func(vals []string) (interface{}, error) {
		dt := help.DateTime{}
		zoneLayout := "Z07:00"
		localZone := time.Now().Format(zoneLayout)
		timeStr := vals[0]
		timeByte := []byte(timeStr)
		timeByte[10] = 'T'
		timeStr = string(timeByte) + localZone
		t, err := time.Parse(time.RFC3339, timeStr)
		dt.Time = t
		return dt, err
	}, help.DateTime{})

	Decoder.RegisterCustomTypeFunc(func(vals []string) (interface{}, error) {
		dateLayout := "2006-01-02"
		dateTimeLayout := "2006-01-02 15:04:05"
		zoneLayout := "Z07:00"
		localZone := time.Now().Format(zoneLayout)
		timeStr := vals[0]

		lenTimeStr := len(timeStr)
		if len(dateLayout) == lenTimeStr {
			timeStr = timeStr + "T00:00:00" + localZone
		} else if len(dateTimeLayout) == lenTimeStr {
			timeByte := []byte(timeStr)
			timeByte[10] = 'T'
			timeStr = string(timeByte) + localZone
		}

		return time.Parse(time.RFC3339, timeStr)
	}, time.Time{})

	Decoder.RegisterCustomTypeFunc(func(vals []string) (interface{}, error) {
		return primitive.ObjectIDFromHex(vals[0])
	}, primitive.ObjectID{})

}
