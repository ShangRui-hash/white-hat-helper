package hackflow

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

//ParseHttpRespConfig 解析http响应配置
type ParseHttpRespConfig struct {
	RoutineCount int
	HttpRespCh   chan *http.Response
}

//ParseHttpResp 解析http响应
func ParseHttpResp(ctx context.Context, config *ParseHttpRespConfig) (chan *ParsedHttpResp, error) {
	resultCh := make(chan *ParsedHttpResp, 1024)
	var wg sync.WaitGroup
	for i := 0; i < config.RoutineCount; i++ {
		wg.Add(1)
		go func() {
		LOOP:
			for {
				select {
				case <-ctx.Done():
					break LOOP
				case resp, ok := <-config.HttpRespCh:
					if !ok {
						break LOOP
					}
					parsedResp := &ParsedHttpResp{
						StatusCode: resp.StatusCode,
						Method:     resp.Request.Method,
						URL:        resp.Request.URL.String(),
						BaseURL:    resp.Request.URL.Scheme + "://" + resp.Request.URL.Host,
						RespHeader: resp.Header,
					}
					body, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						logger.Error("ReadAll failed,err:", err)
						continue
					}
					parsedResp.RespBody = string(body)
					doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
					if err != nil {
						continue
					}
					parsedResp.RespTitle = doc.Find("title").Text()
					logger.Debug("parsedResp.RespTitle:", parsedResp.RespTitle, "parsedResp.URL:", parsedResp.URL)
					resultCh <- parsedResp
					resp.Body.Close()
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
