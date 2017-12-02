package util

import (
	"fmt"
	"reflect"
	"strings"
)

func StructFieldsMap(x interface{}, tag string) map[string]string {
	typ := reflect.TypeOf(x)
	value := reflect.ValueOf(x)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		value = value.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil
	}
	m := make(map[string]string, typ.NumField())
	for i := 0; i < typ.NumField(); i++ {
		v := value.Field(i)
		vs := fmt.Sprintf("%v", v)
		t := typ.Field(i)
		ftag := t.Tag.Get(tag)
		if ftag != "" {
			if strings.Contains(ftag, "omitempty") {
				zero := reflect.Zero(t.Type).Interface()
				current := v.Interface()
				if reflect.DeepEqual(zero, current) {
					continue
				}
			}
			m[strings.Split(ftag, ",")[0]] = vs
		} else { // no tag
			m[SnakeString(t.Name)] = vs
		}
	}
	return m
}

// snake string, XxYy to xx_yy
func SnakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}
