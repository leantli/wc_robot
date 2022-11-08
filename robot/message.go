package robot

import (
	"log"
	"strconv"
	"strings"
	"time"

	"wc_robot/common"
)

// message.go 定义 Message 结构体以及相关的通用方法

// Message 微信消息结构体
type Message struct {
	MsgId            string // 服务端返回的消息id, 可用于撤回消息接口；若消息为图片，还可用于调用微信获取图片接口
	FromUserName     string // 发送消息的人的username
	ToUserName       string // 接受消息的人的username
	MsgType          int    // 消息类型，1为文字，3为图片，34语音消息，43小视频消息，47表情消息，具体参考
	Content          string // 消息内容
	Status           int
	ImgStatus        int
	CreateTime       int64
	VoiceLength      int64
	PlayLength       int64
	FileName         string // 文件名
	FileSize         string // 文件大小
	MediaId          string // 多媒体消息id(图片、视频等)
	Url              string // 多媒体消息访问链接
	AppMsgType       int
	StatusNotifyCode int
	ForwardFlag      int
	AppInfo          *AppInfo
	HasProductId     int
	Ticket           string
	ImgHeight        int
	ImgWidth         int
	SubMsgType       int
	NewMsgId         int64
	OriContent       string
	EncryFileName    string
}

type AppInfo struct {
	AppID string
	Type  int
}

// SendMessage 微信发送消息时消息的结构体
type SendMessage struct {
	Type         int
	Content      string
	FromUserName string
	ToUserName   string
	LocalID      string // 时间戳(unix.milli)左移4位，后面4位为随机数
	ClientMsgId  string
	MediaId      string
}

// NewSendMessage 生成一个SendMessage对象
func NewSendMessage(t int, content, from, to, mediaID string) *SendMessage {
	ti := strconv.FormatInt(time.Now().UnixMilli(), 10)
	return &SendMessage{
		Type:         t,
		Content:      content,
		FromUserName: from,
		ToUserName:   to,
		MediaId:      mediaID,
		LocalID:      ti,
		ClientMsgId:  ti,
	}
}

// 获取消息的发送者
func (m *Message) GetSender() *User {
	return nil
}

// ReplyText 以文本消息的方式回复，返回发送消息的MsgID, err
func (m *Message) ReplyText(content string) (string, error) {
	sm := NewSendMessage(common.MT_TEXT, content, Storage.Self.UserName, m.FromUserName, "")
	smr, err := Caller.SendMsg(Storage.RequiredParams, sm)
	return smr.MsgID, err
}

// 判断消息是否为text类型
func (m *Message) IsText() bool {
	return m.MsgType == common.MT_TEXT && m.Url == ""
}

// 判断消息是否为自身发出
func (m *Message) IsFromSelf() bool {
	return m.FromUserName == Storage.Self.UserName
}

// 判断是否为群组消息
func (m *Message) IsFromGroup() bool {
	return strings.HasPrefix(m.FromUserName, "@@") || (strings.HasPrefix(m.ToUserName, "@@") && m.IsFromSelf())
}

// 判断是否为好友消息
func (m *Message) IsFromMember() bool {
	return strings.HasPrefix(m.FromUserName, "@") && !m.IsFromGroup() && !m.IsFromSelf()
}

// 判断消息发送者nickname是否为指定名称，支持多个nickname，通过","分隔
func (m *Message) IsSentByNickName(n string) bool {
	names := strings.Split(n, ",")
	for _, name := range names {
		u, ok := Storage.MemberMap[m.FromUserName]
		if !ok {
			log.Printf("[WARN]%s 昵称不在好友列表中", m.FromUserName)
			continue
		}
		if u.NickName == name {
			return true
		}
	}
	return false
}

// 判断消息发送者remarkname是否为指定名称(暂未支持群聊备注，目前微信回送的数据中均未显示群聊的remarkname字段)，支持多个remarkname, 通过","分隔
func (m *Message) IsSentByRemarkName(n string) bool {
	// 备注只支持来自好友的，不是好友直接返回
	if !m.IsFromMember() {
		return false
	}
	names := strings.Split(n, ",")
	for _, name := range names {
		u, ok := Storage.MemberMap[m.FromUserName]
		if !ok {
			log.Printf("[WARN]%s 备注不在好友列表中", m.FromUserName)
			continue
		}
		if u.RemarkName == name {
			return true
		}
	}
	return false
}
