package redis

import (
	"fmt"
	"white-hat-helper/settings"
)

const (
	URLSetKeyPrefix            = "urlset::"      //ip的url集合
	URLDetailHashKeyPrefix     = "urldetail::"   //url的详情hash表
	FoundURLSetKeyPrefix       = "foundurlset::" //url的目录集合
	URLFingerprintSetKeyPrefix = "urlfingerprintset::"
	IPSetKeyPrefix             = "ipzset::"
	DomainSetKeyPrefix         = "domainset::"
	IPOSMapKey                 = "ip_os_map"
	IPPortSetKeyPrefix         = "ip_port_set::"
	IPPortDetailKeyPrefix      = "ip_port_detail::"
)

var IPBlackList = []string{
	"127.0.0.1",
	"10.10.10.10",
}

func GetIPSetKey() string {
	//不能放成全局变量，否则settings.CurrentConfig.CompanyID 还未赋值
	return fmt.Sprintf("%s%d", IPSetKeyPrefix, settings.CurrentConfig.CompanyID)
}
