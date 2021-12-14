package main

import (
	"fmt"

	"github.com/ShangRui-hash/hackflow"
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
	stream := hackflow.NewStream()
	subdomainCh, err := hackflow.GetSubfinder().Run(&hackflow.SubfinderRunConfig{
		Proxy:        "socks://127.0.0.1:7890",
		DomainCh:     domainCh,
		RoutineCount: 1000,
	})
	if err != nil {
		logrus.Errorf("subfinder run failed,err:%s", err)
		return
	}
	stream.AddSrc(subdomainCh)
	stream.AddFilter(func(line string) string {
		return "我是过滤器1" + line
	})
	stream.AddFilter(func(line string) string {
		return "我是过滤器2" + line
	})
	for subdomain := range stream.GetDst() {
		fmt.Println(subdomain)
	}
}
