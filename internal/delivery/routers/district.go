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

func RegisterDistrictRouter(r *gin.Engine, db *sqlx.DB, logger *log.Logs, tracer trace.Tracer) *gin.RouterGroup {
	districtRouter := r.Group("/district")

	districtRepo := repository.InitDistrictRepo(db)

	districtService := service.InitDistrictService(districtRepo, logger)
	districtHandler := handlers.InitDistrictHandler(districtService, tracer)

	districtRouter.GET("/by_city_id", districtHandler.GetByID)

	return districtRouter
}
