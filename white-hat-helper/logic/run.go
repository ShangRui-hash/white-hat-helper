package logic

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"white-hat-helper/settings"

	"github.com/ShangRui-hash/hackflow"
	"github.com/sirupsen/logrus"
)

func Run() error {
	//生产者：读取域名列表、ip列表
	domainCh, err := readInput()
	if err != nil {
		logrus.Error("readInput failed,err:", err)
		return err
	}
	//1.被动子域名发现
	hackflow.SetDebug(true)
	subdomainCh, err := hackflow.GetSubfinder().Run(&hackflow.SubfinderRunConfig{
		Proxy:        "socks://127.0.0.1:7890",
		DomainCh:     domainCh,
		RoutineCount: 1000,
	})
	if err != nil {
		logrus.Error("subfinder run failed,err:", err)
		return err
	}
	//2.验证被动发现的域名
	positiveSubdomainCh, err := hackflow.GetKSubdomain().Run(&hackflow.KSubdomainRunConfig{
		Verify:   true,
		DomainCh: subdomainCh,
	})
	if err != nil {
		logrus.Error("ksubdomain run failed,err:", err)
		return err
	}
	//3.获取域名对应的title
	requestCh := hackflow.GetRequest(positiveSubdomainCh)
	responseCh, err := hackflow.RetryHttpSend(&hackflow.RetryHttpSendConfig{
		RequestCh:    requestCh,
		RoutineCount: 1000,
		Proxy:        "socks://127.0.0.1:7890",
	})
	if err != nil {
		logrus.Error("retryHttpSend failed,err:", err)
		return err
	}
	//解析响应报文
	parsedRespCh, err := hackflow.ParseHttpResp(&hackflow.ParseHttpRespConfig{
		RoutineCount: 1000,
		HttpRespCh:   responseCh,
	})
	if err != nil {
		logrus.Error("parseHttpResp failed,err:", err)
		return err
	}
	//4.指纹识别
	fingerprintCh, err := hackflow.DectWhatWeb(&hackflow.DectWhatWebConfig{
		RoutineCount: 1000,
		TargetCh:     parsedRespCh,
	})
	if err != nil {
		logrus.Error("dectWhatWeb failed,err:", err)
		return err
	}
	for fingerprint := range fingerprintCh {
		fmt.Println("fingerprint:", fingerprint)
	}
	return nil
}

//readInput 读取输入
func readInput() (chan string, error) {
	domainCh := make(chan string, 1024)
	var wg sync.WaitGroup
	total := 0
	//1.从域名列表中读
	if settings.CurrentConfig.Domains != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, domain := range strings.Split(settings.CurrentConfig.Domains, ",") {
				domainCh <- domain
				total++
			}
			logrus.Debug("从列表中读取域名完成")
		}()
	}
	//2.从文件中读
	if settings.CurrentConfig.DomainFile != "" {
		file, err := os.Open(settings.CurrentConfig.DomainFile)
		if err != nil {
			logrus.Errorf("open file failed,err:%v\n", err)
			return nil, err
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				domainCh <- scanner.Text()
				total++
			}
			logrus.Debug("从文件中读取域名完成")
		}()
	}
	//3.从标准输入读
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	scanner := bufio.NewScanner(os.Stdin)
	// 	for scanner.Scan() {
	// 		domainCh <- scanner.Text()
	// 		total++
	// 	}
	// 	logrus.Debug("从标准输入读取域名完成")
	// }()
	//4.等待三个协程结束，关闭通道的写权限
	go func() {
		logrus.Debug("启动等待协程，等待输入完成")
		wg.Wait()
		close(domainCh)
		logrus.Debug("输入协程工作完成，管道已经关闭，total:", total)
		if total == 0 {
			logrus.Error("please input domain")
			os.Exit(0)
		}
	}()
	return domainCh, nil
}
