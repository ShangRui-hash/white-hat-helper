package main

import (
	"context"
	"fmt"
	"web_app/pkg/hackflow"

	"github.com/sirupsen/logrus"
)

func main() {
	hackflow.SetDebug(true)
	domainPipe := make(chan string, 10000)
	domainList := []string{"qschou.com", "qsebao.com", "duoerehospital.com", "duoerpharmacy.com"}
	go func() {
		for _, domain := range domainList {
			domainPipe <- domain
		}
		close(domainPipe)
	}()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	subdomainCh, err := hackflow.NewSubfinder(ctx).Run(&hackflow.SubfinderRunConfig{
		Proxy:                          "socks://127.0.0.1:7890",
		DomainCh:                       domainPipe,
		RemoveWildcardAndDeadSubdomain: true,
		OutputInHostIPFormat:           true,
		OutputInJsonLineFormat:         true,
		Silent:                         true,
		RoutineCount:                   1000,
		BaseConfig: &hackflow.BaseConfig{
			CallAfterBegin:    nil,
			CallAfterComplete: nil,
			CallAfterCtxDone:  nil,
			CallAfterFailed:   nil,
		},
	}).Result()
	if err != nil {
		logrus.Errorf("subfinder run failed,err:%s", err)
		return
	}
	for subdomain := range subdomainCh {
		fmt.Println(subdomain)
	}

}
