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

func (m *PrivateMessage) simplify() any {
	m2 := *m
	m2.Message = nil
	m2.RawMessage = ""
	return &m2
}

// Reply 回复
func (m *PrivateMessage) Reply(b *Bot, reply MessageChain) error {
	return b.quickOperation(m, &struct {
		Reply MessageChain `json:"reply"`
	}{reply})
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

func (m *GroupMessage) simplify() any {
	m2 := *m
	m2.Message = nil
	m2.RawMessage = ""
	return &m2
}

// Reply 回复
func (m *GroupMessage) Reply(b *Bot, reply MessageChain, atSender bool) error {
	return b.quickOperation(m, &struct {
		Reply    MessageChain `json:"reply"`
		AtSender bool         `json:"at_sender"`
	}{reply, atSender})
}

// Delete 撤回（需要权限），发送者是匿名用户时无效
func (m *GroupMessage) Delete(b *Bot) error {
	return b.quickOperation(m, &struct {
		Delete bool `json:"delete"`
	}{true})
}

// Kick 把发送者踢出群组（需要权限），不拒绝此人后续加群请求，发送者是匿名用户时无效
func (m *GroupMessage) Kick(b *Bot) error {
	return b.quickOperation(m, &struct {
		Kick bool `json:"kick"`
	}{true})
}

func (m *GroupMessage) Ban(b *Bot, duration int32) error {
	return b.quickOperation(m, &struct {
		Ban         bool  `json:"ban"`
		BanDuration int32 `json:"ban_duration"`
	}{true, duration})
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

// Reply 响应加好友请求，approve是是否同意，remark是添加后的好友备注（仅在同意时有效）
func (r *FriendRequest) Reply(b *Bot, approve bool, remark string) error {
	return b.quickOperation(r, &struct {
		Approve bool   `json:"approve"`
		Remark  string `json:"remark,omitempty"`
	}{approve, remark})
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

// Reply 响应加群请求／邀请，approve是是否同意，reason是拒绝理由（仅在拒绝时有效）
func (r *GroupRequest) Reply(b *Bot, approve bool, reason string) error {
	return b.quickOperation(r, &struct {
		Approve bool   `json:"approve"`
		Reason  string `json:"reason,omitempty"`
	}{approve, reason})
}

// ListenGroupRequest 监听加群请求 / 邀请
func (b *Bot) ListenGroupRequest(l func(request *GroupRequest) bool) {
	listen(b, "request", "group", l)
}
