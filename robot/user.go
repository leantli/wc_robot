package robot

import "wc_robot/common"

// User 抽象的用户结构，用户登陆信息、好友、群组 公众号皆可复用
type User struct {
	// 登陆者用户信息
	Uin               int64  // 用户uin
	UserName          string // 用户username，每次登陆随机分配，个人以@开头，群聊以@@开头
	NickName          string // 微信昵称
	HeadImgUrl        string // 头像图片链接
	RemarkName        string // 备注
	PYInitial         string // 昵称拼音首字母
	PYQuanPin         string // 昵称拼音全拼
	RemarkPYInitial   string // 备注拼音首字母
	RemarkPYQuanPin   string // 备注拼音全拼
	HideInputBarFlag  int    // 是否缩起微信输入框
	StarFriend        int    // 标星朋友数量
	Sex               int    // 性别 1-男 2-女
	Signature         string // 微信个性签名
	AppAccountFlag    int
	VerifyFlag        int // 是否为公众号，0为正常联系人，其余都为公众号(个人公众号/服务号:8；企业服务号24；其余存在部分特殊id)
	ContactFlag       int // 是否为联系人
	WebWxPluginSwitch int // 网页版微信插件开关
	HeadImgFlag       int
	SnsFlag           int
	Province          string // 省份
	City              string // 城市

	// 好友、群组、公众号等额外需要的字段
	IsOwner         int
	MemberCount     int // 群人数
	ChatRoomId      int // 群组id
	UniFriend       int // 共同好友
	OwnerUin        int
	Statues         int
	AttrStatus      int64
	Alias           string
	DisplayName     string // 群成员备注名称
	KeyWord         string
	EncryChatRoomId string
	MemberList      []*User // 群成员
}

// 发送text消息给指定用户(群组), 返回发送消息的MsgID, err
func (u *User) SendTextToUser(toUser *User, content string) (string, error) {
	nsm := NewSendMessage(common.MT_TEXT, content, u.UserName, toUser.UserName, "")
	smr, err := Caller.SendMsg(Storage.RequiredParams, nsm)
	return smr.MsgID, err
}
