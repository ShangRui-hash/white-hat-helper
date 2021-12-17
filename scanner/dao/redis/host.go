package redis

import (
	"fmt"
	"white-hat-helper/pkg/hackflow"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

//saveOneIPDomain 保存一个 ip 和 域名之间的对应关系
func saveOneIPDomain(ip, domain string) error {
	//维护一个ip有序集合,键为公司id
	if _, err := rdb.ZAdd(GetIPSetKey(), redis.Z{Score: 0, Member: ip}).Result(); err != nil {
		logrus.Error("rdb.ZAdd(%s, %s) failed,err:%v", GetIPSetKey(), ip, err)
		return err
	}
	//关联ip和域名，以ip为键，域名为值
	if _, err := rdb.SAdd(DomainSetKeyPrefix+ip, domain).Result(); err != nil {
		logrus.Error("rdb.SAdd(%s, %s) failed,err:%v", DomainSetKeyPrefix+ip, domain, err)
		return err
	}
	return nil
}

//SaveIPDomain 保存ip和域名之间的关系
func SaveIPDomain(inputCh <-chan hackflow.IPDomain) chan interface{} {
	outputCh := make(chan interface{}, 10240)
	go func() {
		for input := range inputCh {
			fmt.Printf("save ip:%s,domain:%s\n", input.IP, input.Domain)
			if err := saveOneIPDomain(input.IP, input.Domain); err != nil {
				logrus.Error("redis saveIPDomain failed,err:", err)
				continue
			}
			outputCh <- input.IP
		}
		close(outputCh)
	}()
	return outputCh
}
