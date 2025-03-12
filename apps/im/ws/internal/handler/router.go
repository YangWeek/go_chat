package handler

import (
	"easy-chat/apps/im/ws/internal/handler/conversation"
	"easy-chat/apps/im/ws/internal/handler/push"
	"easy-chat/apps/im/ws/internal/handler/user"
	"easy-chat/apps/im/ws/internal/svc"
	"easy-chat/apps/im/ws/websocket"
)

func RegisterHandlers(srv *websocket.Server, svc *svc.ServiceContext) {
	srv.AddRoutes([]websocket.Route{
		{
			Method:  "user.online",
			Handler: user.OnLine(svc),
		},
		{
			Method:  "conversation.chat",
			Handler: conversation.Chat(svc),
		},
		{
			Method:  "conversation.markChat",
			Handler: conversation.MarkRead(svc), /// 标记消息已读
		},
		{
			Method:  "push",
			Handler: push.Push(svc), // 推送新消息或状态变更
		},
	})
}
