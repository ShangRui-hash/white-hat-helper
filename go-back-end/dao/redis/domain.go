package redis

import (
	"fmt"
	"web_app/models"
)

func GetDomainListByIP(IP string) ([]string, error) {
	return rdb.SMembers(fmt.Sprintf("%s%s", DomainSetKeyPrefix, IP)).Result()
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
