package utils

import (
	"net/url"
	"reflect"
	"strconv"
)

// StructToQuery 将结构体转换为 url.Values
// 使用 `query` 标签
func StructToQuery(param any) url.Values {
	values := url.Values{}
	v := reflect.ValueOf(param)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		val := v.Field(i)
		if field.Anonymous && val.Kind() == reflect.Struct { // 递归处理匿名字段（嵌套结构体）
			nested := StructToQuery(val.Interface())
			for k, vs := range nested {
				for _, v := range vs {
					values.Add(k, v)
				}
			}
			continue
		}

		tag := field.Tag.Get("query")
		if tag == "" {
			tag = field.Name
		}
		if !val.IsValid() || (val.Kind() == reflect.Ptr && val.IsNil()) {
			continue
		}
		var str string
		switch val.Kind() {
		case reflect.Ptr:
			elem := val.Elem()
			switch elem.Kind() {
			case reflect.String:
				str = elem.String()
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				str = strconv.FormatInt(elem.Int(), 10)
			case reflect.Bool:
				str = strconv.FormatBool(elem.Bool())
			}
		case reflect.String:
			if val.String() == "" {
				continue
			}
			str = val.String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			str = strconv.FormatInt(val.Int(), 10)
		case reflect.Bool:
			str = strconv.FormatBool(val.Bool())
		}
		if str != "" {
			values.Set(tag, str)
		}
	}
	return values
}
