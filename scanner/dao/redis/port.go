package redis

import (
	"white-hat-helper/pkg/hackflow"

	"github.com/sirupsen/logrus"
)

func SaveIPAndPort(inputCh <-chan *hackflow.IPAndPort) chan *hackflow.IPAndPort {
	outputCh := make(chan *hackflow.IPAndPort, 10240)
	go func() {
		for input := range inputCh {
			logrus.Infof("save ip:%s,port:%v\n", input.IP, input.Port)
			if err := saveOneIPPort(input.IP, input.Port); err != nil {
				logrus.Error("redis saveIPPort failed,err:", err)
				continue
			}
			outputCh <- input
		}
		close(outputCh)
	}()
	return outputCh
}

func saveOneIPPort(ip string, port int) error {
	//集合中只存储端口号
	_, err := rdb.SAdd(IPPortSetKeyPrefix+ip, port).Result()
	return err
}
