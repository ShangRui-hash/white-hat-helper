package hackflow

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

//ParsedHttpResp 解析http响应结果
type ParsedHttpResp struct {
	StatusCode int
	URL        string
	BaseURL    string
	Method     string
	RespTitle  string
	RespBody   string
	RespHeader http.Header
}

//ParseHttpRespConfig 解析http响应配置
type ParseHttpRespConfig struct {
	RoutineCount int
	HttpRespCh   chan *http.Response
}

//ParseHttpResp 解析http响应
func ParseHttpResp(config *ParseHttpRespConfig) (chan *ParsedHttpResp, error) {
	resultCh := make(chan *ParsedHttpResp, 1024)
	var wg sync.WaitGroup
	for i := 0; i < config.RoutineCount; i++ {
		wg.Add(1)
		go func() {
			for resp := range config.HttpRespCh {
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
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(resultCh)
	}()
	return resultCh, nil
}
