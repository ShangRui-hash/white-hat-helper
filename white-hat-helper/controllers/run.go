package controllers

import (
	"white-hat-helper/dao/redis"
	"white-hat-helper/logger"
	"white-hat-helper/logic"
	"white-hat-helper/settings"

	"github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"
)

func Run(c *cli.Context) error {
	//0.初始化日志配置
	logger.Init()
	//1.初始化Redis连接
	if err := redis.Init(settings.CurrentConfig.RedisConfigFile); err != nil {
		logrus.Errorf("init redis failed,err:%v\n", err)
		return err
	}
	defer redis.Close()
	//2.开始干活
	if err := logic.Run(); err != nil {
		logrus.Error("logic.Run failed,err:", err)
		return nil
	}
	return nil
}
