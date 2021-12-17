package hackflow

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/serkanalgur/phpfuncs"
	"github.com/sirupsen/logrus"
)

func doGetMoreURL(dict io.Reader, url string) chan interface{} {
	outCh := make(chan interface{}, 1024)
	go func() {
		scanner := bufio.NewScanner(dict)
		for scanner.Scan() {
			outCh <- fmt.Sprintf("%s/%s", url, scanner.Text())
		}
		close(outCh)
	}()
	return outCh
}

//GetMoreURL 读取字典，根据基本的url生成更多的url
func GetMoreURL(dict io.Reader, urlCh chan interface{}) chan interface{} {
	stream := NewStream()
	for url := range urlCh {
		stream.AddSrc(doGetMoreURL(dict, url.(string)))
	}
	return stream.SetDstCount(1).GetDst()[0]
}

type BruteForceURLConfig struct {
	BaseURLCh           chan interface{}
	RoutineCount        int
	RandomAgent         bool
	Proxy               string
	Dictionary          io.Reader
	StatusCodeBlackList []int
}

type BruteForceURLResult struct {
	Method     string
	ParentURL  string
	URL        string
	Location   string
	Title      string
	StatusCode int
}

func BruteForceURL(config *BruteForceURLConfig) (chan *BruteForceURLResult, error) {
	moreURLCh := GetMoreURL(config.Dictionary, config.BaseURLCh)
	requestCh := GenRequest(GenRequestConfig{
		URLCh:       moreURLCh,
		MethodList:  []string{http.MethodGet, http.MethodPost, http.MethodPut},
		RandomAgent: true,
	})
	respCh, err := RetryHttpSend(&RetryHttpSendConfig{
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
	parsedRespCh, err := ParseHttpResp(&ParseHttpRespConfig{
		RoutineCount: 1000,
		HttpRespCh:   respCh,
	})
	if err != nil {
		logrus.Error("parseHttpResp failed,err:", err)
		return nil, err
	}
	outCh := make(chan *BruteForceURLResult, 10240)
	go func() {
		for resp := range parsedRespCh {
			if phpfuncs.InArray(resp.StatusCode, config.StatusCodeBlackList) {
				continue
			}
			outCh <- &BruteForceURLResult{
				ParentURL:  resp.BaseURL,
				URL:        resp.URL,
				Location:   resp.RespHeader.Get("Location"),
				StatusCode: resp.StatusCode,
				Method:     resp.Method,
				Title:      resp.RespTitle,
			}
		}
		close(outCh)
	}()
	return outCh, nil
}
