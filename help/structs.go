package help

import (
	"fmt"
	"reflect"
	"strings"
)

type Struct struct {
	I             interface{}
	compatibleTag string
	useTag        string
}

func NewStruct(i interface{}) Struct {
	return Struct{
		I: i,
	}
}

func (s Struct) GetTypeName() string {
	var typestr string
	typ := reflect.TypeOf(s.I)
	typestr = typ.String()

	lastDotIndex := strings.LastIndex(typestr, ".")
	if lastDotIndex != -1 {
		typestr = typestr[lastDotIndex+1:]
	}

	return typestr
}

func (s Struct) StructName() string {
	v := reflect.TypeOf(s.I)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v.Name()
}

func (s Struct) CompatibleTag(tag string) Struct {
	s.compatibleTag = tag
	return s
}

func (s Struct) UseTag(tag string) Struct {
	s.useTag = tag
	return s
}

// convert struct to map
// s must to be struct, can not be a pointer
func (s Struct) RawStructToMap(snakeCasedKey bool) M {
	v := reflect.ValueOf(s.I)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		panic(fmt.Sprintf("param s must be struct, but got %s", s.I))
	}

	m := M{}
	for i := 0; i < v.NumField(); i++ {
		vf := v.Field(i)
		if vf.CanInterface() {
			vtf := v.Type().Field(i)
			key := ""
			mark := ""
			if snakeCasedKey {
				key = Strings(key).SnakeCasedName()
			} else {
				if s.useTag != "" {
					if vtf.Tag.Get(s.useTag) != "" {
						key = vtf.Tag.Get(s.useTag)
					}
				}

				if key == "" {
					if vtf.Tag.Get("map") != "" {
						key = vtf.Tag.Get("map")
					}
				}

				if key == "" {
					if s.compatibleTag != "" {
						if vtf.Tag.Get(s.compatibleTag) != "" {
							key = vtf.Tag.Get(s.compatibleTag)
						}
					}
				}

				if key == "" {
					key = vtf.Name
				}
			}

			if ks := strings.Split(key, ","); len(ks) > 1 {
				key = ks[0]
				mark = ks[1]
			}

			if key == "-" {
				continue
			}

			if mark == "omitempty" && IsBlank(vf) {
				continue
			}

			val := vf.Interface()
			if vf.Kind() == reflect.Struct {
				valStruct := Struct{I: vf.Interface()}.RawStructToMap(snakeCasedKey)
				if len(valStruct) > 0 {
					val = valStruct
				}
			}

			m[key] = val
		}
	}

	return m
}

// convert struct to map
func (s Struct) StructToMap() M {
	return s.RawStructToMap(false)
}

// convert struct to map
// but struct's field name to snake cased map key
func (s Struct) StructToSnakeKeyMap() M {
	return s.RawStructToMap(true)
}
