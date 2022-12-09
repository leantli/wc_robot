package handlers

import (
	"strings"

	"wc_robot/common"
	"wc_robot/robot"
)

// handlers 包，包含了对消息的各类处理方法

// TODO(leantli): 这些都迁移一下
func InitHandlers(r *robot.Robot) {
	config := common.GetConfig()
	r.Chain.RegisterGlobalCheck(checkOnContact)
	r.Chain.RegisterHandler("菜单|功能|会什么回复", onMenuChecker, onMenu)
	r.Chain.RegisterHandler("存活时间回复", onSurvivalTimeChecker, onSurvivalTime)
	if config.WeatherMsgHandle.SwitchOn {
		r.Chain.RegisterHandler("天气回复", onWeatherChecker, onWeather)
		r.Chain.RegisterHandler("空气质量回复", onAQIChecker, onAQI)
	}
	if config.ALAPI.SwitchOn {
		r.Chain.RegisterHandler("名言回复", onMingYanChecker, onMingYan)
		r.Chain.RegisterHandler("情话回复", onQingHuaChecker, onQingHua)
		r.Chain.RegisterHandler("鸡汤回复", onSoulChecker, onSoul)
	}
	if config.CovidMsgHandle.SwitchOn {
		r.Chain.RegisterHandler("疫情回复", onCovidChecker, onCovid)
	}
	if config.OpenAIHandle.GPTTextSwitchOn {
		r.Chain.RegisterHandler("GPT 文字回复", onGPTTextChecker, onGPTText)
	}
}

// 全局校验，机器人只回复文字、监听的nickname、非自己，其余都不回复，返回 false
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
func checkMatch(msg *robot.Message, keywords []string) bool {
	config := common.GetConfig()
	if msg.IsFromGroup() {
		if !strings.Contains(msg.Content, "@"+config.RobotName) {
			return false
		}
	}
	if len(keywords) == 0 {
		return true
	}
	for _, keyword := range keywords {
		if strings.Contains(msg.Content, keyword) {
			return true
		}
	}
	return false
}
