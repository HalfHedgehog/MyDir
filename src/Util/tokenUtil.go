package Util

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"httpServer/src/global"
	"strconv"
	"time"
)

type Token struct {
	UserId string `json:"userId"`
	jwt.StandardClaims
}

// CreateToken 创建token
func CreateToken(userId int64) string {
	//加密
	c := Token{
		UserId: strconv.FormatInt(userId, 10),
		StandardClaims: jwt.StandardClaims{
			//什么时候生效，现在生效
			NotBefore: time.Now().Unix(),
			//什么时候失效，半个小时后
			ExpiresAt: time.Now().Unix() + global.Config.JWT.ExpiresTime,
			//签发人
			Issuer: global.Config.JWT.Issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	fmt.Println(token)
	newToken, e := token.SignedString([]byte(global.Config.JWT.Key))
	if e != nil {
		fmt.Println(e)
		return ""
	}
	return newToken
}

// DecryToken 解密
func DecryToken(newToken string) (string, error) {
	//解密,在这一步如果出现token过期的问题则会报错
	token, err := jwt.ParseWithClaims(newToken, &Token{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.Config.JWT.Key), nil
	})
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	id := token.Claims.(*Token).UserId
	return id, nil
}
