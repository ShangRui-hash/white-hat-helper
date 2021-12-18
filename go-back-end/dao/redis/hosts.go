package redis

import (
	"fmt"
	"web_app/models"
	"web_app/pkg/hackflow"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

//GetHostsByCompanyID 根据公司ID获取主机列表
func GetHostsByCompanyID(companyID, offset, count int64) (hostList []*models.HostListItem, err error) {
	key := fmt.Sprintf("ipzset::%d", companyID)
	zap.L().Debug("key", zap.String("key", key))
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
	score := len(hostListItem.PortList) * 10
	if _, err := rdb.ZAdd(GetIPSetKey(companyID), redis.Z{Score: float64(score), Member: hostListItem.IP}).Result(); err != nil {
		zap.L().Error("rdb.ZAdd failed,err:", zap.Error(err))
		return err
	}
	return nil
}
