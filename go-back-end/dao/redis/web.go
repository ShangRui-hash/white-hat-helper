package redis

import (
	"encoding/json"
	"strconv"
	"web_app/models"

	"go.uber.org/zap"
)

func GetWebServiceByIP(ip string) ([]models.WebDetail, error) {
	//1.查询该ip的URL集合
	URLList, err := rdb.SMembers(URLSetKeyPrefix + ip).Result()
	if err != nil {
		return nil, err
	}
	//2.查询该ip对应域名的URL集合
	domainList, err := GetDomainListByIP(ip)
	if err != nil {
		return nil, err
	}
	for _, domain := range domainList {
		domainURLList, err := rdb.SMembers(URLSetKeyPrefix + domain).Result()
		if err != nil {
			continue
		}
		URLList = append(URLList, domainURLList...)
	}
	//3.查询URL对应的详细信息
	webs := make([]models.WebDetail, 0, len(URLList))
	for _, url := range URLList {
		webDetail := models.WebDetail{
			WebItem: models.WebItem{
				URL: url,
			},
		}
		//1.获取url的指纹信息
		fingerprintList, err := rdb.SMembers(URLFingerprintSetKeyPrefix + url).Result()
		if err != nil {
			zap.L().Error("rdb.SMembers failed", zap.String("url", url), zap.Error(err))
		}
		webDetail.WebItem.FingerPrint = fingerprintList

		//2.获取url的响应报文信息
		detail, err := rdb.HGetAll(URLDetailHashKeyPrefix + url).Result()
		if err != nil {
			continue
		}

		if title, ok := detail["resp_title"]; ok {
			webDetail.WebItem.Title = title
		}
		if respBody, ok := detail["resp_body"]; ok {
			webDetail.RespBody = respBody
		}
		if respHeader, ok := detail["resp_header"]; ok {
			if err := json.Unmarshal([]byte(respHeader), &webDetail.RespHeader); err != nil {
				zap.L().Error("json.Unmarshal failed ", zap.Error(err))
			}
		}
		if code, ok := detail["status_code"]; ok {
			if statusCode, err := strconv.Atoi(code); err == nil {
				webDetail.WebItem.StatusCode = statusCode
			}
		}
		if location, ok := detail["location"]; ok {
			webDetail.WebItem.Location = location
		}
		webs = append(webs, webDetail)
	}
	return webs, nil
}

func GetWeb(hostList []*models.HostListItem) error {
	for _, host := range hostList {
		webDetailList, err := GetWebServiceByIP(host.IP)
		if err != nil {
			return err
		}
		for _, webDetail := range webDetailList {
			host.WebList = append(host.WebList, webDetail.WebItem)
		}
	}
	return nil
}
