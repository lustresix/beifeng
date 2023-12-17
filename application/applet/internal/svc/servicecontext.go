package svc

import (
	"github.com/lustresix/beifeng/application/applet/internal/config"
	"github.com/lustresix/beifeng/application/user/rpc/user"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type ServiceContext struct {
	Config  config.Config
	UserRPC user.User
	BfRedis *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	conf := redis.RedisConf{
		Host: c.BfRedis.Host,
		Type: "node",
		Pass: c.BfRedis.Pass,
	}

	newRedis, err := redis.NewRedis(conf, redis.WithPass(c.BfRedis.Pass))
	if err != nil {
		panic(err)
	}
	return &ServiceContext{
		Config:  c,
		BfRedis: newRedis,
	}
}
