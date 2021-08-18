package help

import (
	"reflect"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Slice struct {
	I interface{}
}

func NewSlice(i interface{}) Slice {
	return Slice{I: i}
}

func (s Slice) CheckItem(item string) bool {
	if ss, ok := s.I.([]string); ok {
		for _, i := range ss {
			if item == i {
				return true
			}
		}
	}

	return false
}

func (s Slice) CheckMustItem(item string, items ...string) bool {
	b := s.CheckItem(item)
	for _, item := range items {
		if b = s.CheckItem(item); !b {
			break
		}
	}

	return b
}

func (s Slice) CheckPartItem(item string, items ...string) bool {
	b := s.CheckItem(item)
	bis := false
	for _, item := range items {
		if bis = s.CheckItem(item); bis {
			break
		}
	}

	return b || bis
}

func (s Slice) ConvInt() []int {
	ii := []int{}
	if ss, ok := s.I.([]string); ok {
		for _, s := range ss {
			if i, e := strconv.Atoi(s); e == nil {
				ii = append(ii, i)
			}
		}
	}

	return ii
}

func (s Slice) ConvBool() []bool {
	bb := []bool{}
	if ss, ok := s.I.([]string); ok {
		for _, s := range ss {
			if b, e := strconv.ParseBool(s); e == nil {
				bb = append(bb, b)
			}
		}
	}

	return bb
}

func (s Slice) ConvObjectID() []primitive.ObjectID {
	ids := []primitive.ObjectID{}
	if ss, ok := s.I.([]string); ok {
		for _, s := range ss {
			if id, e := primitive.ObjectIDFromHex(s); e == nil {
				ids = append(ids, id)
			}
		}
	}

	return ids
}

// convert slice to map
func (s Slice) SliceToMap(is ...interface{}) M {
	m := M{}
	var inter []interface{}
	var interNum int
	switch si := s.I.(type) {
	case []string:
		if len(is) > 0 {
			var ok bool
			inter, ok = is[0].([]interface{})
			if !ok {
				return m
			}
		}

		interNum = len(inter)
		for key, value := range si {
			if key <= interNum {
				v := reflect.ValueOf(inter[key])
				if v.Kind() == reflect.Ptr {
					v = v.Elem()
				}

				m[value] = v.Interface()
			} else {
				m[value] = nil
			}
		}
	}

	return m
}
