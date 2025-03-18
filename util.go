package validate

import (
	"reflect"
	"unicode"
)

// 获取实例的完整包名
func packageFullName(instance interface{}) string {
	// 获取实例的反射类型
	instanceType := reflect.TypeOf(instance)
	// 如果传入的是指针，需要获取指针指向的实际类型
	if instanceType.Kind() == reflect.Ptr {
		instanceType = instanceType.Elem()
	}
	return instanceType.PkgPath() + "." + instanceType.Name()
}

// 判断传入的值是否为空
func isEmpty(value interface{}) bool {
	// 获取传入值的反射值对象
	v := reflect.ValueOf(value)
	// 检查反射值是否有效
	if !v.IsValid() {
		return true
	}
	// 根据不同的类型进行判断
	switch v.Kind() {
	case reflect.String:
		// 字符串类型，判断长度是否为 0
		return v.Len() == 0
	case reflect.Array, reflect.Slice, reflect.Map:
		// 数组、切片、映射类型，判断长度是否为 0
		return v.Len() == 0
	case reflect.Ptr, reflect.Interface:
		// 指针、接口类型，判断是否为 nil
		return v.IsNil()
	case reflect.Struct:
		// 结构体类型，检查每个字段是否为空
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			if !isEmpty(field.Interface()) {
				return false
			}
		}
		return true
	default:
		// 其他基本类型，判断是否为零值
		return reflect.DeepEqual(value, reflect.Zero(v.Type()).Interface())
	}
}

// 判断是否为数字类型
func isNumericType(i interface{}) bool {
	switch i.(type) {
	case int, int8, int16, int32, int64:
		return true
	case uint, uint8, uint16, uint32, uint64:
		return true
	case float32, float64:
		return true
	default:
		return false
	}
}

// 判断字符串是否全由数字组成
func isNumericString(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// 判断变量是否为数字类型，或者字符串全由数字组成
func isNumeric(v interface{}) bool {
	if num, ok := v.(string); ok {
		return isNumericString(num)
	}
	return isNumericType(v)
}
