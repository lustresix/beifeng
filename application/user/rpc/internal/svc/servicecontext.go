package svc

import (
	"github.com/lustresix/beifeng/application/user/rpc/internal/config"
	"github.com/lustresix/beifeng/application/user/rpc/internal/model"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config    config.Config
	UserModel model.UserModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	// todo: 数据库还没连
	conn := sqlx.NewMysql(c.DataSource)

	return &ServiceContext{
		Config:    c,
		UserModel: model.NewUserModel(conn, c.CacheRedis),
	}
}
