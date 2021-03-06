package memory

import (
	"context"
	"errors"
)

var urlScanCancelFuncMap map[string]context.CancelFunc

func init() {
	urlScanCancelFuncMap = make(map[string]context.CancelFunc)
}

func RegisterURLScanCancelFunc(url string, f context.CancelFunc) {
	urlScanCancelFuncMap[url] = f
}

func StopURLScan(url string) error {
	f, ok := urlScanCancelFuncMap[url]
	if !ok {
		return errors.New("no cancel func")
	}
	f()
	delete(urlScanCancelFuncMap, url)
	return nil
}

func IsURLScanRunning(url string) bool {
	_, ok := urlScanCancelFuncMap[url]
	return ok
}
