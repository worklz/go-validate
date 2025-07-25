package main

import (
	"fmt"
	"regexp"

	"github.com/worklz/go-validate"
)

type UserLogin struct {
	validate.Validator
}

func (u *UserLogin) DefineRules() map[string]interface{} {
	return map[string]interface{}{
		"username": "required|alphaNum",
		"password": "required",
		"captcha":  "required",
	}
}

func (u *UserLogin) DefineTitles() map[string]string {
	return map[string]string{
		"username": "用户名",
		"password": "密码",
		"captcha":  "验证码",
	}
}

func main() {
	// 注册自定义规则
	validate.RegisterRules([]validate.Rule{
		{Name: "alphaNum", Fun: func(value interface{}, param string, datas map[string]interface{}, title string) error {
			valueStr, ok := value.(string)
			if !ok {
				return fmt.Errorf("%s不能为空", title)
			}
			// 定义正则表达式模式
			pattern := "^[a-zA-Z0-9]+$"
			// 编译正则表达式
			match, _ := regexp.MatchString(pattern, valueStr)
			if !match {
				return fmt.Errorf("%s只能为字母或数字", title)
			}
			return nil
		}},
	})
	userLogin := &UserLogin{}
	userLogin.InitInstance(userLogin)
	userLogin.SetDatas(map[string]interface{}{
		"username": "管理员",
		"password": "123456",
	})
	err := userLogin.Check()
	if err != nil {
		fmt.Printf("登录验证失败！%v\r\n", err)
	} else {
		fmt.Println("登录验证通过")
	}
}
