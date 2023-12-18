package logic

import (
	"context"
	"github.com/lustresix/beifeng/application/user/rpc/internal/e"

	"github.com/lustresix/beifeng/application/user/rpc/internal/svc"
	"github.com/lustresix/beifeng/application/user/rpc/service"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindByMobileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindByMobileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindByMobileLogic {
	return &FindByMobileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindByMobileLogic) FindByMobile(in *service.FindByMobileRequest) (*service.FindByMobileResponse, error) {
	if len(in.Mobile) == 0 {
		return nil, e.MobileEmpty
	}
	users, err := l.svcCtx.UserModel.FindOneByMobile(l.ctx, in.Mobile)
	if err != nil || users == nil {
		logx.Errorf("UserModel FindOne error: %v", e.CannotFindUser)
		return nil, nil
	}

	return &service.FindByMobileResponse{
		UserId:   users.Id,
		Username: users.Username,
		Mobile:   users.Mobile,
		Avatar:   users.Avatar,
	}, nil
}
