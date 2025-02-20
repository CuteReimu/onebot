package onebot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/CuteReimu/goutil"
	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
	"golang.org/x/time/rate"
	"log/slog"
	"net/http"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"
)

// WsChannel 连接通道
type WsChannel string

const (
	// WsChannelApi 连接此接口后，向 OneBot 发送如下结构的 JSON 对象，即可调用相应的 API
	WsChannelApi = "api"

	// WsChannelEvent 连接此接口后，OneBot 会在收到事件后推送至客户端
	WsChannelEvent = "event"

	// WsChannelAll 在一条连接上同时提供 api 和 event 的服务
	WsChannelAll = ""
)

// Connect 连接onebot
//
// concurrentEvent 参数如果是true，表示采用并发方式处理事件和消息，由调用者自行解决并发问题。
// 如果是false表示用单线程处理事件和消息，调用者无需关心并发问题。
func Connect(host string, port int, channel WsChannel, accessToken string, qq int64, concurrentEvent bool) (*Bot, error) {
	addr := fmt.Sprintf("ws://%s:%d/%s", host, port, channel)
	log := slog.With("addr", addr)
	log.Info("Dialing")
	var header http.Header
	if len(accessToken) > 0 {
		header = http.Header{"Authorization": []string{"Bearer " + accessToken}}
	}
	c, resp, err := websocket.DefaultDialer.Dial(addr, header)
	if err != nil {
		log.Error("Connect failed")
		return nil, err
	}
	if resp != nil {
		_ = resp.Body.Close() // 仅为了解决lint警告，调用Close方法其实是没有任何效果的
	}
	log.Info("Connected successfully")
	b := &Bot{QQ: qq, c: c, handler: make(map[string]map[string][]listenHandler)}
	if !concurrentEvent {
		b.eventChan = goutil.NewBlockingQueue[func()]()
		go func() {
			for {
				b.eventChan.Take()()
			}
		}()
	}
	go func() {
		for !b.closed.Load() {
			if b.c == nil {
				time.Sleep(3 * time.Second)
				log.Info("trying to reconnect")
				c, resp, err = websocket.DefaultDialer.Dial(addr, header)
				if err != nil {
					log.Error("Connect failed")
					continue
				}
				if resp != nil {
					_ = resp.Body.Close() // 仅为了解决lint警告，调用Close方法其实是没有任何效果的
				}
				log.Info("Connected successfully")
				b.c = c
			}
			for {
				t, message, err := b.c.ReadMessage()
				if err != nil {
					log.Error("read error", "error", err)
					b.c = nil
					break
				}
				if t != websocket.TextMessage {
					continue
				}
				log.Debug("recv", "msg", string(message))
				if !gjson.ValidBytes(message) {
					log.Error("invalid json message")
					continue
				}
				msg := gjson.ParseBytes(message)
				echo := msg.Get("echo")
				if echo.Exists() {
					e := echo.Int()
					retCode := msg.Get("retcode").Int()
					if ch, ok := b.syncIdMap.LoadAndDelete(e); ok {
						ch0 := ch.(chan gjson.Result)
						if retCode != 0 {
							log.Error("request failed", "retcode", retCode, "msg", msg.Get("message"))
						} else {
							ch0 <- msg.Get("data")
						}
						close(ch0)
					}
					continue
				}
				postType := msg.Get("post_type").String()
				func() {
					b.handlerLock.RLock()
					defer b.handlerLock.RUnlock()
					h, ok := b.handler[postType]
					if !ok {
						return
					}
					subType := msg.Get(postType + "_type").String()
					h2, ok := h[subType]
					if !ok {
						return
					}
					if bd := builder[postType][subType]; bd == nil {
						log.Error("cannot find message builder: " + postType)
					} else {
						m := bd()
						err = json.Unmarshal(message, m)
						if err != nil {
							log.Error("json unmarshal failed", "error", err)
							return
						}
						fun := func() {
							defer func() {
								if r := recover(); r != nil {
									log.Error("panic recovered", "error", r, "stack", string(debug.Stack()))
								}
							}()
							for _, f := range h2 {
								if !f(m) {
									break
								}
							}
						}
						b.Run(fun)
					}
				}()
			}
		}
	}()
	return b, nil
}

type Bot struct {
	QQ          int64
	c           *websocket.Conn
	echo        atomic.Int64
	handlerLock sync.RWMutex
	handler     map[string]map[string][]listenHandler
	syncIdMap   sync.Map
	eventChan   *goutil.BlockingQueue[func()]
	limiter     atomic.Pointer[limiter]
	closed      atomic.Bool
}

type limiter struct {
	limiterType string
	limiter     *rate.Limiter
}

func (l *limiter) check() bool {
	if l.limiterType == "wait" {
		if err := l.limiter.Wait(context.Background()); err != nil {
			slog.Error("rate limiter wait error", "error", err)
			return false
		}
		return true
	} else {
		return l.limiter.Allow()
	}
}

func (b *Bot) Close() error {
	b.closed.Store(true)
	c := b.c
	if c != nil {
		return c.Close()
	}
	return nil
}

// SetLimiter 设置限流器，limiterType为"wait"表示等待，为"drop"表示丢弃
func (b *Bot) SetLimiter(limiterType string, l *rate.Limiter) {
	b.limiter.Store(&limiter{limiterType: limiterType, limiter: l})
}

// Run 如果不是并发方式启动，则此方法会将函数放入事件队列。如果是并发方式启动，则此方法等同于go f()。
func (b *Bot) Run(f func()) {
	if b.eventChan == nil {
		go f()
	} else {
		b.eventChan.Put(f)
	}
}

// request 发送请求
func (b *Bot) request(action string, params any) (gjson.Result, error) {
	limiter := b.limiter.Load()
	if limiter != nil && !limiter.check() {
		return gjson.Result{}, errors.New("rate limit exceeded")
	}
	msg := &requestMessage{
		Echo:   b.echo.Add(1),
		Action: action,
		Params: params,
	}
	echo := msg.Echo
	buf, err := json.Marshal(msg)
	if err != nil {
		slog.Error("json marshal failed", "error", err)
		return gjson.Result{}, err
	}
	ch := make(chan gjson.Result, 1)
	b.syncIdMap.Store(echo, ch)
	c := b.c
	if c == nil {
		slog.Error("disconnected, send failed, please wait for reconnecting")
		return gjson.Result{}, err
	}
	err = c.WriteMessage(websocket.TextMessage, buf)
	if err != nil {
		slog.Error("send error", "error", err)
		return gjson.Result{}, err
	}
	slog.Debug("send", "msg", string(buf))
	timeoutTimer := time.AfterFunc(5*time.Second, func() {
		if ch, ok := b.syncIdMap.LoadAndDelete(echo); ok {
			slog.Error("request timeout")
			close(ch.(chan gjson.Result))
		}
	})
	result, ok := <-ch
	if !ok {
		return gjson.Result{}, errors.New("request failed")
	}
	timeoutTimer.Stop()
	code := result.Get("code").Int()
	if code != 0 {
		e := fmt.Sprint("Non-zero code: ", code, ", error message: ", result.Get("msg"))
		slog.Error(e)
		return gjson.Result{}, errors.New(e)
	}
	return result, nil
}

type simplifier interface {
	simplify() any
}

type quickOperationMessage struct {
	Context   any `json:"context"`
	Operation any `json:"operation"`
}

func (b *Bot) quickOperation(context, operation any) error {
	if s, ok := context.(simplifier); ok {
		context = s.simplify()
	}
	_, err := b.request(".handle_quick_operation", &quickOperationMessage{
		Context:   context,
		Operation: operation,
	})
	return err
}

type requestMessage struct {
	Echo   int64  `json:"echo"`
	Action string `json:"action"`
	Params any    `json:"params,omitempty"`
}

var builder = make(map[string]map[string]func() any)

type listenHandler func(message any) bool

func listen[M any](b *Bot, key, subKey string, l func(message M) bool) {
	b.handlerLock.Lock()
	defer b.handlerLock.Unlock()
	if b.handler[key] == nil {
		b.handler[key] = make(map[string][]listenHandler)
	}
	b.handler[key][subKey] = append(b.handler[key][subKey], func(m any) bool { return l(m.(M)) })
}
