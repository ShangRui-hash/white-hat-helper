package redis

import (
	"encoding/json"
	"fmt"
	"web_app/models"
	"web_app/pkg/hackflow"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

func SaveIPAndPort(inputCh <-chan *hackflow.IPAndPort) chan *hackflow.IPAndPort {
	outputCh := make(chan *hackflow.IPAndPort, 10240)
	go func() {
		for input := range inputCh {
			zap.L().Debug("save", zap.String("ip", input.IP), zap.Int("port", input.Port))
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

//GetPortDetailByIP 获取端口详情
func GetPortDetailByIP(IP string) ([]models.PortDetail, error) {
	key := fmt.Sprintf("%s%s", IPPortSetKeyPrefix, IP)
	//1.获取该ip所有的端口号
	portList, err := rdb.SMembers(key).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	//2.从hash表中查询端口的具体信息
	portDetailList := make([]models.PortDetail, 0, len(portList))
	for i := range portList {
		port, err := getPortDetail(IP, portList[i])
		if err != nil {
			continue
		}
		portDetailList = append(portDetailList, port)
	}
	return portDetailList, nil
}

func getPortDetail(IP, portID string) (port models.PortDetail, err error) {
	//查询端口的详细信息
	portDetail, err := rdb.HGet(IPPortDetailKeyPrefix+IP, portID).Result()
	if err != nil {
		zap.L().Error("rdb.HGet failed", zap.Error(err))
		return port, err
	}
	//反序列化端口的详细信息
	if err := json.Unmarshal([]byte(portDetail), &port); err != nil {
		return port, err
	}
	return port, nil
}

//GetPort 获取端口的概要信息
func GetPortByIP(ip string) (*models.PortInfo, error) {
	var portInfo models.PortInfo
	key := fmt.Sprintf("%s%s", IPPortSetKeyPrefix, ip)
	//1.获取端口号列表
	portList, err := rdb.SMembers(key).Result()
	if err != nil {
		return nil, err
	}
	portInfo.Total = len(portList)
	if portInfo.Total > 20 {
		portList = portList[:20]
	}
	portInfo.PortList = make([]models.Port, 0, len(portList))
	//2.获取端口的概要信息
	for i := range portList {
		port, err := getPortDetail(ip, portList[i])
		if err != nil {
			continue
		}
		portInfo.PortList = append(portInfo.PortList, port.Port)
	}
	return &portInfo, nil
}
