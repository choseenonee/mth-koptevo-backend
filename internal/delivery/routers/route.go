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

func RegisterRouteRouter(r *gin.Engine, db *sqlx.DB, logger *log.Logs, tracer trace.Tracer) *gin.RouterGroup {
	routeRouter := r.Group("/route")

	routeRepo := repository.InitRouteRepo(db)
	placeRepo := repository.InitPlaceRepo(db)

	routeService := service.InitRouteService(routeRepo, placeRepo, logger)
	routeHandler := handlers.InitRouteHandler(routeService, tracer)

	routeRouter.POST("/create", routeHandler.Create)
	routeRouter.GET("/by_id", routeHandler.GetRouteByID)
	routeRouter.GET("/by_page", routeHandler.GetRouteByPage)

	return routeRouter
}
