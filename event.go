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
	builder["meta_event"] = map[string]func() any{
		"lifecycle": func() any { return &LifecycleMetaEvent{} },
		"heartbeat": func() any { return &HeartbeatMetaEvent{} },
	}
	builder["notice"] = map[string]func() any{
		"group_upload":   func() any { return &GroupUploadNotice{} },
		"group_admin":    func() any { return &GroupAdminNotice{} },
		"group_decrease": func() any { return &GroupDecreaseNotice{} },
		"group_increase": func() any { return &GroupIncreaseNotice{} },
		"group_ban":      func() any { return &GroupBanNotice{} },
		"friend_add":     func() any { return &FriendAddNotice{} },
		"group_recall":   func() any { return &GroupRecallNotice{} },
		"friend_recall":  func() any { return &FriendRecallNotice{} },
		"notify":         func() any { return &NotifyNotice{} },
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

// Ban 把发送者禁言，对匿名用户也有效
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

type LifecycleMetaEventSubType string

const (
	LifecycleMetaEventEnable  LifecycleMetaEventSubType = "enable"  // OneBot启用
	LifecycleMetaEventDisable LifecycleMetaEventSubType = "disable" // OneBot停用
	LifecycleMetaEventConnect LifecycleMetaEventSubType = "connect" // WebSocket连接成功
)

// LifecycleMetaEvent 生命周期事件
type LifecycleMetaEvent struct {
	Time          int64                     `json:"time"`            // 事件发生的时间戳
	SelfId        int64                     `json:"self_id"`         // 收到事件的机器人 QQ 号
	PostType      string                    `json:"post_type"`       // "meta_event"
	MetaEventType string                    `json:"meta_event_type"` // "lifecycle"
	SubType       LifecycleMetaEventSubType `json:"sub_type"`        // 请求子类型
}

// ListenLifecycleMetaEvent 监听生命周期
func (b *Bot) ListenLifecycleMetaEvent(l func(notice *LifecycleMetaEvent) bool) {
	listen(b, "meta_event", "lifecycle", l)
}

// HeartbeatMetaEvent 心跳事件
type HeartbeatMetaEvent struct {
	Time          int64     `json:"time"`            // 事件发生的时间戳
	SelfId        int64     `json:"self_id"`         // 收到事件的机器人 QQ 号
	PostType      string    `json:"post_type"`       // "meta_event"
	MetaEventType string    `json:"meta_event_type"` // "heartbeat"
	Status        BotStatus `json:"status"`          // 状态信息
	Interval      int64     `json:"interval"`        // 到下次心跳的间隔，单位毫秒
}

// ListenHeartbeatMetaEvent 监听心跳事件
func (b *Bot) ListenHeartbeatMetaEvent(l func(notice *HeartbeatMetaEvent) bool) {
	listen(b, "meta_event", "heartbeat", l)
}

type File struct {
	Id    string `json:"id"`    // 文件 ID
	Name  string `json:"name"`  // 文件名
	Size  int64  `json:"size"`  // 文件大小（字节数）
	Busid int64  `json:"busid"` // busid（目前不清楚有什么作用）
}

// GroupUploadNotice 群文件上传
type GroupUploadNotice struct {
	Time       int64  `json:"time"`        // 事件发生的时间戳
	SelfId     int64  `json:"self_id"`     // 收到事件的机器人 QQ 号
	PostType   string `json:"post_type"`   // "notice"
	NoticeType string `json:"notice_type"` // "group_upload"
	GroupId    int64  `json:"group_id"`    // 群号
	UserId     int64  `json:"user_id"`     // 发送者 QQ 号
	File       File   `json:"file"`        // 文件信息
}

// ListenGroupUploadNotice 监听群文件上传
func (b *Bot) ListenGroupUploadNotice(l func(notice *GroupUploadNotice) bool) {
	listen(b, "notice", "group_upload", l)
}

type GroupAdminNoticeSubType string

const (
	GroupAdminNoticeSet   GroupAdminNoticeSubType = "set"   // 设置管理员
	GroupAdminNoticeUnset GroupAdminNoticeSubType = "unset" // 取消管理员
)

// GroupAdminNotice 群管理员变动
type GroupAdminNotice struct {
	Time       int64                   `json:"time"`        // 事件发生的时间戳
	SelfId     int64                   `json:"self_id"`     // 收到事件的机器人 QQ 号
	PostType   string                  `json:"post_type"`   // "notice"
	NoticeType string                  `json:"notice_type"` // "group_admin"
	SubType    GroupAdminNoticeSubType `json:"sub_type"`    // 事件子类型
	GroupId    int64                   `json:"group_id"`    // 群号
	UserId     int64                   `json:"user_id"`     // 管理员 QQ 号
}

// ListenGroupAdminNotice 监听群管理员变动
func (b *Bot) ListenGroupAdminNotice(l func(notice *GroupAdminNotice) bool) {
	listen(b, "notice", "group_admin", l)
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

type GroupBanNoticeSubType string

const (
	GroupBanNoticeBan     GroupBanNoticeSubType = "ban"      // 禁言
	GroupBanNoticeLiftBan GroupBanNoticeSubType = "lift_ban" // 解除禁言
)

// GroupBanNotice 群禁言
type GroupBanNotice struct {
	Time       int64                 `json:"time"`        // 事件发生的时间戳
	SelfId     int64                 `json:"self_id"`     // 收到事件的机器人 QQ 号
	PostType   string                `json:"post_type"`   // "notice"
	NoticeType string                `json:"notice_type"` // "group_ban"
	SubType    GroupBanNoticeSubType `json:"sub_type"`    // 事件子类型
	GroupId    int64                 `json:"group_id"`    // 群号
	OperatorId int64                 `json:"operator_id"` // 操作者 QQ 号
	UserId     int64                 `json:"user_id"`     // 加入者 QQ 号
	Duration   int64                 `json:"duration"`    // 禁言时长，单位秒
}

// ListenGroupBanNotice 监听群禁言
func (b *Bot) ListenGroupBanNotice(l func(notice *GroupBanNotice) bool) {
	listen(b, "notice", "group_ban", l)
}

// FriendAddNotice 好友添加
type FriendAddNotice struct {
	Time       int64  `json:"time"`        // 事件发生的时间戳
	SelfId     int64  `json:"self_id"`     // 收到事件的机器人 QQ 号
	PostType   string `json:"post_type"`   // "notice"
	NoticeType string `json:"notice_type"` // "friend_add"
	UserId     int64  `json:"user_id"`     // 新添加好友 QQ 号
}

// ListenFriendAddNotice 监听好友添加
func (b *Bot) ListenFriendAddNotice(l func(notice *FriendAddNotice) bool) {
	listen(b, "notice", "friend_add", l)
}

// GroupRecallNotice 群消息撤回
type GroupRecallNotice struct {
	Time       int64  `json:"time"`        // 事件发生的时间戳
	SelfId     int64  `json:"self_id"`     // 收到事件的机器人 QQ 号
	PostType   string `json:"post_type"`   // "notice"
	NoticeType string `json:"notice_type"` // "group_recall"
	GroupId    int64  `json:"group_id"`    // 群号
	UserId     int64  `json:"user_id"`     // 消息发送者 QQ 号
	OperatorId int64  `json:"operator_id"` // 操作者 QQ 号
	MessageId  int64  `json:"message_id"`  // 被撤回的消息 ID
}

// ListenGroupRecallNotice 监听群消息撤回
func (b *Bot) ListenGroupRecallNotice(l func(notice *GroupRecallNotice) bool) {
	listen(b, "notice", "group_recall", l)
}

// FriendRecallNotice 好友消息撤回
type FriendRecallNotice struct {
	Time       int64  `json:"time"`        // 事件发生的时间戳
	SelfId     int64  `json:"self_id"`     // 收到事件的机器人 QQ 号
	PostType   string `json:"post_type"`   // "notice"
	NoticeType string `json:"notice_type"` // "friend_recall"
	UserId     int64  `json:"user_id"`     // 好友 QQ 号
	MessageId  int64  `json:"message_id"`  // 被撤回的消息 ID
}

// ListenFriendRecallNotice 监听好友消息撤回
func (b *Bot) ListenFriendRecallNotice(l func(notice *FriendRecallNotice) bool) {
	listen(b, "notice", "friend_recall", l)
}

type NotifyNoticeSubType string

const (
	NotifyNoticePoke      NotifyNoticeSubType = "poke"       // 群内戳一戳
	NotifyNoticeLuckyKing NotifyNoticeSubType = "lucky_king" // 群红包运气王
	NotifyNoticeHonor     NotifyNoticeSubType = "honor"      // 群成员荣誉变更
)

type NotifyNoticeHonorType string

const (
	NotifyNoticeNone      NotifyNoticeSubType = ""
	NotifyNoticeTalkative NotifyNoticeSubType = "talkative" // 龙王
	NotifyNoticePerformer NotifyNoticeSubType = "performer" // 群聊之火
	NotifyNoticeEmotion   NotifyNoticeSubType = "emotion"   // 快乐源泉
)

// NotifyNotice 其它通知
type NotifyNotice struct {
	Time       int64                 `json:"time"`                 // 事件发生的时间戳
	SelfId     int64                 `json:"self_id"`              // 收到事件的机器人 QQ 号
	PostType   string                `json:"post_type"`            // "notice"
	NoticeType string                `json:"notice_type"`          // "notify"
	SubType    NotifyNoticeSubType   `json:"sub_type"`             // 事件子类型
	GroupId    int64                 `json:"group_id"`             // 群号
	HonorType  NotifyNoticeHonorType `json:"honor_type,omitempty"` // （群成员荣誉变更）荣誉类型
	UserId     int64                 `json:"user_id"`              // （戳一戳）发送者/（红包）发送者/（荣誉变更）成员 QQ 号
	TargetId   int64                 `json:"target_id,omitempty"`  // （戳一戳）被戳者/（红包）运气王 QQ 号
}

// ListenNotifyNotice 监听其它通知
func (b *Bot) ListenNotifyNotice(l func(notice *NotifyNotice) bool) {
	listen(b, "notice", "notify", l)
}
