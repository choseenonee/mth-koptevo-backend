package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"
	"mth/internal/delivery/handlers"
	"mth/internal/repository"
	"mth/internal/service"
	"mth/pkg/log"
)

func RegisterUserRouter(r *gin.Engine, db *sqlx.DB, logger *log.Logs, tracer trace.Tracer) *gin.RouterGroup {
	userRouter := r.Group("/user")

	userRepo := repository.InitUserRepo(db)

	userService := service.InitUserService(userRepo, logger)
	userHandler := handlers.InitUserHandler(userService, tracer)

	userRouter.POST("/check_in", userHandler.CheckIn)
	userRouter.POST("/validate", userHandler.ValidateHash)

	return userRouter
}
