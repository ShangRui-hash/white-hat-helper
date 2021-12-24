package redis

import (
	"fmt"
	"net"
	"net/url"
	"strings"
	"web_app/models"
	"web_app/pkg/hackflow"

	mapset "github.com/deckarep/golang-set"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

//GetDomainInfoByIP 获取ip对应的域名信息
func GetDomainInfoByIP(IP net.IP) (*models.DomainInfo, error) {
	var domainInfo models.DomainInfo
	domains, err := rdb.SMembers(fmt.Sprintf("%s%s", DomainSetKeyPrefix, IP)).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	domainInfo.Total = len(domains)
	domainInfo.DomainList = domains
	return &domainInfo, nil
}

//GetIPsByDomain 获取域名对应的ip地址
func GetIPsByDomain(domain string) ([]string, error) {
	return rdb.SMembers(GetIPSetKey(domain)).Result()
}

//GetDomainURLsByIPUR 根据以ip为host的url获取以域名为host的url列表
func GetDomainURLsByIPURL(IpUrl string) ([]string, error) {
	u, err := url.Parse(IpUrl)
	if err != nil {
		return nil, err
	}
	ip, port := extractIPAndPort(u.Host)
	domainInfo, err := GetDomainInfoByIP(ip)
	if err != nil {
		return nil, err
	}
	domainURLs := make([]string, 0, domainInfo.Total)
	for i := range domainInfo.DomainList {
		domainURLs = append(domainURLs, fmt.Sprintf("%s://%s:%v", u.Scheme, domainInfo.DomainList[i], port))
	}
	return domainURLs, nil
}

//extractIPAndPort 提取ip和端口
func extractIPAndPort(host string) (ip net.IP, port string) {
	ipWithPort := strings.Split(host, ":")
	return net.ParseIP(ipWithPort[0]), ipWithPort[1]
}

//AppendDomainURL 将url中的ip替换为域名，并追加到url管道中
func AppendDomainURL(inURLCh chan interface{}) (outURLCh chan interface{}) {
	outURLCh = make(chan interface{}, 10240)
	go func() {
		for input := range inURLCh {
			outURLCh <- input
			domainURL, err := GetDomainURLsByIPURL(input.(string))
			if err != nil {
				continue
			}
			for i := range domainURL {
				outURLCh <- domainURL[i]
			}
		}
		close(outURLCh)
	}()
	return outURLCh
}

//saveOneIPDomain 保存一个 ip 和 域名之间的对应关系
func saveOneIPDomain(ip, domain string, companyID int64) error {
	//维护一个ip有序集合,键为公司id
	if _, err := rdb.ZAdd(GetCompanyIPZSetKey(companyID), redis.Z{Score: 0, Member: ip}).Result(); err != nil {
		zap.L().Error("rdb.ZAdd failed,err:", zap.Error(err))
		return err
	}
	//维护一个域名有序集合,键为公司id
	if _, err := rdb.ZAdd(GetCompanyDomainZSetKey(companyID), redis.Z{Score: 0, Member: domain}).Result(); err != nil {
		zap.L().Error("rdb.ZAdd failed,err:", zap.Error(err))
		return err
	}
	//关联ip和域名
	//以ip为键，域名为值
	if _, err := rdb.SAdd(DomainSetKeyPrefix+ip, domain).Result(); err != nil {
		zap.L().Error("rdb.SAdd failed,err:%v", zap.Error(err))
		return err
	}
	//以域名为键，ip为值
	if _, err := rdb.SAdd(IPSetKeyPrefix+domain, ip).Result(); err != nil {
		zap.L().Error("rdb.SAdd failed,err:%v", zap.Error(err))
		return err
	}
	return nil
}

//SaveIPDomain 保存ip和域名之间的关系
func SaveIPDomain(inputCh <-chan hackflow.DomainIPs, companyID int64) chan interface{} {
	ipSet := mapset.NewSet()
	outputCh := make(chan interface{}, 10240)
	//读取数据库中已有的ip数据
	go func() {
		ips, err := GetAllIPsByCompanyID(companyID)
		if err != nil {
			zap.L().Error("GetAllIPsByCompanyID failed,err:", zap.Error(err))
			return
		}
		for i := range ips {
			if !ipSet.Contains(ips[i]) {
				outputCh <- ips[i]
				ipSet.Add(ips[i])
			}
		}
	}()
	//存储ip和域名之间的关系
	go func() {
		for input := range inputCh {
			fmt.Printf("save ip:%v,domain:%s\n", input.IP, input.Domain)
			for _, ip := range input.IP {
				if err := saveOneIPDomain(ip, input.Domain, companyID); err != nil {
					zap.L().Error("redis saveIPDomain failed,err:", zap.Error(err))
					continue
				}
				if !ipSet.Contains(ip) {
					outputCh <- ip
					ipSet.Add(ip)
				}
			}
		}
		close(outputCh)
	}()
	return outputCh
}
