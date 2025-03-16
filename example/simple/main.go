package main

import (
	"fmt"

	"github.com/worklz/go-validate"
)

type UserLogin struct {
	validate.Validator
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u *UserLogin) DefineRules() map[string]interface{} {
	return map[string]interface{}{
		"username": "required",
		"password": "required",
	}
}

func (u *UserLogin) DefineMessages() map[string]string {
	return map[string]string{
		"username.required": "用户名不能为空",
		"password.required": "密码不能为空",
	}
}

func (u *UserLogin) DefineScenes() map[string][]string {
	return map[string][]string{
		"login": {"username", "password"},
	}
}

func main() {
	userLogin := &UserLogin{Username: "admin", Password: "123456"}
	validate.Create(userLogin)
	err := userLogin.Check()
	if err != nil {
		fmt.Printf("验证失败！%v", err)
	} else {
		fmt.Println("验证通过")
	}
}
