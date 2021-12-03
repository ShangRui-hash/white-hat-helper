package controllers

import (
	"errors"

	"github.com/gin-gonic/gin"
)

//CtxUserIDKey gin.Context中保存的userID的key
const CtxUserIDKey = "userID"

//CtxUsernameKey gin.Context中保存的username的key
const CtxUsernameKey = "username"

var ErrUserNotLogin = errors.New("用户未登录")

//GetCurrentUser 获取gin.Context中保存的当前登录的用户的userID
func GetCurrentUser(c *gin.Context) (int64, error) {
	uid, exist := c.Get(CtxUserIDKey)
	if !exist {
		return 0, ErrUserNotLogin
	}
	userID, ok := uid.(int64)
	if !ok {
		return 0, ErrUserNotLogin
	}
	return userID, nil
}
