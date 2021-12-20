package redis

import (
	"web_app/models"
	"web_app/pkg/hackflow"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

//GetOS 获取操作系统类型
func GetOS(hostList []*models.HostListItem) error {
	for _, host := range hostList {
		zap.L().Debug("ip", zap.String("ip", host.IP))
		os, err := rdb.HGet(IPOSMapKey, host.IP).Result()
		if err != nil {
			if err == redis.Nil {
				continue
			}
		}
		zap.L().Debug("os", zap.String("os", os))
		host.OS = os
	}
	return nil
}

func GetOSByIP(IP string) (string, error) {
	os, err := rdb.HGet(IPOSMapKey, IP).Result()
	if err != nil && err != redis.Nil {
		return "", err
	}
	return os, nil
}

func SaveIPAndOS(IPAndOSCh chan *hackflow.IPAndOS) hackflow.IPAndOSCh {
	outCh := make(chan *hackflow.IPAndOS, 1024)
	go func() {
		for ipAndOS := range IPAndOSCh {
			outCh <- ipAndOS
			doSaveIPAndOS(ipAndOS.IP, ipAndOS.OS)
		}
	}()
	return outCh
}

func doSaveIPAndOS(IP, OS string) error {
	//维护一个ip和os 之间的哈希表
	if _, err := rdb.HSet(IPOSMapKey, IP, OS).Result(); err != nil {
		logrus.Errorf("rdb.HSet(IPOSMapKey, %s, %s).Result() error:", IP, OS, err)
		return err
	}
	logrus.Info("save ip:", IP, "os:", OS)
	return nil
}
