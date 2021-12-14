package redis

import (
	"encoding/json"
	"fmt"
	"web_app/models"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

//GetHostsByCompanyID 根据公司ID获取主机列表
func GetHostsByCompanyID(param *models.ParamGetHostList) (hostList []*models.HostListItem, err error) {
	key := fmt.Sprintf("ipzset::%d", param.CompanyID)
	zap.L().Debug("key", zap.String("key", key))
	//1.按照分数获取主机列表
	ipList, err := rdb.ZRevRange(key, int64(param.Offset), int64(param.Offset)+int64(param.Count)).Result()
	if err != nil {
		zap.L().Error("get hosts by company id failed", zap.Error(err))
		return nil, err
	}
	for _, ip := range ipList {
		hostList = append(hostList, &models.HostListItem{
			IP: ip,
		})
	}
	return hostList, nil
}

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
	return rdb.HGet(IPOSMapKey, IP).Result()
}

//GetPortDetailByIP 获取端口详情
func GetPortDetailByIP(IP string) ([]models.PortDetail, error) {
	key := fmt.Sprintf("%s%s", IPPortSetKeyPrefix, IP)
	//1.获取该ip所有的端口号
	portList, err := rdb.SMembers(key).Result()
	if err != nil {
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
func GetPort(hostList []*models.HostListItem) error {
	for _, host := range hostList {
		key := fmt.Sprintf("%s%s", IPPortSetKeyPrefix, host.IP)
		//1.获取端口号列表
		portList, err := rdb.SMembers(key).Result()
		if err != nil {
			return err
		}
		zap.L().Debug("portList", zap.Any("portList", portList))
		//2.获取端口的概要信息
		for _, portStr := range portList {
			port, err := getPortDetail(host.IP, portStr)
			if err != nil {
				continue
			}
			host.PortList = append(host.PortList, port.Port)
		}
	}
	return nil
}
