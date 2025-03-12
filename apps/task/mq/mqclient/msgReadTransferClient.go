package mqclient

import (
	"context"
	"easy-chat/apps/task/mq/mq"
	"encoding/json"
	"github.com/zeromicro/go-queue/kq"
)

type MsgReadTransferClient interface {
	Push(msg *mq.MsgMarkRead) error
}

type msgReadTransferClient struct {
	pusher *kq.Pusher
}

func NewMsgReadTransferClient(addr []string, topic string, opts ...kq.PushOption) *msgReadTransferClient {
	return &msgReadTransferClient{
		pusher: kq.NewPusher(addr, topic, opts...),
	}
}

// 把已读的消息推送到kafka
func (c *msgReadTransferClient) Push(msg *mq.MsgMarkRead) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return c.pusher.Push(context.Background(), string(body))
}
