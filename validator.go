package validate

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type ValidatorInterface interface {
	Rules() map[string]interface{}
	Messages() map[string]string
	Titles() map[string]string
	Scenes() map[string][]string
	Check() error
	HandleData(scene string) error
}

type Validator struct {
	Scene     string // 当前验证环境
	ErrPrefix string // 错误前缀
	Err       error  // 错误

	validatorInstance      ValidatorInterface // 验证器实例
	validatorInstanceValue reflect.Value      //验证器实例反射值
}

// 设置验证器实例
func (v *Validator) setValidatorInstance(validator ValidatorInterface) {
	v.validatorInstance = validator

	// 获取指针的反射值
	validatorInstanceValue := reflect.ValueOf(v.validatorInstance)
	// 检查传入的是否为指针
	if validatorInstanceValue.Kind() != reflect.Ptr || validatorInstanceValue.IsNil() {
		v.SetError("未通过正确方法实例当前验证器！")
	} else {
		// 获取指针指向的实际对象的反射值
		v.validatorInstanceValue = validatorInstanceValue.Elem()
	}
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
		err = v.SetError(fmt.Sprintf("属性[%s]不可设置", attr))
		return
	}
	// 将传入的值转换为反射值
	valueToSet := reflect.ValueOf(value)
	// 检查值的类型是否匹配
	if !valueToSet.Type().AssignableTo(attrField.Type()) {
		err = v.SetError(fmt.Sprintf("属性[%s]类型与传入值[%v]类型[%T]不一致", attr, value, value))
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
		err = v.SetError(fmt.Sprintf("属性[%s]无效", attr))
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
		err = v.SetError(attr + " 属性类型需为字符串形式")
	}
	return
}

// 获取验证器实例的错误前缀
func (v *Validator) getErrPrefix() (prefix string, err error) {
	err = v.GetError()
	if err != nil {
		return
	}
	prefix, err = v.getValidatorInstanceStrAttr("ErrPrefix")
	return
}

// 设置验证器实例的错误前缀
func (v *Validator) setErrPrefix(prefix string) (err error) {
	err = v.GetError()
	if err != nil {
		return
	}
	err = v.setValidatorInstanceAttr("ErrPrefix", prefix)
	if err != nil {
		err = v.SetError(err)
	}
	return
}

// 获取验证器错误
func (v *Validator) GetError() (err error) {
	value, _ := v.getValidatorInstanceAttr("Err")
	err, _ = value.(error) // 不需要判断格式
	return
}

// 设置系统内部错误
func (v *Validator) SetSystemError(msg string) error {
	// 错误前缀
	errPrefix, _ := v.getErrPrefix()
	if errPrefix == "" {
		errPrefix = "验证器" + packageFullName(v.validatorInstance)
		v.setErrPrefix(errPrefix)
	}

	if !strings.HasPrefix(msg, errPrefix) {
		msg = errPrefix + msg
	}
	return v.SetError(msg)
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

// 初始化属性
// scene 当前验证场景
func (v *Validator) initAttr(scene string) error {
	// 错误置空
	v.setValidatorInstanceAttr("Err", nil)
	// 当前验证环境
	v.setValidatorInstanceAttr("Scene", scene)
	return nil
}

// 验证指定环境数据
func (v *Validator) CheckScene(scene string) error {
	// 初始化属性
	v.initAttr(scene)
	return nil
}

// 验证
func (v *Validator) Check() error {
	// 初始化属性
	v.initAttr("")
	return nil
}

// 验证后处理数据
// scene 当前验证场景
func (v *Validator) HandleData(scene string) error {
	return nil
}
