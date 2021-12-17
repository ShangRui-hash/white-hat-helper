package hackflow

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
