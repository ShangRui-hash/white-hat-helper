package redis

import (
	"fmt"
	"white-hat-helper/pkg/hackflow"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

//SaveNampResult 保存namp结果,输入为流
func SavePortService(inputCh hackflow.IPAndPortSeviceCh) (outputCh hackflow.IPAndPortSeviceCh) {
	outputCh = make(chan *hackflow.IPAndPortSevice, 10240)
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
func doSaveNampResult(hostListItem *hackflow.IPAndPortSevice) error {
	//2.维护一个ip和 端口+服务的集合
	for _, port := range hostListItem.PortList {
		portStr, err := port.String()
		if err != nil {
			continue
		}
		//维护一个端口号和服务的哈希表,方便更新端口的详细信息
		if _, err := rdb.HSet(IPPortDetailKeyPrefix+hostListItem.IP, fmt.Sprintf("%d", port.Port), portStr).Result(); err != nil {
			logrus.Errorf("rdb.HSet failed,err:", err)
			continue
		}
		logrus.Info("save ip:", hostListItem.IP, "port:", port.Port, "service:", port.Service)
	}
	//3.更新ip 有序集合对应IP的分数
	score := len(hostListItem.PortList) * 10
	if _, err := rdb.ZAdd(GetIPSetKey(), redis.Z{Score: float64(score), Member: hostListItem.IP}).Result(); err != nil {
		logrus.Errorf("rdb.ZAdd failed,err:", err)
		return err
	}
	return nil
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
