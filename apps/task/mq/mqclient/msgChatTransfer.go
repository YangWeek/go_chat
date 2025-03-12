package mqclient

import (
	"context"
	"easy-chat/apps/task/mq/mq"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-queue/kq"
)

// MsgChatTransferClient 提供发送聊天消息的方法。
// 异步的发送消息
type MsgChatTransferClient interface {
	// Push 发送聊天消息
	Push(msg *mq.MsgChatTransfer) error
}

// msgChatTransferClient 实现了 MsgChatTransferClient 接口，用于将聊天消息推送到消息队列中。
type msgChatTransferClient struct {
	pusher *kq.Pusher
}

// NewMsgChatTransferClient 创建一个新的 MsgChatTransferClient 实例。
//
// 该函数用于初始化并返回一个新的消息推送客户端。
func NewMsgChatTransferClient(addr []string, topic string, opts ...kq.PushOption) MsgChatTransferClient {
	return &msgChatTransferClient{
		pusher: kq.NewPusher(addr, topic),
	}
}

// Push 将聊天消息推送到消息队列中。
//
// 该方法将聊天消息序列化为 JSON 格式，并通过 pusher 推送到消息队列中。
func (c *msgChatTransferClient) Push(msg *mq.MsgChatTransfer) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	err = c.pusher.Push(context.Background(), string(body))
	if err != nil {
		return fmt.Errorf("failed to push message to queue: %w", err)
	}
	return nil
}
