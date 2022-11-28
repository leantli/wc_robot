package handlers

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"wc_robot/common"
	"wc_robot/common/alapi"
	"wc_robot/common/covid"
	"wc_robot/common/weather"
	"wc_robot/robot"
)

var (
	begin           = time.Now()                                       // 存活时间计时
	locateWeatherRE = regexp.MustCompile("([\u4e00-\u9fa5]{1,6})天气")   // {城市}天气正则, 匹配位置
	locateAQIRE     = regexp.MustCompile("([\u4e00-\u9fa5]{1,6})空气质量") // {城市}空气质量正则, 匹配位置
	locateCovidRE   = regexp.MustCompile("([\u4e00-\u9fa5]{1,6})疫情")   // {城市}疫情正则，匹配位置
)

func InitHandlers(r *robot.Robot) {
	config := common.GetConfig()
	r.Chain.RegisterGlobalCheck(checkOnContact)
	r.Chain.RegisterHandler("菜单|功能|会什么回复", onMenu)
	r.Chain.RegisterHandler("存活时间回复", onSurvivalTime)
	if config.WeatherMsgHandle.SwitchOn {
		r.Chain.RegisterHandler("天气回复", onWeather)
		r.Chain.RegisterHandler("空气质量回复", onAQI)
	}
	if config.ALAPI.SwitchOn {
		r.Chain.RegisterHandler("名言回复", onMingYan)
		r.Chain.RegisterHandler("情话回复", onQingHua)
		r.Chain.RegisterHandler("鸡汤回复", onSoul)
	}
	if config.CovidMsgHandle.SwitchOn {
		r.Chain.RegisterHandler("疫情回复", onCovid)
	}
}

// 基础校验，机器人只回复文字、监听的nickname、非自己，其余都不回复，返回 false
func checkOnContact(msg *robot.Message) bool {
	if !msg.IsText() {
		return false
	}
	if !msg.IsSentByNickName(common.GetConfig().OnContactNickNames) {
		return false
	}
	if msg.IsFromSelf() {
		return false
	}
	return true
}

// 下面一些匹配：就strings.Contains()和正则匹配二者的性能来说，前者较优

// 判断是否匹配，匹配返回 true, 不匹配返回 false
func checkMatch(msg *robot.Message, keyword string) bool {
	config := common.GetConfig()
	if msg.IsFromGroup() {
		if !(strings.Contains(msg.Content, "@"+config.RobotName) && strings.Contains(msg.Content, keyword)) {
			return false
		}
	}
	if msg.IsFromMember() {
		if !strings.Contains(msg.Content, keyword) {
			return false
		}
	}
	return true
}

// 监听菜单｜功能｜会什么相关的文字进行回复
func onMenu(msg *robot.Message) error {
	config := common.GetConfig()
	if msg.IsFromGroup() {
		if !(strings.Contains(msg.Content, "@"+config.RobotName) &&
			(strings.Contains(msg.Content, "菜单") || strings.Contains(msg.Content, "功能") || strings.Contains(msg.Content, "会什么"))) {
			return nil
		}
	}
	if msg.IsFromMember() {
		if !(strings.Contains(msg.Content, "菜单") || strings.Contains(msg.Content, "功能") || strings.Contains(msg.Content, "会什么")) {
			return nil
		}
	}
	_, err := msg.ReplyText("你好呀👋\n" + `目前只支持"天气"、"空气质量(指标含义)"、"XX(城市、省份、国家)疫情"、"情话"、"鸡汤"、"名言"相关的问题哦`)
	return err
}

// 监听天气相关的文字进行回复
func onWeather(msg *robot.Message) error {
	if !checkMatch(msg, "天气") {
		return nil
	}
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
	if !checkMatch(msg, "空气质量") {
		return nil
	}
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

// 监听心灵鸡汤相关的文字进行回复
func onSoul(msg *robot.Message) error {
	if !checkMatch(msg, "鸡汤") {
		return nil
	}
	s, err := alapi.GetSoul()
	if err != nil {
		return err
	}
	_, err = msg.ReplyText(s)
	return err
}

// 监听情话相关的文字进行回复
func onQingHua(msg *robot.Message) error {
	if !checkMatch(msg, "情话") {
		return nil
	}
	content, err := alapi.GetQinghua()
	if err != nil {
		return err
	}
	_, err = msg.ReplyText(content)
	return err
}

// 监听名言相关的文字进行回复
func onMingYan(msg *robot.Message) error {
	if !checkMatch(msg, "名言") {
		return nil
	}
	content, err := alapi.GetMingYan()
	if err != nil {
		return err
	}
	_, err = msg.ReplyText(content)
	return err
}

// 监听疫情相关的文字进行回复
func onCovid(msg *robot.Message) error {
	if !checkMatch(msg, "疫情") {
		return nil
	}
	hits := locateCovidRE.FindStringSubmatch(msg.Content)
	if len(hits) != 2 {
		return nil
	}
	cr, err := covid.GetCovidResponse(hits[1])
	if err != nil {
		msg.ReplyText("非常抱歉，未检索到该地区疫情数据")
		return err
	}
	_, err = msg.ReplyText(covid.PrintCovidSituation(cr))
	return err
}

func onSurvivalTime(msg *robot.Message) error {
	if !checkMatch(msg, "存活时间") {
		return nil
	}
	now := time.Now()
	nowString := now.Format(common.TimeFormat)
	d := now.Sub(begin)
	second := int(d.Seconds()) % 60
	min := int(d.Minutes()) % 60
	hour := int(d.Hours())
	text := fmt.Sprintf("截止至 %s , 机器人已经存活了 %d 小时 %d 分 %d 秒",
		nowString, hour, min, second)
	_, err := msg.ReplyText(text)
	return err
}
