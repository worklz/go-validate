package validate

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type ValidatorInterface interface {
	SetValidatorInstance(validator ValidatorInterface)
	DefineRules() map[string]interface{}
	GetRules() (rules map[string]interface{}, err error)
	SetRules(rules map[string]interface{}) (err error)
	DefineMessages() map[string]string
	GetMessages() (messages map[string]string, err error)
	SetMessages(messages map[string]string) (err error)
	DefineTitles() map[string]string
	GetTitles() (titles map[string]string, err error)
	SetTitles(titles map[string]string) (err error)
	DefineScenes() map[string][]string
	GetScenes() (scenes map[string][]string, err error)
	SetScenes(scenes map[string][]string) (err error)
	Check() error
	HandleData(scene string) error
}

type Validator struct {
	Rules           map[string]interface{} // 验证规则
	Messages        map[string]string      // 验证提示信息
	Titles          map[string]string      // 验证字段标题
	Datas           map[string]interface{} // 验证数据
	CheckRules      map[string]interface{} // 当前验证规则
	SystemErrPrefix string                 // 系统错误前缀
	Err             error                  // 错误

	validatorInstance      ValidatorInterface // 验证器实例
	validatorInstanceValue reflect.Value      //验证器实例反射值
}

// 设置验证器实例
func (v *Validator) SetValidatorInstance(validator ValidatorInterface) {
	v.validatorInstance = validator

	// 获取指针的反射值
	validatorInstanceValue := reflect.ValueOf(v.validatorInstance)
	// 检查传入的是否为指针
	if validatorInstanceValue.Kind() != reflect.Ptr || validatorInstanceValue.IsNil() {
		v.SetSystemError("未通过正确方法实例当前验证器！")
	} else {
		// 获取指针指向的实际对象的反射值
		v.validatorInstanceValue = validatorInstanceValue.Elem()
	}

	// 设置定义的属性
	v.SetRules(v.DefineRules())
	v.SetMessages(v.DefineMessages())
	v.SetTitles(v.DefineTitles())
	v.SetScenes(v.DefineScenes())
}

// 调用验证器实例方法
func (v *Validator) callValidatorInstanceMethod(methodName string, args []interface{}) (err error) {
	err = v.GetError()
	if err != nil {
		return
	}
	// 获取指定名称的方法
	method := v.validatorInstanceValue.MethodByName(methodName)
	// 判断方法是否有效
	if !method.IsValid() {
		err = v.SetSystemError(fmt.Sprintf("方法[%s]不可调用", methodName))
		return
	}

	// 准备反射调用所需的参数
	var reflectArgs []reflect.Value
	for _, arg := range args {
		reflectArgs = append(reflectArgs, reflect.ValueOf(arg))
	}

	// 调用方法并获取返回值
	results := method.Call(reflectArgs)
	if len(results) > 0 {
		if resErr, ok := results[0].Interface().(error); ok {
			err = resErr
		}
	}
	return
}

// 设置验证器实例属性
func (v *Validator) setValidatorInstanceAttr(attr string, value interface{}) (err error) {
	err = v.GetError()
	if err != nil {
		return
	}
	// 尝试获取指定名称的字段
	attrField := v.validatorInstanceValue.FieldByName(attr)
	// 检查字段是否有效且可设置
	if !attrField.IsValid() || !attrField.CanSet() {
		err = v.SetSystemError(fmt.Sprintf("属性[%s]不可设置", attr))
		return
	}
	// 如果传入的值为nil，则设置为该类型的零值
	if value == nil {
		attrField.Set(reflect.Zero(attrField.Type()))
		return
	}
	// 将传入的值转换为反射值
	valueToSet := reflect.ValueOf(value)
	// 检查值的类型是否匹配
	if !valueToSet.Type().AssignableTo(attrField.Type()) {
		err = v.SetSystemError(fmt.Sprintf("属性[%s]类型与传入值[%v]类型[%T]不一致", attr, value, value))
		return
	}
	// 设置属性值
	attrField.Set(valueToSet)
	return
}

// 获取验证器实例属性
func (v *Validator) getValidatorInstanceAttr(attr string) (res interface{}, err error) {
	if attr != "Err" {
		err = v.GetError()
		if err != nil {
			return
		}
	}
	// 尝试获取指定名称的字段
	attrValue := v.validatorInstanceValue.FieldByName(attr)
	// 检查字段是否有效
	if !attrValue.IsValid() {
		err = v.SetSystemError(fmt.Sprintf("属性[%s]无效", attr))
		return
	}
	// 返回属性值
	res = attrValue.Interface()
	return
}

// 获取验证器实例字符串属性
func (v *Validator) getValidatorInstanceStrAttr(attr string) (res string, err error) {
	err = v.GetError()
	if err != nil {
		return
	}
	attrValue, err := v.getValidatorInstanceAttr(attr)
	if err != nil {
		return
	}
	var ok bool
	res, ok = attrValue.(string)
	if !ok {
		err = v.SetSystemError(attr + " 属性类型需为字符串形式")
	}
	return
}

// 获取系统错误前缀
func (v *Validator) getSystemErrPrefix() (prefix string, err error) {
	err = v.GetError()
	if err != nil {
		return
	}
	prefix, err = v.getValidatorInstanceStrAttr("SystemErrPrefix")
	if err != nil {
		return
	}
	if prefix == "" {
		prefix = "验证器" + packageFullName(v.validatorInstance)
		err = v.setValidatorInstanceAttr("SystemErrPrefix", prefix)
		if err != nil {
			err = v.SetSystemError(err)
		}
	}
	return
}

// 获取验证器错误
func (v *Validator) GetError() (err error) {
	value, _ := v.getValidatorInstanceAttr("Err")
	err, _ = value.(error) // 不需要判断格式
	return
}

// 获取系统错误
func (v *Validator) GetSystemError() (err error) {
	err = v.GetError()
	if err == nil {
		return
	}
	errMsg := err.Error()
	errPrefix, _ := v.getSystemErrPrefix()
	if strings.HasPrefix(errMsg, errPrefix) {
		return
	}
	err = nil
	return
}

// 设置系统内部错误
func (v *Validator) SetSystemError(err interface{}) error {
	var resErrMsg string
	if _errMsg, ok := err.(string); ok {
		resErrMsg = _errMsg
	} else if _err, ok := err.(error); ok {
		resErrMsg = _err.Error()
	} else {
		resErrMsg = "未知错误！"
	}

	// 错误前缀
	errPrefix, _ := v.getSystemErrPrefix()
	if !strings.HasPrefix(resErrMsg, errPrefix) {
		resErrMsg = errPrefix + resErrMsg
	}
	return v.SetSystemError(resErrMsg)
}

// 设置验证器错误
func (v *Validator) SetError(err interface{}) error {
	if _err := v.GetError(); _err != nil {
		return _err
	}
	var resErrMsg string
	if _errMsg, ok := err.(string); ok {
		resErrMsg = _errMsg
	} else if _err, ok := err.(error); ok {
		resErrMsg = _err.Error()
	} else {
		resErrMsg = "未知错误！"
	}
	newErr := errors.New(resErrMsg)
	v.setValidatorInstanceAttr("Err", newErr)
	return newErr
}

// 定义验证规则
func (v *Validator) DefineRules() map[string]interface{} {
	return nil
}

// 获取验证规则
func (v *Validator) GetRules() (rules map[string]interface{}, err error) {
	err = v.GetError()
	if err != nil {
		return
	}
	value, err := v.getValidatorInstanceAttr("Rules")
	if err != nil {
		return
	}
	var ok bool
	rules, ok = value.(map[string]interface{})
	if !ok {
		err = v.SetSystemError("Rules 属性类型需为map[string]interface{}形式")
	}
	return
}

// 设置验证规则
func (v *Validator) SetRules(rules map[string]interface{}) (err error) {
	if rules == nil {
		return
	}
	err = v.GetError()
	if err != nil {
		return
	}

	currRules, err := v.GetRules()
	if err != nil {
		return
	}
	if currRules == nil {
		currRules = map[string]interface{}{}
	}
	for k, v := range rules {
		currRules[k] = v
	}

	err = v.setValidatorInstanceAttr("Rules", currRules)
	if err != nil {
		err = v.SetSystemError(err)
	}
	return
}

// 定义验证提示信息
func (v *Validator) DefineMessages() map[string]string {
	return nil
}

// 获取验证提示信息
func (v *Validator) GetMessages() (messages map[string]string, err error) {
	err = v.GetError()
	if err != nil {
		return
	}
	value, err := v.getValidatorInstanceAttr("Messages")
	if err != nil {
		return
	}
	var ok bool
	messages, ok = value.(map[string]string)
	if !ok {
		err = v.SetSystemError("Messages 属性类型需为map[string]string形式")
	}
	return
}

// 设置验证提示信息
func (v *Validator) SetMessages(messages map[string]string) (err error) {
	if messages == nil {
		return
	}
	err = v.GetError()
	if err != nil {
		return
	}

	currMessages, err := v.GetMessages()
	if err != nil {
		return
	}
	if currMessages == nil {
		currMessages = map[string]string{}
	}
	for k, v := range messages {
		currMessages[k] = v
	}
	err = v.setValidatorInstanceAttr("Messages", currMessages)
	if err != nil {
		err = v.SetSystemError(err)
	}
	return
}

// 定义验证字段标题
func (v *Validator) DefineTitles() map[string]string {
	return nil
}

// 获取验证字段标题
func (v *Validator) GetTitles() (titles map[string]string, err error) {
	err = v.GetError()
	if err != nil {
		return
	}
	value, err := v.getValidatorInstanceAttr("Titles")
	if err != nil {
		return
	}
	var ok bool
	titles, ok = value.(map[string]string)
	if !ok {
		err = v.SetSystemError("Titles 属性类型需为map[string]string形式")
	}
	return
}

// 设置验证字段标题
func (v *Validator) SetTitles(titles map[string]string) (err error) {
	if titles == nil {
		return
	}
	err = v.GetError()
	if err != nil {
		return
	}

	currTitles, err := v.GetTitles()
	if err != nil {
		return
	}
	if currTitles == nil {
		currTitles = map[string]string{}
	}
	for k, v := range titles {
		currTitles[k] = v
	}
	err = v.setValidatorInstanceAttr("Titles", currTitles)
	if err != nil {
		err = v.SetSystemError(err)
	}
	return
}

// 定义验证场景，定义要验证的字段
func (v *Validator) DefineScenes() map[string][]string {
	return nil
}

// 获取验证场景
func (v *Validator) GetScenes() (scenes map[string][]string, err error) {
	err = v.GetError()
	if err != nil {
		return
	}
	value, err := v.getValidatorInstanceAttr("Scenes")
	if err != nil {
		return
	}
	var ok bool
	scenes, ok = value.(map[string][]string)
	if !ok {
		err = v.SetSystemError("Scenes 属性类型需为map[string][]string形式")
	}
	return
}

// 设置验证场景
func (v *Validator) SetScenes(scenes map[string][]string) (err error) {
	if scenes == nil {
		return
	}
	err = v.GetError()
	if err != nil {
		return
	}
	err = v.setValidatorInstanceAttr("Scenes", scenes)
	if err != nil {
		err = v.SetSystemError(err)
	}
	return
}

// 获取验证数据
func (v *Validator) GetDatas() (datas map[string]interface{}, err error) {
	err = v.GetError()
	if err != nil {
		return
	}
	value, err := v.getValidatorInstanceAttr("Datas")
	if err != nil {
		return
	}
	var ok bool
	datas, ok = value.(map[string]interface{})
	if !ok {
		err = v.SetSystemError("Datas 属性类型需为map[string]interface{}形式")
	}
	return
}

// 设置验证数据
func (v *Validator) SetDatas(datas map[string]interface{}) (err error) {
	if datas == nil {
		return
	}
	err = v.GetError()
	if err != nil {
		return
	}
	err = v.setValidatorInstanceAttr("Datas", datas)
	if err != nil {
		err = v.SetSystemError(err)
	}
	return
}

// 初始化属性
// scene 当前验证场景
func (v *Validator) initAttr(scene string) (err error) {
	// 获取系统错误（系统错误直接返回）
	sysErr := v.GetSystemError()
	if sysErr != nil {
		err = sysErr
		return
	}
	// 错误置空
	err = v.setValidatorInstanceAttr("Err", nil)
	if err != nil {
		return
	}
	// 当前验证规则
	checkRules := map[string]interface{}{}
	rules, err := v.GetRules()
	if err != nil {
		return
	}
	if scene != "" {
		var scenes map[string][]string
		scenes, err = v.GetScenes()
		if err != nil {
			return
		}
		if scenes == nil {
			err = v.SetSystemError("未定义验证场景数据！")
			return
		}
		sceneDatas, ok := scenes[scene]
		if !ok {
			err = v.SetSystemError(fmt.Sprintf("未定义验证场景%s数据！", scene))
			return
		}
		for _, dataKey := range sceneDatas {
			dataRules, ok := rules[dataKey]
			if !ok {
				err = v.SetSystemError(fmt.Sprintf("验证场景%s数据%s未定义验证规则！", scene, dataKey))
				return
			}
			checkRules[dataKey] = dataRules
		}
	} else {
		checkRules = rules
	}
	err = v.setValidatorInstanceAttr("CheckRules", checkRules)
	if err != nil {
		return
	}
	return nil
}

// 验证指定环境数据
func (v *Validator) CheckScene(scene string) (err error) {
	// 初始化属性
	err = v.initAttr(scene)
	if err != nil {
		return
	}
	return nil
}

// 验证
func (v *Validator) Check() (err error) {
	// 初始化属性
	err = v.initAttr("")
	if err != nil {
		return
	}
	return
}

// 获取当前验证规则
func (v *Validator) getCheckRules() (rules map[string]interface{}, err error) {
	err = v.GetError()
	if err != nil {
		return
	}
	value, err := v.getValidatorInstanceAttr("CheckRules")
	if err != nil {
		return
	}
	var ok bool
	rules, ok = value.(map[string]interface{})
	if !ok {
		err = v.SetSystemError("CheckRules 属性类型需为map[string]interface{}形式")
	}
	return
}

// 处理验证
func (v *Validator) handleCheck() (err error) {
	checkRules, err := v.getCheckRules()
	if err != nil {
		return
	}
	datas, err := v.GetDatas()
	if err != nil {
		return
	}
	for dataKey, dataRules := range checkRules {
		if dataRules == nil {
			continue
		}
		dataValue, dataExists := datas[dataKey]
		// 定义的规则字符串
		dataRuleStr, isStr := dataRules.(string)
		if isStr {
			if dataRuleStr == "" {
				continue
			}
			dataRuleSlice := strings.Split(dataRuleStr, "|")
			// 判断数据是否为空
			if dataRuleSlice[0] == "required" && (!dataExists || isEmpty(dataValue)) {
				continue
			}
			for _, dataRule := range dataRuleSlice {
				if dataRule == "" {
					continue
				}
				// 获取规则、规则参数
				var ruleName, ruleParam string
				colonIndex := strings.Index(dataRule, ":")
				if colonIndex == -1 {
					ruleName = dataRule
					ruleParam = ""
				}
				ruleName = dataRule[:colonIndex]
				ruleParam = dataRule[colonIndex+1:]
				// 判断是否为注册的规则
				if rule, ok := Rules[ruleName]; ok {
					err = rule.Check(dataValue, ruleParam, datas)
					if err != nil {
						return err
					}
					continue
				}
				// 判断是否为结构体内可调用方法
				err = v.callValidatorInstanceMethod(ruleName, []interface{}{dataValue, ruleParam, datas})
				if err != nil {
					return
				}
			}

			continue
		}
		// 定义的验证方法
		_, isFun := dataRules.(func() error)
		if isFun {
			// 定义的验证方法
			err = v.SetError("定义的验证方法fun形式，正在努力开发中...")
			return
			continue
		}
	}
	return
}

// 验证后处理数据
// scene 当前验证场景
func (v *Validator) HandleData(scene string) error {
	return nil
}
