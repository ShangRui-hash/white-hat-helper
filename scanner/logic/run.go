package logic

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
	"white-hat-helper/dao/redis"
	"white-hat-helper/settings"

	"white-hat-helper/pkg/hackflow"

	"github.com/sirupsen/logrus"
)

func Run() error {
	if settings.CurrentConfig.CompanyID == 0 {
		return fmt.Errorf("company id is empty")
	}
	proxy := "socks://127.0.0.1:7890"
	//生产者：读取域名列表、ip列表
	domainCh, err := readInput()
	if err != nil {
		logrus.Error("readInput failed,err:", err)
		return err
	}
	//1.被动子域名发现,并验证
	hackflow.SetDebug(true)
	subdomainCh, err := hackflow.GetSubfinder().Run(&hackflow.SubfinderRunConfig{
		Proxy:                          proxy,
		Stdin:                          domainCh,
		RemoveWildcardAndDeadSubdomain: true,
		OutputInHostIPFormat:           true,
		OutputInJsonLineFormat:         true,
		Silent:                         true,
		RoutineCount:                   10000,
	}).Result()
	if err != nil {
		logrus.Error("subfinder run failed,err:", err)
		return err
	}
	//2.识别操作系统
	IPAndOSCh := hackflow.GetNmap().OSDection(&hackflow.OSDectionConfig{
		HostCh:    redis.SaveIPDomain(subdomainCh),
		Timeout:   1 * time.Minute,
		BatchSize: 20,
	}).GetIPAndOSCh()
	//2.扫描端口
	// IPAndPortCh, err := hackflow.GetNaabu().Run(&hackflow.NaabuRunConfig{
	// 	Stdin:        hackflow.NewPipe(redis.SaveIPAndOS(IPAndOSCh).GetIPCh()),
	// 	ScanType:     hackflow.CONNECT_SCAN,
	// 	RoutineCount: 1000,
	// }).Result()
	// if err != nil {
	// 	logrus.Error("naabu run failed,err:", err)
	// 	return err
	// }
	IPAndPortCh, err := hackflow.NewPortScanner(20 * time.Second).ConnectScan(
		&hackflow.ScanConfig{
			HostCh:       redis.SaveIPAndOS(IPAndOSCh).GetIPCh(),
			RoutineCount: 1000,
			PortRange:    hackflow.NmapTop1000,
		})
	if err != nil {
		logrus.Error("port scan failed,err:", err)
		return err
	}
	//3.扫描服务
	PortServiceCh := hackflow.GetNmap().ServiceDection(&hackflow.ServiceDectionConfig{
		TargetCh:  redis.SaveIPAndPort(IPAndPortCh),
		Timeout:   2 * time.Minute,
		BatchSize: 30,
	}).GetPortServiceCh()
	urlCh := redis.SavePortService(PortServiceCh).GetWebServiceCh()
	urlChList := hackflow.NewStream().AddSrc(redis.AppendDomainURL(urlCh)).SetDstCount(2).GetDst()
	//4.存储nmap的扫描结果并从端口中提取web服务端口，获取web服务端口的详细信息
	requestCh := hackflow.GenRequest(hackflow.GenRequestConfig{
		URLCh:       urlChList[0],
		MethodList:  []string{http.MethodGet},
		RandomAgent: true,
	})
	responseCh, err := hackflow.RetryHttpSend(&hackflow.RetryHttpSendConfig{
		RequestCh:    requestCh,
		RoutineCount: 1000,
		HttpClientConfig: hackflow.HttpClientConfig{
			Proxy:    proxy,
			RetryMax: 1,
			Redirect: false,
			Checktry: func(ctx context.Context, resp *http.Response, err error) (bool, error) {
				return false, nil
			},
		},
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
	//4.存储响应报文，并对web服务进行指纹识别
	fingerprintCh, err := hackflow.DectWhatWeb(&hackflow.DectWhatWebConfig{
		RoutineCount: 1000,
		TargetCh:     redis.SaveHttpResp(parsedRespCh),
	})
	if err != nil {
		logrus.Error("dectWhatWeb failed,err:", err)
		return err
	}
	//5.对web服务进行目录扫描
	dict, err := os.Open(settings.CurrentConfig.DictPath)
	if err != nil {
		logrus.Error("open dirsearch.txt failed,err:", err)
		return err
	}
	foundURLCh, err := hackflow.BruteForceURL(&hackflow.BruteForceURLConfig{
		BaseURLCh:           redis.SaveFingerprint(fingerprintCh).GetURLCh(),
		RoutineCount:        1000,
		Proxy:               proxy,
		Dictionary:          dict,
		RandomAgent:         true,
		StatusCodeBlackList: []int{404, 405, 403},
	})
	if err != nil {
		logrus.Error("burte force url failed,err:", err)
		return err
	}
	//6.存储目录扫描结果
	redis.SaveFoundURL(foundURLCh)
	logrus.Info("进程运行结束")
	// //2.验证被动发现的域名
	// positiveSubdomainCh, err := hackflow.GetKSubdomain().Run(&hackflow.KSubdomainRunConfig{
	// 	Verify:   true,
	// 	DomainCh: subdomainCh,
	// })
	// if err != nil {
	// 	logrus.Error("ksubdomain run failed,err:", err)
	// 	return err
	// }
	// positiveSubdomainCh2 := hackflow.NewStream().AddSrc(positiveSubdomainCh).SetDstCount(1).AddFilter(func(input string) string {
	// 	fmt.Printf("ksubdomain: %s\n", input)
	// 	return input
	// }).GetDst()[0]

	return nil
}

//readInput 读取输入
func readInput() (Reader io.Reader, err error) {
	domainPipe := hackflow.NewPipe(make(chan interface{}, 1024))
	var wg sync.WaitGroup
	total := 0
	//1.从域名列表中读
	if settings.CurrentConfig.Domains != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, domain := range strings.Split(settings.CurrentConfig.Domains, ",") {
				//注意：这里不能用fmt.Fprintln 这种方法向pipe中写数据，否则会导致读不到数据
				// fmt.Fprintln(domainPipe, strings.TrimSpace(domain))
				domainPipe.Write([]byte(strings.TrimSpace(domain) + "\n"))
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
			_, err := io.Copy(domainPipe, bufio.NewReader(file))
			if err != nil {
				logrus.Errorf("read file failed,err:%v\n", err)
				return
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
		domainPipe.Close()
		logrus.Debug("输入协程工作完成，管道已经关闭，total:", total)
		if total == 0 {
			logrus.Error("please input domain")
			os.Exit(0)
		}
	}()
	return domainPipe, nil
}
