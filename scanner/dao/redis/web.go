package redis

import (
	"encoding/json"
	"net/url"
	"strings"
	"white-hat-helper/pkg/hackflow"

	"github.com/fatih/structs"
	"github.com/sirupsen/logrus"
)

func SaveHttpResp(parsedResp chan *hackflow.ParsedHttpResp) chan *hackflow.ParsedHttpResp {
	outCh := make(chan *hackflow.ParsedHttpResp, 10240)
	go func() {
		for resp := range parsedResp {
			outCh <- resp
			if err := saveHttpResp(resp); err != nil {
				continue
			}
			logrus.Debug("Saved url: ", resp.URL)
		}
		close(outCh)
	}()
	return outCh
}

func saveHttpResp(resp *hackflow.ParsedHttpResp) error {
	u, err := url.Parse(resp.URL)
	if err != nil {
		logrus.Error("Error parsing url: ", err)
		return err
	}
	//TODO 这里可能是ip也可能是域名，需要判断
	ip := strings.Split(u.Host, ":")[0]
	//1.维护一个ip的url集合
	if _, err := rdb.SAdd(URLSetKeyPrefix+ip, resp.URL).Result(); err != nil {
		logrus.Error("Error saving url: ", err)
		return err
	}
	header, err := json.Marshal(resp.RespHeader)
	if err != nil {
		logrus.Error("Error marshaling header: ", err)
		return err
	}
	//2.维护一个url 详细信息的hash表
	data := map[string]interface{}{
		"status_code": resp.StatusCode,
		"resp_title":  resp.RespTitle,
		"resp_body":   resp.RespBody,
		"resp_header": header,
	}
	if _, err := rdb.HMSet(URLDetailHashKeyPrefix+resp.URL, data).Result(); err != nil {
		logrus.Error("Error saving url: ", err)
		return err
	}
	return nil
}

//SaveFingerprint 存储指纹信息
func SaveFingerprint(fingerprintCh chan *hackflow.DectWhatWebResult) {
	for fingerprint := range fingerprintCh {
		if err := saveFingerprint(fingerprint); err != nil {
			continue
		}
		logrus.Debug("Saved fingerprint: ", fingerprint.URL)
	}
}

func saveFingerprint(fingerprint *hackflow.DectWhatWebResult) error {
	//维护一个url的指纹集合
	for key := range fingerprint.FingerPrint {
		if _, err := rdb.SAdd(URLFingerprintSetKeyPrefix+fingerprint.URL, key).Result(); err != nil {
			logrus.Error("Error saving fingerprint: ", err)
			return err
		}
	}
	return nil
}

//SaveFoundURL 存储发现的url
func SaveFoundURL(foundURLCh chan *hackflow.BruteForceURLResult) {
	for foundURL := range foundURLCh {
		if err := saveFoundURL(foundURL); err != nil {
			continue
		}
		logrus.Debug("Saved found url: ", foundURL.URL)
	}
}

func saveFoundURL(foundURL *hackflow.BruteForceURLResult) error {
	//1.维护一个url的目录集合
	if _, err := rdb.SAdd(FoundURLSetKeyPrefix+foundURL.ParentURL, foundURL.URL).Result(); err != nil {
		logrus.Error("Error saving found url: ", err)
		return err
	}
	//2.维护一个url的详细信息的hash表
	if _, err := rdb.HMSet(URLDetailHashKeyPrefix+foundURL.URL, structs.Map(foundURL)).Result(); err != nil {
		logrus.Error("Error saving url: ", err)
		return err
	}
	return nil
}
