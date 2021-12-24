package logic

import (
	"net"
	"web_app/dao/redis"
	"web_app/models"
	"web_app/param"

	"go.uber.org/zap"
)

//GetHostList 获取主机列表
func GetHostList(params *param.ParamGetHostList) ([]*models.HostListItem, error) {
	//1.查询主机列表
	hostList, err := redis.GetHostsByCompanyID(params.CompanyID, int64(params.Offset), int64(params.Count))
	if err != nil {
		zap.L().Error("GetHostList redis.GetHostsByCompanyID error", zap.Error(err))
		return nil, err
	}
	for i := range hostList {
		//1.查询主机的操作系统
		if os, err := redis.GetOSByIP(hostList[i].IP); err == nil {
			hostList[i].OS = os
		}
		//2.查询端口信息
		if portInfo, err := redis.GetPortByIP(hostList[i].IP); err == nil {
			hostList[i].PortInfo = portInfo
		}
		//3.查询主机对应的域名信息
		if domainInfo, err := redis.GetDomainInfoByIP(net.ParseIP(hostList[i].IP)); err == nil {
			hostList[i].DomainInfo = domainInfo
		}
		//4.查询web服务信息
		if webInfo, err := redis.GetWebServiceProfileByIP(hostList[i].IP); err == nil {
			hostList[i].WebInfo = webInfo
		}
	}
	return hostList, nil
}

//GetHostBaseInfo 获取主机的基本信息
func GetHostBaseInfo(ip string) (host models.HostBaseInfo, err error) {
	host.IP = ip
	//1.查询操作系统类型
	if os, err := redis.GetOSByIP(ip); err == nil {
		host.OS = os
	}
	//2.查询域名信息
	if domainInfo, err := redis.GetDomainInfoByIP(net.ParseIP(ip)); err == nil {
		host.DomainList = domainInfo.DomainList
	}
	return host, nil
}

//GetHostPortInfo 获取主机的端口信息
func GetHostPortInfo(ip string) ([]models.PortDetail, error) {
	return redis.GetPortDetailByIP(ip)
}

//GetHostWebInfo 查询主机的web服务信息
func GetHostWebInfo(ip string, offset, count int) ([]*models.WebDetail, error) {
	return redis.GetWebServiceByIP(ip, offset, count)
}
