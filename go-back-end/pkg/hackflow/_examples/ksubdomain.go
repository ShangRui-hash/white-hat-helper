package main

import (
	"fmt"

	"github.com/ShangRui-hash/hackflow"
	"github.com/sirupsen/logrus"
)

func main() {
	hackflow.SetDebug(true)
	ksubdomain, err := hackflow.GetKSubdomain()
	if err != nil {
		logrus.Error("hackflow.GetKSubdomain faled,err:", err)
		return
	}
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
	subdomainCh, err := ksubdomain.Run(hackflow.KSubdomainRunConfig{
		DomainCh:   domainCh,
		BruteLayer: 1,
	})
	if err != nil {
		logrus.Error("ksubdomain.Run failed,err:", err)
		return
	}
	for subdomain := range subdomainCh {
		fmt.Printf("%+v\n", subdomain)
	}
}
