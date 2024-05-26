package onebot

import (
	"encoding/json"
	"fmt"
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
		return 0, fmt.Errorf("invalide message type: %s", messageType)
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
		MessageId int64 `json:"messageId"`
	}{messageId})
	return err
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
		NoCache bool  `json:"no_cache"`
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

type GroupMemberInfo struct {
	GroupId         int64  `json:"group_id"`          // 群号
	UserId          int64  `json:"user_id"`           // QQ 号
	Nickname        string `json:"nickname"`          // 昵称
	Card            string `json:"card"`              // 群名片／备注
	Sex             Sex    `json:"sex"`               // 性别
	Age             int32  `json:"age"`               // 年龄
	Area            string `json:"area"`              // 地区
	JoinTime        int32  `json:"join_time"`         // 加群时间戳
	LastSent        int32  `json:"last_sent"`         // 最后发言时间戳
	Level           string `json:"level"`             // 成员等级
	Role            Role   `json:"role"`              // 角色
	Unfriendly      bool   `json:"unfriendly"`        // 是否不良记录成员
	Title           string `json:"title"`             // 专属头衔
	TitleExpireTime int32  `json:"title_expire_time"` // 专属头衔过期时间戳
	CardChangeable  bool   `json:"card_changeable"`   // 是否允许修改群名片
}

// GetGroupMemberInfo 获取群成员信息，group-群号，qq-QQ号，noCache-是否不使用缓存
func (b *Bot) GetGroupMemberInfo(group, qq int64, noCache bool) (*GroupMemberInfo, error) {
	result, err := b.request("get_group_member_info", &struct {
		GroupId int64 `json:"group_id"`
		UserId  int64 `json:"user_id"`
		NoCache bool  `json:"no_cache"`
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
