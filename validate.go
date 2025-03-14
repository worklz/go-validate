package validate

type Validate struct {
}

func New() *Validate {
	v := &Validate{}
	return v
}

// 校验结构体
func (v *Validate) Check(validator ValidatorInterface) (err error) {

	return nil
}

// 校验单个数据
// data			待验证的数据
// func (v *Validate) check(data interface{}, ruleFuns map[string]interface{}, ruleMessages map[string]string) error {

// }
