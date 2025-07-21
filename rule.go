package validate

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
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
				return errors.New(title + "需由数字组成")
			}
			return nil
		},
	},
	"integer": {
		Name: "integer",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			if _, ok := value.(int); ok {
				return nil
			}
			if str, ok := value.(string); ok {
				_, err := strconv.Atoi(str)
				if err == nil {
					return nil
				}
			}
			return errors.New(title + "需为整数")
		},
	},
	"positiveInt": {
		Name: "positiveInt",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			if !isPositiveInt(value) {
				return errors.New(title + "需为正整数")
			}
			return nil
		},
	},
	"nonnegativeInt": {
		Name: "nonnegativeInt",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			if !isNonnegativeInt(value) {
				return errors.New(title + "需为非负整数")
			}
			return nil
		},
	},
	"float": {
		Name: "float",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			if !isFloat(value) {
				return errors.New(title + "需为浮点数")
			}
			return nil
		},
	},
	"boolean": {
		Name: "boolean",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			if !isBool(value) {
				return errors.New(title + "需为布尔值")
			}
			return nil
		},
	},
	"length": {
		Name: "length",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			valStr, ok := value.(string)
			if !ok {
				return errors.New("需为字符串形式")
			}
			if param == "" {
				return errors.New("验证规则[length]错误")
			}
			valLen := strCharNum(valStr)
			if !strings.Contains(param, ",") {
				limitLen, err := strconv.Atoi(param)
				if err != nil || !isPositiveInt(limitLen) {
					return errors.New("验证规则[length]参数需为正整数")
				}
				if valLen != limitLen {
					return errors.New(title + fmt.Sprintf("限制长度%d", limitLen))
				}
				return nil
			}
			limitLenArr := strings.Split(param, ",")
			if len(limitLenArr) != 2 {
				return errors.New("验证规则[length]参数需为“,”间隔的两个正整数")
			}
			limitMinLen, err1 := strconv.Atoi(limitLenArr[0])
			limitMaxLen, err2 := strconv.Atoi(limitLenArr[1])
			if err1 != nil || err2 != nil || !isPositiveInt(limitMinLen) || !isPositiveInt(limitMaxLen) {
				return errors.New("验证规则[length]参数需为“,”间隔的两个正整数")
			}
			if valLen < limitMinLen || valLen > limitMaxLen {
				return errors.New(title + fmt.Sprintf("限制长度区间%d-%d", limitMinLen, limitMaxLen))
			}
			return nil
		},
	},
	"min": {
		Name: "min",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			valStr, ok := value.(string)
			if !ok {
				return errors.New("需为字符串形式")
			}
			minLen, err := strconv.Atoi(param)
			if err != nil {
				return errors.New("验证规则[min]参数错误")
			}
			valLen := strCharNum(valStr)
			if valLen < minLen {
				return errors.New(title + fmt.Sprintf("限制最小长度%d", minLen))
			}
			return nil
		},
	},
	"max": {
		Name: "max",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			valStr, ok := value.(string)
			if !ok {
				return errors.New("需为字符串形式")
			}
			maxLen, err := strconv.Atoi(param)
			if err != nil {
				return errors.New("验证规则[max]参数错误")
			}
			valLen := strCharNum(valStr)
			if valLen > maxLen {
				return errors.New(title + fmt.Sprintf("限制最大长度%d", maxLen))
			}
			return nil
		},
	},
	"in": {
		Name: "in",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			if param == "" {
				return errors.New("验证规则[in]错误")
			}
			rule := strings.Split(param, ",")
			if len(rule) == 0 {
				return errors.New("验证规则[in]错误")
			}
			valStr := fmt.Sprintf("%v", value)
			for _, r := range rule {
				if r == valStr {
					return nil
				}
			}
			return errors.New(title + "错误")
		},
	},
	"notIn": {
		Name: "notIn",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			if param == "" {
				return errors.New("验证规则[notIn]错误")
			}
			rule := strings.Split(param, ",")
			if len(rule) == 0 {
				return errors.New("验证规则[notIn]错误")
			}
			valStr := fmt.Sprintf("%v", value)
			for _, r := range rule {
				if r == valStr {
					return errors.New(title + "错误")
				}
			}
			return nil
		},
	},
	"between": {
		Name: "between",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			if param == "" {
				return errors.New("验证规则[between]错误")
			}
			rule := strings.Split(param, ",")
			if len(rule) != 2 {
				return errors.New("验证规则[between]错误")
			}
			min, err1 := strconv.Atoi(rule[0])
			max, err2 := strconv.Atoi(rule[1])
			if err1 != nil || err2 != nil {
				return errors.New("验证规则[between]参数错误")
			}
			num, err := strconv.Atoi(fmt.Sprintf("%v", value))
			if err != nil {
				return errors.New(title + "错误")
			}
			if num < min || num > max {
				return errors.New(title + "错误")
			}
			return nil
		},
	},
	"notBetween": {
		Name: "notBetween",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			if param == "" {
				return errors.New("验证规则[notBetween]错误")
			}
			rule := strings.Split(param, ",")
			if len(rule) != 2 {
				return errors.New("验证规则[notBetween]错误")
			}
			min, err1 := strconv.Atoi(rule[0])
			max, err2 := strconv.Atoi(rule[1])
			if err1 != nil || err2 != nil {
				return errors.New("验证规则[notBetween]参数错误")
			}
			num, err := strconv.Atoi(fmt.Sprintf("%v", value))
			if err != nil {
				return errors.New(title + "错误")
			}
			if num >= min && num <= max {
				return errors.New(title + "错误")
			}
			return nil
		},
	},
	"eq": {
		Name: "eq",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			if param == "" {
				return errors.New("验证规则[eq]错误")
			}
			valNum, err1 := strconv.Atoi(fmt.Sprintf("%v", value))
			ruleNum, err2 := strconv.Atoi(param)
			if err1 != nil || err2 != nil {
				return errors.New(title + "错误")
			}
			if valNum != ruleNum {
				return errors.New(title + "错误")
			}
			return nil
		},
	},
	"egt": {
		Name: "egt",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			if param == "" {
				return errors.New("验证规则[egt]错误")
			}
			valNum, err1 := strconv.Atoi(fmt.Sprintf("%v", value))
			ruleNum, err2 := strconv.Atoi(param)
			if err1 != nil || err2 != nil {
				return errors.New(title + "错误")
			}
			if valNum < ruleNum {
				return errors.New(title + "错误")
			}
			return nil
		},
	},
	"gt": {
		Name: "gt",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			if param == "" {
				return errors.New("验证规则[gt]错误")
			}
			valNum, err1 := strconv.Atoi(fmt.Sprintf("%v", value))
			ruleNum, err2 := strconv.Atoi(param)
			if err1 != nil || err2 != nil {
				return errors.New(title + "错误")
			}
			if valNum <= ruleNum {
				return errors.New(title + "错误")
			}
			return nil
		},
	},
	"elt": {
		Name: "elt",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			if param == "" {
				return errors.New("验证规则[elt]错误")
			}
			valNum, err1 := strconv.Atoi(fmt.Sprintf("%v", value))
			ruleNum, err2 := strconv.Atoi(param)
			if err1 != nil || err2 != nil {
				return errors.New(title + "错误")
			}
			if valNum > ruleNum {
				return errors.New(title + "错误")
			}
			return nil
		},
	},
	"lt": {
		Name: "lt",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			if param == "" {
				return errors.New("验证规则[lt]错误")
			}
			valNum, err1 := strconv.Atoi(fmt.Sprintf("%v", value))
			ruleNum, err2 := strconv.Atoi(param)
			if err1 != nil || err2 != nil {
				return errors.New(title + "错误")
			}
			if valNum >= ruleNum {
				return errors.New(title + "错误")
			}
			return nil
		},
	},
	"array": {
		Name: "array",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			if _, ok := isSlice(value); !ok {
				return errors.New(title + "错误")
			}
			return nil
		},
	},
	"arrayIn": {
		Name: "arrayIn",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			arr, ok := isSlice(value)
			if !ok {
				return errors.New(title + "类型错误")
			}
			if len(arr) == 0 {
				return errors.New(title + "不能为空")
			}
			if param == "" {
				return errors.New("验证规则[arrayIn]错误")
			}
			ruleArr := strings.Split(param, ",")
			if !arrayIn(arr, ruleArr) {
				return errors.New(title + "错误")
			}
			return nil
		},
	},
	"arrayEmptyOrIn": {
		Name: "arrayEmptyOrIn",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			arr, ok := isSlice(value)
			if !ok {
				return errors.New(title + "类型错误")
			}
			if len(arr) == 0 {
				return nil
			}
			if param == "" || !strings.Contains(param, ",") {
				return errors.New("验证规则[arrayEmptyOrIn]错误")
			}
			ruleArr := strings.Split(param, ",")
			if !arrayIn(arr, ruleArr) {
				return errors.New(title + "错误")
			}
			return nil
		},
	},
	"arrayPositiveInt": {
		Name: "arrayPositiveInt",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			arr, ok := isSlice(value)
			if !ok {
				return errors.New(title + "类型错误")
			}
			if !ok || len(arr) == 0 {
				return errors.New(title + "不能为空")
			}
			if !isPositiveIntArray(arr) {
				return errors.New(title + "需为正整数数组")
			}
			return nil
		},
	},
	"arrayEmptyOrPositiveInt": {
		Name: "arrayEmptyOrPositiveInt",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			arr, ok := isSlice(value)
			if !ok {
				return errors.New(title + "类型错误")
			}
			if ok && len(arr) == 0 {
				return nil
			}
			if !isPositiveIntArray(arr) {
				return errors.New(title + "需为正整数数组")
			}
			return nil
		},
	},
	"arrayNonnegativeInt": {
		Name: "arrayNonnegativeInt",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			arr, ok := isSlice(value)
			if !ok {
				return errors.New(title + "类型错误")
			}
			if !ok || len(arr) == 0 {
				return errors.New(title + "不能为空")
			}
			if !isNonnegativeIntArray(arr) {
				return errors.New(title + "需为非负正整数数组")
			}
			return nil
		},
	},
	"arrayEmptyOrNonnegativeInt": {
		Name: "arrayEmptyOrNonnegativeInt",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			arr, ok := isSlice(value)
			if !ok {
				return errors.New(title + "类型错误")
			}
			if ok && len(arr) == 0 {
				return nil
			}
			if !isNonnegativeIntArray(arr) {
				return errors.New(title + "需为非负正整数数组")
			}
			return nil
		},
	},
	"mapHas": {
		Name: "mapHas",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			if !isMap(value) {
				return errors.New(title + "格式错误")
			}
			m, ok := value.(map[string]interface{})
			if !ok || len(m) == 0 {
				return errors.New(title + "不能为空")
			}
			if param == "" || !strings.Contains(param, ",") {
				return errors.New("验证规则[mapHas]错误")
			}
			rule := strings.Split(param, ",")
			valueKeys := make([]string, 0, len(m))
			for k := range m {
				valueKeys = append(valueKeys, k)
			}
			if len(arrayDiff(valueKeys, rule)) > 0 {
				return errors.New(title + "错误")
			}
			return nil
		},
	},
	"mapEmptyOrHas": {
		Name: "mapEmptyOrHas",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			if !isMap(value) {
				return errors.New(title + "格式错误")
			}
			m, ok := value.(map[string]interface{})
			if ok && len(m) == 0 {
				return nil
			}
			if param == "" || !strings.Contains(param, ",") {
				return errors.New("验证规则[mapEmptyOrHas]错误")
			}
			rule := strings.Split(param, ",")
			valueKeys := make([]string, 0, len(m))
			for k := range m {
				valueKeys = append(valueKeys, k)
			}
			if len(arrayDiff(valueKeys, rule)) > 0 {
				return errors.New(title + "错误")
			}
			return nil
		},
	},
	"arrayItemHas": {
		Name: "arrayItemHas",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			arr, ok := isSlice(value)
			if !ok {
				return errors.New(title + "类型错误")
			}
			if len(arr) == 0 {
				return errors.New(title + "不能为空")
			}
			if param == "" {
				return errors.New("验证规则[arrayItemHas]错误")
			}
			rule := strings.Split(param, ",")
			for i, vv := range arr {
				if !isMap(vv) {
					return fmt.Errorf("%s[%d]格式错误", title, i)
				}
				m, ok := vv.(map[string]interface{})
				if !ok || len(m) == 0 {
					return fmt.Errorf("%s[%d]不能为空", title, i)
				}
				valueKeys := make([]string, 0, len(m))
				for k := range m {
					valueKeys = append(valueKeys, k)
				}
				if len(arrayDiff(valueKeys, rule)) > 0 {
					return fmt.Errorf("%s[%d]错误", title, i)
				}
			}
			return nil
		},
	},
	"arrayEmptyOrItemHas": {
		Name: "arrayEmptyOrItemHas",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			arr, ok := isSlice(value)
			if !ok {
				return errors.New(title + "类型错误")
			}
			if len(arr) == 0 {
				return nil
			}
			if param == "" {
				return errors.New("验证规则[arrayEmptyOrItemHas]错误")
			}
			rule := strings.Split(param, ",")
			for i, vv := range arr {
				if !isMap(vv) {
					return fmt.Errorf("%s[%d]格式错误", title, i)
				}
				m, ok := vv.(map[string]interface{})
				if !ok || len(m) == 0 {
					return fmt.Errorf("%s[%d]不能为空", title, i)
				}
				valueKeys := make([]string, 0, len(m))
				for k := range m {
					valueKeys = append(valueKeys, k)
				}
				if len(arrayDiff(valueKeys, rule)) > 0 {
					return fmt.Errorf("%s[%d]错误", title, i)
				}
			}
			return nil
		},
	},
	"mobile": {
		Name: "mobile",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			str, ok := value.(string)
			if !ok || !mobileRegex.MatchString(str) {
				return errors.New(title + "需为11位有效手机格式")
			}
			return nil
		},
	},
	"email": {
		Name: "email",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			str, ok := value.(string)
			if !ok || !emailRegex.MatchString(str) {
				return errors.New(title + "需为有效邮箱格式")
			}
			return nil
		},
	},
	"chs": {
		Name: "chs",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			str, ok := value.(string)
			if !ok || !chsRegex.MatchString(str) {
				return errors.New(title + "只能是汉字")
			}
			return nil
		},
	},
	"chsAlphaNum": {
		Name: "chsAlphaNum",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			str, ok := value.(string)
			if !ok || !chsAlphaNumRegex.MatchString(str) {
				return errors.New(title + "只能是汉字/字母/数字")
			}
			return nil
		},
	},
	"chsDash": {
		Name: "chsDash",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			str, ok := value.(string)
			if !ok || !chsDashRegex.MatchString(str) {
				return errors.New(title + "只能是汉字/字母/数字/下划线_/破折号-")
			}
			return nil
		},
	},
	"chsDashSpace": {
		Name: "chsDashSpace",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			str, ok := value.(string)
			if !ok || !chsDashSpaceRegex.MatchString(str) {
				return errors.New(title + "只能是汉字、字母、数字、下划线_、短横线-及空格组合")
			}
			return nil
		},
	},
	"chsDashChar": {
		Name: "chsDashChar",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			str, ok := value.(string)
			fmt.Println(123, str)
			if !ok || !chsDashCharRegex.MatchString(str) {
				return errors.New(title + "只能是汉字、字母、数字、下划线_、短横线-及中文符号组合")
			}
			return nil
		},
	},
	"alphaNum": {
		Name: "alphaNum",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			str, ok := value.(string)
			if !ok || !alphaNumRegex.MatchString(str) {
				return errors.New(title + "只能是字母/数字")
			}
			return nil
		},
	},
	"alphaDash": {
		Name: "alphaDash",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			str, ok := value.(string)
			if !ok {
				return errors.New(title + "格式错误")
			}
			if !alphaDashRegex.MatchString(str) {
				return errors.New(title + "只能是字母/数字/下划线_/短横线-")
			}
			return nil
		},
	},
	"hexColor": {
		Name: "hexColor",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			str, ok := value.(string)
			if !ok {
				return errors.New(title + "格式错误")
			}
			if !hexColorRegex.MatchString(str) {
				return errors.New(title + "错误")
			}
			return nil
		},
	},
	"date": {
		Name: "date",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			str, ok := value.(string)
			if !ok {
				return errors.New(title + "格式错误")
			}
			if !isDate(str) {
				return errors.New(title + "错误")
			}
			return nil
		},
	},
	"datetime": {
		Name: "datetime",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			str, ok := value.(string)
			if !ok {
				return errors.New(title + "格式错误")
			}
			if !isDatetime(str) {
				return errors.New(title + "错误")
			}
			return nil
		},
	},
	"year": {
		Name: "year",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			str, ok := value.(string)
			if !ok {
				return errors.New(title + "格式错误")
			}
			if !isYear(str) {
				return errors.New(title + "错误")
			}
			return nil
		},
	},
	"yearMonth": {
		Name: "yearMonth",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			str, ok := value.(string)
			if !ok {
				return errors.New(title + "格式错误")
			}
			if !isYearMonth(str) {
				return errors.New(title + "错误")
			}
			return nil
		},
	},
	"month": {
		Name: "month",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			str, ok := value.(string)
			if !ok {
				return errors.New(title + "格式错误")
			}
			if !isMonth(str) {
				return errors.New(title + "错误")
			}
			return nil
		},
	},
	"time": {
		Name: "time",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			str, ok := value.(string)
			if !ok {
				return errors.New(title + "格式错误")
			}
			if !isTime(str) {
				return errors.New(title + "错误")
			}
			return nil
		},
	},
	"timeRange": {
		Name: "timeRange",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			return isTimeRange(value, param, title)
		},
	},
	"commaIntervalChsAlphaNum": {
		Name: "commaIntervalChsAlphaNum",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			str, ok := value.(string)
			if !ok {
				return errors.New(title + "格式错误")
			}
			if !commaIntervalChsAlphaNumRegex.MatchString(str) {
				return errors.New(title + "错误")
			}
			return nil
		},
	},
	"commaIntervalPositiveInt": {
		Name: "commaIntervalPositiveInt",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			str, ok := value.(string)
			if !ok {
				return errors.New(title + "格式错误")
			}
			if !commaIntervalPositiveIntRegex.MatchString(str) {
				return errors.New(title + "错误")
			}
			return nil
		},
	},
	"url": {
		Name: "url",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			str, ok := value.(string)
			if !ok {
				return errors.New(title + "格式错误")
			}
			if !urlRegex.MatchString(str) {
				return errors.New(title + "地址格式错误")
			}
			return nil
		},
	},
	"urls": {
		Name: "urls",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			arr, ok := isSlice(value)
			if !ok {
				return errors.New(title + "类型错误")
			}
			for i, v := range arr {
				str, ok := v.(string)
				if !ok {
					return fmt.Errorf("%s第%d个地址格式错误", title, i+1)
				}
				if !urlRegex.MatchString(str) {
					return fmt.Errorf("%s第%d个地址格式错误", title, i+1)
				}
			}
			return nil
		},
	},
	"ip": {
		Name: "ip",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			str, ok := value.(string)
			if !ok {
				return errors.New(title + "格式错误")
			}
			if !ipRegex.MatchString(str) {
				return errors.New(title + "格式错误")
			}
			return nil
		},
	},
	"uri": {
		Name: "uri",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			str, ok := value.(string)
			if !ok {
				return errors.New(title + "格式错误")
			}
			if !uriRegex.MatchString(str) {
				return errors.New(title + "地址格式错误")
			}
			return nil
		},
	},
	"json": {
		Name: "json",
		Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			str, ok := value.(string)
			if !ok {
				return errors.New(title + "格式错误")
			}
			var js interface{}
			if json.Unmarshal([]byte(str), &js) != nil {
				return errors.New(title + "需为json类型字符串")
			}
			js = nil
			return nil
		},
	},
}
