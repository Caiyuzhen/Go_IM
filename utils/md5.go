// 🔐 md5 来进行注册时密码的加密
package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)


// 小写 - 生成 md5 字符串
func Md5Encode(data string) string {
	h := md5.New() // 创建一个 md5 对象
	h.Write([]byte(data)) // 将 data 写入到 h 中
	tempStr := h.Sum(nil)  // 计算 MD5 值, 并存为字符串
	return hex.EncodeToString(tempStr) // 返回加密后的字符串
}

// 转为大写 - 生成 md5 字符串
func MD5Encode(data string) string {
	return strings.ToUpper(Md5Encode(data))
}


// 【🔐 加密】
func MakePassword(plainWd, salt string) string { // 增加多一个随机数(盐值), 增加破解难度, 盐值存在用户表内
	return MD5Encode(plainWd + salt)  // 将【密码【和【盐值】拼接后再进行加密 => 增加破解难度
}

// 【🔐 解密】=> 解密后看看是否 == 密码
func ValidPassword(plainWd, salt string) bool { // 传入【加密后的密码】、【用户账号的密码】
	return Md5Encode(plainWd + salt) // 返回一个布尔值, 表示是否是正确的密码
}