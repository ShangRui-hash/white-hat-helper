package redis

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
	"web_app/models"
	"web_app/pkg/hackflow"

	"github.com/gqcn/structs"

	"go.uber.org/zap"
)

//GetWebServiceByIP 获取某个ip下的所有web服务
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

//GetWeb 获取hostList中所有主机的web服务
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

//SaveHttpResp 保存http响应
func SaveHttpResp(parsedResp chan *hackflow.ParsedHttpResp) chan *hackflow.ParsedHttpResp {
	outCh := make(chan *hackflow.ParsedHttpResp, 10240)
	go func() {
		for resp := range parsedResp {
			outCh <- resp
			if err := saveHttpResp(resp); err != nil {
				continue
			}
			zap.L().Debug("Saved url: ", zap.String("url", resp.URL))
		}
		close(outCh)
	}()
	return outCh
}

func saveHttpResp(resp *hackflow.ParsedHttpResp) error {
	u, err := url.Parse(resp.URL)
	if err != nil {
		zap.L().Error("Error parsing url: ", zap.Error(err))
		return err
	}
	//TODO 这里可能是ip也可能是域名，需要判断
	ip := strings.Split(u.Host, ":")[0]
	//1.维护一个ip的url集合
	if _, err := rdb.SAdd(URLSetKeyPrefix+ip, resp.URL).Result(); err != nil {
		zap.L().Error("Error saving url: ", zap.Error(err))
		return err
	}
	header, err := json.Marshal(resp.RespHeader)
	if err != nil {
		zap.L().Error("Error marshaling header: ", zap.Error(err))
		return err
	}
	//2.维护一个url 详细信息的hash表
	data := map[string]interface{}{
		"status_code": resp.StatusCode,
		"resp_title":  resp.RespTitle,
		"resp_body":   resp.RespBody,
		"resp_header": header,
		"location":    resp.RespHeader.Get("Location"),
	}
	if _, err := rdb.HMSet(URLDetailHashKeyPrefix+resp.URL, data).Result(); err != nil {
		zap.L().Error("Error saving url: ", zap.Error(err))
		return err
	}
	return nil
}

//SaveFingerprint 存储指纹信息
func SaveFingerprint(fingerprintCh hackflow.DectWhatWebResultCh) hackflow.DectWhatWebResultCh {
	outCh := make(chan *hackflow.DectWhatWebResult, 10240)
	go func() {
		defer close(outCh)
		for fingerprint := range fingerprintCh {
			outCh <- fingerprint
			if err := saveFingerprint(fingerprint); err != nil {
				continue
			}
			zap.L().Debug("Saved fingerprint: ", zap.String("url", fingerprint.URL))
		}
	}()
	return outCh
}

//saveFingerprint 存储指纹信息
func saveFingerprint(fingerprint *hackflow.DectWhatWebResult) error {
	//维护一个url的指纹集合
	for key := range fingerprint.FingerPrint {
		if _, err := rdb.SAdd(URLFingerprintSetKeyPrefix+fingerprint.URL, key).Result(); err != nil {
			zap.L().Error("Error saving fingerprint: ", zap.Error(err))
			return err
		}
	}
	return nil
}

//SaveFoundURL 存储发现的url
func SaveFoundURL(foundURLCh chan *hackflow.BruteForceURLResult) {
	go func() {
		for foundURL := range foundURLCh {
			if err := saveFoundURL(foundURL); err != nil {
				continue
			}
			zap.L().Debug("Saved found url: ", zap.String("url:", foundURL.URL))
		}
	}()
}

func saveFoundURL(foundURL *hackflow.BruteForceURLResult) error {
	//1.维护一个url的目录集合
	if _, err := rdb.SAdd(FoundURLSetKeyPrefix+foundURL.ParentURL, foundURL.URL).Result(); err != nil {
		zap.L().Error("Error saving found url: ", zap.Error(err))
		return err
	}
	//2.维护一个url的详细信息的hash表
	if _, err := rdb.HMSet(URLDetailHashKeyPrefix+foundURL.URL, structs.Map(foundURL)).Result(); err != nil {
		zap.L().Error("Error saving url: ", zap.Error(err))
		return err
	}
	return nil
}
