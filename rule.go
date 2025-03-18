package validate

import (
	"errors"
)

// 验证其规则
type Rule struct {
	Name string                                                                                  // 规则名称
	Fun  func(value interface{}, param string, datas map[string]interface{}, title string) error // 校验方法

	validator ValidatorInterface // 验证器实例
}

// 设置验证器实例
func (r *Rule) SetValidator(validator ValidatorInterface) *Rule {
	r.validator = validator
	return r
}

// 校验
func (r *Rule) Check(value interface{}, param string, datas map[string]interface{}, title string) (err error) {
	err = r.Fun(value, param, datas, title)
	return
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

// 注册多个规则
func RegisterRules(rule []Rule) (err error) {
	for _, v := range rule {
		err = RegisterRule(v)
		if err != nil {
			return
		}
	}
	return
}

// 验证器规则
var Rules = map[string]Rule{
	"required": {
		Name: "required",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			if isEmpty(value) {
				return errors.New(title + "不能为空")
			}
			return nil
		},
	},
	"number": {
		Name: "number",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			if !isNumeric(value) {
				return errors.New(title + "需为数字类型或由数字组成")
			}
			return nil
		},
	},
}
