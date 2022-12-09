package handlers

import (
	"wc_robot/common/alapi"
	"wc_robot/robot"
)

func onSoulChecker(msg *robot.Message) bool {
	return checkMatch(msg, []string{"鸡汤"})
}

func onQingHuaChecker(msg *robot.Message) bool {
	return checkMatch(msg, []string{"情话"})
}

func onMingYanChecker(msg *robot.Message) bool {
	return checkMatch(msg, []string{"名言"})
}

// 监听心灵鸡汤相关的文字进行回复
func onSoul(msg *robot.Message) error {
	s, err := alapi.GetSoul()
	if err != nil {
		return err
	}
	_, err = msg.ReplyText(s)
	return err
}

// 监听情话相关的文字进行回复
func onQingHua(msg *robot.Message) error {
	content, err := alapi.GetQinghua()
	if err != nil {
		return err
	}
	_, err = msg.ReplyText(content)
	return err
}

// 监听名言相关的文字进行回复
func onMingYan(msg *robot.Message) error {
	content, err := alapi.GetMingYan()
	if err != nil {
		return err
	}
	_, err = msg.ReplyText(content)
	return err
}
