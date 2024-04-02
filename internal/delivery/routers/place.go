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

func RegisterPlaceRouter(r *gin.Engine, db *sqlx.DB, logger *log.Logs, tracer trace.Tracer) *gin.RouterGroup {
	placeRouter := r.Group("/place")

	placeRepo := repository.InitPlaceRepo(db)

	placeService := service.InitPlaceService(placeRepo, logger)
	placeHandler := handlers.InitPlaceHandler(placeService, tracer)

	placeRouter.POST("/create", placeHandler.Create)
	placeRouter.GET("/by_id", placeHandler.GetByID)
	placeRouter.PUT("/get_all_with_filter", placeHandler.GetAllWithFilter)
	//reviewRouter.POST("/create_on_place", reviewHandler.CreateOnPlace)
	//reviewRouter.GET("/author", reviewHandler.GetByAuthor)
	//reviewRouter.GET("/place", reviewHandler.GetByPlace)
	//reviewRouter.GET("/route", reviewHandler.GetByRoute)
	//reviewRouter.PUT("/update_on_place", reviewHandler.UpdateOnPlace)
	//reviewRouter.PUT("/update_on_route", reviewHandler.UpdateOnRoute)

	return placeRouter
}
