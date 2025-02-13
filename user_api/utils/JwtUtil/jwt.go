package JwtUtil

import (
	"LearningGuide/user_api/global"
	"github.com/golang-jwt/jwt"
	"time"
)

type CustomClaims struct {
	ID          uint
	NickName    string
	AuthorityId uint
	jwt.StandardClaims
}

func CreateJWT(ID uint, nickName string, AuthorityId uint) (string, error) {
	expireTime := time.Now().Add(5 * 24 * time.Hour).Unix()
	claims := CustomClaims{
		ID:          ID,
		NickName:    nickName,
		AuthorityId: AuthorityId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime,
			IssuedAt:  time.Now().Unix(),
			Issuer:    "go-user",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(global.ServerConfig.JwtKey))
}
