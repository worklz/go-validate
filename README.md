# 验证器

[![](https://img.shields.io/badge/Author-worklz-orange.svg)](https://github.com/worklz/go-validate)
[![](https://img.shields.io/badge/version-v1.0.0-brightgreen.svg)](https://github.com/worklz/go-validate)
![GitHub stars](https://img.shields.io/github/stars/worklz/go-validate?style=flat-square)

## 安装/引入

```bash
# 项目根目录下执行
go get github.com/worklz/go-validate
```

```golang
// 项目代码中引入包
import "github.com/worklz/go-validate"
```


## 使用示例

- [简单使用](https://github.com/worklz/go-validate/example/simple/main.go)
- [验证场景](https://github.com/worklz/go-validate/example/scene/main.go)
- [Map数据验证](https://github.com/worklz/go-validate/example/map/main.go)
- [注册验证规则](https://github.com/worklz/go-validate/example/register_rule/main.go)
- [验证器内方法定义为验证规则](https://github.com/worklz/go-validate/example/validator_method_rule/main.go)
- [定义参数验证规则为闭包](https://github.com/worklz/go-validate/example/func_rule/main.go)
- [验证后处理数据](https://github.com/worklz/go-validate/example/handle_datas/main.go)
- [验证单个数据](https://github.com/worklz/go-validate/example/check_var/main.go)

## 验证规则

以下为内置验证规则，可直接使用，更多规则请自行定义（参考上述：注册验证规则示例）。

| 规则 | 描述 | 使用示例 | 解释 | 注意 |
| :-- | :---------------- | :--- | :---------------- | :---------------- |
| required | 验证字段必须 | "name":"required" | 名称必填 | 如果验证规则没有添加required就表示空值，则不会进行后续规则验证 |
| number | 验证字段必须为数字 | "age":"number" | 年龄必须为数字 | 值为数字类型或者数字字符串即可 |
| integer | 验证字段必须为整数 | "count":"integer" | 数量必须为整数 | 值为整数类型或者能转换为整数的字符串 |
| positiveInt | 验证字段必须为正整数 | "score":"positiveInt" | 分数必须为正整数 | 大于0的整数 |
| nonnegativeInt | 验证字段必须为非负整数 | "quantity":"nonnegativeInt" | 数量必须为非负整数 | 大于等于0的整数 |
| float | 验证字段必须为浮点数 | "price":"float" | 价格必须为浮点数 | 值为浮点数类型或者能转换为浮点数的字符串 |
| boolean | 验证字段必须为布尔值 | "isValid":"boolean" | 是否有效必须为布尔值 | 值为布尔类型 |
| length | 验证字段的长度 | "password":"length:6,12" | 密码长度限制在6 - 12位 | 参数为单个正整数时表示固定长度，用逗号分隔的两个正整数表示长度区间 |
| min | 验证字段的最小长度 | "username":"min:3" | 用户名最小长度为3 | 参数必须为正整数 |
| max | 验证字段的最大长度 | "description":"max:200" | 描述最大长度为200 | 参数必须为正整数 |
| in | 验证字段的值必须在指定范围内 | "gender":"in:male,female" | 性别必须为男或女 | 参数用逗号分隔 |
| notIn | 验证字段的值必须不在指定范围内 | "status":"notIn:disabled" | 状态不能为禁用 | 参数用逗号分隔 |
| between | 验证字段的值必须在指定区间内 | "age":"between:18,60" | 年龄必须在18 - 60岁之间 | 参数用逗号分隔，且必须为整数 |
| notBetween | 验证字段的值必须不在指定区间内 | "score":"notBetween:0,50" | 分数不能在0 - 50分之间 | 参数用逗号分隔，且必须为整数 |
| eq | 验证字段的值必须等于指定值 | "code":"eq:123" | 代码必须为123 | 参数必须为整数 |
| egt | 验证字段的值必须大于等于指定值 | "age":"egt:18" | 年龄必须大于等于18岁 | 参数必须为整数 |
| gt | 验证字段的值必须大于指定值 | "price":"gt:100" | 价格必须大于100 | 参数必须为整数 |
| elt | 验证字段的值必须小于等于指定值 | "quantity":"elt:100" | 数量必须小于等于100 | 参数必须为整数 |
| lt | 验证字段的值必须小于指定值 | "score":"lt:60" | 分数必须小于60 | 参数必须为整数 |
| array | 验证字段必须为数组 | "hobbies":"array" | 爱好必须为数组 | 值为切片类型 |
| arrayIn | 验证字段数组中的元素必须在指定范围内 | "roles":"arrayIn:admin,user" | 角色数组中的元素必须为管理员或用户 | 参数用逗号分隔，数组不能为空 |
| arrayEmptyOrIn | 验证字段数组为空或数组中的元素必须在指定范围内 | "tags":"arrayEmptyOrIn:tag1,tag2" | 标签数组为空或标签必须为tag1或tag2 | 参数用逗号分隔 |
| arrayPositiveInt | 验证字段必须为正整数数组 | "ids":"arrayPositiveInt" | ids 字段需为正整数数组 | 字段必须是数组且数组元素都为正整数，数组不能为空 |
| arrayEmptyOrPositiveInt | 验证字段可以为空数组或正整数数组 | "ids":"arrayEmptyOrPositiveInt" | ids 字段可以为空数组或正整数数组 | 若字段为数组且不为空，则数组元素都需为正整数 |
| mapHas | 验证字段必须为包含特定键的非空 map | "info":"mapHas:key1,key2" | info 字段必须是包含 key1 和 key2 的非空 map | 验证规则参数需用逗号分隔，字段必须是 map 且包含规则指定的所有键 |
| mapEmptyOrHas | 验证字段可以为空 map 或包含特定键的 map | "info":"mapEmptyOrHas:key1,key2" | info 字段可以为空 map 或包含 key1 和 key2 的 map | 若字段为非空 map，则必须包含规则指定的所有键，验证规则参数需用逗号分隔 |
| arrayItemHas | 验证字段必须为非空数组，且数组每个元素都是包含特定键的非空 map | "list":"arrayItemHas:key1,key2" | list 字段必须是包含多个 map 的非空数组，每个 map 都要包含 key1 和 key2 | 验证规则参数需用逗号分隔，数组元素必须是 map 且包含规则指定的所有键，数组和 map 都不能为空 |
| arrayEmptyOrItemHas | 验证字段可以为空数组或包含特定键的 map 数组 | "list":"arrayEmptyOrItemHas:key1,key2" | list 字段可以为空数组或包含多个 map 的数组，每个 map 都要包含 key1 和 key2 | 若字段为非空数组，则数组元素必须是 map 且包含规则指定的所有键，验证规则参数需用逗号分隔 |
| mobile | 验证字段必须为 11 位有效手机格式 | "phone":"mobile" | phone 字段需为 11 位有效手机格式 | 字段值需为符合手机格式的字符串 |
| email | 验证字段必须为有效邮箱格式 | "email":"email" | email 字段需为有效邮箱格式 | 字段值需为符合邮箱格式的字符串 |
| chs | 验证字段只能是汉字 | "name":"chs" | name 字段只能是汉字 | 字段值需为纯汉字字符串 |
| chsAlphaNum | 验证字段只能是汉字、字母、数字 | "username":"chsAlphaNum" | username 字段只能是汉字、字母、数字 | 字段值需由汉字、字母、数字组成 |
| chsDash | 验证字段只能是汉字、字母、数字、下划线_、破折号 - | "code":"chsDash" | code 字段只能是汉字、字母、数字、下划线_、破折号 - | 字段值需由指定字符组成 |
| chsDashSpace | 验证字段只能是汉字、字母、数字、下划线_、短横线 - 及空格组合 | "address":"chsDashSpace" | address 字段只能是汉字、字母、数字、下划线_、短横线 - 及空格组合 | 字段值需由指定字符和空格组成 |
| alphaNum | 验证字段只能是字母、数字 | "password":"alphaNum" | password 字段只能是字母、数字 | 字段值需由字母、数字组成 |
| alphaDash | 验证字段只能是字母、数字、下划线_、短横线 - | "slug":"alphaDash" | slug 字段只能是字母、数字、下划线_、短横线 - | 字段值需由指定字符组成 |
| hexColor | 验证字段必须为十六进制颜色格式 | "color":"hexColor" | color 字段必须为十六进制颜色格式 | 字段值需为符合十六进制颜色格式的字符串 |
| date | 验证字段必须为日期格式 | "birthdate":"date" | birthdate 字段必须为日期格式 | 字段值需为符合日期格式的字符串 |
| datetime | 验证字段必须为日期时间格式 | "create_time":"datetime" | create_time 字段必须为日期时间格式 | 字段值需为符合日期时间格式的字符串 |
| year | 验证字段必须为年份格式 | "start_year":"year" | start_year 字段必须为年份格式 | 字段值需为符合年份格式的字符串 |
| yearMonth | 验证字段必须为年月格式 | "period":"yearMonth" | period 字段必须为年月格式 | 字段值需为符合年月格式的字符串 |
| month | 验证字段必须为月份格式 | "due_month":"month" | due_month 字段必须为月份格式 | 字段值需为符合月份格式的字符串 |
| time | 验证字段必须为时间格式 | "start_time":"time" | start_time 字段必须为时间格式 | 字段值需为符合时间格式的字符串 |
| timeRange | 验证字段必须为时间范围格式 | "work_time":"timeRange" | work_time 字段必须为时间范围格式 | 具体格式验证由 isTimeRange 函数决定 |
| commaIntervalChsAlphaNum | 验证字段必须为逗号分隔的汉字、字母、数字组合 | "tags":"commaIntervalChsAlphaNum" | tags 字段必须为逗号分隔的汉字、字母、数字组合 | 字段值需为符合格式的字符串 |
| commaIntervalPositiveInt | 验证字段必须为逗号分隔的正整数组合 | "scores":"commaIntervalPositiveInt" | scores 字段必须为逗号分隔的正整数组合 | 字段值需为符合格式的字符串 |
| url | 验证字段必须为合法的URL地址 | "website":"url" | 网站地址必须是合法的URL格式 | 验证的值必须是字符串类型，否则判定为格式错误 |
| urls | 验证字段必须为包含合法URL地址的数组 | "imageUrls":"urls" | 图片链接列表中的每个链接都必须是合法的URL格式 | 验证的值必须是数组类型，数组中的每个元素必须是字符串类型，否则判定为格式错误 |
| ip | 验证字段必须为合法的IP地址 | "serverIp":"ip" | 服务器IP地址必须是合法的IP格式 | 验证的值必须是字符串类型，否则判定为格式错误 |
| uri | 验证字段必须为合法的URI地址 | "resourceUri":"uri" | 资源的URI地址必须是合法的URI格式 | 验证的值必须是字符串类型，否则判定为格式错误 |