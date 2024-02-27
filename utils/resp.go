package utils

import (
	"encoding/json"
	"net/http"
	"fmt"
)

// 🔥 返回响应的工具类
type H struct {
	Code int
	Msg string
	Data interface {} // 什么数据都可以接收
	Rows interface {} // 什么数据都可以接收
	Total interface {} // 什么数据都可以接收
}



// Resp 通用响应函数 (返回 json 数据)
func Resp(w http.ResponseWriter, code int, data interface{}, msg string) {
	fmt.Println("🚀🚀🚀 响应的数据", data)
	// 设置响应的Content-Type为application/json
	w.Header().Set("Content-Type", "application/json")

	// 返回响应头
	w.WriteHeader(http.StatusOK)

	// 实例化结构体
	h := H {
		Code: code,
		Data: data, // 🔥🔥 返回前段传来的数据（比如上传图片！）
		Msg: msg, // 消息
	}

	// 把结构体实例转为 json  二进制 （序列化）
	res, err := json.Marshal(h)
	if err != nil {
		RespFail(w, "❌ json 转化出错")
		return
	}

	fmt.Println("🚀🚀🚀 序列化为 json 后: ",  string(res))
	fmt.Println("____________________________")

	// 返回转化后的 json
	w.Write(res)
}



// RespList 表示请求处理成功 (返回数据)
func RespList(w http.ResponseWriter, code int, data interface {}, total interface {}) {
	// 设置响应的Content-Type为application/json
	w.Header().Set("Content-Type", "application/json")

	// 返回响应头
	w.WriteHeader(http.StatusOK)

	// 实例化结构体
	h := H {
		Code: code,
		Rows: data, // 行数
		Total: total, // 总数
	}

	// 把结构体实例转为 json
	res, err := json.Marshal(h)
	if err != nil {
		RespFail(w, "❌ json 转化出错")
		return
	}

	// 返回转化后的 json
	w.Write(res)
}




// RespFail 表示请求处理失败的响应 
func RespFail(w http.ResponseWriter, msg string) {
	Resp(w, -1, nil, msg) // -1 表示失败
}




// RespOK 表示请求处理成功的响应
func RespOK(w http.ResponseWriter, data interface {}, msg string) {
	Resp(w, 0, data, msg) // -1 表示失败
}





// RespOkList 表示请求处理成功 (需要返回数据)
func RespOkList(w http.ResponseWriter, data interface {}, total interface {}) {
	RespList(w, 0, data, total) // 传入 0 , 表示成功
}






