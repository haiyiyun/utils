package help

import (
	"time"
)

type Date struct {
	time.Time
}

func (d Date) String() string {
	layout := "2006-01-02"
	return d.Format(layout)
}

type DateTime struct {
	time.Time
}

func (d DateTime) String() string {
	layout := "2006-01-02 15:04:05"
	return d.Format(layout)
}
