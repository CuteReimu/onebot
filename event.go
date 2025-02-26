package onebot

func init() {
	builder["message"] = map[string]func() any{
		"private": func() any { return &PrivateMessage{} },
		"group":   func() any { return &GroupMessage{} },
	}
	builder["request"] = map[string]func() any{
		"friend": func() any { return &FriendRequest{} },
		"group":  func() any { return &GroupRequest{} },
	}
	builder["notice"] = map[string]func() any{
		"group_decrease": func() any { return &GroupDecreaseNotice{} },
		"group_increase": func() any { return &GroupIncreaseNotice{} },
	}
}

type PrivateMessageSubType string

const (
	PrivateMessageFriend PrivateMessageSubType = "friend" // 好友私聊
	PrivateMessageGroup  PrivateMessageSubType = "group"  // 群私聊
	PrivateMessageOther  PrivateMessageSubType = "other"  // 其它
)

// PrivateMessage 私聊消息
type PrivateMessage struct {
	Time        int64                 `json:"time"`         // 事件发生的时间戳
	SelfId      int64                 `json:"self_id"`      // 收到事件的机器人 QQ 号
	PostType    string                `json:"post_type"`    // "message"
	MessageType string                `json:"message_type"` // "private"
	SubType     PrivateMessageSubType `json:"sub_type"`     // 消息子类型
	MessageId   int32                 `json:"message_id"`   // 消息 ID
	UserId      int64                 `json:"user_id"`      // 发送者 QQ 号
	Message     MessageChain          `json:"message"`      // 消息内容
	RawMessage  string                `json:"raw_message"`  // 原始消息内容
	Font        int32                 `json:"font"`         // 字体
	Sender      Profile               `json:"sender"`       // 发送人信息
}

// ListenPrivateMessage 监听私聊消息
func (b *Bot) ListenPrivateMessage(l func(message *PrivateMessage) bool) {
	listen(b, "message", "private", l)
}

type GroupMessageSubType string

const (
	GroupMessageNormal    GroupMessageSubType = "normal"    // 正常消息
	GroupMessageAnonymous GroupMessageSubType = "anonymous" // 匿名消息
	GroupMessageNotice    GroupMessageSubType = "notice"    // 系统提示
)

// GroupMessage 群消息
type GroupMessage struct {
	Time        int64               `json:"time"`         // 事件发生的时间戳
	SelfId      int64               `json:"self_id"`      // 收到事件的机器人 QQ 号
	PostType    string              `json:"post_type"`    // "message"
	MessageType string              `json:"message_type"` // "group"
	SubType     GroupMessageSubType `json:"sub_type"`     // 消息子类型
	MessageId   int32               `json:"message_id"`   // 消息 ID
	GroupId     int64               `json:"group_id"`     // 群号
	UserId      int64               `json:"user_id"`      // 发送者 QQ 号
	Anonymous   *AnonymousMember    `json:"anonymous"`    // 匿名信息，如果不是匿名消息则为 nil
	Message     MessageChain        `json:"message"`      // 消息内容
	RawMessage  string              `json:"raw_message"`  // 原始消息内容
	Font        int32               `json:"font"`         // 字体
	Sender      Member              `json:"sender"`       // 发送人信息
}

// ListenGroupMessage 监听群消息
func (b *Bot) ListenGroupMessage(l func(message *GroupMessage) bool) {
	listen(b, "message", "group", l)
}

// FriendRequest 加好友请求
type FriendRequest struct {
	Time        int64  `json:"time"`         // 事件发生的时间戳
	SelfId      int64  `json:"self_id"`      // 收到事件的机器人 QQ 号
	PostType    string `json:"post_type"`    // "request"
	RequestType string `json:"request_type"` // "friend"
	UserId      int64  `json:"user_id"`      // 发送请求的 QQ 号
	Comment     string `json:"comment"`      // 验证信息
	Flag        string `json:"flag"`         // 请求 flag，在调用处理请求的 API 时需要传入
}

// ListenFriendRequest 监听加好友请求
func (b *Bot) ListenFriendRequest(l func(request *FriendRequest) bool) {
	listen(b, "request", "friend", l)
}

type GroupRequestSubType string

const (
	GroupRequestAdd    GroupRequestSubType = "add"    // 加群请求
	GroupRequestInvite GroupRequestSubType = "invite" // 邀请入群
)

// GroupRequest 加群请求／邀请
type GroupRequest struct {
	Time        int64               `json:"time"`         // 事件发生的时间戳
	SelfId      int64               `json:"self_id"`      // 收到事件的机器人 QQ 号
	PostType    string              `json:"post_type"`    // "request"
	RequestType string              `json:"request_type"` // "group"
	SubType     GroupRequestSubType `json:"sub_type"`     // 请求子类型
	GroupId     int64               `json:"group_id"`     // 群号
	UserId      int64               `json:"user_id"`      // 发送请求的 QQ 号
	Comment     string              `json:"comment"`      // 验证信息
	Flag        string              `json:"flag"`         // 请求 flag，在调用处理请求的 API 时需要传入
}

// ListenGroupRequest 监听加群请求 / 邀请
func (b *Bot) ListenGroupRequest(l func(request *GroupRequest) bool) {
	listen(b, "request", "group", l)
}

type GroupDecreaseNoticeSubType string

const (
	GroupDecreaseNoticeLeave  GroupDecreaseNoticeSubType = "leave"   // 主动退群
	GroupDecreaseNoticeKick   GroupDecreaseNoticeSubType = "kick"    // 成员被踢
	GroupDecreaseNoticeKickMe GroupDecreaseNoticeSubType = "kick_me" // 机器人被踢
)

// GroupDecreaseNotice 群成员减少
type GroupDecreaseNotice struct {
	Time       int64                      `json:"time"`        // 事件发生的时间戳
	SelfId     int64                      `json:"self_id"`     // 收到事件的机器人 QQ 号
	PostType   string                     `json:"post_type"`   // "notice"
	NoticeType string                     `json:"notice_type"` // "group_decrease"
	SubType    GroupDecreaseNoticeSubType `json:"sub_type"`    // 事件子类型
	GroupId    int64                      `json:"group_id"`    // 群号
	OperatorId int64                      `json:"operator_id"` // 操作者 QQ 号（如果是主动退群，则和 user_id 相同）
	UserId     int64                      `json:"user_id"`     // 离开者 QQ 号
}

// ListenGroupDecreaseNotice 监听群成员减少
func (b *Bot) ListenGroupDecreaseNotice(l func(notice *GroupDecreaseNotice) bool) {
	listen(b, "notice", "group_decrease", l)
}

type GroupIncreaseNoticeSubType string

const (
	GroupIncreaseNoticeApprove GroupIncreaseNoticeSubType = "approve" // 管理员已同意入群
	GroupIncreaseNoticeInvite  GroupIncreaseNoticeSubType = "invite"  // 管理员已邀请入群
)

// GroupIncreaseNotice 群成员增加
type GroupIncreaseNotice struct {
	Time       int64                      `json:"time"`        // 事件发生的时间戳
	SelfId     int64                      `json:"self_id"`     // 收到事件的机器人 QQ 号
	PostType   string                     `json:"post_type"`   // "notice"
	NoticeType string                     `json:"notice_type"` // "group_increase"
	SubType    GroupIncreaseNoticeSubType `json:"sub_type"`    // 事件子类型
	GroupId    int64                      `json:"group_id"`    // 群号
	OperatorId int64                      `json:"operator_id"` // 操作者 QQ 号（即管理员 QQ 号）
	UserId     int64                      `json:"user_id"`     // 加入者 QQ 号
}

// ListenGroupIncreaseNotice 监听群成员增加
func (b *Bot) ListenGroupIncreaseNotice(l func(notice *GroupIncreaseNotice) bool) {
	listen(b, "notice", "group_increase", l)
}
