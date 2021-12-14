package main

import (
	"fmt"
	"os"
	"white-hat-helper/pkg/hackflow"

	"github.com/sirupsen/logrus"
)

func main() {
	targetCh := make(chan string, 1024)
	targetCh <- "https://career.huawei.com"
	close(targetCh)
	dict, err := os.Open("./dict.txt")
	if err != nil {
		logrus.Error("open dirsearch.txt failed,err:", err)
		return
	}
	resultCh, err := hackflow.BruteForceURL(&hackflow.BruteForceURLConfig{
		BaseURLCh:           targetCh,
		RoutineCount:        1000,
		Dictionary:          dict,
		RandomAgent:         true,
		StatusCodeBlackList: []int{404, 405, 403},
	})
	for result := range resultCh {
		fmt.Println(result)
	}
}
