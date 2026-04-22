package main

import (
	"log/slog"

	"github.com/CuteReimu/onebot"
	"golang.org/x/time/rate"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	b, err := onebot.Connect("localhost", 8080, onebot.WsChannelAll, "", 123456789, false)
	if err != nil {
		panic(err)
	}
	defer func() { _ = b.Close() }()
	// 设置限流策略为：令牌桶容量为10，每秒放入一个令牌，超过的消息直接丢弃
	b.SetLimiter("drop", rate.NewLimiter(1, 10))
	b.ListenGroupMessage(func(message *onebot.GroupMessage) bool {
		var ret onebot.MessageChain
		ret = append(ret, &onebot.Text{Text: "你说了：\n"})
		ret = append(ret, message.Message...)
		_, err := b.SendGroupMessage(message.GroupId, ret)
		if err != nil {
			slog.Error("发送失败", "error", err)
		}
		return true
	})
	b.ListenPrivateMessage(func(message *onebot.PrivateMessage) bool {
		var ret onebot.MessageChain
		ret = append(ret, &onebot.Text{Text: "你说了：\n"})
		ret = append(ret, message.Message...)
		err := message.Reply(b, ret)
		if err != nil {
			slog.Error("发送失败", "error", err)
		}
		return true
	})
	select {}
}
