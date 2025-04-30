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
	D1       []int  `json:"d1"`
}

func (u *UserLogin) DefineRules() map[string]interface{} {
	return map[string]interface{}{
		"username": "required|chsDashChar",
		"password": "required|",
		"captcha":  "required",
		"d1":       "required|arrayIn:12,23,34",
	}
}

func (u *UserLogin) DefineMessages() map[string]string {
	return map[string]string{
		"username.required": "请输入用户名",
		"password.required": "请输入密码",
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
	userLogin := &UserLogin{Username: "admin（超管）", Password: "123456", Captcha: "1234", D1: []int{12, 23, 26}}
	validate.Create(userLogin)
	err := userLogin.Check()
	if err != nil {
		fmt.Printf("验证失败！%v\r\n", err)
	} else {
		fmt.Println("验证通过")
	}
}
