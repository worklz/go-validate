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
	SetRules(rules map[string]interface{}) (err error)    // 设置验证规则
	AppendRules(rules map[string]interface{}) (err error) // 追加验证规则
	DefineMessages() map[string]string
	GetMessages() (messages map[string]string, err error)
	SetMessages(messages map[string]string) (err error)    // 设置验证提示信息
	AppendMessages(messages map[string]string) (err error) // 追加验证提示信息
	DefineTitles() map[string]string
	GetTitles() (titles map[string]string, err error)
	SetTitles(titles map[string]string) (err error)    // 设置验证字段标题
	AppendTitles(titles map[string]string) (err error) // 追加验证字段标题
	DefineScenes() map[string][]string
	GetScenes() (scenes map[string][]string, err error)
	SetScenes(scenes map[string][]string) (err error)    // 设置验证场景
	AppendScenes(scenes map[string][]string) (err error) // 追加验证场景
	GetDatas() (datas map[string]interface{}, err error)
	SetDatas(datas map[string]interface{}) (err error)
	Check() error
	CheckScene(scene string) error
	GetScene() (scene string, err error)
	HandleDatas(datas map[string]interface{}, scene string) error
}

type Validator struct {
	Rules           map[string]interface{} // 验证规则
	Messages        map[string]string      // 验证提示信息
	Titles          map[string]string      // 验证字段标题
	Scenes          map[string][]string    // 验证场景
	Datas           map[string]interface{} // 验证数据
	Scene           string                 // 当前验证场景
	CheckRules      map[string]interface{} // 当前验证规则
	SystemErrPrefix string                 // 系统错误前缀
	Err             error                  // 错误

	validatorInstance     ValidatorInterface // 验证器实例
	validatorInstancePtr  reflect.Value      // 验证器实例结构体指针的反射值
	validatorInstanceElem reflect.Value      // 验证器实例结构体本身的反射值
}

// 设置验证器实例
func (v *Validator) SetValidatorInstance(validator ValidatorInterface) {
	v.validatorInstance = validator

	// 获取指针的反射值
	validatorInstanceValue := reflect.ValueOf(v.validatorInstance)
	// 检查传入的是否为指针
	if validatorInstanceValue.Kind() != reflect.Ptr || validatorInstanceValue.IsNil() {
		v.SetSystemError("请传入验证器实例指针值！")
	} else {
		// 验证器实例结构体指针的反射值
		v.validatorInstancePtr = validatorInstanceValue
		// 验证器实例结构体本身的反射值
		v.validatorInstanceElem = validatorInstanceValue.Elem()
	}

	// 设置定义的属性
	v.SetRules(v.validatorInstance.DefineRules())
	v.SetMessages(v.validatorInstance.DefineMessages())
	v.SetTitles(v.validatorInstance.DefineTitles())
	v.SetScenes(v.validatorInstance.DefineScenes())
	// 设置验证数据
	v.setDatasByJsonTag()
}

// 设置结构体验证数据，根据json标签
func (v *Validator) setDatasByJsonTag() (err error) {
	err = v.GetError()
	if err != nil {
		return
	}
	datas := make(map[string]interface{})

	// 获取结构体的类型
	typeOf := v.validatorInstanceElem.Type()
	// 遍历结构体的所有字段
	for i := 0; i < v.validatorInstanceElem.NumField(); i++ {
		field := v.validatorInstanceElem.Field(i)
		typeField := typeOf.Field(i)
		// 获取 JSON 标签
		jsonTag := typeField.Tag.Get("json")
		// 解析 JSON 标签，处理可能的选项，如 omitempty
		if commaIndex := strings.Index(jsonTag, ","); commaIndex != -1 {
			jsonTag = jsonTag[:commaIndex]
		}
		// 如果 JSON 标签不为空，则使用该标签作为键
		if jsonTag != "" {
			datas[jsonTag] = field.Interface()
		}
	}

	// 设置验证数据
	v.Datas = datas
	return
}

// 获取系统错误前缀
func (v *Validator) getSystemErrPrefix() (prefix string, err error) {
	err = v.GetError()
	if err != nil {
		return
	}
	if v.SystemErrPrefix == "" {
		v.SystemErrPrefix = "验证器" + packageFullName(v.validatorInstance)
	}
	prefix = v.SystemErrPrefix
	return
}

// 获取验证器错误
func (v *Validator) GetError() error {
	return v.Err
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
	return v.SetError(resErrMsg)
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
	v.Err = newErr
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
	rules = v.Rules
	return
}

// 设置验证规则
func (v *Validator) SetRules(rules map[string]interface{}) (err error) {
	err = v.GetError()
	if err != nil {
		return
	}
	allRules := v.validatorInstance.DefineRules()
	if allRules == nil {
		allRules = map[string]interface{}{}
	}
	for k, v := range rules {
		allRules[k] = v
	}
	v.Rules = allRules
	return
}

// 追加验证规则
func (v *Validator) AppendRules(rules map[string]interface{}) (err error) {
	err = v.GetError()
	if err != nil {
		return
	}
	if rules == nil {
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
	v.Rules = currRules
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
	messages = v.Messages
	return
}

// 设置验证提示信息
func (v *Validator) SetMessages(messages map[string]string) (err error) {
	err = v.GetError()
	if err != nil {
		return
	}
	allMessages := v.validatorInstance.DefineMessages()
	if allMessages == nil {
		allMessages = map[string]string{}
	}
	for k, v := range messages {
		allMessages[k] = v
	}
	v.Messages = allMessages
	return
}

// 追加验证提示信息
func (v *Validator) AppendMessages(messages map[string]string) (err error) {
	err = v.GetError()
	if err != nil {
		return
	}
	if messages == nil {
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
	v.Messages = currMessages
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
	titles = v.Titles
	return
}

// 设置验证字段标题
func (v *Validator) SetTitles(titles map[string]string) (err error) {
	err = v.GetError()
	if err != nil {
		return
	}
	allTitles := v.validatorInstance.DefineTitles()
	if allTitles == nil {
		allTitles = map[string]string{}
	}
	for k, v := range titles {
		allTitles[k] = v
	}
	v.Titles = allTitles
	return
}

// 追加验证字段标题
func (v *Validator) AppendTitles(titles map[string]string) (err error) {
	err = v.GetError()
	if err != nil {
		return
	}
	if titles == nil {
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
	v.Titles = currTitles
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
	scenes = v.Scenes
	return
}

// 获取当前验证场景
func (v *Validator) GetScene() (scene string, err error) {
	err = v.GetError()
	if err != nil {
		return
	}
	scene = v.Scene
	return
}

// 设置验证场景
func (v *Validator) SetScenes(scenes map[string][]string) (err error) {
	err = v.GetError()
	if err != nil {
		return
	}
	allScenes := v.validatorInstance.DefineScenes()
	if allScenes == nil {
		allScenes = map[string][]string{}
	}
	for k, v := range scenes {
		allScenes[k] = v
	}
	v.Scenes = allScenes
	return
}

// 追加验证场景
func (v *Validator) AppendScenes(scenes map[string][]string) (err error) {
	err = v.GetError()
	if err != nil {
		return
	}
	if scenes == nil {
		return
	}

	currScenes, err := v.GetScenes()
	if err != nil {
		return
	}
	if currScenes == nil {
		currScenes = map[string][]string{}
	}
	for k, v := range scenes {
		currScenes[k] = v
	}
	v.Scenes = currScenes
	return
}

// 获取参与验证的数据
func (v *Validator) GetDatas() (datas map[string]interface{}, err error) {
	err = v.GetError()
	if err != nil {
		return
	}
	datas = v.Datas
	return
}

// 设置参与验证的数据（值会同步到对应的json标签属性上）
func (v *Validator) SetDatas(datas map[string]interface{}) (err error) {
	err = v.GetError()
	if err != nil {
		return
	}
	if datas == nil {
		return
	}
	v.Datas = datas

	// 判断属性json标签，并赋值
	// 获取结构体的类型
	t := v.validatorInstanceElem.Type()
	// 遍历结构体的所有字段
	for i := 0; i < v.validatorInstanceElem.NumField(); i++ {
		field := t.Field(i)
		// 获取 JSON 标签
		jsonTag := field.Tag.Get("json")
		// 解析 JSON 标签，处理可能的选项，如 omitempty
		if commaIndex := strings.Index(jsonTag, ","); commaIndex != -1 {
			jsonTag = jsonTag[:commaIndex]
		}
		if jsonTag == "" {
			continue
		}
		jsonTagValue, exists := datas[jsonTag]
		if !exists {
			continue
		}
		// 获取字段的值
		structField := v.validatorInstanceElem.Field(i)
		// 检查字段是否可设置
		if structField.CanSet() {
			// 将传入的值转换为反射值
			val := reflect.ValueOf(jsonTagValue)
			// 检查值的类型是否匹配
			if !val.Type().AssignableTo(structField.Type()) {
				err = v.SetSystemError(fmt.Errorf("属性%s类型%v与传入值类型%v不匹配", field.Name, field.Type, val.Type()))
				return
			}
			// 检查值的类型是否与字段类型匹配
			if val.Type().AssignableTo(structField.Type()) {
				// 设置字段的值
				structField.Set(val)
			}
		}
	}

	return
}

// 设置数据（值会同步到对应的json标签属性上）
func (v *Validator) SetData(key string, value interface{}) (err error) {
	err = v.GetError()
	if err != nil {
		return
	}
	if key == "" {
		return
	}
	datas, err := v.GetDatas()
	if err != nil {
		return
	}
	datas[key] = value
	err = v.SetDatas(datas)
	if err != nil {
		return
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
	v.Err = nil
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
		if scenes != nil {
			if sceneDatas, ok := scenes[scene]; ok {
				for _, dataKey := range sceneDatas {
					dataRules, ok := rules[dataKey]
					if !ok {
						err = v.SetSystemError(fmt.Sprintf("验证场景%s数据%s未定义验证规则！", scene, dataKey))
						return
					}
					checkRules[dataKey] = dataRules
				}
			}
		}
	}
	if len(checkRules) == 0 {
		checkRules = rules
	}
	v.Scene = scene
	v.CheckRules = checkRules
	return nil
}

// 验证指定环境数据
func (v *Validator) CheckScene(scene string) (err error) {
	// 初始化属性
	err = v.initAttr(scene)
	if err != nil {
		return
	}
	// 验证
	err = v.handleCheck()
	return
}

// 验证
func (v *Validator) Check() (err error) {
	// 初始化属性
	err = v.initAttr("")
	if err != nil {
		return
	}
	// 验证
	err = v.handleCheck()
	return
}

// 获取当前验证规则
func (v *Validator) getCheckRules() (rules map[string]interface{}, err error) {
	err = v.GetError()
	if err != nil {
		return
	}
	rules = v.CheckRules
	return
}

// 调用验证器实例的规则方法
func (v *Validator) callValidatorInstanceRuleMethod(methodName string, dataValue interface{}, ruleParam string, datas map[string]interface{}, dataTitle string) (err error) {
	err = v.GetError()
	if err != nil {
		return
	}
	// 获取指定名称的方法
	// 反射值为指针类型，指针可以调用值接收者方法和指针接收者方法
	method := v.validatorInstancePtr.MethodByName(methodName)
	// 判断方法是否有效
	if !method.IsValid() || !method.CanInterface() {
		err = v.SetSystemError(fmt.Sprintf("方法%s不可调用", methodName))
		return
	}

	// 检查参数数量(有一个是接收者)
	methodParamNum := method.Type().NumIn()
	if methodParamNum != 4 {
		err = v.SetSystemError(fmt.Sprintf("方法%s需定义4个参数，但实际有%d个参数", methodName, methodParamNum))
		return
	}
	// 检查参数类型
	no1MethodParamType := method.Type().In(0)
	if no1MethodParamType.Kind() != reflect.Interface || no1MethodParamType.NumMethod() != 0 {
		err = v.SetSystemError(fmt.Sprintf("方法%s的第1个参数类型不正确，需为interface{}", methodName))
		return
	}
	no2MethodParamType := method.Type().In(1)
	if no2MethodParamType.Kind() != reflect.String {
		err = v.SetSystemError(fmt.Sprintf("方法%s的第2个参数类型不正确，需为string", methodName))
		return
	}
	no3MethodParamType := method.Type().In(2)
	if no3MethodParamType.Kind() != reflect.Map || no3MethodParamType.Key().Kind() != reflect.String || no3MethodParamType.Elem().Kind() != reflect.Interface {
		err = v.SetSystemError(fmt.Sprintf("方法%s的第3个参数类型不正确，需为map[string]interface{}", methodName))
		return
	}
	no4MethodParamType := method.Type().In(3)
	if no4MethodParamType.Kind() != reflect.String {
		err = v.SetSystemError(fmt.Sprintf("方法%s的第4个参数类型不正确，需为string", methodName))
		return
	}

	// 检查返回值数量
	methodRetuenNum := method.Type().NumOut()
	if methodRetuenNum != 1 {
		err = v.SetSystemError(fmt.Sprintf("方法%s只需返回1个错误结果，但实际有%d个变量返回", methodName, methodRetuenNum))
		return
	}

	// 检查返回值类型
	no1MethodReturnType := method.Type().Out(0)
	if no1MethodReturnType != reflect.TypeOf((*error)(nil)).Elem() {
		err = v.SetSystemError(fmt.Sprintf("方法%s的返回值类型不正确，需为error", methodName))
		return
	}

	// 准备反射调用所需的参数
	paramValues := []reflect.Value{
		reflect.ValueOf(dataValue),
		reflect.ValueOf(ruleParam),
		reflect.ValueOf(datas),
		reflect.ValueOf(dataTitle),
	}

	// 调用方法并获取返回值
	results := method.Call(paramValues)
	if len(results) > 0 {
		if resErr, ok := results[0].Interface().(error); ok {
			err = resErr
		}
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
	messages, err := v.GetMessages()
	if err != nil {
		return
	}
	titles, err := v.GetTitles()
	if err != nil {
		return
	}
	for dataKey, dataRules := range checkRules {
		if dataRules == nil {
			continue
		}
		dataValue, dataExists := datas[dataKey]
		dataTitle, dataTitleExists := titles[dataKey]
		if !dataTitleExists || dataTitle == "" {
			dataTitle = dataKey
		}
		// 定义的规则字符串
		if dataRuleStr, isStr := dataRules.(string); isStr {
			if dataRuleStr == "" {
				continue
			}
			dataRuleSlice := strings.Split(dataRuleStr, "|")
			// 判断数据是否为空
			if dataRuleSlice[0] != "required" && (!dataExists || isEmpty(dataValue)) {
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
				} else {
					ruleName = dataRule[:colonIndex]
					ruleParam = dataRule[colonIndex+1:]
				}
				// 判断是否为注册的规则
				if rule, ok := Rules[ruleName]; ok {
					err = rule.Check(dataValue, ruleParam, datas, dataTitle)
					if err != nil {
						defineMessage := messages[dataKey+"."+ruleName]
						if defineMessage != "" {
							err = v.SetError(defineMessage)
						}
						return
					}
					continue
				}
				// 判断是否为结构体内可调用方法
				err = v.callValidatorInstanceRuleMethod(ruleName, dataValue, ruleParam, datas, dataTitle)
				if err != nil {
					return
				}
			}

			continue
		}
		// 定义的规则为闭包验证方法
		if dataRuleFun, isFun := dataRules.(func(value interface{}, datas map[string]interface{}, title string) error); isFun {
			err = dataRuleFun(dataValue, datas, dataTitle)
			if err != nil {
				return
			}
			continue
		}
		err = v.SetSystemError(fmt.Sprintf("参数%s验证规则定义需为string或func(value interface{}, datas map[string]interface{}, title string) error类型", dataKey))
		return
	}

	// 验证后处理数据
	scene, err := v.GetScene()
	if err != nil {
		return
	}
	err = v.callValidatorInstanceHandleDatasMethod(datas, scene)
	if err != nil {
		return
	}

	// 设置验证后的数据
	err = v.SetDatas(datas)
	if err != nil {
		return
	}

	return
}

// 调用验证器实例的验证后处理数据方法
func (v *Validator) callValidatorInstanceHandleDatasMethod(datas map[string]interface{}, scene string) (err error) {
	err = v.GetError()
	if err != nil {
		return
	}
	methodName := "HandleDatas"
	// 获取指定名称的方法
	method := v.validatorInstancePtr.MethodByName(methodName)
	// 判断方法是否有效
	if !method.IsValid() || !method.CanInterface() {
		err = v.SetSystemError(fmt.Sprintf("方法%s不可调用", methodName))
		return
	}

	// 检查参数数量
	methodParamNum := method.Type().NumIn()
	if methodParamNum != 2 {
		err = v.SetSystemError(fmt.Sprintf("方法%s需定义2个参数，但实际有%d个参数", methodName, methodParamNum))
		return
	}
	// 检查参数类型
	no1MethodParamType := method.Type().In(0)
	if no1MethodParamType.Kind() != reflect.Map || no1MethodParamType.Key().Kind() != reflect.String || no1MethodParamType.Elem().Kind() != reflect.Interface {
		err = v.SetSystemError(fmt.Sprintf("方法%s的第1个参数类型不正确，需为map[string]interface{}", methodName))
		return
	}
	no2MethodParamType := method.Type().In(1)
	if no2MethodParamType.Kind() != reflect.String {
		err = v.SetSystemError(fmt.Sprintf("方法%s的第2个参数类型不正确，需为string", methodName))
		return
	}

	// 检查返回值数量
	methodRetuenNum := method.Type().NumOut()
	if methodRetuenNum != 1 {
		err = v.SetSystemError(fmt.Sprintf("方法%s只需返回1个错误结果，但实际有%d个变量返回", methodName, methodRetuenNum))
		return
	}

	// 检查返回值类型
	no1MethodReturnType := method.Type().Out(0)
	if no1MethodReturnType != reflect.TypeOf((*error)(nil)).Elem() {
		err = v.SetSystemError(fmt.Sprintf("方法%s的返回值类型不正确，需为error", methodName))
		return
	}

	// 准备反射调用所需的参数
	paramValues := []reflect.Value{
		reflect.ValueOf(datas),
		reflect.ValueOf(scene),
	}

	// 调用方法并获取返回值
	results := method.Call(paramValues)
	if len(results) > 0 {
		if resErr, ok := results[0].Interface().(error); ok {
			err = resErr
		}
	}
	return
}

// 验证后处理数据
func (v *Validator) HandleDatas(datas map[string]interface{}, scene string) error {
	return nil
}
