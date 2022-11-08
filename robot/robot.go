// Package robot 包含了Robot相关的结构体的定义及其方法
package robot

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"wc_robot/common"
)

type Robot struct {
	ctx    context.Context    // 上下文(主要配合cancel终止程序)
	cancel context.CancelFunc // 终止通知函数
	Chain  *MsgHandlerChain   // 消息处理链
}

// 返回一个Robot，并固定使用同一个deviceID
func NewRobot(mode ...Mode) *Robot {
	ctx, cancel := context.WithCancel(context.Background())
	r := &Robot{
		ctx:    ctx,
		cancel: cancel,
		Chain:  &MsgHandlerChain{},
	}
	Storage.RequiredParams = &RequiredParams{DeviceID: getDeviceID()}
	if len(mode) > 0 {
		Caller.SetMode(mode[0])
	}
	return r
}

// 读文件获取设备码，若无则随机生成设备码
func getDeviceID() string {
	f, err := os.Open("deviceID.robot")
	if os.IsExist(err) {
		var b bytes.Buffer
		b.ReadFrom(f)
		return b.String()
	}
	defer f.Close()
	// 若不存在则生成文件及deviceID
	f, err = os.Create("deviceID.robot")
	if err != nil {
		log.Fatalf("[ERROR]创建文件deviceID.robot失败, 错误码err=%v\n", err)
		return ""
	}
	defer f.Close()
	deviceID := generateDeviceID()
	f.WriteString(deviceID)
	return deviceID
}

// 随机生成设备码
func generateDeviceID() string {
	var b strings.Builder
	b.Grow(16)         // deviceID为16位
	b.WriteString("e") // deviceID首字母为e
	for i := 0; i < 15; i++ {
		b.WriteString(strconv.Itoa(rand.Intn(9)))
	}
	return b.String()
}

// isAlive robot 客户端是否存活， true为存活
func (r *Robot) isAlive() bool {
	select {
	case <-r.ctx.Done():
		return false
	default:
		return true
	}
}

// Block 用于保持robot运行，直至robot结束运行(发生异常或主动调用robot.cacel()结束进程)
func (r *Robot) Block() {
	<-r.ctx.Done()
}

// Login 登陆微信--流程如下：
// 1. 获取二维码uuid； 2. 扫描该uuid链接
// 3. 轮询该uuid二维码登陆状态； 4. 成功后获取公共参数(用于后续请求)
// 5. 初始化微信客户端，获取各类信息(用户信息、通讯录信息、首页联系信息等)
func (r *Robot) Login() error {
	uuid, err := Caller.GetLoginUUID()
	if err != nil {
		return fmt.Errorf("GetLoginUUID请求失败, 错误信息err: %v", err)
	}
	err = Caller.OpenQRCode(uuid)
	if err != nil {
		return fmt.Errorf("OpenQRCode执行命令出错, 错误信息err: %v", err)
	}

	var rspRaw []byte
	for {
		rsp, err := Caller.CheckLoginStatus(uuid)
		if err != nil {
			return err
		}
		log.Printf("[INFO]本次登陆状态检查码为%s, 信息为%s", rsp.Code, common.GetLoginCodeDesc(rsp.Code))
		switch rsp.Code {
		case common.LoginStatusWait:
			continue
		case common.LoginStatusScaned:
			continue
		case common.LoginStatusSuccess:
			rspRaw = rsp.Raw
		case common.LoginStatusTimeout:
			return fmt.Errorf("code=%s,msg=%s", rsp.Code, common.GetLoginCodeDesc(rsp.Code))
		default:
			return fmt.Errorf("未知登陆响应码: %s", rsp.Code)
		}
		break
	}

	rsp, err := Caller.GetRequiredParams(rspRaw)
	if err != nil {
		return err
	}

	Storage.RequiredParams = &RequiredParams{
		SKey:       rsp.SKey,
		WxSid:      rsp.WxSid,
		WxUin:      rsp.WxUin,
		PassTicket: rsp.PassTicket,
		DeviceID:   Storage.RequiredParams.DeviceID,
	}
	return r.wechatInit()
}

// wechatInit 微信客户端各类信息初始化(用户信息、通讯录信息、首页联系信息等)
func (r *Robot) wechatInit() error {
	rsp, err := Caller.WebInit(Storage.RequiredParams)
	if err != nil {
		return fmt.Errorf("wechatInit执行失败, err: %v", err)
	}
	Storage.Self = rsp.User
	Storage.SyncKey = rsp.SyncKey
	r.buildMemberMap(rsp.ContactList)

	members, _, err := Caller.GetMemberList(Storage.RequiredParams)
	if err != nil {
		return err
	}
	r.buildMemberMap(members)
	if err = Caller.LoginNotify(Storage.RequiredParams, Storage.Self.UserName); err != nil {
		return err
	}
	// goroutine执行同步消息检查+新消息拉取
	go r.sync()
	return nil
}

// buildMemberMap 构建member map映射, {username : *User}
func (r *Robot) buildMemberMap(members []*User) {
	if Storage.MemberMap == nil {
		Storage.MemberMap = make(map[string]*User)
	}
	for _, user := range members {
		Storage.MemberMap[user.UserName] = user
		log.Printf("[INFO]MemberMap记录,%s--%s", user.UserName, user.NickName)
		// 如果是群聊，还需要将群聊中群成员的username以及相关信息存入MemberMap
		// TODO(leantli): 这里考虑后续对群聊相关的操作有更多需求进行补充，比如只针对群聊中的某个人进行监听回复等(优先级较低)
		if strings.HasPrefix(user.UserName, "@@") {
			r.buildMemberMap(user.MemberList)
		}
	}
}

// sync 同步消息检查+新消息拉取，配合goroutine使用，若发生错误，则传递cancel，使客户端退出
// 收到消息时配合MessageHanlder使用，责任链处理
func (r *Robot) sync() {
	for {
		var catch error
		var retcode int
		// 内层循环主要是轮询synccheck
		for r.isAlive() {
			sr, err := Caller.SyncCheck(Storage.RequiredParams, Storage.SyncKey)
			if err != nil {
				catch = fmt.Errorf("SyncCheck执行中报错,err=%v", err)
				break
			}
			log.Printf("[INFO]Sync得到响应 retCode=%s, selector=%s\n", sr.RetCode, sr.Selector)
			if !sr.IsSuccess() {
				retcode, _ = strconv.Atoi(sr.RetCode)
				break
			}
			if sr.IsNormal() {
				continue
			}
			if err := sr.checkSpecial(); err != nil {
				log.Printf("[ERROR]掉线, err: %v\n", err)
				r.cancel()
				return
			}
			smr, err := Caller.SyncMsg(Storage.RequiredParams, Storage.SyncKey)
			if err != nil {
				catch = fmt.Errorf("SyncMsg执行中报错,err=%v", err)
				break
			}
			// 更新SyncKey
			Storage.SyncKey = smr.SyncKey
			// 更新memberMap
			r.buildMemberMap(smr.ModContactList)
			for _, message := range smr.AddMsgList {
				if _, ok := Storage.MemberMap[message.FromUserName]; !ok {
					continue
				}
				log.Printf("[INFO]收到 %s 的消息 %v\n", Storage.MemberMap[message.FromUserName].NickName, message.Content)
				r.Chain.Handle(message)
			}
		}
		// 如果是用户在手机上主动退出则直接退出，不再循环请求，暂时先写死处理如下
		if retcode == common.Ret_Logout {
			log.Printf("[INFO]用户于手机客户端退出登陆, 响应码:%d, 描述信息%s\n", retcode, common.GetRetDesc(retcode))
			r.cancel()
			return
		}
		// 外层的循环主要保证不因网络问题中断请求, 比如短时间内多次synccheck，会得到HTTP status code "0"
		log.Printf("[ERROR]Sync执行出错, err: %v\n", catch)
		time.Sleep(time.Second * 7)
	}
}
