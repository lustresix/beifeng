package logic

import (
	"context"
	"github.com/lustresix/beifeng/application/user/rpc/internal/e"

	"github.com/lustresix/beifeng/application/user/rpc/internal/svc"
	"github.com/lustresix/beifeng/application/user/rpc/service"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindByIdLogic {
	return &FindByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindByIdLogic) FindById(in *service.FindByIdRequest) (*service.FindByIdResponse, error) {
	if in.UserId == 0 {
		return nil, e.IdEmpty
	}
	findUser, err := l.svcCtx.UserModel.FindOne(l.ctx, in.UserId)
	if err != nil {
		logx.Errorf("UserModel FindOne error: %v", err)
		return nil, err
	}
	if findUser == nil || findUser.Id == 0 {
		return nil, e.CannotFindUser
	}
	return &service.FindByIdResponse{
		UserId:   findUser.Id,
		Username: findUser.Username,
		Avatar:   findUser.Avatar,
		Mobile:   findUser.Mobile,
	}, nil
}
