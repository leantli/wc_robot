package common

import "fmt"

// constant.go文件中包含常量相关的定义

// 微信网页版相关接口
const (
	Jslogin           = "https://login.wx.qq.com/jslogin"                         // 微信网页版获取uuid的接口
	QRCode            = "https://login.weixin.qq.com/qrcode/"                     // 展示微信登陆二维码的页面，给用户扫描
	QRCodeLinux       = "https://login.weixin.qq.com/l/"                          // linux端登陆微信url
	Login             = "https://login.wx.qq.com/cgi-bin/mmwebwx-bin/login"       // 微信网页版检查是否成功登陆的接口
	Webwxnewloginpage = "https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxnewloginpage" // 用户扫描成功后给定的跳转链接，为真正登陆入口，请求后取得登陆公参

	// 以下path配合GetRequiredParams方法中获取的host使用
	Webwxinit            = "/cgi-bin/mmwebwx-bin/webwxinit"         // 微信初始化path
	Webwxstatusnotify    = "/cgi-bin/mmwebwx-bin/webwxstatusnotify" // 微信登陆状态通知path
	Webwxsync            = "/cgi-bin/mmwebwx-bin/webwxsync"         // 微信消息同步path
	Webwxsendmsg         = "/cgi-bin/mmwebwx-bin/webwxsendmsg"      // 微信发送消息path
	Webwxgetcontact      = "/cgi-bin/mmwebwx-bin/webwxgetcontact"
	Webwxsendmsgimg      = "/cgi-bin/mmwebwx-bin/webwxsendmsgimg"
	Webwxsendappmsg      = "/cgi-bin/mmwebwx-bin/webwxsendappmsg"
	Webwxsendvideomsg    = "/cgi-bin/mmwebwx-bin/webwxsendvideomsg"
	Webwxbatchgetcontact = "/cgi-bin/mmwebwx-bin/webwxbatchgetcontact"
	Webwxoplog           = "/cgi-bin/mmwebwx-bin/webwxoplog"
	Webwxverifyuser      = "/cgi-bin/mmwebwx-bin/webwxverifyuser"
	Synccheck            = "/cgi-bin/mmwebwx-bin/synccheck"
	Webwxuploadmedia     = "/cgi-bin/mmwebwx-bin/webwxuploadmedia"
	Webwxgetmsgimg       = "/cgi-bin/mmwebwx-bin/webwxgetmsgimg"
	Webwxgetvoice        = "/cgi-bin/mmwebwx-bin/webwxgetvoice"
	Webwxgetvideo        = "/cgi-bin/mmwebwx-bin/webwxgetvideo"
	Webwxlogout          = "/cgi-bin/mmwebwx-bin/webwxlogout"
	Webwxgetmedia        = "/cgi-bin/mmwebwx-bin/webwxgetmedia"
	Webwxupdatechatroom  = "/cgi-bin/mmwebwx-bin/webwxupdatechatroom"
	Webwxrevokemsg       = "/cgi-bin/mmwebwx-bin/webwxrevokemsg"
	Webwxcheckupload     = "/cgi-bin/mmwebwx-bin/webwxcheckupload"
	Webwxpushloginurl    = "/cgi-bin/mmwebwx-bin/webwxpushloginurl"
	Webwxgeticon         = "/cgi-bin/mmwebwx-bin/webwxgeticon"
	Webwxcreatechatroom  = "/cgi-bin/mmwebwx-bin/webwxcreatechatroom"
)

// 设置为uos桌面版，绕过设备被限制登录微信网页端
// 参考：https://github.com/wechaty/puppet-wechat/issues/127
const (
	UOSPatchClientVersion = "2.0.0"
	UOSPatchExtspam       = "Go8FCIkFEokFCggwMDAwMDAwMRAGGvAESySibk50w5Wb3uTl2c2h64jVVrV7gNs06GFlWplHQbY/5FfiO++1yH4ykC" +
		"yNPWKXmco+wfQzK5R98D3so7rJ5LmGFvBLjGceleySrc3SOf2Pc1gVehzJgODeS0lDL3/I/0S2SSE98YgKleq6Uqx6ndTy9yaL9qFxJL7eiA/R" +
		"3SEfTaW1SBoSITIu+EEkXff+Pv8NHOk7N57rcGk1w0ZzRrQDkXTOXFN2iHYIzAAZPIOY45Lsh+A4slpgnDiaOvRtlQYCt97nmPLuTipOJ8Qc5p" +
		"M7ZsOsAPPrCQL7nK0I7aPrFDF0q4ziUUKettzW8MrAaiVfmbD1/VkmLNVqqZVvBCtRblXb5FHmtS8FxnqCzYP4WFvz3T0TcrOqwLX1M/DQvcHa" +
		"GGw0B0y4bZMs7lVScGBFxMj3vbFi2SRKbKhaitxHfYHAOAa0X7/MSS0RNAjdwoyGHeOepXOKY+h3iHeqCvgOH6LOifdHf/1aaZNwSkGotYnYSc" +
		"W8Yx63LnSwba7+hESrtPa/huRmB9KWvMCKbDThL/nne14hnL277EDCSocPu3rOSYjuB9gKSOdVmWsj9Dxb/iZIe+S6AiG29Esm+/eUacSba0k8" +
		"wn5HhHg9d4tIcixrxveflc8vi2/wNQGVFNsGO6tB5WF0xf/plngOvQ1/ivGV/C1Qpdhzznh0ExAVJ6dwzNg7qIEBaw+BzTJTUuRcPk92Sn6QDn" +
		"2Pu3mpONaEumacjW4w6ipPnPw+g2TfywJjeEcpSZaP4Q3YV5HG8D6UjWA4GSkBKculWpdCMadx0usMomsSS/74QgpYqcPkmamB4nVv1JxczYIT" +
		"IqItIKjD35IGKAUwAA=="
)

// 登陆响应码
const (
	// 扫码登陆成功
	LoginStatusSuccess = "200"
	// 未扫码
	LoginStatusWait = "408"
	// 已扫码，待手机上确认
	LoginStatusScaned = "201"
	// 未扫码超时
	LoginStatusTimeout = "400"
)

var mapLoginDesc = map[string]string{
	LoginStatusSuccess: "扫码登陆成功",
	LoginStatusWait:    "未扫码",
	LoginStatusScaned:  "已扫码，待手机上确认",
	LoginStatusTimeout: "未扫码超时",
}

// GetLoginCodeDesc 根据登陆校验时响应码返回对应的错误信息描述
func GetLoginCodeDesc(loginCode string) string {
	if desc, exist := mapLoginDesc[loginCode]; exist {
		return desc
	}
	return fmt.Sprintf("未知的登陆状态响应码: %s", loginCode)
}

const (
	Ret_Success = 0
	Ret_Logout  = 1101
)

// mapRetCodeDesc 基本响应码/描述
var mapRetDesc = map[int]string{
	Ret_Success: "成功",
	-14:         "ticket错误",
	1:           "传入参数错误",
	1100:        "未登录提示",
	Ret_Logout:  "未检测到登录",
	1102:        "cookie值无效",
	1203:        "当前登录环境异常, 为了安全起见请不要在web端进行登录",
	1205:        "操作频繁",
}

// 获取基本响应码retCode的信息描述
func GetRetDesc(retCode int) string {
	if desc, exist := mapRetDesc[retCode]; exist {
		return desc
	}
	return fmt.Sprintf("未知的基本响应码retCode: %d", retCode)
}

const (
	Selector_Normal           = "0"
	Selector_Msg              = "2"
	Selector_ENTER_LEAVE_CHAT = "7"
)

// mapSelectorDesc selector 描述
var mapSelectorDesc = map[string]string{
	Selector_Normal:           "正常",
	Selector_Msg:              "有新消息",
	"4":                       "有人修改了自己的昵称或你修改了别人的备注",
	"6":                       "存在删除或新增的好友信息",
	Selector_ENTER_LEAVE_CHAT: "进入或离开聊天界面",
}

// 获取基本响应码retCode的信息描述
func GetSelectorDesc(selector string) string {
	if desc, exist := mapSelectorDesc[selector]; exist {
		return desc
	}
	return fmt.Sprintf("未知的selector: %s", selector)
}

// 登陆响应码
const (
	// 文本消息
	MT_TEXT = 1
	// 图片消息
	MT_IMAGE = 3
	// 语音消息
	MT_VOICE = 34
	// 好友请求消息
	MT_VERIFY = 37
	// 好友推荐消息
	MT_POSSIBLE_FRIEND = 40
	// 分享名片消息
	MT_SHARE_CARD = 42
	// 视频消息
	MT_VIDEO = 43
	// 表情消息
	MT_EMOTICON = 47
	// 位置消息
	MT_LOCATION = 48
	// 多媒体消息(分享链接)
	MT_MEDIA = 49
	// VOIP消息
	MT_VOIPMSG = 50
	// 状态通知，比如自身访问了某一个聊天页面
	MT_STATUS_NOTIFY = 51
	// VOIP结束通知
	MT_VOIP_NOTIFY = 52
	// VOIP邀请
	MT_VOIP_INVITE = 53
	// 小视频
	MT_MICRO_VIDEO = 62
	// 系统通知消息
	MT_SYS_NOTICE = 9999
	// 系统消息
	MT_SYS = 10000
	// 撤回消息
	MT_RECALLED = 10002
)

var mapMsgTypeDesc = map[int]string{
	MT_TEXT:            "文本消息",
	MT_IMAGE:           "图片消息",
	MT_VOICE:           "语音消息",
	MT_VERIFY:          "好友请求消息",
	MT_POSSIBLE_FRIEND: "好友推荐消息",
	MT_SHARE_CARD:      "分享名片消息",
	MT_VIDEO:           "视频消息",
	MT_EMOTICON:        "表情消息",
	MT_LOCATION:        "位置消息",
	MT_MEDIA:           "多媒体消息(分享链接)",
	MT_VOIPMSG:         "VOIP消息",
	MT_STATUS_NOTIFY:   "状态通知，比如自身访问了某一个聊天页面",
	MT_VOIP_NOTIFY:     "VOIP结束通知",
	MT_VOIP_INVITE:     "VOIP邀请",
	MT_MICRO_VIDEO:     "小视频",
	MT_SYS_NOTICE:      "系统通知消息",
	MT_SYS:             "系统消息",
	MT_RECALLED:        "撤回消息",
}

// GetMsgTypeDesc 根据登陆校验时响应码返回对应的错误信息描述
func GetMsgTypeDesc(MsgType int) string {
	if desc, exist := mapMsgTypeDesc[MsgType]; exist {
		return desc
	}
	return fmt.Sprintf("未知的消息类型MsgType: %d", MsgType)
}

const TimeFormat string = "2006-01-02 15:04:05"
