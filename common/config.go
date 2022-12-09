package common

import (
	"io"
	"log"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

type Config struct {
	RobotName          string `yaml:"robot_name"`           // 机器人在微信中的昵称
	OnContactNickNames string `yaml:"on_contact_nicknames"` // 机器人只回复的人/群的nickname, 未设置则不会回复任何人

	WeatherMsgHandle *WeatherMsgHandle  `yaml:"weather_msg_handle"` // 天气、空气质量消息自动回复
	ALAPI            *ALAPI             `yaml:"alapi"`              // 情话、鸡汤、名言消息自动回复
	CovidMsgHandle   *CovidMsgHandle    `yaml:"covid_msg_handle"`   // 疫情消息自动回复
	OpenAIHandle     *OpenAI            `yaml:"openai"`             // OpenAI API 相关配置
	WeatherSchedule  []*WeatherSchedule `yaml:"weather_schedules"`  // 每天定时发送天气提醒
	ClockInSchedule  []*ClockInSchedule `yaml:"clock_in_schedules"` // 每天定时发送信息
	DaysMatters      []*DaysMatter      `yaml:"days_matters"`       // 重要的日子， 设置则会每天定时提醒距离该日子的时间
}

// 天气、空气质量消息自动回复相关参数
type WeatherMsgHandle struct {
	SwitchOn bool   `yaml:"switch_on"` // "天气、空气质量回复"开关，true为开，false为关闭
	CityCode string `yaml:"city_code"` // [Deprecated] 回复天气时回复的地区
}

// 情话、鸡汤、名言自动回复相关参数
type ALAPI struct {
	SwitchOn bool   `yaml:"switch_on"` // "情话、鸡汤、名言回复"开关，true为开，false为关闭
	Token    string `yaml:"token"`     // alapi 调用的 token
}

// 疫情消息自动回复相关参数
type CovidMsgHandle struct {
	SwitchOn bool `yaml:"switch_on"` // "疫情回复"开关
}

// OpenAI API 相关配置
type OpenAI struct {
	ApiKey           string `yaml:"api_key"`                   // OpenAI 的 API_KEY
	GPTTextSwitchOn  bool   `yaml:"gpt_text_switch_on"`        // GPT 文字回复开关, true 为开, false 为关闭
	GPTTextIsDefault bool   `yaml:"gpt_text_is_default_reply"` // GPT 文字回复是否作为默认回复(当不触发其他关键词时)，不是默认回复时需要通过 "gpt xxx" 才能触发
}

// 重要的日子， 设置则会每天定时提醒距离该日子的时间
type DaysMatter struct {
	SwitchOn      bool   `yaml:"switch_on"`      // 重要的日子 开关
	ToNickNames   string `yaml:"to_nicknames"`   // 重要的日子 要提醒的用户昵称，支持多个，通过英文","分隔
	ToRemarkNames string `yaml:"to_remarknames"` // 发送给哪些人，发送多人则通过英文","分隔,比如a,b,c(备注，微信没有返回群聊备注)
	Times         string `yaml:"times"`          // 重要的日子 要提醒的时间，支持多个时间, 通过英文","分隔，比如00:00,9:00(注意合法时间)
	Date          string `yaml:"date"`           // 重要的日子 日子，格式为"2022-9-29"
	Content       string `yaml:"content"`        // 重要的日子 事项
}

// 每天定时发送天气提醒
type WeatherSchedule struct {
	SwitchOn      bool   `yaml:"switch_on"`      // 该功能开关，true为开，false为关闭
	ToNickNames   string `yaml:"to_nicknames"`   // 发送给哪些人/群，发送多人/群则通过英文","分隔,比如a,b,c(昵称)
	ToRemarkNames string `yaml:"to_remarknames"` // 发送给哪些人，发送多人则通过英文","分隔,比如a,b,c(备注，微信没有返回群聊备注)
	Times         string `yaml:"times"`          // 发送的时间，支持多个时间，通过英文","分隔，比如00:00,9:00(注意合法时间)
	CityCode      string `yaml:"city_code"`      // 发送哪个城市/地区的天气情况，基于citycode，详见../common/weather/cityID.xlsx
}

// 基于昵称，定时发送信息
type ClockInSchedule struct {
	SwitchOn      bool   `yaml:"switch_on"`      // 该定时任务开关，true为开，false为关
	ToNickNames   string `yaml:"to_nicknames"`   // 发送给哪些人/群，发送多人/群则通过英文","分隔,比如a,b,c(昵称)
	ToRemarkNames string `yaml:"to_remarknames"` // 发送给哪些人，发送多人则通过英文","分隔,比如a,b,c(备注，微信没有返回群聊备注)
	Times         string `yaml:"times"`          // 发送的时间，支持多个时间，通过英文","分隔，比如00:00,9:00(注意合法时间)
	Text          string `yaml:"text"`           // 发送的信息
}

const configPath = "config.yaml"

var config Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		f, err := os.Open(configPath)
		if err != nil {
			log.Fatalf("[ERROR]打开文件config.yaml失败, err: %v", err)
		}
		defer f.Close()
		b, err := io.ReadAll(f)
		if err != nil {
			log.Fatalf("[ERROR]读取文件config.yaml失败, err: %v", err)
		}
		if err := yaml.Unmarshal(b, &config); err != nil {
			log.Fatalf("[ERROR]解析文件失败, 请查看是否格式有问题, err: %v", err)
		}
		for _, temp := range config.WeatherSchedule {
			log.Printf("[INFO]%v\n", temp)
		}
		for _, temp := range config.ClockInSchedule {
			log.Printf("[INFO]%v\n", temp)
		}
		for _, temp := range config.DaysMatters {
			log.Printf("[INFO]%v\n", temp)
		}
	})
	return &config
}
