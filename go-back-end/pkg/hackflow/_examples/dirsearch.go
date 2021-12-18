package main

import (
	"fmt"
	"white-hat-helper/pkg/hackflow"

	"github.com/sirupsen/logrus"
)

func main() {
	hackflow.SetDebug(true)
	urlCh := make(chan []byte, 1)
	urlCh <- []byte("https://365.lenovo.com.cn/")
	resultCh, err := hackflow.GetDirSearch().Run(hackflow.DirSearchConfig{
		Stdin:       hackflow.NewPipe(urlCh),
		FullURL:     true,
		RandomAgent: true,
	}).ParsedResult()
	if err != nil {
		logrus.Error("dirsearch run failed,err:", err)
		return
	}
	for result := range resultCh {
		fmt.Println(result)
	}
}
