package redis

import (
	"fmt"
	"web_app/models"
	"web_app/pkg/hackflow"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

//GetAllIPsByCompanyID 获取公司下所有的IP
func GetAllIPsByCompanyID(companyID int64) (ips []string, err error) {
	key := GetCompanyIPZSetKey(companyID)
	ips, err = rdb.ZRevRange(key, 0, -1).Result()
	if err != nil {
		if err == redis.Nil {
			return []string{}, nil
		}
		return nil, err
	}
	return ips, nil
}

//GetHostsByCompanyID 根据公司ID获取主机列表
func GetHostsByCompanyID(companyID, offset, count int64) (hostList []*models.HostListItem, err error) {
	key := GetCompanyIPZSetKey(companyID)
	//1.按照分数获取主机列表
	ipList, err := rdb.ZRevRange(key, offset, offset+count).Result()
	if err != nil {
		zap.L().Error("get hosts by company id failed", zap.Error(err))
		return nil, err
	}
	for _, ip := range ipList {
		hostList = append(hostList, &models.HostListItem{
			IP: ip,
		})
	}
	//2.启动一个协程，悄咪咪地更新ip的分数
	go func() {
		for _, ip := range ipList {
			if err := UpdateIPScore(companyID, ip); err != nil {
				zap.L().Error("UpdateIPScore failed,err:", zap.Error(err))
				continue
			}
		}
	}()
	return hostList, nil
}

//SaveNampResult 保存namp结果,输入为流
func SavePortService(inputCh hackflow.IPAndPortSeviceCh, companyID int64) (outputCh hackflow.IPAndPortSeviceCh) {
	outputCh = make(chan *hackflow.IPAndPortSevice, 10240)
	go func() {
		for hostListItem := range inputCh {
			if err := doSaveNampResult(hostListItem, companyID); err != nil {
				zap.L().Error("doSaveNampResult error:", zap.Error(err))
				continue
			}
			outputCh <- hostListItem
			fmt.Println("save namp result success,item:", hostListItem)
		}
		close(outputCh)
	}()
	return outputCh
}

//doSaveNampResult 保存一条nmap的结果
func doSaveNampResult(hostListItem *hackflow.IPAndPortSevice, companyID int64) error {
	//2.维护一个ip和 端口+服务的集合
	for _, port := range hostListItem.PortList {
		portStr, err := port.String()
		if err != nil {
			continue
		}
		//维护一个端口号和服务的哈希表,方便更新端口的详细信息
		if _, err := rdb.HSet(IPPortDetailKeyPrefix+hostListItem.IP, fmt.Sprintf("%d", port.Port), portStr).Result(); err != nil {
			zap.L().Error("rdb.HSet failed,err:", zap.Error(err))
			continue
		}
		zap.L().Info("save", zap.String("ip", hostListItem.IP), zap.Int("port:", port.Port), zap.String("service:", port.Service))
	}
	//3.更新ip 有序集合对应IP的分数
	if err := UpdateIPScore(companyID, hostListItem.IP); err != nil {
		zap.L().Error("UpdateIPScore failed,err:", zap.Error(err))
		return err
	}
	return nil
}
