package main

import (
	"fmt"

	"github.com/worklz/go-validate"
)

type UserParams struct {
	validate.Validator
	Username string `json:"username"`
	Password string `json:"password"`
	Captcha  string `json:"captcha"`
	UserId   uint   `json:"user_id"`
}

func (u *UserParams) DefineRules() map[string]interface{} {
	return map[string]interface{}{
		"username": "required",
		"password": "required",
		"captcha":  "required",
		"user_id":  "required",
	}
}

func (u *UserParams) DefineTitles() map[string]string {
	return map[string]string{
		"username": "用户名",
		"password": "密码",
		"captcha":  "验证码",
	}
}

func (u *UserParams) DefineScenes() map[string][]string {
	return map[string][]string{
		"login":    {"username", "password", "captcha"},
		"register": {"username", "password"},
		"info":     {"user_id"},
	}
}

// 验证后处理数据
func (u *UserParams) HandleDatas(datas map[string]interface{}, scene string) (err error) {
	switch scene {
	case "login":
		err = u.SetData("user_id", uint(1))
		if err != nil {
			return
		}
		fmt.Println("登录验证后处理数据")
	case "register":
		err = u.SetData("user_id", uint(2))
		if err != nil {
			return
		}
		fmt.Println("注册验证后处理数据")
	case "info":
		err = u.SetError("用户不存在")
		return
	}
	return
}

func main() {
	userLogin := &UserParams{Username: "admin", Password: "123456", Captcha: "1234"}
	validate.Create(userLogin)
	err := userLogin.CheckScene("login")
	if err != nil {
		fmt.Printf("登录验证失败！%v\r\n", err)
	} else {
		fmt.Println("登录验证通过")
	}
	fmt.Printf("登录后的数据datas：%v 结构体UserId：%d\r\n", userLogin.Datas, userLogin.UserId)

	userRegister := &UserParams{Username: "admin", Password: "123456"}
	validate.Create(userRegister)
	err = userRegister.CheckScene("register")
	if err != nil {
		fmt.Printf("注册验证失败！%v\r\n", err)
	} else {
		fmt.Println("注册验证通过")
	}
	fmt.Printf("注册后的数据datas：%v 结构体UserId：%d\r\n", userRegister.Datas, userRegister.UserId)

	userInfo := &UserParams{UserId: 3}
	validate.Create(userInfo)
	err = userInfo.CheckScene("info")
	if err != nil {
		fmt.Printf("获取信息验证失败！%v\r\n", err)
	} else {
		fmt.Println("获取信息验证通过")
	}
	fmt.Printf("获取信息的数据datas：%v 结构体UserId：%d\r\n", userInfo.Datas, userInfo.UserId)

}
