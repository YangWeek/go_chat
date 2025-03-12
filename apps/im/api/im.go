package main

import (
	"flag"
	"fmt"

	"easy-chat/apps/im/api/internal/config"
	"easy-chat/apps/im/api/internal/handler"
	"easy-chat/apps/im/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/im.yaml", "the config file")

// 向上的接口层
// im-api 服务功能如何:
// 1 根据用户获取聊天记录
// 2 建立会话
// 3 获取会话
// 4 更新会话

// 配置应该可以从配置中心拉取
func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	//创建 REST 服务器 并 解决跨域问题
	server := rest.MustNewServer(c.RestConf, rest.WithCors())
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
