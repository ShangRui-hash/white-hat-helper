package hackflow

import (
	"sync"

	wappalyzer "github.com/projectdiscovery/wappalyzergo"
)

//DectWhatWebResult 是DectWhatWeb的结果
type DectWhatWebResult struct {
	URL         string
	FingerPrint map[string]struct{}
}

//DecWhatWebConfig 是DectWhatWeb的配置
type DectWhatWebConfig struct {
	TargetCh     chan *ParsedHttpResp
	RoutineCount int
}

// DectWhatWeb 根据响应报文来探测网站的指纹信息
func DectWhatWeb(config *DectWhatWebConfig) (chan *DectWhatWebResult, error) {
	resultCh := make(chan *DectWhatWebResult, 1024)
	wappalyzerClient, err := wappalyzer.New()
	if err != nil {
		logger.Error("wappalyzer.New failed,err:", err)
		return nil, err
	}
	//消费者
	var wg sync.WaitGroup
	for i := 0; i < config.RoutineCount; i++ {
		wg.Add(1)
		go func() {
			for target := range config.TargetCh {
				fingerprints := wappalyzerClient.Fingerprint(target.RespHeader, []byte(target.RespBody))
				resultCh <- &DectWhatWebResult{
					URL:         target.URL,
					FingerPrint: fingerprints,
				}
			}
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(resultCh)
	}()
	return resultCh, nil
}
