package logic

import (
	"context"
	"github.com/lustresix/beifeng/application/applet/internal/e"
	"github.com/lustresix/beifeng/application/applet/internal/svc"
	"github.com/lustresix/beifeng/application/applet/internal/types"
	"github.com/lustresix/beifeng/application/user/rpc/user"
	"github.com/lustresix/beifeng/pkg/encrypt"
	"github.com/lustresix/beifeng/pkg/jwt"
	"github.com/lustresix/beifeng/pkg/xcode"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	// 清理字符串的空格等
	req.Mobile = strings.TrimSpace(req.Mobile)
	req.VerificationCode = strings.TrimSpace(req.VerificationCode)

	if len(req.Mobile) == 0 {
		return nil, e.LoginMobileEmpty
	}
	if len(req.VerificationCode) == 0 {
		return nil, e.VerificationCodeEmpty
	}
	// 先验证验证码是否正确
	err = checkVerificationCode(req.Mobile, req.VerificationCode, l.svcCtx.BfRedis)
	if err != nil {
		logx.Errorf("checkVerificationCode error: %v", err)
		return nil, err
	}

	// 再验证号码是否存在
	mobile, err := encrypt.DecMobile(req.Mobile)
	if err != nil {
		logx.Errorf("DecMobile error: %v", err)
		return nil, err
	}
	u, err := l.svcCtx.UserRPC.FindByMobile(l.ctx, &user.FindByMobileRequest{
		Mobile: mobile,
	})
	if err != nil {
		logx.Errorf("UserRPC FindByMobile error: %v", err)
		return nil, err
	}
	if u == nil || u.UserId == 0 {
		return nil, xcode.AccessDenied
	}

	// 如果查到了，发jwtToken
	token, err := jwt.CreatToken(jwt.TokenOptions{
		AccessSecret: l.svcCtx.Config.Auth.AccessSecret,
		AccessExpire: l.svcCtx.Config.Auth.AccessExpire,
		Fields: map[string]interface{}{
			"userId": u.UserId,
		},
	})
	if err != nil {
		logx.Errorf("CreatToken error: %v", err)
		return nil, err
	}
	_ = delActivationCache(req.Mobile, l.svcCtx.BfRedis)

	return &types.LoginResponse{
		UserId: u.UserId,
		Token: types.Token{
			AccessToken:  token.AccessToken,
			AccessExpire: token.AccessExpire,
		},
	}, nil
}
