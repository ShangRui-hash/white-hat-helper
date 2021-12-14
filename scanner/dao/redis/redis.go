package redis

import (
	"fmt"
	"white-hat-helper/settings"

	"github.com/go-redis/redis"
)

// 声明一个全局的rdb变量
var rdb *redis.Client

// Init 初始化连接
func Init(redisConfigFile string) (err error) {
	//1.读取redis配置文件
	redisConfig := settings.NewRedisConfig()
	if err := redisConfig.Init(redisConfigFile); err != nil {
		return err
	}
	//2.创建redis连接
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
		Password: redisConfig.Passwd, // no password set
		DB:       redisConfig.Db,     // use default DB
		PoolSize: redisConfig.PoolSize,
	})

	_, err = rdb.Ping().Result()
	return err
}

//Close 关闭连接
func Close() {
	rdb.Close()
}
