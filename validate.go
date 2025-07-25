package validate

import (
	"errors"
	"fmt"
	"strings"
)

// 创建验证器
func Create(validator ValidatorInterface) {
	// 设置验证器实例
	validator.InitInstance(validator)
}

// 校验单个变量
// data: 待验证的数据
// rule: 规则
// title: 标题
// messages: 自定义错误信息
func CheckVar(data interface{}, rule string, title string, messages map[string]string) (err error) {
	ruleSlice := strings.Split(rule, "|")
	// 判断数据是否为空
	if ruleSlice[0] != "required" && isEmpty(data) {
		return
	}
	for _, ruleItemStr := range ruleSlice {
		if ruleItemStr == "" {
			continue
		}
		// 获取规则、规则参数
		var ruleName, ruleParam string
		colonIndex := strings.Index(ruleItemStr, ":")
		if colonIndex == -1 {
			ruleName = ruleItemStr
			ruleParam = ""
		} else {
			ruleName = ruleItemStr[:colonIndex]
			ruleParam = ruleItemStr[colonIndex+1:]
		}
		// 判断是否为注册的规则
		if rule, ok := Rules[ruleName]; ok {
			err = rule.Check(data, ruleParam, nil, title)
			if err != nil {
				defineMessage := messages[ruleName]
				if defineMessage != "" {
					err = errors.New(defineMessage)
				}
				return
			}
			continue
		}
		err = fmt.Errorf("验证规则%s未定义", ruleName)
		return
	}
	return
}
