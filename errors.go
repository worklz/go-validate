package validate

import "reflect"

// 无效的validate参数传递
type InvalidValidateError struct {
	Type reflect.Type
}

// Error returns InvalidValidationError message
func (e *InvalidValidateError) Error() string {

	if e.Type == nil {
		return "validate: (nil)"
	}

	return "validate: (nil " + e.Type.String() + ")"
}
