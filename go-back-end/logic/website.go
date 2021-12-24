package logic

import (
	"context"
	"os"
	"web_app/dao/memory"
	"web_app/dao/redis"
	"web_app/models"
	"web_app/param"
	"web_app/pkg/hackflow"
	"web_app/settings"

	"go.uber.org/zap"
)

func GetWebSiteList(companyID int64, offset, count int) ([]*models.WebDetail, error) {
	return redis.GetWebServiceByCompanyID(companyID, offset, count)
}

func StartURLDirScan(param *param.ParamStartURLDirScan) error {
	ctx, cancel := context.WithCancel(context.Background())
	dict, err := os.Open(settings.Conf.DictPath)
	if err != nil {
		zap.L().Error("open dirsearch.txt failed,err:", zap.Error(err))
		return err
	}
	urlCh := make(chan interface{}, 1)
	urlCh <- param.URL
	close(urlCh)
	respCh, err := hackflow.BruteForceURL(ctx, &hackflow.BruteForceURLConfig{
		BaseURLCh:           urlCh,
		RoutineCount:        100,
		Proxy:               settings.Conf.Proxy,
		Dictionary:          dict,
		RandomAgent:         true,
		StatusCodeBlackList: hackflow.DefaultStatusCodeBlackList,
	})
	if err != nil {
		zap.L().Error("burte force url failed,err:", zap.Error(err))
		return err
	}
	//维护一个url和退出函数的映射
	memory.RegisterURLScanCancelFunc(param.URL, cancel)
	//存储扫描到的结果
	redis.SaveHttpResp(respCh, 0)
	return nil
}

func StopURLDirScan(param *param.ParamStopURLDirScan) error {
	return memory.StopURLScan(param.URL)
}

//DeleteURLSubDir 删除URL的指定子目录
func DeleteURLSubDir(param *param.ParamDeleteURLSubDir) error {
	return redis.DeleteURLSubDir(param.ParentURL, param.SubURL)
}

func GetSubDir(url string, offset, count int) ([]hackflow.BruteForceURLResult, error) {
	return redis.GetURLByParentURL(url, offset, count)
}
