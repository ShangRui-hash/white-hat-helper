package hackflow

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
)

func doGetMoreURL(dict io.Reader, url string, outCh chan interface{}) chan interface{} {
	scanner := bufio.NewScanner(dict)
	for scanner.Scan() {
		outCh <- fmt.Sprintf("%s/%s", url, scanner.Text())
	}
	return outCh
}

//GetMoreURL 读取字典，根据基本的url生成更多的url
func GetMoreURL(dict io.Reader, urlCh chan interface{}) chan interface{} {
	outCh := make(chan interface{}, 1024)
	go func() {
		for url := range urlCh {
			doGetMoreURL(dict, url.(string), outCh)
		}
		close(outCh)
	}()
	return outCh
}

var DefaultStatusCodeBlackList string = "400,401,402,403,404,405,500,501,502,503,504"

type BruteForceURLConfig struct {
	BaseURLCh           chan interface{}
	RoutineCount        int
	RandomAgent         bool
	Proxy               string
	StatusCodeBlackList string
	Dictionary          io.Reader
}
type dirSearchGo struct {
	baseTool
}

func NewDirSearchGo(ctx context.Context) *dirSearchGo {
	return &dirSearchGo{
		baseTool{
			ctx: ctx,
		},
	}
}

func (d *dirSearchGo) Run(config *BruteForceURLConfig) (chan *ParsedHttpResp, error) {
	moreURLCh := GetMoreURL(config.Dictionary, config.BaseURLCh)
	requestCh := GenRequest(d.ctx, GenRequestConfig{
		URLCh:       moreURLCh,
		MethodList:  []string{http.MethodGet, http.MethodPost, http.MethodPut},
		RandomAgent: true,
	})
	respCh, err := RetryHttpSend(d.ctx, &RetryHttpSendConfig{
		RequestCh:    requestCh,
		RoutineCount: config.RoutineCount,
		HttpClientConfig: HttpClientConfig{
			Proxy:    config.Proxy,
			Redirect: false,
			Checktry: func(ctx context.Context, resp *http.Response, err error) (bool, error) {
				return false, nil
			},
			RetryMax: 1,
		},
	})
	if err != nil {
		return nil, err
	}
	//解析响应报文
	return ParseHttpResp(d.ctx, &ParseHttpRespConfig{
		RoutineCount: 1000,
		HttpRespCh:   respCh,
	})
}
