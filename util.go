package validate

import "reflect"

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
