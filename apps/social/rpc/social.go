package main

import (
	"easy-chat/pkg/interceptor"
	"easy-chat/pkg/interceptor/rpcserver"
	"flag"
	"fmt"

	"easy-chat/apps/social/rpc/internal/config"
	"easy-chat/apps/social/rpc/internal/server"
	"easy-chat/apps/social/rpc/internal/svc"
	"easy-chat/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/social.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		social.RegisterSocialServer(grpcServer, server.NewSocialServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})

	// 为gRPC服务添加一元拦截器，以增强服务的功能和性能。
	// 这里分别添加了日志拦截器、幂等性拦截器和同步限流拦截器。
	s.AddUnaryInterceptors(rpcserver.LogInterceptor)
	s.AddUnaryInterceptors(interceptor.NewIdempotenceServer(interceptor.NewDefaultIdempotent(c.Cache[0].RedisConf)))
	s.AddUnaryInterceptors(rpcserver.SyncXLimitInterceptor(100))

	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
