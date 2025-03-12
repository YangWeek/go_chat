package websocket

import (
	"easy-chat/apps/im/ws/websocket/auth"
	"time"
)

// WebSocket 服务配置框架，采用 ‌函数式选项模式
// ServerOptions 定义 WebSocket 服务器的选项配置函数。
type ServerOptions func(opt *websocketOption)

type websocketOption struct {
	auth.Authentication        // WebSocket 服务器的身份认证设置
	patten              string // WebSocket 路由模式

	ack          AckType       // 消息确认类型
	ackTimeout   time.Duration // 消息确认超时时间
	sendErrCount int           // 发送错误次数限制

	maxConnectionIdle time.Duration // 最大连接空闲时间

	concurrency int // 群消息并发处理量级
}

// newWebsocketServerOption 创建一个新的 websocketOption 实例。
//
// 该函数初始化 WebSocket 服务器的选项，并应用提供的 ServerOptions 函数。
//
// 参数:
//   - opts: 可选的 ServerOptions 函数，用于配置 WebSocket 服务器选项。
//
// 返回:
//   - websocketOption: 配置好的 WebSocket 服务器选项。ServerOptions
func newWebsocketServerOption(opts ...ServerOptions) websocketOption {
	o := websocketOption{
		Authentication:    new(auth.WebSocketAuth),
		maxConnectionIdle: defaultMaxConnectionIdle,
		ackTimeout:        defaultAckTimeout,
		sendErrCount:      defaultSendErrCount,
		patten:            "/ws",
		concurrency:       defaultConcurrency,
	}
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

// WithWebsocketAuthentication 配置 WebSocket 服务器的身份认证。
//
// 该函数返回一个 ServerOptions 函数，用于设置 WebSocket 服务器的身份认证。
func WithWebsocketAuthentication(auth auth.Authentication) ServerOptions {
	return func(opt *websocketOption) {
		opt.Authentication = auth
	}
}

// WithWebsocketPatten 配置 WebSocket 的路由模式。
//
// 该函数返回一个 ServerOptions 函数，用于设置 WebSocket 服务器的路由模式。
func WithWebsocketPatten(patten string) ServerOptions {
	return func(opt *websocketOption) {
		opt.patten = patten
	}
}

// WithServerAck 配置消息确认类型。
//
// 该函数返回一个 ServerOptions 函数，用于设置 WebSocket 服务器的消息确认类型。
func WithServerAck(ack AckType) ServerOptions {
	return func(opt *websocketOption) {
		opt.ack = ack
	}
}

// WithServerSendErrCount 配置发送错误次数限制。
//
// 该函数返回一个 ServerOptions 函数，用于设置 WebSocket 服务器的发送错误次数限制。
func WithServerSendErrCount(sendErrCount int) ServerOptions {
	return func(opt *websocketOption) {
		opt.sendErrCount = sendErrCount
	}
}

// WithWebsocketMaxConnectionIdle 配置最大连接空闲时间。
//
// 该函数返回一个 ServerOptions 函数，用于设置 WebSocket 服务器的最大连接空闲时间。
func WithWebsocketMaxConnectionIdle(duration time.Duration) ServerOptions {
	return func(opt *websocketOption) {
		if opt.maxConnectionIdle > 0 {
			opt.maxConnectionIdle = duration
		}
	}
}
