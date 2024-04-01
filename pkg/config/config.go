package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

const (
	DBName     = "DB_NAME"
	DBUser     = "DB_USER"
	DBPassword = "DB_PASSWORD"
	DBPort     = "DB_PORT"
	DBHost     = "DB_HOST"
	JWTExpire  = "JWT_EXPIRE"
	Secret     = "SECRET"
	JaegerHost = "JAEGER_HOST"
	JaegerPort = "JAEGER_PORT"
)

func InitConfig() {
	envPath, _ := os.Getwd()
	envPath = filepath.Join(envPath, "..") // workdir is cmd
	envPath = filepath.Join(envPath, "/deploy")

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(envPath)

	viper.AutomaticEnv()

	err := viper.ReadInConfig()

	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to init config. Error:%v, readed config: %v, %v", err.Error(),
			viper.GetString(DBName), viper.GetString(DBPort)))
	}
}
