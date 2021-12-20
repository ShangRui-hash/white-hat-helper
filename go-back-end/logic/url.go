package logic

import (
	"context"
	"net/http"
	"web_app/dao/memory"
	"web_app/dao/redis"
	"web_app/param"
	"web_app/pkg/hackflow"

	"go.uber.org/zap"
)

func URLDirScan(param *param.ParamURLDirScan) error {
	ctx, cancel := context.WithCancel(context.Background())
	resultCh, err := hackflow.NewDirSearch(ctx).Run(hackflow.DirSearchConfig{
		URL:                 param.URL,
		FullURL:             true,
		RandomAgent:         true,
		HTTPMethod:          http.MethodGet,
		MinRespContentSize:  2,
		StatusCodeBlackList: "403,404,405,500",
	}).Result()
	if err != nil {
		zap.L().Error("dirsearch run failed", zap.Error(err))
		cancel()
		return err
	}
	//维护一个url和退出函数的映射
	memory.RegisterURLScanCancelFunc(param.URL, cancel)
	//存储扫描到的结果
	redis.SaveFoundURL(resultCh)
	return nil
}
