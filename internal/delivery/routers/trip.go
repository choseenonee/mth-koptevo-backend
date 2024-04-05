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

func RegisterTripRouter(r *gin.Engine, db *sqlx.DB, logger *log.Logs, tracer trace.Tracer) *gin.RouterGroup {
	tripRouter := r.Group("/trip")

	tripRepo := repository.InitTripRepo(db)

	tripService := service.InitTripService(tripRepo, logger)
	tripHandler := handlers.InitTripHandler(tripService, tracer)

	tripRouter.POST("/create", tripHandler.Create)
	tripRouter.GET("/by_id", tripHandler.GetByID)
	tripRouter.GET("/by_user_id", tripHandler.GetByUser)

	tripRouter.PUT("/route/add", tripHandler.AddRoute)
	tripRouter.PUT("/route/change/day", tripHandler.ChangeRouteDay)
	tripRouter.PUT("/route/change/position", tripHandler.ChangeRoutePosition)
	tripRouter.DELETE("/route", tripHandler.DeleteRoute)

	tripRouter.PUT("/place/add", tripHandler.AddPlace)
	tripRouter.PUT("/place/change/day", tripHandler.ChangePlaceDay)
	tripRouter.PUT("/place/change/position", tripHandler.ChangePlacePosition)
	tripRouter.DELETE("/place", tripHandler.DeletePlace)

	return tripRouter
}
