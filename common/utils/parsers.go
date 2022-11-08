package utils

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"net/http"
)

// parsers.go 主要包含了一些通用的格式转换、解析操作

// ScanXml 按xml格式解析resp的payload（Body）
func ScanXml(resp *http.Response, v interface{}) error {
	return xml.NewDecoder(resp.Body).Decode(v)
}

// ScanJson 按json格式解析resp的payload（Body）
func ScanJson(resp *http.Response, v interface{}) error {
	return json.NewDecoder(resp.Body).Decode(v)
}

// ToJsonBuff 读取
func ToJsonBuff(v any) (*bytes.Buffer, error) {
	var b bytes.Buffer
	e := json.NewEncoder(&b)
	// 设置禁止html转义
	e.SetEscapeHTML(false)
	err := e.Encode(v)
	return &b, err
}
