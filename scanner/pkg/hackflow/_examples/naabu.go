package main

import (
	"fmt"
	"white-hat-helper/pkg/hackflow"

	"github.com/sirupsen/logrus"
)

func main() {
	// Raw()
	Result()
}

func Raw() {
	hackflow.SetDebug(false)
	hosts := []string{
		"52.82.107.19",
		"8.25.82.225",
		"47.246.44.227",
		"47.93.151.132",
		"127.0.0.1",
	}

	hostPipe := hackflow.NewPipe(make(chan interface{}, 1024))
	go func() {
		for _, host := range hosts {
			hostPipe.Write([]byte(host + "\n"))
		}
		if err := hostPipe.Close(); err != nil {
			logrus.Errorf("close pipe failed,err:%v", err)
		}
	}()

	naabu := hackflow.GetNaabu().Run(&hackflow.NaabuRunConfig{
		RoutineCount: 1,
		Stdin:        hostPipe,
	})
	go func() {
		stderrCh, err := naabu.GetStderrPipe()
		if err != nil {
			logrus.Error(err)
			return
		}
		for err := range stderrCh.Chan() {
			logrus.Error(err)
		}
	}()
	resultCh, err := naabu.GetStdoutPipe()
	if err != nil {
		logrus.Error(err)
		return
	}
	for result := range resultCh.Chan() {
		fmt.Print(result.(string))
	}
}

func Result() {
	hackflow.SetDebug(true)
	hosts := []string{
		"52.82.107.19",
		"8.25.82.225",
		"47.246.44.227",
		"47.93.151.132",
		"127.0.0.1",
	}

	hostPipe := hackflow.NewPipe(make(chan interface{}, 1024))
	go func() {
		for _, host := range hosts {
			hostPipe.Write([]byte(host))
		}
		// if err := hostPipe.Close(); err != nil {
		// 	logrus.Errorf("close pipe failed,err:%v", err)
		// }
	}()

	naabu := hackflow.GetNaabu().Run(&hackflow.NaabuRunConfig{
		RoutineCount: 10000,
		Stdin:        hostPipe,
		ScanType:     hackflow.CONNECT_SCAN,
	})
	go func() {
		stderrCh, err := naabu.GetStderrPipe()
		if err != nil {
			logrus.Error(err)
			return
		}
		for err := range stderrCh.Chan() {
			logrus.Error("from stderr:", err)
		}
	}()
	resultCh, err := naabu.Result()
	if err != nil {
		logrus.Error(err)
		return
	}
	for result := range resultCh {
		fmt.Println(result)
	}
}
