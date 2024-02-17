package utils // 导出为 utils 包

import (
	"encoding/json"
)


func IsValidJSON(data []byte) bool { // 首字母要大写, 才能被其他包调用
	var js json.RawMessage
	return json.Unmarshal(data, &js) == nil
}