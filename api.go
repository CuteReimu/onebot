package onebot

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
)

// SendPrivateMessage 发送私聊消息，消息ID
func (b *Bot) SendPrivateMessage(userId int64, message MessageChain) (int64, error) {
	result, err := b.request("send_private_msg", &struct {
		UserId  int64        `json:"user_id"`
		Message MessageChain `json:"message"`
	}{userId, message})
	if err != nil {
		return 0, err
	}
	return result.Get("message_id").Int(), nil
}

// SendGroupMessage 发送群消息，group-群号，message-发送的内容，返回消息id
func (b *Bot) SendGroupMessage(group int64, message MessageChain) (int64, error) {
	result, err := b.request("send_group_msg", &struct {
		GroupId int64        `json:"group_id"`
		Message MessageChain `json:"message"`
	}{group, message})
	if err != nil {
		return 0, err
	}
	return result.Get("message_id").Int(), nil
}

type MessageType string

const (
	MessageTypePrivate MessageType = "private" // 私聊消息
	MessageTypeGroup   MessageType = "group"   // 群消息
)

// SendMessage 发送消息，返回消息id
func (b *Bot) SendMessage(messageType MessageType, targetId int64, message MessageChain) (int64, error) {
	m := map[string]any{
		"message_type": string(messageType),
		"message":      message,
	}
	switch messageType {
	case MessageTypePrivate:
		m["user_id"] = targetId
	case MessageTypeGroup:
		m["group_id"] = targetId
	default:
		return 0, fmt.Errorf("invalid message type: %s", messageType)
	}
	result, err := b.request("send_msg", m)
	if err != nil {
		return 0, err
	}
	return result.Get("message_id").Int(), nil
}

// DeleteMessage 撤回消息，messageId-需要撤回的消息的ID
func (b *Bot) DeleteMessage(messageId int64) error {
	_, err := b.request("delete_msg", &struct {
		MessageId int64 `json:"message_id"`
	}{messageId})
	return err
}

type Message struct {
	Time        int32        `json:"time"`         // 发送时间
	MessageType MessageType  `json:"message_type"` // 消息类型
	MessageId   int32        `json:"message_id"`   // 消息ID
	RealId      int32        `json:"real_id"`      // 消息真实ID
	Sender      Profile      `json:"sender"`       // 发送人信息
	Message     MessageChain `json:"message"`      // 消息内容
}

// GetMessage 获取消息
func (b *Bot) GetMessage(messageId int32) (*Message, error) {
	result, err := b.request("get_msg", &struct {
		MessageId int32 `json:"message_id"`
	}{messageId})
	if err != nil {
		return nil, err
	}
	var msg *Message
	err = json.Unmarshal([]byte(result.Raw), &msg)
	return msg, err
}

// GetForwardMessage 获取合并转发消息
func (b *Bot) GetForwardMessage(id string) ([]any, error) {
	result, err := b.request("get_forward_msg", &struct {
		Id string `json:"id"`
	}{id})
	if err != nil {
		return nil, err
	}
	var ret []any
	for _, msg := range result.Array() {
		postType := msg.Get("post_type").String()
		if postType == "message_sent" {
			postType = "message"
		}
		subType := msg.Get(postType + "_type").String()
		if bd := builder[postType][subType]; bd == nil {
			slog.Error("cannot find message builder: " + postType)
			return nil, errors.New("cannot find message builder: " + postType)
		} else {
			m := bd()
			err = json.Unmarshal([]byte(result.Raw), m)
			if err != nil {
				slog.Error("json unmarshal failed", "error", err)
				return nil, err
			}
			ret = append(ret, m)
		}
	}
	return ret, nil
}

// SendLike 发送好友赞，userId-对方QQ号，times-赞的次数（每个好友每天最多10次）
func (b *Bot) SendLike(userId int64, times int32) error {
	_, err := b.request("send_like", &struct {
		UserId int64 `json:"user_id"`
		Times  int32 `json:"times,omitempty"`
	}{userId, times})
	return err
}

// SetGroupKick 群组踢人，rejectAddRequest-拒绝此人的加群请求
func (b *Bot) SetGroupKick(groupId, userId int64, rejectAddRequest bool) error {
	_, err := b.request("set_group_kick", &struct {
		GroupId          int64 `json:"group_id"`
		UserId           int64 `json:"user_id"`
		RejectAddRequest bool  `json:"reject_add_request,omitempty"`
	}{groupId, userId, rejectAddRequest})
	return err
}

// SetGroupBan 群组单人禁言，duration-禁言时长，单位秒，0表示取消禁言
func (b *Bot) SetGroupBan(groupId, userId int64, duration int32) error {
	_, err := b.request("set_group_ban", &struct {
		GroupId  int64 `json:"group_id"`
		UserId   int64 `json:"user_id"`
		Duration int32 `json:"duration"`
	}{groupId, userId, duration})
	return err
}

// SetGroupAnonymousBan 群组匿名用户禁言，flag-需从群消息上报的 AnonymousMember 中获得，duration-禁言时长，单位秒，无法取消匿名用户禁言
func (b *Bot) SetGroupAnonymousBan(groupId int64, flag string, duration int32) error {
	_, err := b.request("set_group_anonymous_ban", &struct {
		GroupId       int64  `json:"group_id"`
		AnonymousFlag string `json:"anonymous_flag"`
		Duration      int32  `json:"duration"`
	}{groupId, flag, duration})
	return err
}

// SetGroupWholeBan 群组全员禁言
func (b *Bot) SetGroupWholeBan(groupId int64, enable bool) error {
	_, err := b.request("set_group_whole_ban", &struct {
		GroupId int64 `json:"group_id"`
		Enable  bool  `json:"enable"`
	}{groupId, enable})
	return err
}

// SetGroupAdmin 群组设置管理员
func (b *Bot) SetGroupAdmin(groupId, userId int64, enable bool) error {
	_, err := b.request("set_group_admin", &struct {
		GroupId int64 `json:"group_id"`
		UserId  int64 `json:"user_id"`
		Enable  bool  `json:"enable"`
	}{groupId, userId, enable})
	return err
}

// SetGroupAnonymous 群组匿名
func (b *Bot) SetGroupAnonymous(groupId int64, enable bool) error {
	_, err := b.request("set_group_anonymous", &struct {
		GroupId int64 `json:"group_id"`
		Enable  bool  `json:"enable"`
	}{groupId, enable})
	return err
}

// SetGroupCard 设置群名片（群备注），card-群名片内容，不填或空字符串表示删除群名片
func (b *Bot) SetGroupCard(groupId, userId int64, card string) error {
	_, err := b.request("set_group_card", &struct {
		GroupId int64  `json:"group_id"`
		UserId  int64  `json:"user_id"`
		Card    string `json:"card,omitempty"`
	}{groupId, userId, card})
	return err
}

// SetGroupName 设置群名
func (b *Bot) SetGroupName(groupId int64, groupName string) error {
	_, err := b.request("set_group_leave", &struct {
		GroupId   int64  `json:"group_id"`
		GroupName string `json:"group_name"`
	}{groupId, groupName})
	return err
}

// SetGroupLeave 退出群，group-群号，isDismiss-是否解散
func (b *Bot) SetGroupLeave(group int64, isDismiss bool) error {
	_, err := b.request("set_group_leave", &struct {
		GroupId   int64 `json:"group_id"`
		IsDismiss bool  `json:"is_dismiss,omitempty"`
	}{group, isDismiss})
	return err
}

// SetGroupSpecialTitle 设置群组专属头衔，specialTitle-专属头衔，不填或空字符串表示删除专属头衔，duration-专属头衔有效期，单位秒，-1表示永久，不过此项似乎没有效果，可能是只有某些特殊的时间长度有效，有待测试
func (b *Bot) SetGroupSpecialTitle(groupId, userId int64, specialTitle string, duration int32) error {
	_, err := b.request("set_group_special_title", &struct {
		GroupId      int64  `json:"group_id"`
		UserId       int64  `json:"user_id"`
		SpecialTitle string `json:"specialTitle,omitempty"`
		Duration     int32  `json:"duration"`
	}{groupId, userId, specialTitle, duration})
	return err
}

// SetFriendAddRequest 处理加好友请求，flag-从请求中获取，approve-是否同意，remark-同意好友后的备注
func (b *Bot) SetFriendAddRequest(flag string, approve bool, remark string) error {
	_, err := b.request("set_friend_add_request", &struct {
		Flag    string `json:"flag"`
		Approve bool   `json:"approve"`
		Remark  string `json:"remark,omitempty"`
	}{flag, approve, remark})
	return err
}

// SetGroupAddRequest 处理加群请求／邀请，flag-从请求中获取，subType-从请求中获取，approve-是否同意，reason-拒绝理由
func (b *Bot) SetGroupAddRequest(flag string, subType GroupRequestSubType, approve bool, reason string) error {
	_, err := b.request("set_group_add_request", &struct {
		Flag    string              `json:"flag"`
		SubType GroupRequestSubType `json:"subType"`
		Approve bool                `json:"approve"`
		Reason  string              `json:"reason,omitempty"`
	}{flag, subType, approve, reason})
	return err
}

type LoginInfo struct {
	UserId   int64  `json:"user_id"`  // QQ号
	Nickname string `json:"nickname"` // QQ昵称
}

// GetLoginInfo 获取登录号信息
func (b *Bot) GetLoginInfo() (*LoginInfo, error) {
	result, err := b.request("get_login_info", nil)
	if err != nil {
		return nil, err
	}
	var loginInfo *LoginInfo
	err = json.Unmarshal([]byte(result.Raw), &loginInfo)
	return loginInfo, err
}

// GetStrangerInfo 获取陌生人信息，userId-QQ号，noCache-是否不使用缓存
func (b *Bot) GetStrangerInfo(userId int64, noCache bool) (*Profile, error) {
	result, err := b.request("get_stranger_info", &struct {
		UserId  int64 `json:"user_id"`
		NoCache bool  `json:"no_cache,omitempty"`
	}{userId, noCache})
	if err != nil {
		return nil, err
	}
	var strangerInfo *Profile
	err = json.Unmarshal([]byte(result.Raw), &strangerInfo)
	return strangerInfo, err
}

// GetFriendList 获取陌生人信息，获取好友列表
func (b *Bot) GetFriendList() ([]*Friend, error) {
	result, err := b.request("get_friend_list", nil)
	if err != nil {
		return nil, err
	}
	var friends []*Friend
	err = json.Unmarshal([]byte(result.Raw), &friends)
	return friends, err
}

type GroupInfo struct {
	GroupId        int64  `json:"group_id"`         // 群号
	GroupName      string `json:"group_name"`       // 群名称
	MemberCount    int32  `json:"member_count"`     // 成员数
	MaxMemberCount int32  `json:"max_member_count"` // 最大成员数（群容量）
}

// GetGroupInfo 获取群信息，group-群号，noCache-是否不使用缓存
func (b *Bot) GetGroupInfo(group int64, noCache bool) (*GroupInfo, error) {
	result, err := b.request("get_group_info", &struct {
		GroupId int64 `json:"group_id"`
		NoCache bool  `json:"no_cache,omitempty"`
	}{group, noCache})
	if err != nil {
		return nil, err
	}
	var groupInfo *GroupInfo
	err = json.Unmarshal([]byte(result.Raw), &groupInfo)
	return groupInfo, err
}

// GetGroupList 获取群列表
func (b *Bot) GetGroupList() ([]*GroupInfo, error) {
	result, err := b.request("get_group_list", nil)
	if err != nil {
		return nil, err
	}
	var groupInfo []*GroupInfo
	err = json.Unmarshal([]byte(result.Raw), &groupInfo)
	return groupInfo, err
}

type String string

func (s *String) UnmarshalJSON(bytes []byte) error {
	if len(bytes) > 2 && bytes[0] == '"' && bytes[len(bytes)-1] == '"' {
		return json.Unmarshal(bytes, (*string)(s))
	}
	*s = String(bytes)
	return nil
}

func (s *String) String() string {
	return string(*s)
}

type GroupMemberInfo struct {
	GroupId         int64  `json:"group_id"`          // 群号
	UserId          int64  `json:"user_id"`           // QQ 号
	Nickname        string `json:"nickname"`          // 昵称
	Card            string `json:"card"`              // 群名片／备注
	Sex             Sex    `json:"sex"`               // 性别
	Age             int32  `json:"age"`               // 年龄
	Area            string `json:"area"`              // 地区
	JoinTime        int32  `json:"join_time"`         // 加群时间戳
	LastSentTime    int32  `json:"last_sent_time"`    // 最后发言时间戳
	Level           String `json:"level"`             // 成员等级（onebot11标准这里是字符串，但很多框架把它返回了一个整型数值，这里就用这个方式兼容一下）
	Role            Role   `json:"role"`              // 角色
	Unfriendly      bool   `json:"unfriendly"`        // 是否不良记录成员
	Title           string `json:"title"`             // 专属头衔
	TitleExpireTime int32  `json:"title_expire_time"` // 专属头衔过期时间戳
	CardChangeable  bool   `json:"card_changeable"`   // 是否允许修改群名片
}

func (m *GroupMemberInfo) CardOrNickname() string {
	if m.Card != "" {
		return m.Card
	}
	return m.Nickname
}

// GetGroupMemberInfo 获取群成员信息，group-群号，qq-QQ号，noCache-是否不使用缓存
func (b *Bot) GetGroupMemberInfo(group, qq int64, noCache bool) (*GroupMemberInfo, error) {
	result, err := b.request("get_group_member_info", &struct {
		GroupId int64 `json:"group_id"`
		UserId  int64 `json:"user_id"`
		NoCache bool  `json:"no_cache,omitempty"`
	}{group, qq, noCache})
	if err != nil {
		return nil, err
	}
	var groupMemberInfo *GroupMemberInfo
	err = json.Unmarshal([]byte(result.Raw), &groupMemberInfo)
	return groupMemberInfo, err
}

// GetGroupMemberList 获取群成员列表，group-群号
//
// 注意：获取群成员列表时，不保证能获取到每个成员的所有信息，有些信息（例如area、title等）可能无法获得。
// 想要获取所有信息，请调用 Bot.GetGroupMemberInfo 方法获取单个成员信息。
func (b *Bot) GetGroupMemberList(group int64) ([]*GroupMemberInfo, error) {
	result, err := b.request("get_group_member_list", &struct {
		GroupId int64 `json:"group_id"`
	}{group})
	if err != nil {
		return nil, err
	}
	var groupInfo []*GroupMemberInfo
	err = json.Unmarshal([]byte(result.Raw), &groupInfo)
	return groupInfo, err
}

// GetCookies 获取Cookies，domain-需要获取cookies的域名
func (b *Bot) GetCookies(domain string) (string, error) {
	result, err := b.request("get_cookies", &struct {
		Domain string `json:"domain"`
	}{domain})
	if err != nil {
		return "", err
	}
	var ret struct {
		Cookies string `json:"cookies"`
	}
	err = json.Unmarshal([]byte(result.Raw), &ret)
	return ret.Cookies, err
}

// GetCsrfToken 获取 CSRF Token
func (b *Bot) GetCsrfToken() (int32, error) {
	result, err := b.request("get_csrf_token", nil)
	if err != nil {
		return 0, err
	}
	var ret struct {
		Token int32 `json:"token"`
	}
	err = json.Unmarshal([]byte(result.Raw), &ret)
	return ret.Token, err
}

// GetCredentials 获取QQ相关接口凭证，即上面两个接口的合并
func (b *Bot) GetCredentials(domain string) (string, int32, error) {
	result, err := b.request("get_credentials", &struct {
		Domain string `json:"domain"`
	}{domain})
	if err != nil {
		return "", 0, err
	}
	var ret struct {
		Cookies string `json:"cookies"`
		Token   int32  `json:"token"`
	}
	err = json.Unmarshal([]byte(result.Raw), &ret)
	return ret.Cookies, ret.Token, err
}

// GetRecord 获取语音，file-收到的语音文件名（消息段的file参数），outFormat-要转换到的格式，返回转换后的语音文件路径
//
// 提示：要使用此接口，通常需要安装 ffmpeg，请参考 OneBot 实现的相关说明。
// outFormat目前支持mp3、amr、wma、m4a、spx、ogg、wav、flac
func (b *Bot) GetRecord(file, outFormat string) (string, error) {
	result, err := b.request("get_record", &struct {
		File      string `json:"file"`
		OutFormat string `json:"out_format"`
	}{file, outFormat})
	if err != nil {
		return "", err
	}
	var ret struct {
		File string `json:"file"`
	}
	err = json.Unmarshal([]byte(result.Raw), &ret)
	return ret.File, err
}

// GetImage 获取图片，file-收到的图片文件名（消息段的file参数），返回下载后的图片文件路径
func (b *Bot) GetImage(file string) (string, error) {
	result, err := b.request("get_image", &struct {
		File string `json:"file"`
	}{file})
	if err != nil {
		return "", err
	}
	var ret struct {
		File string `json:"file"`
	}
	err = json.Unmarshal([]byte(result.Raw), &ret)
	return ret.File, err
}

// CanSendImage 检查是否可以发送图片
func (b *Bot) CanSendImage() (bool, error) {
	result, err := b.request("can_send_image", nil)
	if err != nil {
		return false, err
	}
	var ret struct {
		Yes bool `json:"yes"`
	}
	err = json.Unmarshal([]byte(result.Raw), &ret)
	return ret.Yes, err
}

// CanSendRecord 检查是否可以发送语音
func (b *Bot) CanSendRecord() (bool, error) {
	result, err := b.request("can_send_record", nil)
	if err != nil {
		return false, err
	}
	var ret struct {
		Yes bool `json:"yes"`
	}
	err = json.Unmarshal([]byte(result.Raw), &ret)
	return ret.Yes, err
}

type BotStatus struct {
	Online *bool `json:"online"` // 当前QQ在线，true|false|nil
	Good   bool  `json:"good"`   // 状态符合预期，意味着各模块正常运行、功能正常，且 QQ 在线
}

// GetStatus 获取运行状态
func (b *Bot) GetStatus() (*BotStatus, error) {
	result, err := b.request("get_status", nil)
	if err != nil {
		return nil, err
	}
	var ret *BotStatus
	err = json.Unmarshal([]byte(result.Raw), &ret)
	return ret, err
}

type BotVersionInfo struct {
	AppName         string `json:"app_name"`         // 应用标识，如"mirai-native"
	AppVersion      string `json:"app_version"`      // 应用版本，如"1.2.3"
	ProtocolVersion string `json:"protocol_version"` // OneBot标准版本，如"v11"
}

// GetVersionInfo 获取版本信息
func (b *Bot) GetVersionInfo() (*BotVersionInfo, error) {
	result, err := b.request("get_version_info", nil)
	if err != nil {
		return nil, err
	}
	var ret *BotVersionInfo
	err = json.Unmarshal([]byte(result.Raw), &ret)
	return ret, err
}

// SetRestart 重启OneBot实现，delay-要延迟的毫秒数，如果默认情况下无法重启，可以尝试设置延迟为2000左右。
//
// 由于重启 OneBot 实现同时需要重启 API 服务，这意味着当前的 API 请求会被中断，因此需要异步地重启，接口返回的 status 是 async。
func (b *Bot) SetRestart(delay int32) error {
	_, err := b.request("set_restart", &struct {
		Delay int32 `json:"delay,omitempty"`
	}{delay})
	return err
}

// CleanCache 用于清理积攒了太多的缓存文件
func (b *Bot) CleanCache() error {
	_, err := b.request("clean_cache", nil)
	return err
}
