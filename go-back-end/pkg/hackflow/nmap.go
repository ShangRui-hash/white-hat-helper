package hackflow

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Ullaakut/nmap"
	"github.com/sirupsen/logrus"
)

type nmapV2 struct {
	baseTool
	osResultCh          chan *nmap.Run
	portServiceResultCh chan *nmap.Run
}

func NewNmap(ctx context.Context) *nmapV2 {
	return &nmapV2{
		baseTool: baseTool{
			ctx:  ctx,
			name: "nmap",
			desp: "端口扫描、服务识别、操作系统识别",
		},
		osResultCh:          make(chan *nmap.Run, 1024),
		portServiceResultCh: make(chan *nmap.Run, 1024),
	}
}

type OSDectionConfig struct {
	HostCh    chan interface{}
	Timeout   time.Duration
	BatchSize int
}

func (n *nmapV2) OSDection(config *OSDectionConfig) *nmapV2 {
	var count int32 = 0
	var wg sync.WaitGroup
	go func() {
	LOOP:
		for {
			if count < int32(config.BatchSize) {
				select {
				case <-n.ctx.Done():
					break LOOP
				case target, ok := <-config.HostCh:
					if !ok {
						break LOOP
					}
					atomic.AddInt32(&count, 1)
					wg.Add(1)
					go func() {
						defer wg.Done()
						n.doOSDection(target.(string), config.Timeout)
						atomic.AddInt32(&count, -1)
					}()
				}
			}
		}
		//等待所有进程干完后
		wg.Wait()
		close(n.osResultCh)
	}()
	return n
}

func (n *nmapV2) doOSDection(target string, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(n.ctx, timeout)
	defer cancel()
	logger.Debug("nmap run:", target)
	scanner, err := nmap.NewScanner(
		nmap.WithTargets(target),
		nmap.WithOSDetection(),
		nmap.WithContext(ctx),
		nmap.WithSkipHostDiscovery(), // -Pn
	)
	if err != nil {
		logger.Error("nmap.NewScanner faield: ", err)
		return
	}
	result, warnings, err := scanner.Run()
	if len(warnings) > 0 {
		logger.Debugf("nmap warnings: %s", warnings)
	}
	if err != nil {
		logger.Error("nmap.Run faield: ", err)
		return
	}
	n.osResultCh <- result
}

type ServiceDectionConfig struct {
	TargetCh  chan *IPAndPort
	Timeout   time.Duration
	BatchSize int
}

func (n *nmapV2) ServiceDection(config *ServiceDectionConfig) *nmapV2 {
	var count int32 = 0
	var wg sync.WaitGroup
	lackCh := make(chan struct{}, 1024)
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
						logger.Error("nmapV2 run faield: ", err)
					}
					atomic.AddInt32(&count, -1)
				}()
			}
		}
		wg.Wait()
		close(n.portServiceResultCh)
	}()
	return n
}

func (n *nmapV2) doServiceDection(target *IPAndPort, timeout time.Duration) error {
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

func (n *nmapV2) GetIPAndOSCh() IPAndOSCh {
	ipAndOsCh := make(chan *IPAndOS, 1024)
	go func() {
		for result := range n.osResultCh {
			for _, host := range result.Hosts {
				logger.Debugf("os:%+v\n", host.OS)
				os := UNKNOWN_OS
				if len(host.OS.Matches) > 0 {
					os = host.OS.Matches[0].Name
				}
				ipAndOsCh <- &IPAndOS{
					IP: host.Addresses[0].Addr,
					OS: os,
				}
			}
		}
		close(ipAndOsCh)
	}()
	return ipAndOsCh
}

func (hostCh IPAndPortSeviceCh) GetWebServiceCh() (urlCh chan interface{}) {
	urlCh = make(chan interface{}, 10240)
	go func() {
		for host := range hostCh {
			for _, port := range host.PortList {
				var protocol string
				if port.Service == "http" {
					protocol = "http"
				} else if port.Service == "ssl" {
					protocol = "https"
				}
				if protocol != "" {
					//IP形式的url
					urlCh <- fmt.Sprintf("%s://%s:%d", protocol, host.IP, port.Port)
					logrus.Debug("web service:", fmt.Sprintf("%s://%s:%d", protocol, host.IP, port.Port))
				}
			}
		}
		close(urlCh)
	}()
	return urlCh
}

func (n *nmapV2) GetPortServiceCh() chan *IPAndPortSevice {
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
	}()
	return resultCh
}
