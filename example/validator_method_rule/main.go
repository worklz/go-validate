package main

import (
	"errors"
	"fmt"

	"github.com/worklz/go-validate"
)

type UserLogin struct {
	validate.Validator
}

func (u *UserLogin) DefineRules() map[string]interface{} {
	return map[string]interface{}{
		"username": "required",
		"password": "required",
		"captcha":  "required|IsCaptcha",
	}
}

func (u *UserLogin) DefineTitles() map[string]string {
	return map[string]string{
		"username": "用户名",
		"password": "密码",
		"captcha":  "验证码",
	}
}

// 验证器定义的验证规则
// 要为值接收者方法，才能被父结构体反射执行匹配到
// 首字母需大写，为导出方法
func (u UserLogin) IsCaptcha(value interface{}, param string, datas map[string]interface{}, title string) error {
	if value != "1234" {
		return errors.New(title + "只能为1234")
	}
	return nil
}

func main() {
	userLogin := &UserLogin{}
	validate.Create(userLogin)
	userLogin.SetDatas(map[string]interface{}{
		"username": "admin",
		"password": "123456",
		"captcha":  "123456",
	})
	err := userLogin.Check()
	if err != nil {
		fmt.Printf("登录验证失败！%v\r\n", err)
	} else {
		fmt.Println("登录验证通过")
	}
}
