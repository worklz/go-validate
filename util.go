package validate

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// 预定义的正则表达式
var (
	mobileRegex                   = regexp.MustCompile(`^1[3-9]\d{9}$`)
	emailRegex                    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	chsRegex                      = regexp.MustCompile(`^[\u4e00-\u9fa5]+$`)
	chsAlphaNumRegex              = regexp.MustCompile(`^[\u4e00-\u9fa5a-zA-Z0-9]+$`)
	chsDashRegex                  = regexp.MustCompile(`^[\u4e00-\u9fa5a-zA-Z0-9_-]+$`)
	chsDashSpaceRegex             = regexp.MustCompile(`^[\u4e00-\u9fa5a-zA-Z0-9_ -]+$`)
	alphaNumRegex                 = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	alphaDashRegex                = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	hexColorRegex                 = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)
	dateRegex                     = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	datetimeRegex                 = regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$`)
	yearRegex                     = regexp.MustCompile(`^\d{4}$`)
	yearMonthRegex                = regexp.MustCompile(`^\d{4}-\d{2}$`)
	monthRegex                    = regexp.MustCompile(`^(0?[1-9]|1[0-2])$`)
	timeRegex                     = regexp.MustCompile(`^\d{2}:\d{2}:\d{2}$`)
	commaIntervalChsAlphaNumRegex = regexp.MustCompile(`^[\u4e00-\u9fa5a-zA-Z0-9]+(,[\u4e00-\u9fa5a-zA-Z0-9]+)*$`)
	commaIntervalPositiveIntRegex = regexp.MustCompile(`^[1-9]\d*(,[1-9]\d*)*$`)
	urlRegex                      = regexp.MustCompile(`^https?://.*$`)
	ipRegex                       = regexp.MustCompile(`^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$|^([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}$`)
	uriRegex                      = regexp.MustCompile(`^/.*$`)
)

// 获取实例的完整包名
func packageFullName(instance interface{}) string {
	// 获取实例的反射类型
	instanceType := reflect.TypeOf(instance)
	// 如果传入的是指针，需要获取指针指向的实际类型
	if instanceType.Kind() == reflect.Ptr {
		instanceType = instanceType.Elem()
	}
	return instanceType.PkgPath() + "." + instanceType.Name()
}

// 判断传入的值是否为空
func isEmpty(value interface{}) bool {
	// 获取传入值的反射值对象
	v := reflect.ValueOf(value)
	// 检查反射值是否有效
	if !v.IsValid() {
		return true
	}
	// 根据不同的类型进行判断
	switch v.Kind() {
	case reflect.String:
		// 字符串类型，判断长度是否为 0
		return v.Len() == 0
	case reflect.Array, reflect.Slice, reflect.Map:
		// 数组、切片、映射类型，判断长度是否为 0
		return v.Len() == 0
	case reflect.Ptr, reflect.Interface:
		// 指针、接口类型，判断是否为 nil
		return v.IsNil()
	case reflect.Struct:
		// 结构体类型，检查每个字段是否为空
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			if !isEmpty(field.Interface()) {
				return false
			}
		}
		return true
	default:
		// 其他基本类型，判断是否为零值
		return reflect.DeepEqual(value, reflect.Zero(v.Type()).Interface())
	}
}

// 判断变量是否为数字类型，或者字符串全由数字组成
func isNumeric(value interface{}) bool {
	switch v := value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return true
	case string:
		_, err := strconv.ParseFloat(v, 64)
		return err == nil
	default:
		return false
	}
}

// 检查值是否为正整数
func isPositiveInt(value interface{}) bool {
	if num, ok := value.(int); ok {
		return num > 0
	}
	if str, ok := value.(string); ok {
		num, err := strconv.Atoi(str)
		return err == nil && num > 0
	}
	return false
}

// 检查值是否为非负整数
func isNonnegativeInt(value interface{}) bool {
	if num, ok := value.(int); ok {
		return num >= 0
	}
	if str, ok := value.(string); ok {
		num, err := strconv.Atoi(str)
		return err == nil && num >= 0
	}
	return false
}

// 检查值是否为浮点数
func isFloat(value interface{}) bool {
	if _, ok := value.(float32); ok {
		return true
	}
	if _, ok := value.(float64); ok {
		return true
	}
	if str, ok := value.(string); ok {
		_, err := strconv.ParseFloat(str, 64)
		return err == nil
	}
	return false
}

// 检查值是否为布尔值
func isBool(value interface{}) bool {
	_, ok := value.(bool)
	return ok
}

// 检查数组中的值是否都在规则数组中
func arrayIn(value []interface{}, rule []string) bool {
	for _, v := range value {
		found := false
		for _, r := range rule {
			if fmt.Sprintf("%v", v) == r {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// 计算两个数组的差集
func arrayDiff(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

// 判断值是否为数组
func isArray(value interface{}) bool {
	return reflect.ValueOf(value).Kind() == reflect.Slice
}

// 判断数组元素是否为正整数
func isPositiveIntArray(value []interface{}) bool {
	for _, v := range value {
		if num, ok := v.(int); !ok || num <= 0 {
			return false
		}
	}
	return true
}

// 判断值是否为map
func isMap(value interface{}) bool {
	return reflect.ValueOf(value).Kind() == reflect.Map
}

// 验证日期格式
func isDate(value string) bool {
	_, err := time.Parse("2006-01-02", value)
	return err == nil
}

// 验证日期时间格式
func isDatetime(value string) bool {
	_, err := time.Parse("2006-01-02 15:04:05", value)
	return err == nil
}

// 验证年份格式
func isYear(value string) bool {
	_, err := strconv.Atoi(value)
	return err == nil && len(value) == 4
}

// 验证年月格式
func isYearMonth(value string) bool {
	_, err := time.Parse("2006-01", value)
	return err == nil
}

// 验证月份格式
func isMonth(value string) bool {
	month, err := strconv.Atoi(value)
	return err == nil && month >= 1 && month <= 12
}

// 验证时间格式
func isTime(value string) bool {
	_, err := time.Parse("15:04:05", value)
	return err == nil
}

// 验证时间范围
func isTimeRange(value interface{}, rule string, title string) error {
	valueMap, ok := value.(map[string]interface{})
	if !ok {
		return errors.New(title + "格式错误")
	}
	start, startExists := valueMap["start"].(string)
	end, endExists := valueMap["end"].(string)
	if !startExists && !endExists {
		return errors.New(title + "错误")
	}
	if rule == "" {
		return errors.New("验证规则[timeRange]的参数数据缺失")
	}
	checkFunc := map[string]func(string) bool{
		"date":      isDate,
		"datetime":  isDatetime,
		"year":      isYear,
		"yearMonth": isYearMonth,
		"month":     isMonth,
		"time":      isTime,
	}
	check, exists := checkFunc[rule]
	if !exists {
		return fmt.Errorf("验证规则[timeRange:%s]错误", rule)
	}
	if start != "" && !check(start) {
		return errors.New(title + "开始时间错误")
	}
	if end != "" && !check(end) {
		return errors.New(title + "结束时间错误")
	}
	if start != "" && end != "" {
		startNum, _ := strconv.Atoi(strings.ReplaceAll(start, "[^0-9]", ""))
		endNum, _ := strconv.Atoi(strings.ReplaceAll(end, "[^0-9]", ""))
		if startNum > endNum {
			return errors.New(title + "开始时间不能大于结束时间")
		}
	}
	return nil
}
