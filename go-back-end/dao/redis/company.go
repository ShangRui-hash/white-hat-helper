package redis

import "fmt"

//GetAssetCount 获取公司的资产数
func GetAssetCount(CompanyID int64) (count int64, err error) {
	return rdb.ZCard(fmt.Sprintf("%s%d", IPSetKeyPrefix, CompanyID)).Result()
}
