package main

import (
	"fmt"

	"github.com/ShangRui-hash/hackflow"
	"github.com/sirupsen/logrus"
)

func main() {
	hosts := []string{
		"yihutest.lenovo.com",
		"mp.csw.lenovo.com",
		"tools.dds.lenovo.com",
		"vaultstudio.lenovo.com",
		"api-us.csw.lenovo.com",
		"sun.cnno1.holo-one.lenovo.com",
		"turn.cnno1.holo-one.lenovo.com",
		"vantagefeed-origin.csw.lenovo.com",
		"chs.lenovo.com",
	}
	hostCh := make(chan string, 1024)
	for _, host := range hosts {
		hostCh <- host
	}
	close(hostCh)
	hackflow.SetDebug(true)
	resultCh, err := hackflow.GetNaabu().Run(&hackflow.NaabuRunConfig{
		RoutineCount: 2000,
		HostCh:       hostCh,
	})
	if err != nil {
		logrus.Error(err)
		return
	}
	for result := range resultCh {
		fmt.Println(result)
	}
}
