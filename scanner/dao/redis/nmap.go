package redis

import (
	"fmt"
	"white-hat-helper/pkg/hackflow"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

//SaveNampResult 保存namp结果,输入为流
func SaveNampResult(inputCh chan hackflow.HostListItem) (outputCh chan hackflow.HostListItem) {
	outputCh = make(chan hackflow.HostListItem, 10240)
	go func() {
		for hostListItem := range inputCh {
			if err := doSaveNampResult(hostListItem); err != nil {
				logrus.Errorf("doSaveNampResult error:", err)
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
func doSaveNampResult(hostListItem hackflow.HostListItem) error {
	//1.维护一个ip和os 之间的哈希表
	if _, err := rdb.HSet(IPOSMapKey, hostListItem.IP, hostListItem.OS).Result(); err != nil {
		logrus.Errorf("rdb.HSet(IPOSMapKey, %s, %s).Result() error:", hostListItem.IP, hostListItem.OS, err)
		return err
	}
	//2.维护一个ip和 端口+服务的集合
	for _, port := range hostListItem.PortList {
		portStr, err := port.String()
		if err != nil {
			continue
		}
		//集合中只存储端口号
		if _, err := rdb.SAdd(IPPortSetKeyPrefix+hostListItem.IP, port.Port).Result(); err != nil {
			logrus.Errorf("rdb.SAdd failed,err:", err)
			continue
		}
		//维护一个端口号和服务的哈希表,方便更新端口的详细信息
		if _, err := rdb.HSet(IPPortDetailKeyPrefix+hostListItem.IP, fmt.Sprintf("%d", port.Port), portStr).Result(); err != nil {
			logrus.Errorf("rdb.HSet failed,err:", err)
			continue
		}
	}
	//3.更新ip 有序集合对应IP的分数
	score := len(hostListItem.PortList) * 10
	if hostListItem.OS != hackflow.UNKNOWN_OS {
		score += 20
	}
	if _, err := rdb.ZAdd(GetIPSetKey(), redis.Z{Score: float64(score), Member: hostListItem.IP}).Result(); err != nil {
		logrus.Errorf("rdb.ZAdd failed,err:", err)
		return err
	}
	return nil
}
