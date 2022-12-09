package openai

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"wc_robot/common"
	"wc_robot/common/utils"
)

// OpenAI 文档见 https://beta.openai.com/docs/api-reference/introduction

const baseUrl = "https://api.openai.com/v1"
const contentType = "application/json"

func getOpenAICompletions(req *GPTCompletionsReq) (*GPTCompletionsResp, error) {
	url, err := url.Parse(baseUrl + "/completions")
	if err != nil {
		return nil, err
	}
	body, err := utils.ToJsonBuff(req)
	if err != nil {
		return nil, err
	}
	r, err := http.NewRequest(http.MethodPost, url.String(), body)
	if err != nil {
		return nil, err
	}
	r.Header.Add("Authorization", "Bearer "+common.GetConfig().OpenAIHandle.ApiKey)
	r.Header.Add("Content-Type", contentType)
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("响应失败，状态码为 %s", resp.Status)
	}
	var gcresp GPTCompletionsResp
	if err := utils.ScanJson(resp, &gcresp); err != nil {
		return nil, err
	}
	return &gcresp, nil
}

// 获取 GPT 的文字模型回复
func GetGPTTextReply(msg string) (string, error) {
	req := &GPTCompletionsReq{
		Model:       "text-davinci-003",
		Prompt:      msg,
		MaxTokens:   2048, // 设置问题+回答最大长度 2048 字节
		Temperature: 0.8,
	}
	resp, err := getOpenAICompletions(req)
	if err != nil {
		return "", err
	}
	if len(resp.Choices) == 0 {
		log.Printf("choice 为空的情况, resp=%+v", resp)
		return "不好意思，我不知如何回答该问题", nil
	}
	if resp.Choices[0].FinishReason == "length" {
		log.Printf("顺便记录一下回答过长的情况, Q:%s A:%+v", msg, resp)
		return "不好意思，这个问题讨论起来太过于广泛，我无法简短的回答", nil
	}
	return resp.Choices[0].Text, nil
}
