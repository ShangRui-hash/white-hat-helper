package main

import (
	"fmt"

	"github.com/ShangRui-hash/hackflow"
	"github.com/sirupsen/logrus"
)

func main() {
	hackflow.SetDebug(true)
	URLCh := make(chan string, 1024)
	URLCh <- "segmentfault.com"
	close(URLCh)
	httpx := hackflow.GetHttpx()
	resultCh, err := httpx.Run(&hackflow.HttpRunConfig{
		DisplayTitle:     true,
		DisplaySatusCode: true,
		RandomAgent:      true,
		Proxy:            "socks://127.0.0.1:7890",
		Threads:          100,
		URLCh:            URLCh,
	})
	if err != nil {
		logrus.Error("httpx.Run failed,err:", err)
	}
	for result := range resultCh {
		fmt.Println(result)
	}
}
