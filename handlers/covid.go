package handlers

import (
	"regexp"

	"wc_robot/common/covid"
	"wc_robot/robot"
)

var locateCovidRE = regexp.MustCompile("([\u4e00-\u9fa5]{1,6})疫情") // {城市}疫情正则，匹配位置

func onCovidChecker(msg *robot.Message) bool {
	return checkMatch(msg, []string{"疫情"})
}

// 监听疫情相关的文字进行回复
func onCovid(msg *robot.Message) error {
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
