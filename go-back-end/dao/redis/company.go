package redis

//GetAssetCount 获取公司的资产数
func GetAssetCount(companyID int64) (count int64, err error) {
	return rdb.ZCard(GetCompanyIPZSetKey(companyID)).Result()
}

//GetWebSiteCount 获取公司的站点数
func GetWebSiteCount(companyID int64) (count int64, err error) {
	return rdb.ZCard(GetCompanyWebSiteZSetKey(companyID)).Result()
}

func GetDomainCount(companyID int64) (count int64, err error) {
	return rdb.ZCard(GetCompanyDomainZSetKey(companyID)).Result()
}
