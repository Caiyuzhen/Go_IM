package utils

import (
	"encoding/json"
	"net/http"
)

// 🔥 分页的工具类（好友列表分页）
type H struct {
	Code int
	Msg string
	Data interface {} // 什么数据都可以接收
	Rows interface {} // 什么数据都可以接收
	Total interface {} // 什么数据都可以接收
}

// Resp 通用响应函数
func Resp(w http.ResponseWriter, code int, data interface{}, msg string) {

}


// RespFail 表示请求处理失败的响应
func RespFail(w http.ResponseWriter, msg string) {
	Resp(w, -1, nil, msg) // -1 表示失败
}



// RespOk 表示请求处理成功的响应，不带数据返回
func RespOk(w http.ResponseWriter, data interface {}, msg string) {
	Resp(w, 0, nil, msg) // -1 表示失败
}





// RespOkList 表示请求处理成功，且需要返回列表数据
func RespOkList(w http.ResponseWriter, data interface {}, total interface {}) {
	RespList(w, 0, data, total) // 传入 0 表示成功
}




// RespList 表示请求处理成功, 并返回列表数据
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

