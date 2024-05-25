package middleware

import (
	"net/http"

	"github.com/palp1tate/FlowFederate/api/internal/consts"
	"github.com/palp1tate/FlowFederate/api/internal/errorx"

	"github.com/gin-gonic/gin"
)

func AdminAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role := ctx.GetInt("role")
		if role != consts.Admin {
			ctx.JSON(http.StatusForbidden, gin.H{
				"code": errorx.ErrMustAdmin,
				"msg":  errorx.GetMsg(errorx.ErrMustAdmin),
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
