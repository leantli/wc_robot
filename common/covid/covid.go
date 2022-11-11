// Package covid 定义该三方接口供应商相关的请求、解析、模型等
package covid

import (
	"strings"
)

// 百度 opendata 实时疫情数据响应结构体
type CovidResponse struct {
	ResultNum string `json:"ResultNum"`
	Result    []struct {
		DisplayData struct {
			ResultData struct {
				TplData struct {
					Desc     string `json:"desc"`
					DataList []struct {
						TotalDesc string `json:"total_desc"`
						TotalNum  string `json:"total_num"`
					} `json:"data_list"`
					Location string `json:"location"`
				} `json:"tplData"`
			} `json:"resultData"`
		} `json:"DisplayData"`
	} `json:"Result"`
}

// PrintCovidSituation 根据 CovidResponse 中的数据描述打印疫情信息
// 第二个 Result 一般是疫情情况的事件脉络，但内容并不够及时，直接去掉不用
func PrintCovidSituation(cr *CovidResponse) string {
	sb := strings.Builder{}
	td := cr.Result[0].DisplayData.ResultData.TplData
	sb.WriteString(td.Location + "今日疫情:\n")
	for _, d := range td.DataList {
		sb.WriteString(d.TotalDesc)
		sb.WriteString(" ")
		sb.WriteString(d.TotalNum)
		sb.WriteString("\n")
	}
	sb.WriteString(td.Desc)
	return sb.String()
}
