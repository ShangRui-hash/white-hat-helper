package main

import (
	"fmt"
	"white-hat-helper/pkg/hackflow"

	"github.com/sirupsen/logrus"
)

func main() {
	hackflow.SetDebug(true)
	domainPipe := hackflow.NewPipe(make(chan []byte, 10000))
	domainList := []string{
		"lenovo.com",
		"lenovo.com.cn",
		"lenovomm.com",
		"lenovo.cn",
		"lenovo.net",
		"motorola.com",
		"motorola.com.cn",
	}
	go func() {
		for _, domain := range domainList {
			domainPipe.Chan() <- []byte(domain + "\n")
		}
		domainPipe.Close()
	}()
	subdomainCh, err := hackflow.GetSubfinder().Run(&hackflow.SubfinderRunConfig{
		Proxy:                          "socks://127.0.0.1:7890",
		Stdin:                          domainPipe,
		RemoveWildcardAndDeadSubdomain: true,
		OutputInHostIPFormat:           true,
		OutputInJsonLineFormat:         true,
		Silent:                         true,
		RoutineCount:                   1000,
	}).Result()
	if err != nil {
		logrus.Errorf("subfinder run failed,err:%s", err)
		return
	}
	for subdomain := range subdomainCh {
		fmt.Println(subdomain)
	}

}
