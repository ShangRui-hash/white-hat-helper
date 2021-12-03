package jwt

import (
	"errors"
	"time"
	"web_app/dao/redis"

	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
)

// PriKey 用来对token签名的私钥（一定要妥善保管）
var (
	PriKey          = []byte("夏天夏天悄悄过去")
	AdminIssuer     = "admin"
	UserIssuer      = "user"
	ErrInvalidToken = errors.New("invalid token")
)

//TokenExpireDuration token的过期时间
const TokenExpireDuration = time.Hour * 100

// MyClaims 自定义声明结构体并内嵌jwt.StandardClaims
// jwt包自带的jwt.StandardClaims只包含了官方字段
// 我们这里需要额外记录一个username字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
type MyClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenToken 生成JWT
func GenToken(userID int64, username string, issuer string) (string, error) {
	// 创建一个我们自己的声明
	c := MyClaims{
		userID,
		username, // 自定义字段
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), // 过期时间
			Issuer:    issuer,                                     // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(PriKey)
}

//GenAdminToken 生成管理员用的token
func GenAdminToken(userID int64, username string) (string, error) {
	return GenToken(userID, username, AdminIssuer)
}

//GenUserToken 生成用户使用的token
func GenUserToken(userID int64, username string) (string, error) {
	return GenToken(userID, username, UserIssuer)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	var mc = new(MyClaims)
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, mc,
		func(token *jwt.Token) (i interface{}, err error) {
			return PriKey, nil
		})
	if err != nil {
		return nil, err
	}

	if token.Valid { // 校验token
		return mc, nil
	}
	return nil, ErrInvalidToken
}

//IsValidToken 效验token是否合法
func IsValidToken(token string, issuer string) (*MyClaims, bool) {
	//检查token是否在黑名单中
	exist, err := redis.NewTokenBlackList().IsMember(token)
	if err != nil {
		zap.L().Error("redis.NewTokenBlackList().IsMember(token) failed", zap.Error(err))
		return nil, false
	}
	if exist {
		return nil, false
	}
	//解析token
	mc, err := ParseToken(token)
	if err != nil || mc.Issuer != issuer {
		return nil, false
	}
	return mc, true
}
