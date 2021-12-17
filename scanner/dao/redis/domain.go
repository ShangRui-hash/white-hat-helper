package redis

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/sirupsen/logrus"
)

func GetDomainListByIP(IP string) ([]string, error) {
	return rdb.SMembers(fmt.Sprintf("%s%s", DomainSetKeyPrefix, IP)).Result()
}

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
				outURLCh <- fmt.Sprintf("%s://%s:%v", u.Scheme, domain, port)
				logrus.Info("web service:", fmt.Sprintf("%s://%s:%v", u.Scheme, domain, port))
			}
		}
	}()

	return outURLCh

}
