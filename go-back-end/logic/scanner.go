package logic

import (
	"context"
	"io"
	"net/http"
	"strings"
	"time"
	"web_app/dao/redis"
	"web_app/pkg/hackflow"

	"go.uber.org/zap"
)

type scanner struct {
	ctx       context.Context
	proxy     string
	companyID int64
}

func NewScanner(ctx context.Context, proxy string, companyID int64) *scanner {
	return &scanner{
		ctx:       ctx,
		proxy:     proxy,
		companyID: companyID,
	}
}

//TransformationStage 转化阶段，将目标转化为更多的资产
func (s *scanner) TransformationStage(scanArea string) (<-chan hackflow.IPDomain, error) {
	//生产者：读取域名列表、ip列表
	domainPipe := hackflow.NewPipe(make(chan interface{}, 1024))
	go func() {
		for _, domain := range strings.Split(scanArea, ",") {
			domainPipe.Write([]byte(strings.TrimSpace(domain) + "\n")) //注意：这里不能用fmt.Fprintln 这种方法向pipe中写数据，否则会导致读不到数据
		}
		domainPipe.Close()
		zap.L().Debug("从列表中读取域名完成")
	}()
	//1.被动子域名发现,并验证
	subfinder := hackflow.NewSubfinder(s.ctx)
	subdomainCh, err := subfinder.Run(&hackflow.SubfinderRunConfig{
		Proxy:                          s.proxy,
		Stdin:                          domainPipe,
		RemoveWildcardAndDeadSubdomain: true,
		OutputInHostIPFormat:           true,
		OutputInJsonLineFormat:         true,
		Silent:                         true,
		RoutineCount:                   10000,
	}).Result()
	if err != nil {
		zap.L().Error("subfinder run failed,err:", zap.Error(err))
		return nil, err
	}
	return subdomainCh, nil
}

//HostScanStage 主机扫描阶段
func (s *scanner) HostScanStage(subdomainCh <-chan hackflow.IPDomain) (hackflow.IPAndPortSeviceCh, error) {
	//1.识别操作系统
	nmap := hackflow.NewNmap(s.ctx)
	IPAndOSCh := nmap.OSDection(&hackflow.OSDectionConfig{
		HostCh:    redis.SaveIPDomain(subdomainCh, s.companyID),
		Timeout:   1 * time.Minute,
		BatchSize: 20,
	}).GetIPAndOSCh()
	//2.扫描端口
	portScanner := hackflow.NewPortScanner(s.ctx, 20*time.Second)
	IPAndPortCh, err := portScanner.ConnectScan(
		&hackflow.ScanConfig{
			HostCh:       redis.SaveIPAndOS(IPAndOSCh).GetIPCh(),
			RoutineCount: 1000,
			PortRange:    hackflow.NmapTop1000,
		})
	if err != nil {
		zap.L().Error("port scan failed,err:", zap.Error(err))
		return nil, err
	}
	//3.扫描服务
	portServiceCh := nmap.ServiceDection(&hackflow.ServiceDectionConfig{
		TargetCh:  redis.SaveIPAndPort(IPAndPortCh),
		Timeout:   2 * time.Minute,
		BatchSize: 30,
	}).GetPortServiceCh()
	return portServiceCh, nil
}

//WebScanStage web服务扫描阶段
func (s *scanner) WebScanStage(portServiceCh hackflow.IPAndPortSeviceCh, webDirDictionary io.Reader) error {
	//1.提取web服务
	urlCh := redis.SavePortService(portServiceCh, s.companyID).GetWebServiceCh()
	urlChList := hackflow.NewStream().AddSrc(redis.AppendDomainURL(urlCh)).SetDstCount(2).GetDst()
	requestCh := hackflow.GenRequest(s.ctx, hackflow.GenRequestConfig{
		URLCh:       urlChList[0],
		MethodList:  []string{http.MethodGet},
		RandomAgent: true,
	})
	zap.L().Debug("hackflow.GenRequest return")
	responseCh, err := hackflow.RetryHttpSend(s.ctx, &hackflow.RetryHttpSendConfig{
		RequestCh:    requestCh,
		RoutineCount: 1000,
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
	zap.L().Debug("hackflow.RetryHttpSend return ")
	//解析响应报文
	parsedRespCh, err := hackflow.ParseHttpResp(s.ctx, &hackflow.ParseHttpRespConfig{
		RoutineCount: 1000,
		HttpRespCh:   responseCh,
	})
	if err != nil {
		zap.L().Error("parseHttpResp failed,err:", zap.Error(err))
		return err
	}
	zap.L().Debug("hackflow.ParseHttpResp return")
	//4.存储响应报文，并对web服务进行指纹识别
	fingerprintCh, err := hackflow.DectWhatWeb(s.ctx, &hackflow.DectWhatWebConfig{
		RoutineCount: 1000,
		TargetCh:     redis.SaveHttpResp(parsedRespCh),
	})
	if err != nil {
		zap.L().Error("dectWhatWeb failed,err:", zap.Error(err))
		return err
	}
	zap.L().Debug("hackflow.dectWhatWeb return")
	//5.对web服务进行目录扫描
	foundURLCh, err := hackflow.BruteForceURL(s.ctx, &hackflow.BruteForceURLConfig{
		BaseURLCh:           redis.SaveFingerprint(fingerprintCh).GetURLCh(),
		RoutineCount:        1000,
		Proxy:               s.proxy,
		Dictionary:          webDirDictionary,
		RandomAgent:         true,
		StatusCodeBlackList: []int{404, 405, 403},
	})
	if err != nil {
		zap.L().Error("burte force url failed,err:", zap.Error(err))
		return err
	}
	zap.L().Debug("hackflow.bruteForceURL return")
	//6.存储目录扫描结果
	redis.SaveFoundURL(foundURLCh)
	zap.L().Debug("redis.SaveFoundURL return")
	return nil
}

//run 开始扫描
func (s *scanner) Run(scanArea string, webDirDictionary io.Reader) error {
	subdomainCh, err := s.TransformationStage(scanArea)
	if err != nil {
		return err
	}
	zap.L().Debug("s.TransformationStage return")
	ipAndPortServiceCh, err := s.HostScanStage(subdomainCh)
	if err != nil {
		return err
	}
	zap.L().Debug("s.HostScanStage return")
	if err := s.WebScanStage(ipAndPortServiceCh, webDirDictionary); err != nil {
		return err
	}
	zap.L().Debug("s.WebScanStage return")
	return nil
}
