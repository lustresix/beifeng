package logic

import (
	"context"
	"errors"
	"fmt"
	"github.com/lustresix/beifeng/application/applet/internal/e"
	"github.com/lustresix/beifeng/application/user/rpc/user"
	"github.com/lustresix/beifeng/pkg/encrypt"
	"github.com/lustresix/beifeng/pkg/jwt"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"strings"

	"github.com/lustresix/beifeng/application/applet/internal/svc"
	"github.com/lustresix/beifeng/application/applet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterRequest) (resp *types.RegisterResponse, err error) {
	// 清理字符串的空格等
	req.Name = strings.TrimSpace(req.Name)
	req.Mobile = strings.TrimSpace(req.Mobile)
	req.Password = strings.TrimSpace(req.Password)
	req.VerificationCode = strings.TrimSpace(req.VerificationCode)

	// 内容判断，不为空
	if len(req.Mobile) == 0 {
		return nil, e.RegisterMobileEmpty
	}
	if len(req.Password) == 0 {
		return nil, e.RegisterPasswdEmpty
	} else {
		req.Password = encrypt.EncPassword(req.Password)
	}
	if len(req.VerificationCode) == 0 {
		return nil, e.VerificationCodeEmpty
	}

	// 验证码检测
	err = checkVerificationCode(req.Mobile, req.VerificationCode, l.svcCtx.BfRedis)
	if err != nil {
		logx.Errorf("checkVerificationCode error: %v", err)
		return nil, err
	}

	// 加密电话号码
	encMobile, err := encrypt.EncMobile(req.Mobile)
	if err != nil {
		logx.Errorf("EncMobile error: %v", err)
		return nil, err
	}

	// 检查是否已经注册过了
	users, err := l.svcCtx.UserRPC.FindByMobile(l.ctx, &user.FindByMobileRequest{
		Mobile: encMobile,
	})
	fmt.Println(err)
	if err != nil {
		logx.Errorf("UserRPC FindByMobile error: %v", err)
		return nil, err
	}
	fmt.Println(22222)
	if users != nil && users.UserId > 0 {
		return nil, e.MobileHasRegistered
	}
	fmt.Println(users, err)
	// 注册
	registerResponse, err := l.svcCtx.UserRPC.Register(l.ctx, &user.RegisterRequest{
		Username: req.Name,
		Mobile:   encMobile,
	})
	if err != nil {
		logx.Errorf("UserRPC Register error: %v", err)
		return nil, err
	}

	// 编写jwt
	token, err := jwt.CreatToken(jwt.TokenOptions{
		AccessSecret: l.svcCtx.Config.Auth.AccessSecret,
		AccessExpire: l.svcCtx.Config.Auth.AccessExpire,
		Fields: map[string]interface{}{
			"userId": registerResponse.UserId,
		},
	})
	if err != nil {
		logx.Errorf("CreatToken error: %v", err)
		return nil, err
	}
	_ = delActivationCache(req.Mobile, l.svcCtx.BfRedis)

	return &types.RegisterResponse{
		UserId: registerResponse.UserId,
		Token:  types.Token(token),
	}, nil
}

// 验证码检测
func checkVerificationCode(mobile, VerificationCode string, rds *redis.Redis) error {
	cache, err := getActivationCache(mobile, rds)
	if err != nil {
		return err
	}
	if cache == "" {
		return errors.New("verification code expired")
	} else if cache != VerificationCode {
		return errors.New("verification code is not true")
	}
	return nil
}

// 登录之后删除验证码
func delActivationCache(mobile string, rds *redis.Redis) error {
	key := fmt.Sprintf(prefixActivation, mobile)
	_, err := rds.Del(key)
	return err
}
