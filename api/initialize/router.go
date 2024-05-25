package initialize

import (
	"github.com/palp1tate/FlowFederate/api/global"
	"github.com/palp1tate/FlowFederate/api/internal/utils"
	"github.com/palp1tate/FlowFederate/api/middleware"
	"github.com/palp1tate/FlowFederate/api/router"

	"github.com/gin-gonic/gin"
	sentinel "github.com/sentinel-group/sentinel-go-adapters/gin"
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
	})
	r.Use(
		middleware.Cors(),
		sentinel.SentinelMiddleware(sentinel.WithResourceExtractor(utils.ResourceExtractor)),
	)
	ApiGroup := r.Group("/api")
	{
		router.InitThirdPartyRouter(ApiGroup)
		router.InitUserRouter(ApiGroup)
		router.InitAdminRouter(ApiGroup)
		router.InitTrainRouter(ApiGroup)
	}
	return r
}
