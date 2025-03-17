package validate

import "errors"

// 验证其规则
type Rule struct {
	Name  string                                                                    // 规则名称
	Check func(value interface{}, param string, datas map[string]interface{}) error // 校验方法
}

// 注册规则
func RegisterRule(rule Rule) (err error) {
	if rule.Name == "" {
		err = errors.New("rule name is empty")
		return
	}
	Rules[rule.Name] = rule
	return
}

// 验证器规则
var Rules = map[string]Rule{}
