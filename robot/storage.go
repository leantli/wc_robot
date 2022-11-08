package robot

// storage.go 存储 robot、caller、messagehandlechain、message 共同需要的一些数据

// RobotStorage robot客户端提供的内存存储结构
type RobotStorage struct {
	RequiredParams *RequiredParams  // 微信交互操作时常用的公参
	UUID           string           // 每次登陆都会变更的uuid
	SyncKey        *SyncKey         // 微信消息同步所使用的key
	Self           *User            // 该Storage对应robot登陆用户的信息
	MemberMap      map[string]*User // 用户的好友列表，key为username, value为对应的user，便于查找
}

// RequiredParams 微信交互操作时常用的公参，在扫码成功后获取
type RequiredParams struct {
	SKey       string // 公参
	WxSid      string // 公参
	WxUin      int64  // 公参
	PassTicket string // 公参
	DeviceID   string // 一个robot有一个deviceID，并且不能时常变动
}

// SyncKey 微信消息同步使用的key
type SyncKey struct {
	Count int
	List  []*struct{ Key, Val int64 }
}

// Storage 全局存储 robot、caller、messagehandlechain、message 共同需要的一些数据
var Storage RobotStorage

// 根据传入的remarkname检索用户通讯录
func (s *RobotStorage) SearchMembersByRemarkName(limit int, remarkname string) []*User {
	return s.SearchMembers(limit, func(user *User) bool { return user.RemarkName == remarkname })
}

// 根据传入的nickname检索用户通讯录
func (s *RobotStorage) SearchMembersByNickName(limit int, nickname string) []*User {
	return s.SearchMembers(limit, func(user *User) bool { return user.NickName == nickname })
}

// 检索用户通讯录，通过传入func(user *User) bool进行判断，可参考
func (s *RobotStorage) SearchMembers(limit int, condFns ...func(user *User) bool) []*User {
	result := make([]*User, 0, limit)
	for _, m := range s.MemberMap {
		if len(result) >= limit {
			break
		}
		var passed int
		for _, condFn := range condFns {
			if !condFn(m) {
				break
			}
			passed++
		}
		if passed == len(condFns) {
			result = append(result, m)
		}
	}
	return result
}
