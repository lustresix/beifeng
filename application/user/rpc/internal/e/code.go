package e

import "github.com/lustresix/beifeng/pkg/xcode"

var (
	// RegisterNameEmpty 注册名字为空
	RegisterNameEmpty = xcode.New(20001, "注册名字不能为空")

	// MobileEmpty 手机号为空
	MobileEmpty = xcode.New(20001, "手机号为空")

	// CannotFindUser 找不到用户
	CannotFindUser = xcode.New(20001, "找不到用户")

	// IdEmpty 用户Id为空
	IdEmpty = xcode.New(20001, "用户Id为空")
)
