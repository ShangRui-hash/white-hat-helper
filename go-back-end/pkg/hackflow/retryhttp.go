package hackflow

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/url"
	"sync"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/sirupsen/logrus"
)

type HttpClientConfig struct {
	Logger   interface{}
	Proxy    string
	Redirect bool
	RetryMax int
	Checktry retryablehttp.CheckRetry
}

//RetryHttpSendConfig 用来配置RetryHttpSend
type RetryHttpSendConfig struct {
	RequestCh    chan *http.Request
	RoutineCount int
	HttpClientConfig
}

func NewHttpClient(config *HttpClientConfig) (*http.Client, error) {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = config.RetryMax
	if config.Checktry != nil {
		retryClient.CheckRetry = config.Checktry
	}
	if config.Logger == nil {
		retryClient.Logger = logrus.StandardLogger()
	} else {
		retryClient.Logger = config.Logger
	}
	tr := &http.Transport{
		//跳过证书验证
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	//代理
	if config.Proxy != "" {
		proxy, err := url.Parse(config.Proxy)
		if err != nil {
			logrus.Error("url.Parse failed,err:", err)
			return nil, err
		}
		tr.Proxy = http.ProxyURL(proxy)
	}
	retryClient.HTTPClient.Transport = tr
	//不进行重定向，便于获取http头中 location 字段
	if !config.Redirect {
		retryClient.HTTPClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}
	return retryClient.StandardClient(), nil
}

//RetryHttpSend 用来发送http请求，如果发送失败，会自动重试
func RetryHttpSend(ctx context.Context, config *RetryHttpSendConfig) (chan *http.Response, error) {
	resultCh := make(chan *http.Response, 1024)
	client, err := NewHttpClient(&config.HttpClientConfig)
	if err != nil {
		logrus.Error("NewHttpClient failed,err:", err)
		return nil, err
	}
	var wg sync.WaitGroup
	for i := 0; i < config.RoutineCount; i++ {
		wg.Add(1)
		go func(ctx context.Context) {
		LOOP:
			for {
				select {
				case <-ctx.Done():
					break LOOP
				case request, ok := <-config.RequestCh:
					if !ok {
						break LOOP
					}
					resp, err := client.Do(request)
					if err != nil {
						logger.Errorf("client.Do(%s) failded,err:%v", request.URL, err)
						continue
					}
					resultCh <- resp
				}
			}
			wg.Done()
		}(ctx)
	}
	go func() {
		wg.Wait()
		close(resultCh)
	}()
	return resultCh, nil
}
