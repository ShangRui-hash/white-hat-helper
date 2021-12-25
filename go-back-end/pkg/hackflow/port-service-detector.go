package hackflow

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Ullaakut/nmap"
)

type portServiceDetector struct {
	baseTool
	portServiceResultCh chan *nmap.Run
}

func NewPortServiceDetector(ctx context.Context) *portServiceDetector {
	return &portServiceDetector{
		baseTool: baseTool{
			ctx:  ctx,
			name: "port-service-detector",
			desp: "端口服务识别",
		},
		portServiceResultCh: make(chan *nmap.Run, 1024),
	}
}

type ServiceDectionConfig struct {
	*BaseConfig
	TargetCh  chan *IPAndPort
	Timeout   time.Duration
	BatchSize int
}

func (n *portServiceDetector) Run(config *ServiceDectionConfig) *portServiceDetector {
	var count int32 = 0
	var wg sync.WaitGroup
	lackCh := make(chan struct{}, 1024)
	//监控当前进程数量，如果小于最大进程，就向lackCh中发送信号
	go func() {
	LOOP:
		for {
			select {
			case <-n.ctx.Done():
				break LOOP
			default:
				if count < int32(config.BatchSize) {
					lackCh <- struct{}{}
				}
			}
		}
		close(lackCh)
	}()

	go func() {
	LOOP:
		for {
			select {
			case <-n.ctx.Done():
				if config.CallAfterCtxDone != nil {
					config.CallAfterCtxDone(n)
				}
				break LOOP
			case <-lackCh:
				target, ok := <-config.TargetCh
				if !ok {
					break LOOP
				}
				atomic.AddInt32(&count, 1)
				wg.Add(1)
				go func() {
					defer wg.Done()
					if err := n.doServiceDection(target, config.Timeout); err != nil {
						logger.Error("portServiceDetector run faield: ", err)
					}
					atomic.AddInt32(&count, -1)
				}()
			}
		}
		wg.Wait()
		close(n.portServiceResultCh)
		if config.CallAfterComplete != nil {
			config.CallAfterComplete(n)
		}
	}()
	if config.CallAfterBegin != nil {
		config.CallAfterBegin(n)
	}
	return n
}

func (n *portServiceDetector) doServiceDection(target *IPAndPort, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(n.ctx, timeout)
	defer cancel()
	logger.Debug("nmap do service dection:", target)
	scanner, err := nmap.NewScanner(
		nmap.WithTargets(target.IP),
		nmap.WithPorts(fmt.Sprintf("%d", target.Port)),
		nmap.WithServiceInfo(),
		nmap.WithContext(ctx),
		nmap.WithSkipHostDiscovery(), // -Pn
	)
	if err != nil {
		logger.Error("nmap.NewScanner faield: ", err)
		return err
	}
	result, warnings, err := scanner.Run()
	logger.Debugf("nmap warnings: %s", warnings)
	if err != nil {
		logger.Error("nmap.Run faield: ", err)
		return err
	}
	n.portServiceResultCh <- result
	return nil
}

func (hostCh IPAndPortSeviceCh) GetWebServiceCh() (urlCh chan interface{}) {
	urlCh = make(chan interface{}, 10240)
	go func() {
		for host := range hostCh {
			for i := range host.PortList {
				var protocol string
				if host.PortList[i].Service == "http" {
					protocol = "http"
				} else if host.PortList[i].Service == "ssl" {
					protocol = "https"
				}
				if protocol != "" {
					//IP形式的url
					urlCh <- fmt.Sprintf("%s://%s:%d", protocol, host.IP, host.PortList[i].Port)
				}
			}
		}
		close(urlCh)
	}()
	return urlCh
}

func (n *portServiceDetector) GetPortServiceCh() chan *IPAndPortSevice {
	resultCh := make(chan *IPAndPortSevice, 10240)
	go func() {
		for result := range n.portServiceResultCh {
			for _, host := range result.Hosts {
				portList := make([]Port, 0)
				for _, port := range host.Ports {
					if port.State.State == "closed" {
						continue
					}
					portList = append(portList, Port{
						Port:     int(port.ID),
						Service:  port.Service.Name,
						Version:  port.Service.Version,
						Status:   port.State.State,
						Protocol: port.Protocol,
						Script:   port.Scripts,
					})
					fmt.Printf("port:%+v\n", port)
				}
				resultCh <- &IPAndPortSevice{
					IP:       host.Addresses[0].Addr,
					PortList: portList,
				}
			}
		}
		close(resultCh)
	}()
	return resultCh
}
