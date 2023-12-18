package logic

import (
	"context"
	"github.com/lustresix/beifeng/application/user/rpc/internal/e"
	"github.com/lustresix/beifeng/application/user/rpc/internal/model"
	"github.com/lustresix/beifeng/application/user/rpc/internal/svc"
	"github.com/lustresix/beifeng/application/user/rpc/service"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(i *service.RegisterRequest) (*service.RegisterResponse, error) {
	if len(i.Username) == 0 {
		return nil, e.RegisterNameEmpty
	}
	u := &model.User{
		Username:   i.Username,
		Mobile:     i.Mobile,
		CreateTime: time.Now(),
	}

	insert, err := l.svcCtx.UserModel.Insert(l.ctx, u)
	if err != nil {
		logx.Errorf("Insert UserModel error: %v", err)
		return nil, err
	}
	lastInsertId, err := insert.LastInsertId()
	if err != nil {
		logx.Errorf("LastInsertId error: %v", err)
		return nil, err
	}

	return &service.RegisterResponse{UserId: lastInsertId}, nil
}
