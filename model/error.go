package model

import "errors"

var (
	// ErrUserNoExist 用户不存在.
	ErrUserNoExist = errors.New("The user does not exist")
	// ErrorUserPasswordError 密码错误.
	ErrorUserPasswordError = errors.New("Incorrect user name or password")
	// ErrServerAlreadyExist 指定的服务已存在.
	ErrServerAlreadyExist = errors.New("The service already exists")
	// ErrInvalidParameter 无效的参数
	ErrInvalidParameter = errors.New("Invalid parameter")
)
