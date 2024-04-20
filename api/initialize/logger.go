package initialize

import (
	"github.com/palp1tate/FlowFederate/api/global"

	"go.uber.org/zap"
)

func InitLogger() {
	logger, _ := zap.NewProduction()
	if global.ServerConfig.Debug {
		logger, _ = zap.NewDevelopment()
	}
	zap.ReplaceGlobals(logger)
}
