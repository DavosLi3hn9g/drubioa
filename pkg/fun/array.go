package fun

import (
	"reflect"
)

func Struct2Map(obj interface{}, tag string) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	if t.Kind() == reflect.Ptr { //识别指针
		t = t.Elem()
		v = v.Elem()
	}
	var data = make(map[string]interface{})
	if t.Kind() == reflect.Struct { //结构体
		for i := 0; i < t.NumField(); i++ {
			if tag != "" {
				data[t.Field(i).Tag.Get(tag)] = v.Field(i).Interface()
			} else {
				data[t.Field(i).Name] = v.Field(i).Interface()
			}
		}
	}
	return data
}
