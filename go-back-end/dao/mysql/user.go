package mysql

import (
	"web_app/models"

	"go.uber.org/zap"
)

//UserLogin 用户登录
func UserLogin(user *models.User) error {
	sql := `SELECT 
	id,username,password 
	FROM user 
	WHERE username=? AND password = md5(md5(?))`
	err := db.Get(user, sql, user.Username, user.Password)
	if err != nil {
		zap.L().Error("mysql.UserLogin db.Get failed", zap.Error(err))
		return err
	}
	return nil
}

//GetUserInfo 获取用户详细信息
func GetUserInfo(id int64) (models.UserInfo, error) {
	var userinfo models.UserInfo
	sql := `SELECT username,avatar FROM user WHERE id=?`
	err := db.Get(&userinfo, sql, id)
	return userinfo, err
}

//IsUsernameExist 查询用户名是否存在
func IsUsernameExist(username string) (bool, error) {
	sql := `SELECT count(id) FROM user WHERE username=?`
	var count int
	err := db.Get(&count, sql, username)
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}
