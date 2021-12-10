package redis

import (
	"encoding/json"
	"fmt"
	"web_app/models"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

var (
	IPOSMapKey         = "ip_os_map"
	IPPortSetKeyPrefix = "ip_port_set::"
)

//GetHostsByCompanyID 根据公司ID获取主机列表
func GetHostsByCompanyID(param *models.ParamGetHostList) (hostList []*models.HostListItem, err error) {
	key := fmt.Sprintf("ipzset::%d", param.CompanyID)
	zap.L().Debug("key", zap.String("key", key))
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
	portList, err := rdb.SMembers(key).Result()
	if err != nil {
		return nil, err
	}
	portDetailList := make([]models.PortDetail, 0, len(portList))
	for i := range portList {
		var port models.PortDetail
		if err := json.Unmarshal([]byte(portList[i]), &port); err != nil {
			return nil, err
		}
		portDetailList = append(portDetailList, port)
	}
	return portDetailList, nil
}

func GetPort(hostList []*models.HostListItem) error {
	for _, host := range hostList {
		key := fmt.Sprintf("%s%s", IPPortSetKeyPrefix, host.IP)
		portList, err := rdb.SMembers(key).Result()
		if err != nil {
			return err
		}
		for _, portStr := range portList {
			var port models.Port
			if err := json.Unmarshal([]byte(portStr), &port); err != nil {
				return err
			}
			host.PortList = append(host.PortList, port)
		}
	}
	return nil
}
