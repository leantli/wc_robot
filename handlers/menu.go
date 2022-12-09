package handlers

import (
	"strings"

	"wc_robot/common"
	"wc_robot/robot"
)

func onMenuChecker(msg *robot.Message) bool {
	config := common.GetConfig()
	if msg.IsFromGroup() {
		if !(strings.Contains(msg.Content, "@"+config.RobotName) &&
			(strings.Contains(msg.Content, "菜单") || strings.Contains(msg.Content, "功能") || strings.Contains(msg.Content, "会什么"))) {
			return false
		}
	}
	if msg.IsFromMember() {
		if !(strings.Contains(msg.Content, "菜单") || strings.Contains(msg.Content, "功能") || strings.Contains(msg.Content, "会什么")) {
			return false
		}
	}
	return true
}

// 监听菜单｜功能｜会什么相关的文字进行回复
func onMenu(msg *robot.Message) error {
	_, err := msg.ReplyText("你好呀👋\n" +
		`支持自动回复"XX(城市/地区)天气","XX(城市/地区)空气质量"关键词(天气数据来源：小米天气)\n` +
		`支持自动回复"XX(城市/省份/国家)疫情"关键词(疫情数据来源：百度实时疫情)\n` +
		`每日定时发送天气预报\n` +
		`每日定时发送消息\n` +
		`重要的日子提醒(类似倒数日)\n` +
		`GPT 语言自动回复`)
	return err
}
