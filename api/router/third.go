package router

import (
	"github.com/palp1tate/FlowFederate/api/handler"
	"github.com/palp1tate/FlowFederate/api/middleware"

	"github.com/gin-gonic/gin"
)

func InitThirdPartyRouter(Router *gin.RouterGroup) {
	ThirdPartyRouter := Router.Group("/third_party")
	{
		ThirdPartyRouter.GET("/get_captcha", handler.GetPicCaptcha)
		ThirdPartyRouter.POST("/send_sms", handler.SendSms)
		ThirdPartyRouter.POST("upload_file", middleware.JWTAuth(), handler.UploadFile)
		ThirdPartyRouter.DELETE("delete_file", middleware.JWTAuth(), handler.DeleteFile)
	}
}
