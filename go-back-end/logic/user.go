package logic

import (
	"errors"
	"web_app/dao/mysql"
	"web_app/dao/redis"
	"web_app/models"
	"web_app/pkg/jwt"

	"go.uber.org/zap"
)

//UserLogin 用户登录
func UserLogin(params models.ParamLogin) (token string, err error) {
	//1.效验用户名的错误次数
	counter := redis.NewCounter(params.Username)
	val, err := counter.Get()
	if err != nil {
		zap.L().Error("counter.Get failed", zap.Error(err))
		return "", err
	}
	if val >= 10 {
		return "", errors.New("用户名密码错误次数30分钟内已超过10次")
	}
	//2.效验用户名是否存在
	exist, err := mysql.IsUsernameExist(params.Username)
	if err != nil {
		return "", err
	}
	if !exist {
		return "", errors.New("用户名不存在")
	}
	//3.查询密码是否正确
	user := models.User{
		Username: params.Username,
		Password: params.Password,
	}
	if err = mysql.UserLogin(&user); err != nil {
		zap.L().Error("mysql.UserLogin failed", zap.Error(err))
		counter.Incr() //增加错误次数
		return "", err
	}
	//4.更新登录时间
	_, err = redis.NewLoginTime().UpdateLoginTime(user.Username)
	if err != nil {
		zap.L().Error("redis.NewLoginTime().UpdateLoginTime failed", zap.Error(err))
		return "", err
	}
	//5.生成token
	return jwt.GenUserToken(user.ID, user.Username)
}

//Logout 退出登录
func Logout(token string) error {
	_, err := redis.NewTokenBlackList().Add(token)
	return err
}

//GetUserInfo 获取用户信息
func GetUserInfo(userid int64) (user models.UserInfo, err error) {
	return mysql.GetUserInfo(userid)
}
