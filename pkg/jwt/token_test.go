package jwt

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var opt = TokenOptions{
	AccessSecret: "lustresix",
	AccessExpire: 36000,
	Fields:       nil,
}

func TestCreatToken(t *testing.T) {
	token, err := CreatToken(opt)
	assert.Nil(t, err)
	t.Log(token)
}

func TestGenToken(t *testing.T) {
	now := time.Now().Add(-time.Minute).Unix()
	token, err := genToken(now, opt.AccessSecret, opt.Fields, opt.AccessExpire)
	assert.Nil(t, err)
	t.Log(token)
}
