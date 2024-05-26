package router

import (
	"github.com/palp1tate/FlowFederate/api/handler"

	"github.com/gin-gonic/gin"
)

func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("/user")
	{
		UserRouter.POST("/register", handler.Register)
		//	UserRouter.POST("/login", handler.UserLogin)
		//	UserRouter.GET("/get_user", middleware.JWTAuth(), handler.GetUser)
		//	UserRouter.PUT("/reset_password", handler.ResetPassword)
		//	UserRouter.PUT("/update_user", middleware.JWTAuth(), handler.UpdateUser)
		//}
	}
}
