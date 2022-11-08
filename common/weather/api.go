package weather

import (
	"net/http"
	"net/url"
	"strings"

	"wc_robot/common/utils"
)

// 小米天气接口(小米部分数据来源于彩云天气)
const weather_url = "https://weatherapi.market.xiaomi.com/wtr-v3/weather/all"

// http请求获取天气数据
func GetWeather(cityCode string) (*Weather, error) {
	uri, err := url.Parse(weather_url)
	if err != nil {
		return nil, err
	}
	params := url.Values{}
	params.Add("locationKey", strings.Join([]string{"weathercn", cityCode}, ":"))
	// 以下皆为固定值
	params.Add("latitude", "0")
	params.Add("longitude", "0")
	params.Add("sign", "zUFJoAR2ZVrDy1vF3D07")
	params.Add("isGlobal", "false")
	params.Add("locale", "zh_cn")
	params.Add("appKey", "weather20151024")
	uri.RawQuery = params.Encode()
	resp, err := http.Get(uri.String())
	if err != nil {
		return nil, err
	}

	var w Weather
	err = utils.ScanJson(resp, &w)
	return &w, err
}
