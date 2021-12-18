package hackflow

import (
	"encoding/json"

	"github.com/Ullaakut/nmap"
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

type IPAndPort struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}
type IPAndPortCh chan *IPAndPort

func (i IPAndPortCh) GetIPCh() chan interface{} {
	IPCh := make(chan interface{}, 1024)
	go func() {
		for IPAndPort := range i {
			IPCh <- IPAndPort.IP
		}
		close(IPCh)
	}()
	return IPCh
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

type IPAndPortSevice struct {
	IP       string `json:"ip"`
	PortList []Port `json:"ports"`
}

type IPAndPortSeviceCh chan *IPAndPortSevice
