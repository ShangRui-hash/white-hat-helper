package settings

import (
	"fmt"

	"github.com/spf13/viper"
)

//CurrentConfig 当前配置
var CurrentConfig Config

//Config 配置
type Config struct {
	RedisConfigFile string
	Domains         string
	DomainFile      string
	DictPath        string
	Debug           bool
	CompanyID       int
}

//RedisConfig redis配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Db       int    `mapstructure:"db"`
	Passwd   string `mapstructure:"passwd"`
	PoolSize int    `mapstructure:"pool_size"`
}

func NewRedisConfig() *RedisConfig {
	return &RedisConfig{}
}

//Init 读取Redis配置文件
func (r *RedisConfig) Init(filePath string) error {
	//指定配置文件
	viper.SetConfigFile(filePath)
	//读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("viper.ReadInConfig() failed,err:%v\n", err)
		return err
	}
	//反序列化配置信息
	err = viper.Unmarshal(r)
	if err != nil {
		fmt.Printf("viper.Unmarshal(&Conf) failed,err:%v\n", err)
		return err
	}
	return nil
}
