package validate

import (
	"reflect"
	"strings"
)

// 获取实例的完整包名
func packageFullName(instance interface{}) string {
	// 获取实例的反射类型
	instanceType := reflect.TypeOf(instance)
	// 如果传入的是指针，需要获取指针指向的实际类型
	if instanceType.Kind() == reflect.Ptr {
		instanceType = instanceType.Elem()
	}
	return instanceType.PkgPath()
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

// 根据 JSON 标签将结构体转换为 map[string]interface{}
func structToMapByJsonTag(obj interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	// 获取传入对象的反射值
	value := reflect.ValueOf(obj)
	// 如果传入的是指针，获取指针指向的值
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	// 检查传入的是否为结构体
	if value.Kind() != reflect.Struct {
		return result
	}
	// 获取结构体的类型
	typeOf := value.Type()
	// 遍历结构体的所有字段
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		typeField := typeOf.Field(i)
		// 获取 JSON 标签
		jsonTag := typeField.Tag.Get("json")
		// 解析 JSON 标签，处理可能的选项，如 omitempty
		if commaIndex := strings.Index(jsonTag, ","); commaIndex != -1 {
			jsonTag = jsonTag[:commaIndex]
		}
		// 如果 JSON 标签不为空，则使用该标签作为键
		if jsonTag != "" {
			result[jsonTag] = field.Interface()
		}
	}
	return result
}
