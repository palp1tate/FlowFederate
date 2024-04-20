package global

import (
	"github.com/palp1tate/FlowFederate/api/config"
	"gorm.io/gorm"

	ut "github.com/go-playground/universal-translator"
)

var (
	DB         *gorm.DB
	Translator ut.Translator

	ServerConfig *config.ServerConfig
)
