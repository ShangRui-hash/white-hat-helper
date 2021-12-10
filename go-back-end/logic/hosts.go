package logic

import (
	"web_app/dao/redis"
	"web_app/models"

	"go.uber.org/zap"
)

//GetHostList 获取主机列表
func GetHostList(params *models.ParamGetHostList) ([]*models.HostListItem, error) {
	hostList, err := redis.GetHostsByCompanyID(params)
	if err != nil {
		zap.L().Error("GetHostList redis.GetHostsByCompanyID error", zap.Error(err))
		return nil, err
	}
	if err := redis.GetOS(hostList); err != nil {
		zap.L().Error("GetHostList redis.GetOS error", zap.Error(err))
		return nil, err
	}
	if err := redis.GetPort(hostList); err != nil {
		zap.L().Error("GetHostList redis.GetPort error", zap.Error(err))
		return nil, err
	}
	return hostList, nil
}

func GetHostDetail(ip string) (host models.HostDetail, err error) {
	host.IP = ip
	os, err := redis.GetOSByIP(ip)
	if err != nil {
		zap.L().Error("GetOSByIP failed ", zap.Error(err))
		return host, nil
	}
	host.OS = os
	portDetailList, err := redis.GetPortDetailByIP(ip)
	if err != nil {
		zap.L().Error("GetHostList redis.GetPort error", zap.Error(err))
		return host, err
	}
	host.PortList = portDetailList
	return host, nil
}
