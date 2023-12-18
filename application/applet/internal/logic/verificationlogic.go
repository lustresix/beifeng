package logic

import (
	"context"
	"fmt"
	"github.com/lustresix/beifeng/application/user/rpc/user"
	"github.com/lustresix/beifeng/pkg/util"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"strconv"
	"time"

	"github.com/lustresix/beifeng/application/applet/internal/svc"
	"github.com/lustresix/beifeng/application/applet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	// 加盐
	prefixVerificationCount = "bf#verification#count#%s"
	prefixActivation        = "biz#activation#%s"

	// 每个用户限制一天只能发十次验证码
	verificationLimitPerDay = 10

	// 验证码有效时长为半个小时
	expireActivation = 60 * 30
)

type VerificationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVerificationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerificationLogic {
	return &VerificationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VerificationLogic) Verification(req *types.VerificationRequest) (resp *types.VerificationResponse, err error) {
	count, err := l.getVerificationCount(req.Mobile)
	if err != nil {
		logx.Errorf("getVerificationCount mobile: %s error: %v", req.Mobile, err)
	}
	if count > verificationLimitPerDay {
		return nil, err
	}
	// 是否已经获取过code了
	code, _ := getActivationCache(req.Mobile, l.svcCtx.BfRedis)
	if code == "" {
		code = util.RandomNumeric(6)
	}
	// 发送短信
	_, err = l.svcCtx.UserRPC.SendSms(l.ctx, &user.SendSmsRequest{
		Mobile: req.Mobile,
	})
	if err != nil {
		logx.Errorf("sendSms mobile: %s error: %v", req.Mobile, err)
		return nil, err
	}

	// 把内容存在redis里面
	err = saveMobileInfo(req.Mobile, code, l.svcCtx.BfRedis)
	if err != nil {
		logx.Errorf("saveActivationCache mobile: %s error: %v", req.Mobile, err)
		return nil, err
	}

	// 记录这个手机号一天发过几次信息
	err = l.incrVerificationCount(req.Mobile)
	if err != nil {
		logx.Errorf("incrVerificationCount mobile: %s error: %v", req.Mobile, err)
		return nil, err
	}

	return &types.VerificationResponse{}, nil
}

// 从缓存之中获取电话号码
func (l *VerificationLogic) getVerificationCount(mobile string) (int, error) {
	key := fmt.Sprintf(prefixVerificationCount, mobile)
	val, err := l.svcCtx.BfRedis.Get(key)
	if err != nil {
		return 0, err
	}
	if len(val) == 0 {
		return 0, nil
	}

	return strconv.Atoi(val)
}

func (l *VerificationLogic) incrVerificationCount(mobile string) error {
	key := fmt.Sprintf(prefixVerificationCount, mobile)
	_, err := l.svcCtx.BfRedis.Incr(key)
	if err != nil {
		return err
	}

	return l.svcCtx.BfRedis.Expireat(key, util.EndOfDay(time.Now()).Unix())
}

func getActivationCache(mobile string, rds *redis.Redis) (string, error) {
	key := fmt.Sprintf(prefixActivation, mobile)
	return rds.Get(key)
}

func saveMobileInfo(mobile, code string, rds *redis.Redis) error {
	key := fmt.Sprintf(prefixActivation, mobile)
	return rds.Setex(key, code, expireActivation)
}
