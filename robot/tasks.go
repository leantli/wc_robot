package robot

import (
	"fmt"
	"log"
	"strings"
	"time"

	"wc_robot/common"
	"wc_robot/common/weather"

	"github.com/jasonlvhit/gocron"
)

// tasks.go 定义了robot相关的定时任务

const (
	separator  = ","
	formatTime = "2006-1-2" // 用于time.Format(), 格式化go的时间
)

// 初始化定时任务并运行
func InitTasks(config *common.Config) {
	// 每天定时发送今日天气情况给指定好友/群(昵称)
	for _, wn := range config.WeatherSchedule {
		if !wn.SwitchOn {
			continue
		}
		times := strings.Split(wn.Times, separator)
		for _, time := range times {
			if err := gocron.Every(1).Day().At(time).Do(getWeatherToMembers, wn.CityCode, wn.ToNickNames, wn.ToRemarkNames); err != nil {
				log.Printf("[ERROR]添加定时任务%v出错,err: %v\n", wn, err)
			}
		}
	}

	// 每天定时发送文本信息给指定好友/群(昵称)
	for _, cin := range config.ClockInSchedule {
		if !cin.SwitchOn {
			continue
		}
		times := strings.Split(cin.Times, separator)
		for _, time := range times {
			if err := gocron.Every(1).Day().At(time).Do(sendText, cin.Text, cin.ToNickNames, cin.ToRemarkNames); err != nil {
				log.Printf("[ERROR]添加定时任务%v出错,err: %v\n", cin, err)
			}
		}
	}

	// 重要的日子定时任务
	for _, dm := range config.DaysMatters {
		if !dm.SwitchOn {
			continue
		}
		times := strings.Split(dm.Times, separator)
		date, _ := time.Parse(formatTime, dm.Date)
		for _, time := range times {
			if err := gocron.Every(1).Day().At(time).Do(daysMatter, dm.Content, date, dm.ToNickNames, dm.ToRemarkNames); err != nil {
				log.Printf("[ERROR]添加定时任务%v出错,err: %v\n", dm, err)
			}
		}
	}

	gocron.Start()
}

// 获取指定天气并发送到指定好友/群聊
func getWeatherToMembers(cityCode, toNickNames, ToRemarkNames string) {
	users := getToUsers(toNickNames, ToRemarkNames)
	if len(users) == 0 {
		return
	}

	w, err := weather.GetWeather(cityCode)
	if err != nil {
		log.Printf("[ERROR]获取天气信息失败 err:%v", err)
	}
	for _, user := range users {
		_, err = Storage.Self.SendTextToUser(user, weather.CurrentWeatherInfo(w))
		if err != nil {
			log.Printf("[ERROR]发送消息给 %s 失败, err:%v", user.NickName, err)
		}
	}
}

// 发送消息给指定好友/群聊
func sendText(content, toNickNames, ToRemarkNames string) {
	users := getToUsers(toNickNames, ToRemarkNames)
	if len(users) == 0 {
		return
	}

	for _, user := range users {
		log.Printf("[INFO]要发送消息给%s\n", user.NickName)
		_, err := Storage.Self.SendTextToUser(user, content)
		if err != nil {
			log.Printf("[ERROR]发送消息给 %s 失败, err:%v", user.NickName, err)
		}
	}
}

// 重要的日子
func daysMatter(content string, day time.Time, toNickNames, ToRemarkNames string) {
	users := getToUsers(toNickNames, ToRemarkNames)
	if len(users) == 0 {
		return
	}

	// 只取当前时间的年月日
	now, _ := time.Parse(formatTime, time.Now().Format(formatTime))
	interval := int(day.Sub(now).Hours() / 24)
	var s string
	if interval == 0 {
		s = fmt.Sprintf("今天就是%s!", content)
	}
	if interval > 0 {
		s = fmt.Sprintf("还有%d天就是%s", interval, content)
	}
	if interval < 0 {
		interval = -interval
		s = fmt.Sprintf("%s已经%d天", content, interval)
	}

	for _, user := range users {
		_, err := Storage.Self.SendTextToUser(user, s)
		if err != nil {
			log.Printf("[ERROR]发送消息给 %s 失败, err:%v", user.NickName, err)
		}
	}
}

// 根据nicknames和remarknames返回得到user
func getToUsers(toNickNames, ToRemarkNames string) []*User {
	var nicknames []string
	var remarknames []string
	if len(toNickNames) != 0 {
		nicknames = strings.Split(toNickNames, separator)
	}
	if len(ToRemarkNames) != 0 {
		remarknames = strings.Split(ToRemarkNames, separator)
	}

	var users []*User
	for _, name := range nicknames {
		ms := Storage.SearchMembersByNickName(1, name)
		if len(ms) == 0 {
			log.Printf("[ERROR]未找到该用户/群: %s, 请确认用户/群聊昵称是否输入正确, 不支持备注\n", name)
			continue
		}
		log.Printf("[INFO]根据toNickNames: %s, 检索到用户如下\n", toNickNames)
		for _, m := range ms {
			log.Printf("[INFO]用户昵称：%s\n", m.NickName)
		}
		users = append(users, ms...)
	}

	for _, name := range remarknames {
		ms := Storage.SearchMembersByRemarkName(1, name)
		if len(ms) == 0 {
			log.Printf("[ERROR]未找到该用户: %s, 请确认用户备注是否输入正确, 不支持群聊备注, 因为微信未返回群聊备注名\n", name)
			continue
		}
		log.Printf("[INFO]根据ToRemarkNames:%s, 检索到用户如下\n", ToRemarkNames)
		for _, m := range ms {
			log.Printf("[INFO]用户昵称 %s, 备注 %s\n", m.NickName, m.RemarkName)
		}
		users = append(users, ms...)
	}
	return users
}
