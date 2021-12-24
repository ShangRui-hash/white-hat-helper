package redis

import (
	"fmt"
)

const (
	DomainSetKeyPrefix          = "domainset::ip::" //ip为键，域名为值
	IPSetKeyPrefix              = "ipset::domain::" //域名为键，ip为值
	IPOSMapKey                  = "ip_os_map"
	IPPortSetKeyPrefix          = "ip_port_set::"
	IPPortDetailKeyPrefix       = "ip_port_detail::"
	TaskPidHashKey              = "task_pid_hash"
	URLSetKeyPrefix             = "urlset::"     //ip的url集合
	URLDetailHashKeyPrefix      = "urldetail::"  //url的详情hash表
	SubDirZSetKeyPrefix         = "subdirzset::" //url的目录集合
	URLFingerprintSetKeyPrefix  = "urlfingerprintset::"
	CompanyIPSetKeyPrefix       = "ipzset::company_id::"      //公司的ip集合
	CompanyDomainSetKeyPrefix   = "domainzset::company_id::"  //公司的域名集合
	CompanyWebSiteZSetKeyPrefix = "websitezset::company_id::" //公司的web站点集合
)

//GetIPSetKey 获取公司的IP集合的key
func GetCompanyIPZSetKey(companyID int64) string {
	return fmt.Sprintf("%s%d", CompanyIPSetKeyPrefix, companyID)
}

//GetCompanyDomainSetKey 获取公司的域名集合的key
func GetCompanyDomainZSetKey(companyID int64) string {
	return fmt.Sprintf("%s%d", CompanyDomainSetKeyPrefix, companyID)
}

func GetCompanyWebSiteZSetKey(companyID int64) string {
	return fmt.Sprintf("%s%d", CompanyWebSiteZSetKeyPrefix, companyID)
}

//GetWebSiteZSetKeyOfIP 获取ip所对应的站点集合的key
func GetWebSiteZSetKeyOfIP(ip string) string {
	return fmt.Sprintf("%s%s", URLSetKeyPrefix, ip)
}

//GetURLDetailHashKey 获取url的详情hash表的key
func GetURLDetailHashKey(url string) string {
	return fmt.Sprintf("%s%s", URLDetailHashKeyPrefix, url)
}

//GetIPSetKey 获取域名对应的ip集合的key
func GetIPSetKey(domain string) string {
	return fmt.Sprintf("%s%s", IPSetKeyPrefix, domain)
}

func GetSubDirZSetKey(url string) string {
	return fmt.Sprintf("%s%s", SubDirZSetKeyPrefix, url)
}
