package redis

import "web_app/models"

func GetCompanyStat(companyID int64) (*models.CompanyStat, error) {
	var stat models.CompanyStat
	//1.获取主机数量
	hostCount, err := GetAssetCount(companyID)
	if err != nil {
		return nil, err
	}
	stat.HostCount = uint32(hostCount)
	//2.获取域名数量
	domainCount, err := GetDomainCount(companyID)
	if err != nil {
		return nil, err
	}
	stat.DomainCount = uint32(domainCount)
	//3.获取站点数量
	webSiteCount, err := GetWebSiteCount(companyID)
	if err != nil {
		return nil, err
	}
	stat.WebSiteCount = uint32(webSiteCount)
	return &stat, nil
}
