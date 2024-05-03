package global

import (
	"github.com/palp1tate/FlowFederate/api/config"
	"github.com/palp1tate/FlowFederate/service/edge/pb"
	"gorm.io/gorm"

	ut "github.com/go-playground/universal-translator"
)

var (
	DB *gorm.DB

	Translator ut.Translator

	ServerConfig *config.ServerConfig

	EdgeServiceClient pb.EdgeServiceClient
)
