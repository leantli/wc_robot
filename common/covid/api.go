package covid

import (
	"fmt"
	"net/http"
	"wc_robot/common/utils"
)

const uri string = "https://opendata.baidu.com/data/inner?resource_id=5653&query=%s新型肺炎最新动态&dsp=iphone&tn=wisexmlnew&alr=1&is_opendata=1"

// 调百度实时疫情 API 获取疫情数据响应
func GetCovidResponse(location string) (*CovidResponse, error) {
	rsp, err := http.Get(fmt.Sprintf(uri, location))
	if err != nil {
		return nil, err
	}
	var cs CovidResponse
	if err := utils.ScanJson(rsp, &cs); err != nil {
		return nil, err
	}
	return &cs, nil
}
