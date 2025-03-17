package validate

// 创建验证器
func Create(validator ValidatorInterface) {
	// 设置验证器实例
	validator.SetValidatorInstance(validator)
}
