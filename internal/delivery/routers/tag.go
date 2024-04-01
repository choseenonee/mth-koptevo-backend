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

func RegisterTagRouter(r *gin.Engine, db *sqlx.DB, logger *log.Logs, tracer trace.Tracer) *gin.RouterGroup {
	tagRouter := r.Group("/tag")

	tagRepo := repository.InitTagRepo(db)

	tagService := service.InitTagService(tagRepo, logger)
	tagHandler := handlers.InitTagHandler(tagService, tracer)

	tagRouter.POST("/create", tagHandler.Create)
	tagRouter.GET("/get_all", tagHandler.GetAll)

	return tagRouter
}
