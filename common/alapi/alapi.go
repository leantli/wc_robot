// Package alapi 定义该三方接口供应商相关的请求、解析、模型等
package alapi

// ALAPI 的响应结构体
type AlapiResp struct {
	Code int                              `json:"code"`
	Msg  string                           `json:"msg"`
	Data struct{ Content, Author string } `json:"data"`
}
