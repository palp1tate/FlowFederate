package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/palp1tate/FlowFederate/api/global"
	"github.com/palp1tate/FlowFederate/api/initialize"
	"github.com/palp1tate/FlowFederate/internal/util"

	"go.uber.org/zap"
)

func main() {
	initialize.InitConfig()
	initialize.InitLogger()
	initialize.InitMySQL()
	router := initialize.Router()
	if err := initialize.InitTranslator("zh"); err != nil {
		zap.S().Warn(err)
		panic(err)
	}
	initialize.InitServiceConn()

	host := global.ServerConfig.Api.Host
	port := flag.Int("p", 0, "port number")
	flag.Parse()
	if *port == 0 {
		*port, _ = util.GetFreePort()
	}
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: router,
	}
	go func() {
		zap.S().Infof("Starting service at %s:%d", host, *port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zap.S().Panic(err)
		}
	}()

	client := initialize.NewRegistryClient(global.ServerConfig.Consul.Host, global.ServerConfig.Consul.Port)
	apiName := global.ServerConfig.Api.Name
	apiTags := global.ServerConfig.Api.Tags
	apiId := util.GenerateUUID()
	err := client.Register(host, *port, apiName, apiTags, apiId)
	if err != nil {
		zap.S().Panic(err)
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = srv.Shutdown(ctx); err != nil {
		zap.S().Error(err)
	}
	if err = client.DeRegister(apiId); err != nil {
		zap.S().Warnf(err.Error())
	} else {
		zap.S().Infof("%s logged off successfully", apiName)
	}
}
