package onebot

import (
	"encoding/json"
	"testing"
)

func TestMessage(t *testing.T) {
	type PrivateMessage struct {
		UserId  int64        `json:"user_id"`
		Message MessageChain `json:"message"`
	}
	privateMessage := &PrivateMessage{1000, MessageChain{&Text{Text: "123"}}}
	buf, err := json.Marshal(privateMessage)
	if err != nil {
		t.Fatal(err)
	}
	if string(buf) != `{"user_id":1000,"message":[{"type":"text","data":{"text":"123"}}]}` {
		t.Fatal(string(buf))
	}
	privateMessage = nil
	err = json.Unmarshal(buf, &privateMessage)
	if err != nil {
		t.Fatal(err)
	}
	if privateMessage.UserId != 1000 || len(privateMessage.Message) != 1 {
		t.Fatal(privateMessage)
	}
	if text, ok := privateMessage.Message[0].(*Text); !ok || text.Text != "123" {
		t.Fatal(privateMessage)
	}
}
