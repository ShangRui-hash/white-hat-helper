package hackflow

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/Ullaakut/nmap"
)

//UNKNOWN_OS 未知操作系统
const UNKNOWN_OS = "unknown"

type HostListItem struct {
	IP       string `json:"ip"`
	OS       string `json:"os"`
	PortList []Port `json:"ports"`
}

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
	resultCh chan HostListItem
}

func newNmap() Tool {
	return &nmapV2{
		baseTool: baseTool{
			name: "nmap",
			desp: "端口扫描、服务识别、操作系统识别",
		},
		resultCh: make(chan HostListItem, 10240),
	}
}

func GetNmap() *nmapV2 {
	return container.Get(NMAP).(*nmapV2)
}

type NmapRunConfig struct {
	TargetCh  chan string
	Timeout   time.Duration
	BatchSize int
}

func (n *nmapV2) Run(config *NmapRunConfig) (*nmapV2, error) {
	batchSize := config.BatchSize
	batchTargetCh := make(chan []string, 1024)
	batch := make([]string, 0)
	go func() {
		for target := range config.TargetCh {
			batch = append(batch, target)
			if len(batch) == batchSize {
				batchTargetCh <- batch
				batch = make([]string, 0)
			}
		}
		if len(batch) > 0 {
			batchTargetCh <- batch
		}
		close(batchTargetCh)
	}()
	go func() {
		for targetList := range batchTargetCh {
			var wg sync.WaitGroup
			for _, target := range targetList {
				wg.Add(1)
				go func(target string) {
					defer wg.Done()
					if err := n.run(target, config.Timeout); err != nil {
						logger.Error("nmapV2 run faield: ", err)
					}
				}(target)
			}
			wg.Wait()
		}
		close(n.resultCh)
	}()
	return n, nil
}

func (n *nmapV2) run(target string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	logger.Debug("nmap run:", target)
	scanner, err := nmap.NewScanner(
		nmap.WithTargets(target),
		nmap.WithServiceInfo(),
		nmap.WithOSDetection(),
		nmap.WithContext(ctx),
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
		os := UNKNOWN_OS
		if len(host.OS.Matches) > 0 {
			os = host.OS.Matches[0].Name
		}
		n.resultCh <- HostListItem{
			IP:       target,
			OS:       os,
			PortList: portList,
		}
	}
	return nil
}

func (n *nmapV2) Result() chan HostListItem {
	return n.resultCh
}

func (n *nmapV2) Print() chan HostListItem {
	resultCh := make(chan HostListItem, 10240)
	go func() {
		for result := range n.resultCh {
			fmt.Printf("%+v\n", result)
			resultCh <- result
		}
		close(resultCh)
	}()
	return resultCh
}
