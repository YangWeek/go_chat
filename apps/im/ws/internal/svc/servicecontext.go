package svc

import (
	"easy-chat/apps/im/immodels"
	"easy-chat/apps/im/ws/internal/config"
	"easy-chat/apps/task/mq/mqclient"
)

type ServiceContext struct {
	Config config.Config

	immodels.ChatLogModel
	mqclient.MsgChatTransferClient
	mqclient.MsgReadTransferClient
}

// NewServiceContext 创建并返回一个新的 ServiceContext 实例。
//
// 参数:
//   - c: 配置结构体，包含服务所需的各种配置信息。
//
// 返回值:
//   - *ServiceContext: 一个指向新创建的 ServiceContext 的指针。
func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:                c,
		MsgChatTransferClient: mqclient.NewMsgChatTransferClient(c.MsgChatTransfer.Addrs, c.MsgChatTransfer.Topic),
		MsgReadTransferClient: mqclient.NewMsgReadTransferClient(c.MsgReadTransfer.Addrs, c.MsgReadTransfer.Topic),
		ChatLogModel:          immodels.MustChatLogModel(c.Mongo.Url, c.Mongo.Db),
	}
}
