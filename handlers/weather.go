package handlers

import (
	"fmt"
	"regexp"
	"strings"

	"wc_robot/common/weather"
	"wc_robot/robot"
)

var (
	locateWeatherRE = regexp.MustCompile("([\u4e00-\u9fa5]{1,6})天气")   // {城市}天气正则, 匹配位置
	locateAQIRE     = regexp.MustCompile("([\u4e00-\u9fa5]{1,6})空气质量") // {城市}空气质量正则, 匹配位置
)

func onWeatherChecker(msg *robot.Message) bool {
	return checkMatch(msg, []string{"天气"})
}

func onAQIChecker(msg *robot.Message) bool {
	return checkMatch(msg, []string{"空气质量"})
}

// 监听天气相关的文字进行回复
func onWeather(msg *robot.Message) error {
	hits := locateWeatherRE.FindStringSubmatch(msg.Content)
	if len(hits) != 2 {
		return nil
	}
	city := hits[1]
	runeCity := []rune(city)
	if len(runeCity) < 2 {
		_, err := msg.ReplyText("地区匹配过于宽泛，请规范输入，如\"深圳南山天气\"")
		return err
	}
	// 只取匹配到的城市的最后两个字作模糊查询
	wr, err := weather.GetCityLike(string(runeCity[len(runeCity)-2:]))
	if err != nil {
		return err
	}
	if len(wr.Data) == 1 {
		for k, v := range wr.Data {
			w, err := weather.GetWeather(k)
			if err != nil {
				return err
			}
			v = strings.Join(strings.Split(v, ", "), "-")
			_, err = msg.ReplyText(fmt.Sprintf("%s天气情况\n%s", v, w.GetCurrentWeatherInfo()))
			return err
		}
	}
	citys := wr.GetCityLike()
	var w *weather.WeatherResp
	for c, id := range citys {
		if strings.Contains(c, city) {
			if w != nil {
				_, err := msg.ReplyText("地区匹配过于宽泛，请规范输入，如\"深圳南山天气\"")
				return err
			}
			if w, err = weather.GetWeather(id); err != nil {
				return err
			}
		}
	}
	if w != nil {
		_, err = msg.ReplyText(fmt.Sprintf("%s天气情况\n%s", city, w.GetCurrentWeatherInfo()))
		return err
	}
	_, err = msg.ReplyText("很抱歉，无法获取到该地区的天气")
	return err
}

// 监听空气质量(指标含义) 的文字进行回复
func onAQI(msg *robot.Message) error {
	if strings.Contains(msg.Content, "指标含义") {
		msg.ReplyText(weather.AQIIndicesDesc())
		return nil
	}
	hits := locateAQIRE.FindStringSubmatch(msg.Content)
	if len(hits) != 2 {
		return nil
	}
	city := hits[1]
	runeCity := []rune(city)
	if len(runeCity) < 2 {
		_, err := msg.ReplyText("地区匹配过于宽泛，请规范输入，如\"深圳南山空气质量\"")
		return err
	}
	// 只取匹配到的城市的最后两个字作模糊查询
	wr, err := weather.GetCityLike(string(runeCity[len(runeCity)-2:]))
	if err != nil {
		return err
	}
	if len(wr.Data) == 1 {
		for k, v := range wr.Data {
			w, err := weather.GetWeather(k)
			if err != nil {
				return err
			}
			v = strings.Join(strings.Split(v, ", "), "-")
			_, err = msg.ReplyText(fmt.Sprintf("%s空气质量情况\n%s", v, w.GetAQIInfo()))
			return err
		}
	}
	citys := wr.GetCityLike()
	var w *weather.WeatherResp
	for c, id := range citys {
		if strings.Contains(c, city) {
			if w != nil {
				_, err := msg.ReplyText("地区匹配过于宽泛，请规范输入，如\"深圳南山空气质量\"")
				return err
			}
			if w, err = weather.GetWeather(id); err != nil {
				return err
			}
		}
	}
	if w != nil {
		_, err = msg.ReplyText(fmt.Sprintf("%s空气质量情况\n%s", city, w.GetAQIInfo()))
		return err
	}
	_, err = msg.ReplyText("很抱歉，无法获取到该地区的空气质量")
	return err
}
