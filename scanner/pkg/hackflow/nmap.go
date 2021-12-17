package hackflow

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/Ullaakut/nmap"
	"github.com/sirupsen/logrus"
)

//UNKNOWN_OS 未知操作系统
const UNKNOWN_OS = "unknown"

type Port struct {
	Port     int           `json:"port"`
	Service  string        `json:"service"`
	Status   string        `json:"status"`
	Version  string        `json:"version"`
	Protocol string        `json:"protocol"`
	Script   []nmap.Script `json:"script"`
}

func (p *Port) String() (string, error) {
	b, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(b), err
}

type nmapV2 struct {
	baseTool
	osResultCh          chan *nmap.Run
	portServiceResultCh chan *nmap.Run
}

func newNmap() Tool {
	return &nmapV2{
		baseTool: baseTool{
			name: "nmap",
			desp: "端口扫描、服务识别、操作系统识别",
		},
		osResultCh:          make(chan *nmap.Run, 1024),
		portServiceResultCh: make(chan *nmap.Run, 1024),
	}
}

func GetNmap() *nmapV2 {
	return container.Get(NMAP).(*nmapV2)
}

type OSDectionConfig struct {
	HostCh    chan interface{}
	Timeout   time.Duration
	BatchSize int
}

func (n *nmapV2) OSDection(config *OSDectionConfig) *nmapV2 {
	go func() {
		count := 0
		var wg sync.WaitGroup
	LOOP:
		for {
			if count < config.BatchSize {
				target, ok := <-config.HostCh
				if !ok {
					break LOOP
				}
				count++
				wg.Add(1)
				go func() {
					defer wg.Done()
					n.doOSDection(target.(string), config.Timeout)
					count--
				}()
			}
		}
		wg.Wait()
		close(n.osResultCh)
	}()
	return n
}

func (n *nmapV2) doOSDection(target string, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
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
	logger.Debugf("nmap warnings: %s", warnings)
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
	go func() {
		count := 0
		var wg sync.WaitGroup
	LOOP:
		for {
			if count < config.BatchSize {
				target, ok := <-config.TargetCh
				if !ok {
					break LOOP
				}
				count++
				wg.Add(1)
				go func() {
					defer wg.Done()
					if err := n.doServiceDection(target, config.Timeout); err != nil {
						logger.Error("nmapV2 run faield: ", err)
					}
					count--
				}()
			}
		}
		wg.Wait()
		close(n.portServiceResultCh)
	}()
	return n
}

func (n *nmapV2) doServiceDection(target *IPAndPort, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
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

type IPAndOS struct {
	IP string `json:"ip"`
	OS string `json:"os"`
}
type IPAndOSCh chan *IPAndOS

func (i IPAndOSCh) GetIPCh() chan interface{} {
	IPCh := make(chan interface{}, 1024)
	go func() {
		for ipAndOs := range i {
			IPCh <- ipAndOs.IP
		}
		close(IPCh)
	}()
	return IPCh
}

func (n *nmapV2) GetIPAndOSCh() IPAndOSCh {
	ipAndOsCh := make(chan *IPAndOS, 1024)
	go func() {
		for result := range n.osResultCh {
			for _, host := range result.Hosts {
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

type IPAndPortSevice struct {
	IP       string `json:"ip"`
	PortList []Port `json:"ports"`
}

type IPAndPortSeviceCh chan *IPAndPortSevice

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

// func (n *nmapV2) Print() chan HostListItem {
// 	resultCh := make(chan HostListItem, 10240)
// 	go func() {
// 		for result := range n.resultCh {
// 			fmt.Printf("%+v\n", result)
// 			resultCh <- result
// 		}
// 		close(resultCh)
// 	}()
// 	return resultCh
// }
