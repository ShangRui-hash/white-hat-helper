package redis

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

//LoginTime 登录时间
type LoginTime struct {
	key string
}

//NewLoginTime 新建一个记录登录时间的Hash
func NewLoginTime() *LoginTime {
	return &LoginTime{
		key: fmt.Sprintf("%s::login_time", viper.GetString("redis.key_prefix")),
	}
}

//UpdateLoginTime 更新登录时间
func (l *LoginTime) UpdateLoginTime(username string) (bool, error) {
	return rdb.HSet(l.key, username, time.Now().Format("2006-01-02T15:04:05Z")).Result()
}

//Get 获取登录时间
func (l *LoginTime) Get(username string) (string, error) {
	return rdb.HGet(l.key, username).Result()
}
