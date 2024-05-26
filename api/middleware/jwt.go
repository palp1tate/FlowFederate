package middleware

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/palp1tate/FlowFederate/api/global"
	"github.com/palp1tate/FlowFederate/api/internal/errorx"
	"github.com/palp1tate/FlowFederate/api/internal/utils"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": errorx.ErrTokenNeed,
				"msg":  errorx.GetMsg(errorx.ErrTokenNeed),
			})
			c.Abort()
			return
		}
		j := utils.NewJWT(global.ServerConfig.JWT.SigningKey, global.ServerConfig.JWT.Expiration)
		claims, err := j.ParseToken(token)
		if err != nil {
			if errors.Is(err, errorx.TokenExpired) {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code": errorx.ErrTokenExpired,
					"msg":  errorx.GetMsg(errorx.ErrTokenExpired),
				})
				c.Abort()
				return
			}

			c.JSON(http.StatusUnauthorized, gin.H{
				"code": errorx.ErrTokenParseFailed,
				"msg":  errorx.GetMsg(errorx.ErrTokenParseFailed),
			})
			c.Abort()
			return
		}
		id, _ := strconv.Atoi(claims.ID)
		role, _ := strconv.Atoi(claims.Role)
		c.Set("id", id)
		c.Set("role", role)
		c.Set("token", token)
		c.Next()
	}
}
