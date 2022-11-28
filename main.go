package main

import (
	"log"
	"time"

	"wc_robot/common"
	"wc_robot/handlers"
	"wc_robot/robot"
	"wc_robot/tasks"
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
	r := robot.NewRobot()
	handlers.InitHandlers(r)
	if err := r.Login(); err != nil {
		log.Println(err)
	}
	tasks.InitTasks(common.GetConfig())
	r.Block()
}
