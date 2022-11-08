package robot

import (
	"fmt"
	"strconv"
	"time"

	"wc_robot/common"
)

// dto.go文件中主要包含了一些响应结构体，主要用于数据流转，因为响应的部分字段并无存储价值

// LoginStatusResponse 是CheckLoginStatus请求返回结果的封装
// 这里本来想直接封装code和redirect_uri的，但是如果检测未扫码时只会返回code
type LoginStatusResponse struct {
	Code string
	Raw  []byte
}

// 在扫码成功后，请求公参的响应xml结构体
type RequiredParamsResponse struct {
	Ret         int    `xml:"ret"`         // 响应码
	Message     string `xml:"message"`     // 响应消息
	SKey        string `xml:"skey"`        // 公参
	WxSid       string `xml:"wxsid"`       // 公参
	WxUin       int64  `xml:"wxuin"`       // 公参
	PassTicket  string `xml:"pass_ticket"` // 公参
	IsGrayScale int    `xml:"isgrayscale"`
}

// 请求的微信host，携带返回各类domain的方法
type Host string

// 微信正常请求域名的前缀
func (h Host) BaseDomain() string {
	return fmt.Sprintf("https://%s", h)
}

// 微信进行文件传输时域名的前缀
func (h Host) FileDomain() string {
	return fmt.Sprintf("https://file.%s", h)
}

// 微信进行同步检查时域名的前缀
func (h Host) SyncDomain() string {
	return fmt.Sprintf("https://webpush.%s", h)
}

// 微信初始化响应结构体
type WebInitResponse struct {
	BaseResponse        *BaseResponse // 基本响应数据
	Count               int           // 微信首页联系人数量
	ContactList         []*User       // 微信首页联系人信息列表
	SyncKey             *SyncKey      // 微信消息同步key
	User                *User         // 登陆者的用户信息
	ChatSet             string        // 首页联系人的uin，通过‘,’分隔
	Skey                string
	ClientVersion       int
	SystemTime          int64 // unix时间，单位为秒
	GrayScale           int
	InviteStartCount    int
	MPSubscribeMsgCount int               // 公众号推送消息数量
	MPSubscribeMsgList  []*MPSubscribeMsg // 公众号推送消息信息列表
	ClickReportInterval int
}

// 基本响应数据
type BaseResponse struct {
	Ret    int    // 响应码
	ErrMsg string // 错误信息
}

// 公众号的订阅信息
type MPSubscribeMsg struct {
	MPArticleCount int    // 公众号文章数量
	Time           int64  // 推送unix时间(秒)
	UserName       string // 公众号username(@xxxx)
	NickName       string // 公众号名
	MPArticleList  []struct {
		Title  string // 文章标题
		Digest string // 文章摘要
		Cover  string // 文章封面链接
		Url    string // 文章链接
	}
}

// 同步检查时微信的响应结构
type SyncCheckResponse struct {
	RetCode  string // 表示请求是否成功，映射参考mapSyncCheckRetDesc
	Selector string // 0 表示正常无消息，其余都表示有新消息
}

// 请求成功且正常
func (s *SyncCheckResponse) IsSuccess() bool {
	return s.RetCode == "0"
}

// 正常响应，无新消息
func (s *SyncCheckResponse) IsNormal() bool {
	return s.IsSuccess() && s.Selector == "0"
}

var count int      // 用于计算同秒内 selector=7 的次数，判断是否掉线
var mark int       // 记录上一次出现 selector=7 的秒数
const retryMax = 4 // 最大重试次数
// 特别的响应校验, selector = 7，此时可能存在掉线情况，一定时间后微信端会无限返送该响应
func (s *SyncCheckResponse) checkSpecial() error {
	if !(s.IsSuccess() && s.Selector == common.Selector_ENTER_LEAVE_CHAT) {
		return nil
	}
	now := time.Now().Second()
	if mark != now {
		mark = now
		return nil
	}
	count++
	// 判断是否掉线
	if count < retryMax {
		return nil
	}
	count = 0
	return fmt.Errorf("已掉线, selector 出现次数频繁，超过重试次数 %d", retryMax)
}

// 实现Error接口，返回retcode响应码对应的错误信息
func (s *SyncCheckResponse) Error() string {
	c, err := strconv.Atoi(s.RetCode)
	if err != nil {
		return fmt.Sprintf("strconv.Atoi(s.RetCode)执行失败, retCode=%v, 报错信息err=%v", s.RetCode, err)
	}
	return common.GetRetDesc(c)
}

// SyncMsgResp 拉取新消息的响应结构体
type SyncMsgResp struct {
	BaseResponse           *BaseResponse // 基本响应
	AddMsgCount            int           // 新增消息数量
	AddMsgList             []*Message    // 新增消息列表
	ModContactCount        int           // 修改联系人的数量
	ModContactList         []*User       // 修改联系人列表
	DelContactCount        int           // 删除联系人的数量
	DelContactList         []*User       // 删除联系人列表
	ModChatRoomMemberCount int           // 群成员变动的数量
	ModChatRoomMemberList  []*User       // 群成员变动列表
	ContinueFlag           int
	SyncKey                *SyncKey // 新一轮消息更新需要使用的SyncKey(实际都用这个也没问题)
	SKey                   string
	SyncCheckKey           *SyncKey // 新一轮同步检测需要使用的SyncKey
}

type SendMessageResp struct {
	BaseResponse *BaseResponse // 基本响应
	MsgID        string        // 服务端返回的消息ID，可用于撤回接口
	LocalId      string        // 本地消息ID，是我们自己请求时的参数
}
