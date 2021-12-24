package redis

import (
	"web_app/pkg/hackflow"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

//CalculateIPScore 计算IP的分数
func CalculateHostScore(ip string) int {
	score := 0 //每次计算前清空
	//1.获取ip对应的端口号列表
	portList, err := rdb.HKeys(IPPortDetailKeyPrefix + ip).Result()
	if err != nil && err != redis.Nil {
		zap.L().Error("rdb.HKeys failed,err:", zap.Error(err))
		return score
	}
	for _, port := range portList {
		switch port {
		case "22":
			score += 50
		case "23":
			score += 40
		default:
			score += 10
		}
	}
	//2.获取ip对应的web服务列表
	webs, err := GetAllWebServiceByIP(ip)
	if err != nil && err != redis.Nil {
		zap.L().Error("GetWebServiceByIP failed,err:", zap.Error(err))
		return score
	}
	for i := range webs {
		switch webs[i].StatusCode {
		case 200:
			score += 30
		case 301, 302:
			score += 20
		default:
			score += 10
		}
		for j := range webs[i].FingerPrint {
			switch webs[i].FingerPrint[j] {
			case "PHP":
				score += 30
			case "ThinkPHP":
				score += 40
			default:
				score += 10
			}
		}
		score += len(webs[i].Dirs) * 20
	}
	return score
}

//UpdateIPScore 更新公司的ip集合中指定ip的分数
func UpdateIPScore(companyID int64, ip string) error {
	score := CalculateHostScore(ip)
	if _, err := rdb.ZAdd(GetCompanyIPZSetKey(companyID), redis.Z{Score: float64(score), Member: ip}).Result(); err != nil {
		zap.L().Error("rdb.ZAdd failed,err:", zap.Error(err))
		return err
	}
	return nil
}

func CalculateURLScore(url string, resp *hackflow.ParsedHttpResp) int {
	score := 0
	switch resp.StatusCode {
	case 200:
		score += 30
	case 301, 302:
		score += 20
	default:
		score += 10
	}
	if len(resp.RespTitle) > 0 {
		score += 10
	}
	return score
}
