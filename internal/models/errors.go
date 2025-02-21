package models

import "errors"

var (
	// 没有在数据库中找到对应的记录
	ErrNoRecord = errors.New("models:no matching record found")
	// 使用错误的email进行登录
	ErrInvalidCredentials = errors.New("models:invalid credential")
	// 尝试通过一个重复的邮箱进行注册
	ErrDuplicateEmail = errors.New("models:duplicate email")
)
