package openai

// 同包名下存放实体

// GPTCompletionsResp gpt/completions api 的请求体
type GPTCompletionsReq struct {
	Model            string  `json:"model"`             // 请求采用的模型，官网这里给的例子是"text-davinci-003"
	Prompt           string  `json:"prompt"`            // 我们传入的消息
	MaxTokens        int     `json:"max_tokens"`        // 回复的字数限制，最大设为 4096，单位看官网感觉应该是字节，一个中文 2 个token
	Temperature      float32 `json:"temperature"`       // 模型参数，值越高答案越不稳定，官方文档说更具创造性的程序可以采用 0.9，默认为 1
	TopP             int     `json:"top_p"`             // 模型参数，使用核心采样替代上面的温度采样，官方文档建议只改二者其一即可
	FrequencyPenalty int     `json:"frequency_penalty"` // 正值增加模型讨论新主题的可能性，取值-2.0～2.0，默认为 0
	PresencePenalty  int     `json:"presence_penalty"`  // 正值降低模型逐字重复同一行的可能性，取值-2.0～2.0，默认为 0
}

// GPTCompletionsResp gpt/completions api 的响应结构体
type GPTCompletionsResp struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int      `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Text         string `json:"text"`
	Index        int    `json:"index"`
	Logprobs     int    `json:"logprobs"`
	FinishReason string `json:"finish_reason"`
}
