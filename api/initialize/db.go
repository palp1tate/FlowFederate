package initialize

import (
	"fmt"

	"gorm.io/gorm/schema"

	"github.com/palp1tate/FlowFederate/api/global"
	"go.uber.org/zap"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitMySQL() {
	var err error
	m := global.ServerConfig.MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		m.User, m.Password, m.Host, m.Port, m.Database)
	ormLogger := logger.Default.LogMode(logger.Silent)
	if global.ServerConfig.Debug {
		ormLogger = ormLogger.LogMode(logger.Info)
	}
	global.DB, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                      dsn,
		DefaultStringSize:        256,
		DisableDatetimePrecision: true,
	}), &gorm.Config{
		Logger: ormLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		zap.S().Panic(err.Error())
	}
}
