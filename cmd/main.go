package main

import (
	"fmt"
	"github.com/spf13/viper"
	"mth/pkg/config"
	"mth/pkg/database"
	"mth/pkg/log"
	"mth/pkg/trace"
)

const serviceName = "mth backend"

func main() {
	logger, loggerInfoFile, loggerErrorFile := log.InitLogger()
	defer loggerInfoFile.Close()
	defer loggerErrorFile.Close()

	logger.Info("Logger Initialized")

	config.InitConfig()
	logger.Info("Config Initialized")

	jaegerURL := fmt.Sprintf("http://%v:%v/api/traces", viper.GetString(config.JaegerHost), viper.GetString(config.JaegerPort))
	tracer := trace.InitTracer(jaegerURL, serviceName)
	logger.Info("Tracer Initialized")

	db := database.GetDB()
	logger.Info("Database Initialized")

	_ = db
	_ = tracer
}
