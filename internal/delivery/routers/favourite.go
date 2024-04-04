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

func RegisterFavouriteRouter(r *gin.Engine, db *sqlx.DB, logger *log.Logs, tracer trace.Tracer) *gin.RouterGroup {
	favouriteRouter := r.Group("/favourite")

	placeRepo := repository.InitPlaceRepo(db)
	routeRepo := repository.InitRouteRepo(db)
	favouriteRepo := repository.InitFavouriteRepo(db)

	favouriteService := service.InitFavouriteService(favouriteRepo, placeRepo, routeRepo, logger)
	favouriteHandler := handlers.InitFavouriteHandler(favouriteService, tracer)

	favouriteRouter.POST("/like_place", favouriteHandler.LikePlace)
	favouriteRouter.POST("/like_route", favouriteHandler.LikeRoute)
	favouriteRouter.GET("/by_user_id", favouriteHandler.GetLikedByUser)

	return favouriteRouter
}
