package onebot

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"log/slog"
	"strings"
)

type messageSegment struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

type MessageChain []SingleMessage

func (c *MessageChain) MarshalJSON() ([]byte, error) {
	segments := make([]messageSegment, 0, len(*c))
	for _, m := range *c {
		segments = append(segments, messageSegment{
			Type: m.GetMessageType(),
			Data: m,
		})
	}
	return json.Marshal(segments)
}

func (c *MessageChain) UnmarshalJSON(data []byte) error {
	if !gjson.ValidBytes(data) {
		return errors.New("invalid json data")
	}
	result := gjson.ParseBytes(data)
	if !result.IsArray() {
		return errors.New("result is not array")
	}
	*c = parseMessageChain(result.Array())
	return nil
}

type SingleMessage interface {
	GetMessageType() string
}

// Text 纯文本
type Text struct {
	Text string `json:"text"` // 文字消息
}

func (m *Text) GetMessageType() string {
	return "text"
}

func (m *Text) String() string {
	return m.Text
}

// Face QQ表情
type Face struct {
	Id string `json:"id"` // QQ 表情 ID
}

func (m *Face) GetMessageType() string {
	return "face"
}

func (m *Face) String() string {
	return fmt.Sprintf("[CQ:face,id=%s]", m.Id)
}

// Image 图片
type Image struct {
	// 图片文件名
	//
	// 当发送图片时，除了可以直接使用收到的文件名直接发送以外，还支持：
	//
	// 绝对路径，格式使用file URI；
	// 网络URL，格式为http://xxx；
	// Base64编码，格式为base64://<base64编码>。
	File string `json:"file"`

	// 图片类型，"flash"表示闪照，无此参数表示普通图片
	Type string `json:"type,omitempty"`

	// 图片的URL，只有在收图片时才有这个字段
	Url string `json:"url,omitempty"`

	// 以下字段只在通过网络URL发送时有效
	Cache   string `json:"cache,omitempty"`   // 1和0表示是否使用已缓存的文件，默认是
	Proxy   string `json:"proxy,omitempty"`   // 1和0表示是否通过代理下载文件（需通过环境变量或配置文件配置代理），默认是
	Timeout string `json:"timeout,omitempty"` // 单位秒，表示下载网络文件的超时时间，默认不超时
}

func (m *Image) GetMessageType() string {
	return "image"
}

func (m *Image) String() string {
	return fmt.Sprintf("[CQ:image,file=%s]", m.File)
}

// Record 语音
type Record struct {
	// 语音文件名
	//
	// 当发送语音时，除了可以直接使用收到的文件名直接发送以外，还支持：
	//
	// 绝对路径，格式使用file URI；
	// 网络URL，格式为http://xxx；
	// Base64编码，格式为base64://<base64编码>。
	File string `json:"file"`

	// 发送时可选，是否变声
	Magic bool `json:"magic,omitempty"`

	// 语音的URL，只有在收语音时才有这个字段
	Url string `json:"url,omitempty"`

	// 以下字段只在通过网络URL发送时有效
	Cache   string `json:"cache,omitempty"`   // 1和0表示是否使用已缓存的文件，默认是
	Proxy   string `json:"proxy,omitempty"`   // 1和0表示是否通过代理下载文件（需通过环境变量或配置文件配置代理），默认是
	Timeout string `json:"timeout,omitempty"` // 单位秒，表示下载网络文件的超时时间，默认不超时
}

func (m *Record) GetMessageType() string {
	return "record"
}

func (m *Record) String() string {
	return fmt.Sprintf("[CQ:record,file=%s]", m.File)
}

// Video 短视频
type Video struct {
	// 短视频文件名
	//
	// 当发送短视频时，除了可以直接使用收到的文件名直接发送以外，还支持：
	//
	// 绝对路径，格式使用file URI；
	// 网络URL，格式为http://xxx；
	// Base64编码，格式为base64://<base64编码>。
	File string `json:"file"`

	// 短视频的URL，只有在收短视频时才有这个字段
	Url string `json:"url,omitempty"`

	// 以下字段只在通过网络URL发送时有效
	Cache   string `json:"cache,omitempty"`   // 1和0表示是否使用已缓存的文件，默认是
	Proxy   string `json:"proxy,omitempty"`   // 1和0表示是否通过代理下载文件（需通过环境变量或配置文件配置代理），默认是
	Timeout string `json:"timeout,omitempty"` // 单位秒，表示下载网络文件的超时时间，默认不超时
}

func (m *Video) GetMessageType() string {
	return "video"
}

func (m *Video) String() string {
	return fmt.Sprintf("[CQ:video,file=%s]", m.File)
}

// At @某人
type At struct {
	QQ string `json:"qq"` // 群员QQ号，或者"all"表示全体成员
}

func (m *At) GetMessageType() string {
	return "at"
}

func (m *At) String() string {
	return fmt.Sprintf("[CQ:at,qq=%s]", m.QQ)
}

// RPS 猜拳魔法表情
type RPS struct {
}

func (m *RPS) GetMessageType() string {
	return "rps"
}

func (m *RPS) String() string {
	return "[CQ:rps]"
}

// Dice 掷骰子魔法表情
type Dice struct {
}

func (m *Dice) GetMessageType() string {
	return "dice"
}

func (m *Dice) String() string {
	return "[CQ:dice]"
}

// Shake 窗口抖动
type Shake struct {
}

func (m *Shake) GetMessageType() string {
	return "shake"
}

func (m *Shake) String() string {
	return "[CQ:shake]"
}

// Poke 戳一戳，字段含义参考文档
//
// https://github.com/botuniverse/onebot-11/blob/master/message/segment.md#%E6%88%B3%E4%B8%80%E6%88%B3
type Poke struct {
	Type string `json:"type"`
	Id   string `json:"id"`
	Name string `json:"name,omitempty"` // 表情名，发送时无需此字段
}

func (m *Poke) GetMessageType() string {
	return "poke"
}

func (m *Poke) String() string {
	return fmt.Sprintf("[CQ:poke,type=%s,id=%s]", m.Type, m.Id)
}

// Anonymous 匿名发消息
type Anonymous struct {
	Ignore string `json:"ignore,omitempty"` // 发消息时使用，1和0表示无法匿名时是否继续发送
}

func (m *Anonymous) GetMessageType() string {
	return "anonymous"
}

func (m *Anonymous) String() string {
	return "[CQ:anonymous]"
}

// Share 链接分享
type Share struct {
	Url     string `json:"url"`               // URL
	Title   string `json:"title"`             // 标题
	Content string `json:"content,omitempty"` // 内容描述，发送时可不填
	Image   string `json:"image,omitempty"`   // 图片URL，发送时可不填
}

func (m *Share) GetMessageType() string {
	return "share"
}

func (m *Share) String() string {
	return fmt.Sprintf("[CQ:share,url=%s,title=%s]", m.Url, m.Title)
}

type ContactType string

const (
	ContactTypeQQ    ContactType = "qq"    // 推荐好友
	ContactTypeGroup ContactType = "group" // 推荐群
)

// Contact 推荐好友、推荐群
type Contact struct {
	Type ContactType `json:"type"` // 类型
	Id   string      `json:"id"`   // QQ号或群号
}

func (m *Contact) GetMessageType() string {
	return "contact"
}

func (m *Contact) String() string {
	return fmt.Sprintf("[CQ:contact,type=%s,id=%s]", m.Type, m.Id)
}

// Location 位置
type Location struct {
	Lat     string `json:"lat"`               // 纬度
	Lon     string `json:"lon"`               // 经度
	Title   string `json:"title,omitempty"`   // 标题，发送时可不填
	Content string `json:"content,omitempty"` // 内容描述，发送时可不填
}

func (m *Location) GetMessageType() string {
	return "location"
}

func (m *Location) String() string {
	return fmt.Sprintf("[CQ:location,lat=%s,lon=%s]", m.Lat, m.Lon)
}

// Music 音乐分享
type Music struct {
	// "qq": QQ音乐, "163": 网易云音乐, "xm": 虾米音乐, "custom": 音乐自定义分享
	Type string `json:"type"`

	// 歌曲ID，用于非自定义分享
	Id string `json:"id,omitempty"`

	// 以下字段用于自定义分享

	Url     string `json:"url,omitempty"`     // 点击后跳转目标 URL
	Audio   string `json:"audio,omitempty"`   // 音乐 URL
	Title   string `json:"title,omitempty"`   // 标题
	Content string `json:"content,omitempty"` // 内容描述，发送时可不填
	Image   string `json:"image,omitempty"`   // 图片URL，发送时可不填
}

func (m *Music) GetMessageType() string {
	return "music"
}

func (m *Music) String() string {
	if m.Type == "custom" {
		return fmt.Sprintf("[CQ:music,type=custom,url=%s,audio=%s,title=%s]", m.Url, m.Audio, m.Title)
	}
	return fmt.Sprintf("[CQ:music,type=%s,id=%s]", m.Type, m.Id)
}

// Reply 回复
type Reply struct {
	Id string `json:"id"` // 回复时引用的消息 ID
}

func (m *Reply) GetMessageType() string {
	return "reply"
}

func (m *Reply) String() string {
	return fmt.Sprintf("[CQ:reply,id=%s]", m.Id)
}

// Forward 合并转发
type Forward struct {
	Id string `json:"id"` // 合并转发 ID，需调用 Bot.GetForwardMsg 方法获取具体内容
}

func (m *Forward) GetMessageType() string {
	return "forward"
}

func (m *Forward) String() string {
	return fmt.Sprintf("[CQ:forward,id=%s]", m.Id)
}

// Node 合并转发节点
type Node struct {
	Id string `json:"id,omitempty"` // 转发的消息 ID

	// 也可以不要 ID 字段，改为用以下字段进行自定义节点

	UserId   string       `json:"user_id,omitempty"`  // 发送者 QQ 号
	Nickname string       `json:"nickname,omitempty"` // 发送者昵称
	Content  MessageChain `json:"content,omitempty"`  // 消息内容
}

func (m *Node) GetMessageType() string {
	return "node"
}

func (m *Node) String() string {
	if len(m.Id) == 0 {
		var sb strings.Builder
		for _, m := range m.Content {
			sb.WriteString(fmt.Sprint(m))
		}
		return fmt.Sprintf("[CQ:node,user_id=%s,nickname=%s,content=%s]", m.UserId, m.Nickname, sb.String())
	}
	return fmt.Sprintf("[CQ:node,id=%s]", m.Id)
}

type Xml struct {
	Data string `json:"data"` // XML文本
}

func (m *Xml) GetMessageType() string {
	return "xml"
}

func (m *Xml) String() string {
	return fmt.Sprintf("[CQ:xml,data=%s]", m.Data)
}

type Json struct {
	Data string `json:"data"` // Json文本
}

func (m *Json) GetMessageType() string {
	return "json"
}

func (m *Json) String() string {
	return fmt.Sprintf("[CQ:json,data=%s]", m.Data)
}

var singleMessageBuilder = make(map[string]func() SingleMessage)

func init() {
	initMessageBuilder := func(f func() SingleMessage) {
		singleMessageBuilder[f().GetMessageType()] = f
	}
	initMessageBuilder(func() SingleMessage { return &Text{} })
	initMessageBuilder(func() SingleMessage { return &Face{} })
	initMessageBuilder(func() SingleMessage { return &Image{} })
	initMessageBuilder(func() SingleMessage { return &Record{} })
	initMessageBuilder(func() SingleMessage { return &Video{} })
	initMessageBuilder(func() SingleMessage { return &At{} })
	initMessageBuilder(func() SingleMessage { return &RPS{} })
	initMessageBuilder(func() SingleMessage { return &Dice{} })
	initMessageBuilder(func() SingleMessage { return &Shake{} })
	initMessageBuilder(func() SingleMessage { return &Poke{} })
	initMessageBuilder(func() SingleMessage { return &Anonymous{} })
	initMessageBuilder(func() SingleMessage { return &Share{} })
	initMessageBuilder(func() SingleMessage { return &Contact{} })
	initMessageBuilder(func() SingleMessage { return &Location{} })
	initMessageBuilder(func() SingleMessage { return &Music{} })
	initMessageBuilder(func() SingleMessage { return &Reply{} })
	initMessageBuilder(func() SingleMessage { return &Forward{} })
	initMessageBuilder(func() SingleMessage { return &Node{} })
	initMessageBuilder(func() SingleMessage { return &Xml{} })
	initMessageBuilder(func() SingleMessage { return &Json{} })
}

func parseMessageChain(results []gjson.Result) MessageChain {
	if len(results) == 0 {
		return nil
	}
	ret := make(MessageChain, 0, len(results))
	for i := range results {
		if results[i].Type != gjson.JSON {
			slog.Error("single message is not json: " + results[i].Type.String())
			continue
		}
		singleMessageType := results[i].Get("type").String()
		if builder, ok := singleMessageBuilder[singleMessageType]; ok {
			m := builder()
			if err := json.Unmarshal([]byte(results[i].Get("data").Raw), m); err == nil {
				ret = append(ret, m)
			} else {
				slog.Error("json unmarshal failed", "buf", results[i].Raw, "error", err)
			}
		} else {
			slog.Error("unknown single message type: " + results[i].String())
		}
	}
	return ret
}
