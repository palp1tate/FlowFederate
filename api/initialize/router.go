package initialize

import (
	"github.com/palp1tate/FlowFederate/api/global"
	"github.com/palp1tate/FlowFederate/api/middleware"
	"github.com/palp1tate/FlowFederate/api/router"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()
	if !global.ServerConfig.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "ok",
		})
		return
	})
	r.Use(
		middleware.Cors(),
	)
	ApiGroup := r.Group("/api")
	{
		router.InitTrainRouter(ApiGroup)
	}
	return r
}
