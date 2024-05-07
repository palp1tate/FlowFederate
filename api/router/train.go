package router

import (
	"github.com/palp1tate/FlowFederate/api/handler"

	"github.com/gin-gonic/gin"
)

func InitTrainRouter(Router *gin.RouterGroup) {
	TrainRouter := Router.Group("/train")
	{
		TrainRouter.POST("/start", handler.Train)
		TrainRouter.GET("/list", handler.GetTaskList)
		TrainRouter.GET("/detail", handler.GetTask)
	}
}
