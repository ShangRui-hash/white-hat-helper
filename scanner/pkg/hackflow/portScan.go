package hackflow

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type portScanner struct {
	dialer      net.Dialer
	ipAndPortCh chan *IPAndPort
}

func NewPortScanner(timeout time.Duration) *portScanner {
	return &portScanner{
		dialer:      net.Dialer{Timeout: timeout},
		ipAndPortCh: make(chan *IPAndPort, 1024),
	}
}

type ScanConfig struct {
	HostCh       chan interface{}
	RoutineCount int
}

func genTargets(HostCh chan interface{}) chan *IPAndPort {
	targets := make(chan *IPAndPort, 1024)
	go func() {
		for host := range HostCh {

			targets <- &IPAndPort{
				IP:   host.(string),
				Port: 80,
			}

		}
		close(targets)
	}()
	return targets
}

//ConnectScan TCP全连接扫描
func (p *portScanner) ConnectScan(config *ScanConfig) IPAndPortCh {
	targetCh := genTargets(config.HostCh)
	var wg sync.WaitGroup
	for i := 0; i < config.RoutineCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for host := range targetCh {
				logrus.Info("connect scan: ", host)
				p.doConnectScan(host.IP, host.Port)
			}
		}()
	}
	go func() {
		wg.Wait()
		close(p.ipAndPortCh)
	}()
	return p.ipAndPortCh
}

func (p *portScanner) doConnectScan(host string, port int) {
	con, err := p.dialer.Dial("tcp4", fmt.Sprintf(`%s:%d`, host, port))
	if err == nil { //连接成功
		con.Close()
		p.ipAndPortCh <- &IPAndPort{
			IP:   host,
			Port: port,
		}
	}
}
