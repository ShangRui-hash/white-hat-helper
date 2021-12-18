package main

import (
	"fmt"

	"github.com/ShangRui-hash/hackflow"
	"github.com/sirupsen/logrus"
)

func main() {
	hackflow.SetDebug(true)
	//1.采集url
	urlCh, err := hackflow.GetUrlCollector().Run(&hackflow.UrlCollectorCofnig{
		Keyword:      ".php?id=1",
		SearchEngine: "baidu",
	})
	if err != nil {
		logrus.Error("GetUrlCollector().Run failed,err:", err)
		return
	}
	//2.将url转化为request
	requestCh := hackflow.GetRequest(urlCh)
	//3.请求request
	responseCh, err := hackflow.RetryHttpSend(&hackflow.RetryHttpSendConfig{
		RequestCh:    requestCh,
		RoutineCount: 100,
	})
	if err != nil {
		logrus.Error("hackflow.RetryHttpSend failed,err:", err)
		return
	}
	//4.解析response
	parsedRespCh, err := hackflow.ParseHttpResp(&hackflow.ParseHttpRespConfig{
		HttpRespCh:   responseCh,
		RoutineCount: 100,
	})

	//5.识别指纹信息
	fingerprintCh, err := hackflow.DectWhatWeb(&hackflow.DectWhatWebConfig{
		TargetCh:     parsedRespCh,
		RoutineCount: 100,
	})
	if err != nil {
		logrus.Error("hackflow.DectWhatWeb failed,err:", err)
		return
	}
	for fingerprint := range fingerprintCh {
		fmt.Println(fingerprint)
	}
}
