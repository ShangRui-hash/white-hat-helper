package hackflow

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Ullaakut/nmap"
	"github.com/sirupsen/logrus"
)

type osDetector struct {
	baseTool
	osResultCh chan *nmap.Run
}

func NewOSDetector(ctx context.Context) *osDetector {
	return &osDetector{
		baseTool: baseTool{
			ctx:  ctx,
			name: "os-detector",
			desp: "操作系统识别",
		},
		osResultCh: make(chan *nmap.Run, 1024),
	}
}

type OSDectionConfig struct {
	*BaseConfig
	HostCh    chan interface{}
	Timeout   time.Duration
	BatchSize int
}

func (o *osDetector) Run(config *OSDectionConfig) *osDetector {
	var count int32 = 0
	var wg sync.WaitGroup
	go func() {
	LOOP:
		for {
			if count < int32(config.BatchSize) {
				select {
				case <-o.ctx.Done():
					if config.CallAfterCtxDone != nil {
						config.CallAfterCtxDone(o)
					}
					break LOOP
				case target, ok := <-config.HostCh:
					if !ok {
						break LOOP
					}
					atomic.AddInt32(&count, 1)
					wg.Add(1)
					go func() {
						defer wg.Done()
						o.doOSDection(target.(string), config.Timeout)
						atomic.AddInt32(&count, -1)
					}()
				}
			}
		}
		//等待所有进程干完后
		wg.Wait()
		close(o.osResultCh)
		if config.CallAfterComplete != nil {
			config.CallAfterComplete(o)
		}
	}()
	if config.CallAfterBegin != nil {
		config.CallAfterBegin(o)
	}
	return o
}

func (o *osDetector) doOSDection(target string, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(o.ctx, timeout)
	defer cancel()
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
	o.osResultCh <- result
}

func (o *osDetector) GetIPAndOSCh() IPAndOSCh {
	ipAndOsCh := make(chan *IPAndOS, 1024)
	go func() {
		for result := range o.osResultCh {
			for i := range result.Hosts {
				os := UNKNOWN_OS
				if len(result.Hosts[i].OS.Matches) > 0 {
					os = result.Hosts[i].OS.Matches[0].Name
				}
				ipAndOsCh <- &IPAndOS{
					IP: result.Hosts[i].Addresses[0].Addr,
					OS: os,
				}
			}
		}
		close(ipAndOsCh)
		logrus.Error("ipAndOsCh closed")
	}()
	return ipAndOsCh
}
