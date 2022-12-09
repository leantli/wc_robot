package handlers

import (
	"fmt"
	"time"

	"wc_robot/common"
	"wc_robot/robot"
)

var begin = time.Now() // 存活时间计时

func onSurvivalTimeChecker(msg *robot.Message) bool {
	return checkMatch(msg, []string{"存活时间"})
}

func onSurvivalTime(msg *robot.Message) error {
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
