package redis

import (
	"fmt"
)

const (
	DomainSetKeyPrefix         = "domainset::"
	IPOSMapKey                 = "ip_os_map"
	IPPortSetKeyPrefix         = "ip_port_set::"
	IPPortDetailKeyPrefix      = "ip_port_detail::"
	TaskPidHashKey             = "task_pid_hash"
	URLSetKeyPrefix            = "urlset::"      //ip的url集合
	URLDetailHashKeyPrefix     = "urldetail::"   //url的详情hash表
	FoundURLSetKeyPrefix       = "foundurlset::" //url的目录集合
	URLFingerprintSetKeyPrefix = "urlfingerprintset::"
	IPSetKeyPrefix             = "ipzset::"
)

var IPBlackList = []string{
	"127.0.0.1",
	"10.10.10.10",
}

//GetIPSetKey 获取公司的IP集合的key
func GetIPSetKey(companyID int64) string {
	return fmt.Sprintf("%s%d", IPSetKeyPrefix, companyID)
}
