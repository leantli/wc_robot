package main

import (
	"log"
	"strings"
	"time"

	"wc_robot/common"
	"wc_robot/common/alapi"
	"wc_robot/common/weather"
	"wc_robot/robot"
)

// 日志设置初始化
func init() {
	log.SetFlags(log.Llongfile | log.Ldate | log.Ltime)

	// 部署在 linux 上可直接通过 nohup ./wc_robot > robot.log & 运行并打印日志
	// 本机测试运行可取消下方注释，记录 log 便于观察

	// // 打印日志到本地 wc_robot.log
	// outputLogPath := "wc_robot.log"
	// f, err := os.Create(outputLogPath)
	// if err != nil {
	// 	log.Println("[WARN]创建日志文件失败, 日志仅输出在控制台")
	// }
	// w := io.MultiWriter(os.Stdout, f)
	// log.SetOutput(w)
}

func main() {
	begin := time.Now()
	defer func() {
		log.Printf("[INFO]本次机器人运行时间为: %s", time.Since(begin).String())
	}()
	config := common.GetConfig()

	r := robot.NewRobot()
	r.Chain.RegisterHandler("菜单|功能|会什么回复", onMenu)
	if config.WeatherMsgHandle.SwitchOn {
		r.Chain.RegisterHandler("天气回复", onWeather)
		r.Chain.RegisterHandler("空气质量回复", onAQI)
	}
	if config.ALAPI.SwitchOn {
		r.Chain.RegisterHandler("名言回复", onMingYan)
		r.Chain.RegisterHandler("情话回复", onQingHua)
		r.Chain.RegisterHandler("鸡汤回复", onSoul)
	}
	if err := r.Login(); err != nil {
		log.Println(err)
	}
	robot.InitTasks(config)
	r.Block()
}

// 基础校验，机器人只回复文字、监听的nickname、非自己
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

// 监听菜单｜功能｜会什么相关的文字进行回复
func onMenu(msg *robot.Message) error {
	config := common.GetConfig()
	if !checkOnContact(msg) {
		return nil
	}
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
	_, err := msg.ReplyText("你好呀👋\n" + `目前只支持"天气"、"空气质量(指标含义)"、"情话"、"鸡汤"、"名言"相关的问题哦`)
	return err
}

// 监听天气相关的文字进行回复
func onWeather(msg *robot.Message) error {
	config := common.GetConfig()
	if !checkOnContact(msg) {
		return nil
	}
	if msg.IsFromGroup() {
		if !(strings.Contains(msg.Content, "@"+config.RobotName) && strings.Contains(msg.Content, "天气")) {
			return nil
		}
	}
	if msg.IsFromMember() {
		if !strings.Contains(msg.Content, "天气") {
			return nil
		}
	}

	w, err := weather.GetWeather(config.WeatherMsgHandle.CityCode)
	if err != nil {
		return err
	}
	_, err = msg.ReplyText(weather.CurrentWeatherInfo(w))
	return err
}

// 监听空气质量(指标含义) 的文字进行回复
func onAQI(msg *robot.Message) error {
	config := common.GetConfig()
	if !checkOnContact(msg) {
		return nil
	}
	if msg.IsFromGroup() {
		if !(strings.Contains(msg.Content, "@"+config.RobotName) && (strings.Contains(msg.Content, "空气质量"))) {
			return nil
		}
		if strings.Contains(msg.Content, "指标含义") {
			msg.ReplyText(weather.AQIIndicesDesc())
			return nil
		}
	}
	if msg.IsFromMember() {
		if !strings.Contains(msg.Content, "空气质量") {
			return nil
		}
		if strings.Contains(msg.Content, "指标含义") {
			msg.ReplyText(weather.AQIIndicesDesc())
			return nil
		}
	}

	w, err := weather.GetWeather(config.WeatherMsgHandle.CityCode)
	if err != nil {
		return err
	}
	_, err = msg.ReplyText(weather.AQIInfo(w))
	return err
}

// 监听心灵鸡汤相关的文字进行回复
func onSoul(msg *robot.Message) error {
	config := common.GetConfig()
	if !checkOnContact(msg) {
		return nil
	}
	if msg.IsFromGroup() {
		if !(strings.Contains(msg.Content, "@"+config.RobotName) && strings.Contains(msg.Content, "鸡汤")) {
			return nil
		}
	}
	if msg.IsFromMember() {
		if !strings.Contains(msg.Content, "鸡汤") {
			return nil
		}
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
	config := common.GetConfig()
	if !checkOnContact(msg) {
		return nil
	}
	if msg.IsFromGroup() {
		if !(strings.Contains(msg.Content, "@"+config.RobotName) && strings.Contains(msg.Content, "情话")) {
			return nil
		}
	}
	if msg.IsFromMember() {
		if !strings.Contains(msg.Content, "情话") {
			return nil
		}
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
	config := common.GetConfig()
	if !checkOnContact(msg) {
		return nil
	}
	if msg.IsFromGroup() {
		if !(strings.Contains(msg.Content, "@"+config.RobotName) && strings.Contains(msg.Content, "名言")) {
			return nil
		}
	}
	if msg.IsFromMember() {
		if !strings.Contains(msg.Content, "名言") {
			return nil
		}
	}

	content, err := alapi.GetMingYan()
	if err != nil {
		return err
	}
	_, err = msg.ReplyText(content)
	return err
}