package main

import (
	"fmt"

	"github.com/worklz/go-validate"
)

func main() {
	// 定义变量标题
	var1 := "qwe"
	err := validate.CheckVar(var1, "number", "参数1", nil)
	if err != nil {
		fmt.Printf("var1验证失败！%v\r\n", err)
	} else {
		fmt.Println("var1验证通过")
	}

	// 自定义错误信息
	var2 := "qwe"
	err = validate.CheckVar(var2, "number", "", map[string]string{"number": "参数2错误"})
	if err != nil {
		fmt.Printf("var2验证失败！%v\r\n", err)
	} else {
		fmt.Println("var2验证通过")
	}
}
