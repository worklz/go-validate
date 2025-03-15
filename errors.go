package validate

import "reflect"

// 无效的参数传递
type InvalidParamError struct {
	Message string
	Type    reflect.Type
}

func (e *InvalidParamError) Error() string {
	if e.Type == nil {
		return "validate: (nil)"
	}
	return "validate: (nil " + e.Type.String() + ")"
}
