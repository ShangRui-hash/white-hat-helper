package redis

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

//Counter 计数器
type Counter struct {
	key string
}

//NewCounter 创建一个计数器
func NewCounter(username string) *Counter {
	return &Counter{
		key: fmt.Sprintf("login_counter::%s", username),
	}
}

//Incr 增加该用户名错误登录的次数
func (c *Counter) Incr() error {
	_, err := rdb.Get(c.key).Result()
	if err == redis.Nil {
		//键不存在
		rdb.Set(c.key, 1, 30*time.Minute)
		return nil
	}
	if err != nil {
		return err
	}
	//键存在
	_, err = rdb.IncrBy(c.key, 1).Result()
	return err
}

//Get 获取计数器的值
func (c *Counter) Get() (int, error) {
	val, err := rdb.Get(c.key).Result()
	if err == redis.Nil {
		//键不存在
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(val)
}
