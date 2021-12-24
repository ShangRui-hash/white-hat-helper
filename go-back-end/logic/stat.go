package logic

import (
	"web_app/dao/redis"
	"web_app/models"
)

func GetCompanyStat(companyID int64) (*models.CompanyStat, error) {
	return redis.GetCompanyStat(companyID)
}
