package validate

type ValidatorInterface interface {
	Rules() map[string]interface{}
	Messages() map[string]string
	Titles() map[string]string
	Scenes() map[string][]string
	HandleData(scene string) error
}

type Validator struct {
}

// 验证规则
func (v *Validator) Rules() map[string]interface{} {
	return nil
}

// 验证提示信息
func (v *Validator) Messages() map[string]string {
	return nil
}

// 验证字段标题
func (v *Validator) Titles() map[string]string {
	return nil
}

// 验证场景，定义要验证的字段
func (v *Validator) Scenes() map[string][]string {
	return nil
}

// 验证后处理数据
// scene 当前验证场景
func (v *Validator) HandleData(scene string) error {
	return nil
}
