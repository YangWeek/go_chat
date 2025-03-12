package main

import (
	"easy-chat/apps/im/ws/internal/config"
	"easy-chat/apps/im/ws/internal/handler"
	"easy-chat/apps/im/ws/internal/svc"
	"easy-chat/apps/im/ws/websocket"
	"easy-chat/apps/im/ws/websocket/auth"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"time"
)

var configFile = flag.String("f", "etc/social.yaml", "the config file")

// 处理聊天消息
func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	// 配置初始化，确保服务配置正确无误。

	// 创建服务上下文。
	ctx := svc.NewServiceContext(c)
	// 创建服务
	srv := websocket.NewServer(c.ListenOn,
		// 通过函数式选项模式（Functional Options）将 JWT 认证逻辑注入服务端
		websocket.WithWebsocketAuthentication(auth.NewJwtAuth(ctx)), // 实现这个 jwt的校验
		websocket.WithServerAck(websocket.OnlyAck),                  // 设置ack 机制
		websocket.WithWebsocketMaxConnectionIdle(7*time.Hour),       // 设置最长的连接时间
		websocket.WithServerSendErrCount(3),
	)
	defer srv.Stop()

	// 注册路由

	handler.RegisterHandlers(srv, ctx)

	fmt.Println("start websocket server at ", c.ListenOn, " ..... ")
	srv.Start()
}
