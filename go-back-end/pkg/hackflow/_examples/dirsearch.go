package main

import (
	"context"
	"fmt"
	"net/http"
	"web_app/pkg/hackflow"

	"github.com/sirupsen/logrus"
)

func main() {
	hackflow.SetDebug(true)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	resultCh, err := hackflow.NewDirSearch(ctx).Run(hackflow.DirSearchConfig{
		URL:                 "http://49.4.14.213:80/",
		FullURL:             true,
		RandomAgent:         true,
		HTTPMethod:          http.MethodGet,
		MinRespContentSize:  2,
		StatusCodeBlackList: "403,404,500",
	}).Result()
	if err != nil {
		logrus.Error("dirsearch run failed,err:", err)
		return
	}
	for result := range resultCh {
		fmt.Println(result)
	}
}
