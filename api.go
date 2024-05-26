package onebot

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

// DeleteMessage 撤回消息，messageId-需要撤回的消息的ID
func (b *Bot) DeleteMessage(messageId int64) error {
	_, err := b.request("delete_msg", &struct {
		MessageId int64 `json:"messageId"`
	}{messageId})
	return err
}
