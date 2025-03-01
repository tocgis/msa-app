package jwtutil

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

var JwtSecret = []byte("jwtSecret_v3")

const JWT_CONTEXT_KEY = "jwt_context_key"

type Token struct {
	Name string
	Uid  int
	jwt.StandardClaims
}

func CreateJwtToken(name string, id int) (string, error) {
	var token Token
	token.StandardClaims = jwt.StandardClaims{
		Audience:  "",                                             // 受众群体
		ExpiresAt: time.Now().Add(3600 * 24 * time.Second).Unix(), // 到期时间
		Id:        "",                                             // 编号
		IssuedAt:  time.Now().Unix(),                              // 签发时间
		Issuer:    "kit_v3",                                       // 签发人
		NotBefore: time.Now().Unix(),                              // 生效时间
		Subject:   "login",                                        // 主题
	}
	token.Name = name
	token.Uid = id
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, token)
	return tokenClaims.SignedString(JwtSecret)
}

func ParseToken(token string) (jwt.MapClaims, error) {
	jwtToken, err := jwt.ParseWithClaims(token, jwt.MapClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return JwtSecret, nil
	})
	if err != nil || jwtToken == nil {
		return nil, err
	}
	claim, ok := jwtToken.Claims.(jwt.MapClaims)
	if ok && jwtToken.Valid {
		return claim, nil
	} else {
		return nil, nil
	}
}
