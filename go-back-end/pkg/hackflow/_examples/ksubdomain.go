package main

import (
	"context"
	"fmt"
	"web_app/pkg/hackflow"

	"github.com/sirupsen/logrus"
)

func main() {
	hackflow.SetDebug(true)
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	subdomainCh, err := hackflow.NewKSubdomain(ctx).Run(&hackflow.KSubdomainRunConfig{
		DomainCh:   domainCh,
		BruteLayer: 1,
	}).Result()
	if err != nil {
		logrus.Error("ksubdomain.Run failed,err:", err)
		return
	}
	for subdomain := range subdomainCh {
		fmt.Printf("%+v\n", subdomain)
	}
}
