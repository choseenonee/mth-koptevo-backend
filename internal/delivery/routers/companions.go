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

func RegisterCompanionsRouter(r *gin.Engine, db *sqlx.DB, logger *log.Logs, tracer trace.Tracer) *gin.RouterGroup {
	companionsRouter := r.Group("/companions")

	companionsRepo := repository.InitCompanionsRepo(db)

	companionsService := service.InitCompanionsService(companionsRepo, logger)
	companionsHandler := handlers.InitCompanionsHandler(companionsService, tracer)

	companionsRouter.POST("/create_place_companion", companionsHandler.CreateCompanionPlace)
	companionsRouter.POST("/create_route_companion", companionsHandler.CreateCompanionRoute)

	return companionsRouter
}
