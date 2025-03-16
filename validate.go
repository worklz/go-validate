package validate

// type Validate *validate

// type validate struct {
// }

// 校验单个数据
// data			待验证的数据
// func (v *Validate) check(data interface{}, ruleFuns map[string]interface{}, ruleMessages map[string]string) error {

// }

// 创建验证器
func Create(validator ValidatorInterface) {
	// 设置验证器实例
	validator.SetValidatorInstance(validator)
}
