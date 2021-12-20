package redis

import (
	"fmt"
	"net/url"
	"strings"
	"web_app/models"
	"web_app/pkg/hackflow"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

func GetDomainListByIP(IP string) ([]string, error) {
	domains, err := rdb.SMembers(fmt.Sprintf("%s%s", DomainSetKeyPrefix, IP)).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	return domains, nil
}

func GetDomainList(hostList []*models.HostListItem) error {
	for _, host := range hostList {
		domainList, err := GetDomainListByIP(host.IP)
		if err != nil {
			return err
		}
		host.DomainList = domainList
	}
	return nil
}

//AppendDomainURL 将url中的ip替换为域名，并追加到url管道中
func AppendDomainURL(inURLCh chan interface{}) (outURLCh chan interface{}) {
	outURLCh = make(chan interface{}, 10240)
	go func() {
		for input := range inURLCh {
			outURLCh <- input
			u, err := url.Parse(input.(string))
			if err != nil {
				continue
			}
			ipWithPort := strings.Split(u.Host, ":")
			ip := ipWithPort[0]
			port := ipWithPort[1]
			domainList, err := GetDomainListByIP(ip)
			if err != nil {
				continue
			}
			for _, domain := range domainList {
				u := fmt.Sprintf("%s://%s:%v", u.Scheme, domain, port)
				outURLCh <- u
				zap.L().Info("append a url", zap.String("web service", u))
			}
		}
	}()
	return outURLCh
}

//saveOneIPDomain 保存一个 ip 和 域名之间的对应关系
func saveOneIPDomain(ip, domain string, companyID int64) error {
	//维护一个ip有序集合,键为公司id
	if _, err := rdb.ZAdd(GetIPSetKey(companyID), redis.Z{Score: 0, Member: ip}).Result(); err != nil {
		zap.L().Error("rdb.ZAdd failed,err:", zap.Error(err))
		return err
	}
	//关联ip和域名，以ip为键，域名为值
	if _, err := rdb.SAdd(DomainSetKeyPrefix+ip, domain).Result(); err != nil {
		zap.L().Error("rdb.SAdd failed,err:%v", zap.Error(err))
		return err
	}
	return nil
}

//SaveIPDomain 保存ip和域名之间的关系
func SaveIPDomain(inputCh <-chan hackflow.DomainIPs, companyID int64) chan interface{} {
	outputCh := make(chan interface{}, 10240)
	go func() {
		for input := range inputCh {
			fmt.Printf("save ip:%v,domain:%s\n", input.IP, input.Domain)
			for _, ip := range input.IP {
				if err := saveOneIPDomain(ip, input.Domain, companyID); err != nil {
					zap.L().Error("redis saveIPDomain failed,err:", zap.Error(err))
					continue
				}
				outputCh <- ip
			}
		}
		close(outputCh)
	}()
	return outputCh
}
