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

func RegisterReviewRouter(r *gin.Engine, db *sqlx.DB, logger *log.Logs, tracer trace.Tracer) *gin.RouterGroup {
	reviewRouter := r.Group("/review")

	reviewRepo := repository.InitReviewRepo(db)

	reviewService := service.InitReviewService(reviewRepo, logger)
	reviewHandler := handlers.InitReviewHandler(reviewService, tracer)

	reviewRouter.POST("/create_on_route", reviewHandler.CreateOnRoute)
	reviewRouter.POST("/create_on_place", reviewHandler.CreateOnPlace)
	reviewRouter.GET("/author", reviewHandler.GetByAuthor)
	reviewRouter.GET("/place", reviewHandler.GetByPlace)
	reviewRouter.GET("/route", reviewHandler.GetByRoute)
	reviewRouter.PUT("/update_on_place", reviewHandler.UpdateOnPlace)
	reviewRouter.PUT("/update_on_route", reviewHandler.UpdateOnRoute)

	return reviewRouter
}
