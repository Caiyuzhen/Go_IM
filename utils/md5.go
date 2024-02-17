// 🔐 md5 来进行注册时密码的加密
package utils // 导出为 utils 包

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
	_"fmt"
)


// 小写 - 生成 md5 字符串  -> 🌟注意大小写!!
func Md5Encode(data string) string {
	h := md5.New() // 创建一个 md5 对象
	h.Write([]byte(data)) // 将 data 写入到 h 中
	tempStr := h.Sum(nil)  // 计算 MD5 值, 并存为字符串
	return hex.EncodeToString(tempStr) // 返回加密后的字符串
}

// 转为大写 - 生成 md5 字符串 -> 🌟注意大小写!!
func MD5Encode(data string) string {
	return strings.ToUpper(Md5Encode(data))
}


// 【🔐 加密】->  在 userServer 内进行调用
func MakePassword(plainWd, salt string) string { // 增加多一个随机数(盐值), 增加破解难度, 盐值存在用户表内
	return MD5Encode(plainWd + salt)  // 将【密码【和【盐值】拼接后再进行加密 => 增加破解难度
}

// 【🔐 解密】=> 解密后看看是否 == 密码  ->  在 userServer 内进行调用
func ValidPassword(userInputPwd, salt string, dataBaseUserPassword string) bool { // 传入【加密后的密码】| 盐值 ｜【用户账号的密码】
	// fmt.Println("✅ 拿到了用户输入的密码: ", userInputPwd)
	// fmt.Println("✅ 拿到了用户的盐值: ", salt)
	// fmt.Println("✅ 拿到了用户的加密密码: ", dataBaseUserPassword)
	// fmt.Println("✅ 拿到了用户的解密密码: ", Md5Encode(userInputPwd + salt))
	return MD5Encode(userInputPwd + salt) == dataBaseUserPassword// 返回一个布尔值, 表示是否是正确的密码
}