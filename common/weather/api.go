package weather

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"wc_robot/common/utils"
)

const (
	weather_url   = "https://weatherapi.market.xiaomi.com/wtr-v3/weather/all" // 小米天气接口(小米部分数据来源于彩云天气)
	city_like_url = "https://wis.qq.com/city/like"                            // 城市列表接口，返回 map 结构 -> "城市id" : "对应的城市"
)

// http 请求获取天气数据
func GetWeather(cityCode string) (*WeatherResp, error) {
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
	var r WeatherResp
	err = utils.ScanJson(resp, &r)
	return &r, err
}

// http 请求模糊查询城市，获取 "城市 id":"城市" map
func GetCityLike(city string) (*CityLikeResp, error) {
	uri, err := url.Parse(city_like_url)
	if err != nil {
		return nil, err
	}
	params := url.Values{}
	params.Add("source", "pc")
	params.Add("city", city)
	uri.RawQuery = params.Encode()
	resp, err := http.Get(uri.String())
	if err != nil {
		return nil, err
	}
	var r CityLikeResp
	if err := utils.ScanJson(resp, &r); err != nil {
		return nil, err
	}
	if r.Status != http.StatusOK {
		return nil, fmt.Errorf("request GetCityLike failed, status: %v, message: %s", r.Status, r.Message)
	}
	return &r, nil
}
