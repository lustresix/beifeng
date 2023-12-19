package main

import (
	"flag"
	"fmt"
	"github.com/lustresix/beifeng/pkg/xcode"
	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/lustresix/beifeng/application/applet/internal/config"
	"github.com/lustresix/beifeng/application/applet/internal/handler"
	"github.com/lustresix/beifeng/application/applet/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/applet-api.yaml", "the config file")

func main() {
	flag.Parse()

	// 读取配置文件
	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	// 自定义错误处理方法
	httpx.SetErrorHandler(xcode.ErrHandler)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
