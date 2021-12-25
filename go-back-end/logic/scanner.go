package logic

import (
	"context"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
	"web_app/dao/redis"
	"web_app/models"
	"web_app/pkg/hackflow"

	"go.uber.org/zap"
)

func GetDefaultBaseConfig(taskID int64) *hackflow.BaseConfig {
	return &hackflow.BaseConfig{
		CallAfterBegin: func(t hackflow.Tool) {
			if err := redis.BeginTask(taskID, t.GetName()); err != nil {
				zap.L().Error("redis.BeginTask failed,err:", zap.Error(err))
				return
			}
		},
		CallAfterComplete: func(t hackflow.Tool) {
			if err := redis.CompletedTask(taskID, t.GetName()); err != nil {
				zap.L().Error("redis.CompleteTask failed,err:", zap.Error(err))
				return
			}
		},
		CallAfterCtxDone: func(t hackflow.Tool) {
			if err := redis.StopTask(taskID, t.GetName()); err != nil {
				zap.L().Error("redis.StopTask failed,err:", zap.Error(err))
				return
			}
		},
		CallAfterFailed: func(t hackflow.Tool) {
			if err := redis.FailedTask(taskID, t.GetName()); err != nil {
				zap.L().Error("redis.FailedTask failed,err:", zap.Error(err))
				return
			}
		},
	}
}

type scanner struct {
	ctx   context.Context
	proxy string
	task  *models.Task
}

func NewScanner(ctx context.Context, proxy string, task *models.Task) *scanner {
	return &scanner{
		ctx:   ctx,
		proxy: proxy,
		task:  task,
	}
}

//TransformationStage 转化阶段，将目标转化为更多的资产
func (s *scanner) TransformationStage() (chan interface{}, error) {
	//生产者：读取域名列表、ip列表
	domainPipe := make(chan string, 1024)
	// domainChForKSubdomain := make(chan string, 1024)
	go func() {
		for _, domain := range strings.Split(s.task.ScanArea, ",") {
			zap.L().Info("domain:", zap.String("domain", domain))
			domainPipe <- domain
			// domainChForKSubdomain <- domain
		}
		close(domainPipe)
		// close(domainChForKSubdomain)
		zap.L().Debug("从列表中读取域名完成")
	}()
	//1.被动子域名发现,并验证
	subdomainCh, err := hackflow.NewSubfinder(s.ctx).Run(&hackflow.SubfinderRunConfig{
		BaseConfig:                     GetDefaultBaseConfig(s.task.ID),
		Proxy:                          s.proxy,
		DomainCh:                       domainPipe,
		RemoveWildcardAndDeadSubdomain: true,
		OutputInHostIPFormat:           true,
		OutputInJsonLineFormat:         true,
		Silent:                         true,
		RoutineCount:                   1000,
	}).Result()
	if err != nil {
		zap.L().Error("subfinder run failed,err:", zap.Error(err))
		return nil, err
	}
	//2.子域名爆破
	// subdomainCh2, err := hackflow.NewKSubdomain(s.ctx).Run(&hackflow.KSubdomainRunConfig{
	// 	BruteLayer: 1,
	// 	DomainCh:   domainChForKSubdomain,
	// }).Result()
	// if err != nil {
	// 	zap.L().Error("ksubdomain run failed,err:", zap.Error(err))
	// 	return nil, err
	// }
	outCh := make(chan hackflow.DomainIPs, 1024)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for item := range subdomainCh {
			zap.L().Info("xxxsubdomain:", zap.String("subdomain", item.Domain), zap.String("ip", item.IP[0]))
			outCh <- item
		}
	}()
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	for item := range subdomainCh2 {
	// 		outCh <- item
	// 	}
	// }()
	go func() {
		wg.Wait()
		close(outCh)
	}()
	return redis.SaveIPDomain(outCh, s.task.CompanyID), nil
}

//HostScanStage 主机扫描阶段
func (s *scanner) HostScanStage(ipCh chan interface{}) (hackflow.IPAndPortSeviceCh, error) {
	//1.识别操作系统
	IPAndOSCh := hackflow.NewOSDetector(s.ctx).Run(&hackflow.OSDectionConfig{
		BaseConfig: GetDefaultBaseConfig(s.task.ID),
		HostCh:     ipCh,
		Timeout:    1 * time.Minute,
		BatchSize:  20,
	}).GetIPAndOSCh()
	//2.扫描端口
	IPAndPortCh, err := hackflow.NewPortScanner(s.ctx, 20*time.Second).ConnectScan(
		&hackflow.ScanConfig{
			BaseConfig:   GetDefaultBaseConfig(s.task.ID),
			HostCh:       redis.SaveIPAndOS(IPAndOSCh).GetIPCh(),
			RoutineCount: 1000,
			PortRange:    hackflow.NmapTop1000,
		})
	if err != nil {
		zap.L().Error("port scan failed,err:", zap.Error(err))
		return nil, err
	}
	//3.扫描服务
	portServiceCh := hackflow.NewPortServiceDetector(s.ctx).Run(&hackflow.ServiceDectionConfig{
		BaseConfig: GetDefaultBaseConfig(s.task.ID),
		TargetCh:   redis.SaveIPAndPort(IPAndPortCh),
		Timeout:    2 * time.Minute,
		BatchSize:  30,
	}).GetPortServiceCh()
	return portServiceCh, nil
}

//WebScanStage web服务扫描阶段
func (s *scanner) WebScanStage(urlCh chan interface{}, webDirDictionary io.Reader) error {
	requestCh := hackflow.GenRequest(s.ctx, hackflow.GenRequestConfig{
		URLCh:       urlCh,
		MethodList:  []string{http.MethodGet},
		RandomAgent: true,
	})
	responseCh, err := hackflow.RetryHttpSend(s.ctx, &hackflow.RetryHttpSendConfig{
		RequestCh:    requestCh,
		RoutineCount: 100,
		HttpClientConfig: hackflow.HttpClientConfig{
			Proxy:    s.proxy,
			RetryMax: 1,
			Redirect: false,
			Checktry: func(ctx context.Context, resp *http.Response, err error) (bool, error) {
				return false, nil
			},
		},
	})
	if err != nil {
		zap.L().Error("retryHttpSend failed,err:", zap.Error(err))
		return err
	}
	//解析响应报文
	parsedRespCh, err := hackflow.ParseHttpResp(s.ctx, &hackflow.ParseHttpRespConfig{
		RoutineCount: 100,
		HttpRespCh:   responseCh,
	})
	if err != nil {
		zap.L().Error("parseHttpResp failed,err:", zap.Error(err))
		return err
	}
	//4.存储响应报文，并对web服务进行指纹识别
	fingerprintCh, err := hackflow.NewWhatWeb(s.ctx).Run(&hackflow.DectWhatWebConfig{
		BaseConfig:   GetDefaultBaseConfig(s.task.ID),
		RoutineCount: 100,
		TargetCh:     redis.SaveHttpResp(parsedRespCh, s.task.CompanyID),
	})
	if err != nil {
		zap.L().Error("dectWhatWeb failed,err:", zap.Error(err))
		return err
	}
	zap.L().Debug("hackflow.dectWhatWeb return")
	//5.对web服务进行目录扫描
	_ = fingerprintCh
	// respCh, err := hackflow.NewDirSearchGo(s.ctx).Run(&hackflow.BruteForceURLConfig{
	// 	BaseURLCh:           redis.SaveFingerprint(fingerprintCh).GetURLCh(),
	// 	RoutineCount:        100,
	// 	Proxy:               s.proxy,
	// 	Dictionary:          webDirDictionary,
	// 	RandomAgent:         true,
	// 	StatusCodeBlackList: hackflow.DefaultStatusCodeBlackList,
	// })
	// if err != nil {
	// 	zap.L().Error("burte force url failed,err:", zap.Error(err))
	// 	return err
	// }
	// //6.存储目录扫描结果
	// redis.SaveHttpResp(respCh, s.task.CompanyID)
	return nil
}

//run 开始扫描
func (s *scanner) Run(scanArea string, webDirDictionary io.Reader) error {
	ipCh, err := s.TransformationStage()
	if err != nil {
		return err
	}
	zap.L().Debug("s.TransformationStage return")
	ipAndPortServiceCh, err := s.HostScanStage(ipCh)
	if err != nil {
		return err
	}
	//1.提取web服务
	urlCh := redis.SavePortService(ipAndPortServiceCh, s.task.CompanyID).GetWebServiceCh()
	zap.L().Debug("s.HostScanStage return")
	if err := s.WebScanStage(urlCh, webDirDictionary); err != nil {
		return err
	}
	zap.L().Debug("s.WebScanStage return")
	return nil
}
