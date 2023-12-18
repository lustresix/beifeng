package jwt

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type (
	TokenOptions struct {
		// 用于令牌签名的密钥
		AccessSecret string

		// 令牌过期时间（秒）
		AccessExpire int64

		// 要包含在令牌负载中的附加字段
		Fields map[string]interface{}
	}

	Token struct {
		// 生成的访问令牌
		AccessToken string `json:"access_token"`

		// 访问令牌的过期时间
		AccessExpire int64 `json:"access_expire"`
	}
)

func CreatToken(opt TokenOptions) (Token, error) {
	var token Token

	now := time.Now().Add(-time.Minute).Unix()
	accessToken, err := genToken(now, opt.AccessSecret, opt.Fields, opt.AccessExpire)
	if err != nil {
		return token, err
	}
	token.AccessToken = accessToken
	token.AccessExpire = now + opt.AccessExpire

	return token, nil
}

func genToken(iat int64, secretKey string, payloads map[string]interface{}, seconds int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	for k, v := range payloads {
		claims[k] = v
	}
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims

	// 使用密钥对令牌进行签名，并返回签名后的字符
	return token.SignedString([]byte(secretKey))
}
