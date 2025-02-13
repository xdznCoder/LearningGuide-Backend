package middlewares

import (
	"LearningGuide/user_api/global"
	customjwt "LearningGuide/user_api/utils/JwtUtil"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
	"net/http"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": "Token 为空",
			})
			c.Abort()
			return
		}
		token, err := jwt.ParseWithClaims(authHeader, &customjwt.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(global.ServerConfig.JwtKey), nil
		})
		if err != nil {
			switch err.(*jwt.ValidationError).Errors {
			case jwt.ValidationErrorExpired:
				c.JSON(http.StatusUnauthorized, gin.H{
					"msg": "Token 已过期",
				})
				c.Abort()
				return
			default:
				zap.S().Infof("Parse JWT Failed: %v", err)
				c.JSON(http.StatusUnauthorized, gin.H{
					"msg": "无效 Token",
				})
				c.Abort()
				return
			}
		}
		c.Set("claims", token.Claims.(*customjwt.CustomClaims))
		c.Set("userId", token.Claims.(*customjwt.CustomClaims).ID)
		return
	}
}

func IsAdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := c.Get("claims")
		if ok && user.(*customjwt.CustomClaims).AuthorityId != 1 {
			return
		}
		c.JSON(http.StatusForbidden, gin.H{
			"msg": "权限不足",
		})
		c.Abort()
	}
}
