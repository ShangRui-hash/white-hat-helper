package redis

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"web_app/dao/memory"
	"web_app/models"
	"web_app/pkg/hackflow"

	"github.com/go-redis/redis"

	"go.uber.org/zap"
)

func GetURLListByCompanyID(companyID int64, offset, count int) ([]string, error) {
	URLList, err := rdb.ZRevRange(GetCompanyWebSiteZSetKey(companyID), int64(offset), int64(offset+count)).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	return URLList, nil
}

func GetURLListByIP(ip string, offset, count int) ([]string, error) {
	urls, err := rdb.ZRevRange(GetWebSiteZSetKeyOfIP(ip), int64(offset), int64(offset+count)).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	return urls, nil
}

//GetALLURLByIP 查询该IP下的所有URL
func GetALLURLByIP(ip string) ([]string, error) {
	URLList, err := rdb.ZRevRange(GetWebSiteZSetKeyOfIP(ip), 0, -1).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	return URLList, nil
}

//GetWebDetailList 根据提供的url列表获取这些url的详细信息
func GetWebDetailList(URLList []string) []*models.WebDetail {
	webs := make([]*models.WebDetail, 0, len(URLList))
	for i := range URLList {
		webDetail := &models.WebDetail{
			WebItem: models.WebItem{
				URL: URLList[i],
			},
			ToolStatusList: make([]models.ToolStatus, 0),
		}
		//1.获取url的指纹信息
		if fingerprintList, err := rdb.SMembers(URLFingerprintSetKeyPrefix + URLList[i]).Result(); err == nil {
			webDetail.WebItem.FingerPrint = fingerprintList
		}
		//2.获取url的响应报文信息
		if resp, err := GetRespByURL(URLList[i]); err == nil {
			webDetail.WebItem.Title = resp.RespTitle
			webDetail.WebItem.StatusCode = resp.StatusCode
			webDetail.WebItem.Location = resp.RespHeader.Get("Location")
			webDetail.RespHeader = resp.RespHeader
			webDetail.RespBody = resp.RespBody
		}
		//4.查询目录扫描器的运行状态
		dirScanStatus := models.ToolStatus{
			Name:   "目录扫描",
			Status: memory.IsURLScanRunning(URLList[i]),
		}
		webDetail.ToolStatusList = append(webDetail.ToolStatusList, dirScanStatus)
		webs = append(webs, webDetail)
	}
	return webs
}

//GetWebServiceCountByCompanyID 获取某个公司下的web站点列表
func GetWebServiceByCompanyID(companyID int64, offset, count int) ([]*models.WebDetail, error) {
	URLList, err := GetURLListByCompanyID(companyID, offset, count)
	if err != nil {
		return nil, err
	}
	return GetWebDetailList(URLList), nil
}

//GetWebServiceByIP 分页获取某个ip下的所有web服务(web站点)
func GetWebServiceByIP(ip string, offset, count int) ([]*models.WebDetail, error) {
	URLList, err := GetURLListByIP(ip, offset, count)
	if err != nil {
		return nil, err
	}
	return GetWebDetailList(URLList), nil
}

//GetAllWebServiceByIP 获取某个ip下的所有web服务
func GetAllWebServiceByIP(ip string) ([]*models.WebDetail, error) {
	//1.查询该ip的URL集合
	URLList, err := GetALLURLByIP(ip)
	if err != nil {
		return nil, err
	}
	//2.查询URL对应的详细信息
	return GetWebDetailList(URLList), nil
}

func GetRespByURL(url string) (hackflow.ParsedHttpResp, error) {
	var resp hackflow.ParsedHttpResp
	detail, err := rdb.HGetAll(URLDetailHashKeyPrefix + url).Result()
	if err != nil {
		return resp, err
	}
	if title, ok := detail["resp_title"]; ok {
		resp.RespTitle = title
	}
	if respHeader, ok := detail["resp_header"]; ok {
		if err := json.Unmarshal([]byte(respHeader), &resp.RespHeader); err != nil {
			zap.L().Error("json.Unmarshal failed ", zap.Error(err))
		}
	}
	if code, ok := detail["status_code"]; ok {
		if statusCode, err := strconv.Atoi(code); err == nil {
			resp.StatusCode = statusCode
		}
	}
	if respBody, ok := detail["resp_body"]; ok {
		resp.RespBody = respBody
	}
	return resp, nil
}

func GetRespHeader(url string) (http.Header, error) {
	var respHeader http.Header
	header, err := rdb.HGet(URLDetailHashKeyPrefix+url, "resp_header").Result()
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(header), &respHeader); err != nil {
		zap.L().Error("json.Unmarshal failed ", zap.Error(err))
	}
	return respHeader, nil
}

func GetRespTitle(url string) (string, error) {
	return rdb.HGet(URLDetailHashKeyPrefix+url, "resp_title").Result()
}
func GetRespStatusCode(url string) (int, error) {
	return rdb.HGet(URLDetailHashKeyPrefix+url, "status_code").Int()
}
func GetRespLocation(url string) (string, error) {
	return rdb.HGet(URLDetailHashKeyPrefix+url, "location").Result()
}

//GetWebServiceProfileByIP 获取对应ip的web服务概要信息
func GetWebServiceProfileByIP(ip string) (*models.WebInfo, error) {
	var webInfo models.WebInfo
	//1.查询该ip的URL集合
	URLList, err := GetALLURLByIP(ip)
	if err != nil {
		return nil, err
	}
	webInfo.Total = len(URLList)
	if webInfo.Total > 6 {
		URLList = URLList[:6]
	}
	//2.查询所有URL的概要信息
	webs := make([]models.WebItem, 0, len(URLList))
	for i := range URLList {
		webItem := models.WebItem{
			URL: URLList[i],
		}
		//1.获取url的指纹信息
		if fingerprintList, err := rdb.SMembers(URLFingerprintSetKeyPrefix + URLList[i]).Result(); err == nil {
			webItem.FingerPrint = fingerprintList
		}
		//2.获取url的响应报文信息
		if location, err := GetRespLocation(URLList[i]); err == nil {
			webItem.Location = location
		}
		if title, err := GetRespTitle(URLList[i]); err == nil {
			webItem.Title = title
		}
		if code, err := GetRespStatusCode(URLList[i]); err == nil {
			webItem.StatusCode = code
		}
		webs = append(webs, webItem)
	}
	webInfo.WebList = webs
	return &webInfo, nil

}

//SaveHttpResp 保存http响应
func SaveHttpResp(parsedResp chan *hackflow.ParsedHttpResp, companyID int64) chan *hackflow.ParsedHttpResp {
	outCh := make(chan *hackflow.ParsedHttpResp, 10240)
	go func() {
		for resp := range parsedResp {
			outCh <- resp
			if err := saveHttpResp(resp, companyID); err != nil {
				continue
			}
			zap.L().Debug("Saved url: ", zap.String("url", resp.URL))
		}
		close(outCh)
	}()
	return outCh
}

//saveHttpResp 保存Http响应报文
func saveHttpResp(resp *hackflow.ParsedHttpResp, companyID int64) error {
	if resp.BaseURL == resp.URL { //说明是站点
		return saveOneWebSite(resp, companyID)
	}
	//说明是站点的子目录
	return saveOneSubDir(resp)
}

//向域名的站点集合中添加站点
func addWebSiteToDomainZSet(domain, url string) error {
	//查询该域名对应的ip地址
	ips, err := GetIPsByDomain(domain)
	if err != nil {
		return err
	}
	//向域名对应的ip的url集合中添加url
	for i := range ips {
		if err := rdb.ZAdd(GetWebSiteZSetKeyOfIP(ips[i]), redis.Z{0, url}).Err(); err != nil {
			return err
		}
	}
	//向域名对应的url集合中添加url
	if err := rdb.ZAdd(GetWebSiteZSetKeyOfIP(domain), redis.Z{Score: 0, Member: url}).Err(); err != nil {
		zap.L().Error("Error saving url: ", zap.Error(err))
		return err
	}
	return nil
}

//向ip的站点集合中添加站点
func addWebSiteToIPZSet(ip net.IP, url string) error {
	//查询ip地址对应的域名
	domainInfo, err := GetDomainInfoByIP(ip)
	if err != nil {
		return err
	}
	//向域名对应的url的有序集合中添加数据
	for i := range domainInfo.DomainList {
		if err := rdb.ZAdd(GetWebSiteZSetKeyOfIP(domainInfo.DomainList[i]), redis.Z{0, url}).Err(); err != nil {
			return err
		}
	}
	//向ip对应的url集合中添加url
	if err := rdb.ZAdd(GetWebSiteZSetKeyOfIP(ip.String()), redis.Z{Score: 0, Member: url}).Err(); err != nil {
		zap.L().Error("Error saving url: ", zap.Error(err))
		return err
	}
	return nil
}

//addWebSiteToCompanyZSet 向公司的站点集合中添加站点
func addWebSiteToCompanyZSet(companyID int64, url string) error {
	return rdb.ZAdd(GetCompanyWebSiteZSetKey(companyID), redis.Z{0, url}).Err()
}

//saveURLDetail 保存url的详细信息
func saveURLDetail(resp *hackflow.ParsedHttpResp) error {
	header, err := json.Marshal(resp.RespHeader)
	if err != nil {
		zap.L().Error("Error marshaling header: ", zap.Error(err))
		return err
	}
	data := map[string]interface{}{
		"method":      resp.Method,
		"status_code": resp.StatusCode,
		"resp_title":  resp.RespTitle,
		"resp_body":   resp.RespBody,
		"resp_header": header,
		"location":    resp.RespHeader.Get("Location"),
	}
	if _, err := rdb.HMSet(GetURLDetailHashKey(resp.BaseURL), data).Result(); err != nil {
		zap.L().Error("Error saving url: ", zap.Error(err))
		return err
	}
	return nil
}

//saveOneWebSite 保存一个Web站点的信息
func saveOneWebSite(resp *hackflow.ParsedHttpResp, companyID int64) error {
	if err := addWebSiteToCompanyZSet(companyID, resp.URL); err != nil {
		return err
	}
	u, err := url.Parse(resp.BaseURL)
	if err != nil {
		zap.L().Error("Error parsing url: ", zap.Error(err))
		return err
	}
	IPorDomain := strings.Split(u.Host, ":")[0]
	if ip := net.ParseIP(IPorDomain); ip != nil {
		//ip下的站点集合
		if err := addWebSiteToIPZSet(ip, resp.URL); err != nil {
			return err
		}
	} else {
		//域名下的站点集合
		if err := addWebSiteToDomainZSet(IPorDomain, resp.URL); err != nil {
			return err
		}
	}
	return saveURLDetail(resp)
}

func saveOneSubDir(resp *hackflow.ParsedHttpResp) error {
	if strings.Contains(hackflow.DefaultStatusCodeBlackList, fmt.Sprintf("%v", resp.StatusCode)) {
		return nil
	}
	//维护一个站点的子目录的有序集合
	if err := rdb.ZAdd(GetSubDirZSetKey(resp.BaseURL), redis.Z{Score: 0, Member: resp.URL}).Err(); err != nil {
		return err
	}
	return saveURLDetail(resp)
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

//DeleteURLSubDir 删除url的指定子目录
func DeleteURLSubDir(parentURL, subURL string) error {
	//1.从集合中删除url
	if _, err := rdb.SRem(GetSubDirZSetKey(parentURL), subURL).Result(); err != nil {
		zap.L().Error("Error saving found url: ", zap.Error(err))
		return err
	}
	//2.删除url的详细信息
	if _, err := rdb.Del(URLDetailHashKeyPrefix + subURL).Result(); err != nil {
		zap.L().Error("Error saving url: ", zap.Error(err))
		return err
	}
	return nil
}

//GetURLByParentURL 获取父url下的所有子url
func GetURLByParentURL(parentURL string, offset, count int) ([]hackflow.BruteForceURLResult, error) {
	//1.获取子目录
	urls, err := rdb.SMembers(GetSubDirZSetKey(parentURL)).Result()
	if err != nil && err != redis.Nil {
		zap.L().Error("Error getting found url: ", zap.Error(err))
		return nil, err
	}
	//2.获取子目录的详细信息
	webList := make([]hackflow.BruteForceURLResult, 0, len(urls))
	for i := range urls {
		var webItem hackflow.BruteForceURLResult
		webItem.URL = urls[i]
		if title, err := GetRespTitle(urls[i]); err == nil {
			webItem.Title = title
		}
		if location, err := GetRespLocation(urls[i]); err == nil {
			webItem.Location = location
		}
		if statusCode, err := GetRespStatusCode(urls[i]); err == nil {
			webItem.StatusCode = statusCode
		}
		webList = append(webList, webItem)
	}
	return webList, nil
}
