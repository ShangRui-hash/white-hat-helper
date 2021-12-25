package hackflow

import (
	"context"
	"sync"

	wappalyzer "github.com/projectdiscovery/wappalyzergo"
)

//DectWhatWebResult 是DectWhatWeb的结果
type DectWhatWebResult struct {
	URL         string
	FingerPrint map[string]struct{}
}

type DectWhatWebResultCh chan *DectWhatWebResult

func (d DectWhatWebResultCh) GetURLCh() chan interface{} {
	urlCh := make(chan interface{}, 1024)
	go func() {
		for item := range d {
			urlCh <- item.URL
		}
		close(urlCh)
	}()
	return urlCh
}

type whatWeb struct {
	baseTool
}

func NewWhatWeb(ctx context.Context) *whatWeb {
	return &whatWeb{
		baseTool: baseTool{
			ctx:  ctx,
			name: "whatweb",
			desp: "探测网站的指纹信息",
		},
	}
}

//DecWhatWebConfig 是DectWhatWeb的配置
type DectWhatWebConfig struct {
	*BaseConfig
	TargetCh     chan *ParsedHttpResp
	RoutineCount int
}

// DectWhatWeb 根据响应报文来探测网站的指纹信息
func (w *whatWeb) Run(config *DectWhatWebConfig) (chan *DectWhatWebResult, error) {
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
		LOOP:
			for {
				select {
				case <-w.ctx.Done():
					if config.CallAfterCtxDone != nil {
						config.CallAfterCtxDone(w)
					}
					break LOOP
				case target, ok := <-config.TargetCh:
					if !ok {
						break LOOP
					}
					fingerprints := wappalyzerClient.Fingerprint(target.RespHeader, []byte(target.RespBody))
					resultCh <- &DectWhatWebResult{
						URL:         target.URL,
						FingerPrint: fingerprints,
					}
				}
			}
			wg.Done()
		}()
	}
	if config.CallAfterBegin != nil {
		config.CallAfterBegin(w)
	}
	go func() {
		wg.Wait()
		close(resultCh)
		if config.CallAfterComplete != nil {
			config.CallAfterComplete(w)
		}
	}()
	return resultCh, nil
}
