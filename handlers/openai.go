package handlers

import (
	"regexp"

	"wc_robot/common"
	"wc_robot/common/openai"
	"wc_robot/robot"
)

var (
	gptFnRE      = regexp.MustCompile(`(?i)gpt\s(.*)`)
	gptDefaultRE = regexp.MustCompile(`@` + common.GetConfig().RobotName + `[  \s](.*)`)
)

func onGPTTextChecker(msg *robot.Message) bool {
	if common.GetConfig().OpenAIHandle.GPTTextIsDefault {
		// 作为基本消息处理逻辑(不匹配上其他功能就会自动用 GPT 处理)
		return checkMatch(msg, []string{})
	}
	// 作为功能提供，只回复"gpt xxx"格式
	return checkMatch(msg, []string{"gpt", "GPT"})
}

func onGPTText(msg *robot.Message) error {
	var hits []string
	// 使用功能形式的 gpt，则无需考虑是群还是个人的消息
	if !common.GetConfig().OpenAIHandle.GPTTextIsDefault {
		hits = gptFnRE.FindStringSubmatch(msg.Content)
	}
	if common.GetConfig().OpenAIHandle.GPTTextIsDefault {
		if msg.IsFromGroup() {
			hits = gptDefaultRE.FindStringSubmatch(msg.Content)
		} else {
			hits = make([]string, 2)
			hits[1] = msg.Content
		}
	}
	if len(hits) != 2 {
		return nil
	}
	reply, err := openai.GetGPTTextReply(hits[1])
	if err != nil {
		return err
	}
	_, err = msg.ReplyText(hits[1] + reply)
	return err
}
