package main

import (
	"fmt"

	"github.com/ShangRui-hash/hackflow"
	"github.com/sirupsen/logrus"
)

func main() {
	domainCh := make(chan string, 1024)
	domainList := []string{
		"lenovo.com",
		"lenovo.com.cn",
		"lenovomm.com",
		"lenovo.cn",
		"lenovo.net",
		"motorola.com",
		"motorola.com.cn",
		"baiying.cn",
	}
	go func() {
		for _, domain := range domainList {
			domainCh <- domain
		}
		close(domainCh)
	}()
	//1.被动子域名发现
	hackflow.SetDebug(true)
	subdomainCh, err := hackflow.GetSubfinder().Run(&hackflow.SubfinderRunConfig{
		Proxy:    "socks://127.0.0.1:7890",
		DomainCh: domainCh,
	})
	if err != nil {
		logrus.Error("subfinder run failed,err:", err)
		return
	}
	//2.验证被动发现的域名
	positiveSubdomainCh, err := hackflow.GetKSubdomain().Run(&hackflow.KSubdomainRunConfig{
		Verify:   true,
		DomainCh: subdomainCh,
	})
	if err != nil {
		logrus.Error("ksubdomain run failed,err:", err)
		return
	}
	positiveSubdomainCh = SaveDomain(positiveSubdomainCh)
	//3.获取域名对应的title
	resultCh, err := hackflow.GetHttpx().Run(&hackflow.HttpxRunConfig{
		DisplayTitle: true,
		Proxy:        "socks://127.0.0.1:7890",
		URLCh:        positiveSubdomainCh,
	})
	if err != nil {
		logrus.Error("httpx run failed,err:", err)
		return
	}
	for result := range resultCh {
		fmt.Println(result)
	}
}

func SaveDomain(inCh chan string) chan string {
	outCh := make(chan string, 1024)
	go func() {
		for item := range inCh {
			fmt.Println(item)
			outCh <- item
		}
		close(outCh)
	}()
	return outCh
}
