package svc

import (
	"easy-chat/apps/im/immodels"
	"easy-chat/apps/im/ws/websocket"
	"easy-chat/apps/social/rpc/socialclient"
	"easy-chat/apps/task/mq/internal/config"
	"easy-chat/pkg/constants"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"net/http"
)

type ServiceContext struct {
	config.Config

	WsClient websocket.Client
	*redis.Redis

	socialclient.Social

	immodels.ChatLogModel
	immodels.ConversationModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	svc := &ServiceContext{
		Config:            c,
		Redis:             redis.MustNewRedis(c.Redisx),
		ChatLogModel:      immodels.MustChatLogModel(c.Mongo.Url, c.Mongo.Db),
		ConversationModel: immodels.MustConversationModel(c.Mongo.Url, c.Mongo.Db),
		Social:            socialclient.NewSocial(zrpc.MustNewClient(c.SocialRpc)),
	}
	// 获取token redis 超级token 给这个mq 使用
	token, err := svc.GetSystemToken()
	if err != nil {
		panic(err)
	}

	// 设置token http 请求
	// http 请求头中设置token
	header := http.Header{}
	header.Set("Authorization", token)

	// 创建Websocket客户端，并设置 JWT 认证信息
	svc.WsClient = websocket.NewClient(c.Ws.Host, websocket.WithClientHeader(header))
	return svc
}

func (svc *ServiceContext) GetSystemToken() (string, error) {
	return svc.Redis.Get(constants.RedisSystemRootToken)
}
