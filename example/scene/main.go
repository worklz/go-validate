package main

import (
	"fmt"

	"github.com/worklz/go-validate"
)

type UserLogin struct {
	validate.Validator
	Username string `json:"username"`
	Password string `json:"password"`
	Captcha  string `json:"captcha"`
}

func (u *UserLogin) DefineRules() map[string]interface{} {
	return map[string]interface{}{
		"username": "required",
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

func (u *UserLogin) DefineScenes() map[string][]string {
	return map[string][]string{
		"login":    {"username", "password", "captcha"},
		"register": {"username", "password"},
	}
}

func main() {
	userLogin := &UserLogin{Username: "admin", Password: "123456"}
	userLogin.InitInstance(userLogin)
	err := userLogin.Check()
	if err != nil {
		fmt.Printf("登录验证失败！%v\r\n", err)
	} else {
		fmt.Println("登录验证通过")
	}

	err = userLogin.CheckScene("register")
	if err != nil {
		fmt.Printf("注册验证失败！%v\r\n", err)
	} else {
		fmt.Println("注册验证通过")
	}

}
