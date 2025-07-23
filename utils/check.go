package utils

import "reflect"

func CheckValue[T any](value T) bool {
	v := reflect.ValueOf(value)
	// 判断是否为指针类型
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return false
		}
		// 非nil指针，判断指针指向的值是否为零值
		elem := v.Elem()
		zero := reflect.Zero(elem.Type())
		return !reflect.DeepEqual(elem.Interface(), zero.Interface())
	}
	// 非指针，直接判断是否为零值
	zero := reflect.Zero(v.Type())
	return !reflect.DeepEqual(v.Interface(), zero.Interface())
}
