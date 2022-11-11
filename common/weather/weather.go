// Package weather 定义了对小米天气的请求和解析接口
package weather

import (
	"fmt"
	"strconv"
	"strings"
)

// 天气
type Weather struct {
	Current        *Current        `json:"current"`        // 当前天气预报
	ForecastDaily  *ForecastDaily  `json:"forecastDaily"`  // 未来15日天气预报(含今日)
	ForecastHourly *ForecastHourly `json:"forecastHourly"` // 未来24h天气预报(不含当前小时)
	AQI            *AQI            `json:"aqi"`            // 当前空气质量

	errCode string
	errDesc string
}

// 当前天气
type Current struct {
	Weather     string         `json:"weather"`     // 天气码，对应信息见constant.go
	Temperature *ModuleCurrent `json:"temperature"` // 温度
	Humidity    *ModuleCurrent `json:"humidity"`    // 湿度
	Pressure    *ModuleCurrent `json:"pressure"`    // 气压
	PubTime     string         `json:"pubTime"`     // 发布时间
}

func (c *Current) String() string {
	return fmt.Sprintf("当前天气: %s\n当前温度: %s%s\n当前湿度: %s%s\n当前气压: %s%s\n本次数据更新时间: %s",
		GetWeatherCodeDesc(c.Weather),
		c.Temperature.Value, c.Temperature.Unit,
		c.Humidity.Value, c.Humidity.Unit,
		c.Pressure.Value, c.Pressure.Unit,
		c.PubTime,
	)
}

// 当前天气的模块，包含体感、湿度、气压、温度
type ModuleCurrent struct {
	Unit  string // 单位
	Value string // 值
}

// 15天天气预报(含今日)
type ForecastDaily struct {
	PubTime     string       `json:"pubTime"`     // 发布时间
	Temperature *ModuleDaily `json:"temperature"` // 15天天气预报的温度模块
	Weather     *ModuleDaily `json:"weather"`     // 15天天气预报的天气模块
}

// 15天天气预报的模块，包含温度和天气
type ModuleDaily struct {
	Unit  string                      `json:"unit"`  // 单位
	Value []struct{ From, To string } `json:"value"` // 值，一般有15个
}

// 未来24小时天气预报(不含当前小时)
type ForecastHourly struct {
	Temperature *ModuleHourly `json:"temperature"` // 未来24小时天气预报的温度模块
	Weather     *ModuleHourly `json:"weather"`     // 未来24小时天气预报的天气模块
	AQI         *ModuleHourly `json:"aqi"`         // 未来24小时天气预报的AQI模块
}

// 未来24小时天气预报的模块，包含温度、天气
type ModuleHourly struct {
	PubTime string `json:"pubTime"` // 发布时间
	Value   []int  `json:"value"`   // 23个值
}

// 当前的空气质量
type AQI struct {
	Aqi     string `json:"aqi"`     // 空气质量
	CO      string `json:"co"`      // 一氧化碳
	NO2     string `json:"no2"`     // 二氧化氮
	O3      string `json:"o3"`      // 臭氧
	PM10    string `json:"pm10"`    // 10微米以下可吸入颗粒
	PM25    string `json:"pm25"`    // 25微米以下可吸入颗粒
	SO2     string `json:"so2"`     // 二氧化硫
	PubTime string `json:"pubTime"` // 发布时间
	Suggest string `json:"suggest"` // 发布时间
}

// 获取天气信息中的当前天气信息
func (w *Weather) GetCurrentWeatherInfo() string {
	next1 := strconv.Itoa(w.ForecastHourly.Weather.Value[0])
	next2 := strconv.Itoa(w.ForecastHourly.Weather.Value[1])
	next3 := strconv.Itoa(w.ForecastHourly.Weather.Value[2])
	return fmt.Sprintf("你好呀👋\n当前天气: %s\n"+
		"当前温度: %s%s\n"+
		"当前湿度: %s%s\n"+
		"当前空气质量: %s\n"+
		"预期未来三小时天气: %s, %s, %s\n"+
		"今日最低/高温度: %s/%s\n"+
		"今日天气预期: %s->%s\n"+
		"本次数据更新时间: %s",
		GetWeatherCodeDesc(w.Current.Weather),
		w.Current.Temperature.Value, w.Current.Temperature.Unit,
		w.Current.Humidity.Value, w.Current.Humidity.Unit,
		w.AQI.Aqi,
		GetWeatherCodeDesc(next1), GetWeatherCodeDesc(next2), GetWeatherCodeDesc(next3),
		w.ForecastDaily.Temperature.Value[0].To, w.ForecastDaily.Temperature.Value[0].From,
		GetWeatherCodeDesc(w.ForecastDaily.Weather.Value[0].From), GetWeatherCodeDesc(w.ForecastDaily.Weather.Value[0].To),
		strings.TrimSuffix(w.Current.PubTime, "+08:00"), //去除pubTime后面的时区显示+08:00
	)
}

// 获取天气信息中的AQI空气质量信息
func (w *Weather) GetAQIInfo() string {
	return fmt.Sprintf("你好呀👋\n当前空气质量: %s %s\n"+
		"PM2.5细颗粒物: %sμg/m³\n"+
		"PM10可吸入颗粒物: %sμg/m³\n"+
		"SO2二氧化硫: %sμg/m³\n"+
		"NO2二氧化氮: %sμg/m³\n"+
		"O3臭氧: %sμg/m³\n"+
		"CO一氧化碳: %smg/m³\n"+
		"本次数据更新时间: %s",
		w.AQI.Aqi, GetAQIQuality(w.AQI.Aqi),
		w.AQI.PM25,
		w.AQI.PM10,
		w.AQI.SO2,
		w.AQI.NO2,
		w.AQI.O3,
		w.AQI.CO,
		strings.TrimSuffix(w.Current.PubTime, "+08:00"), //去除pubTime后面的时区显示+08:00
	)
}

// 获取AQI空气质量指标的描述
func AQIIndicesDesc() string {
	return strings.Join([]string{pm25Desc, pm10Desc, so2Desc, no2Desc, o3Desc, coDesc}, "\n")
}
