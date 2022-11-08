package alapi

import "fmt"

const (
	host = "https://v2.alapi.cn"

	soulPath    = "/api/soul"
	qinghuaPath = "/api/qinghua"
	mingyanPath = "/api/mingyan"
	jokePath    = "/api/joke/random"
)

// 响应码
const (
	// 请求成功
	Success = 200
	// 请求超过限制次数
	OverLimit = 102
	// 请求 QPS 超过限制
	OverQPS = 429
)

var mapCodeDesc = map[int]string{
	Success:   "请求成功",
	OverLimit: "请求超过限制次数",
	OverQPS:   "请求 QPS 超过限制",
	404:       "接口地址不存在",
	422:       "接口请求失败",
	400:       "接口请求失败",
	405:       "请求方法不被允许",
	100:       "token错误",
	101:       "账号过期",
	104:       "来源或者ip不在白名单",
	406:       "没有更多数据了",
}

// GetCodeDesc 根据响应码返回对应的信息
func GetCodeDesc(code int) string {
	if desc, exist := mapCodeDesc[code]; exist {
		return desc
	}
	return fmt.Sprintf("未知的响应码: %d", code)
}
