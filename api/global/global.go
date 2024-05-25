package global

import (
	"github.com/palp1tate/FlowFederate/api/config"
	"github.com/palp1tate/FlowFederate/service/edge/pb"
	"github.com/redis/go-redis/v9"

	ut "github.com/go-playground/universal-translator"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB

	Translator ut.Translator

	NacosConfig *config.NacosConfig

	RedisClient *redis.Client

	ServerConfig *config.ServerConfig

	EdgeServiceClient pb.EdgeServiceClient
)
