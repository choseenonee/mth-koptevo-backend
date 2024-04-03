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

func RegisterNoteRouter(r *gin.Engine, db *sqlx.DB, logger *log.Logs, tracer trace.Tracer) *gin.RouterGroup {
	noteRouter := r.Group("/note")

	noteRepo := repository.InitNoteRepo(db)

	noteService := service.InitNoteService(noteRepo, logger)
	noteHandler := handlers.InitNoteHandler(noteService, tracer)

	noteRouter.POST("/create", noteHandler.Create)
	noteRouter.GET("/by_user_and_place_ids", noteHandler.GetByIDs)
	noteRouter.GET("/by_user_id", noteHandler.GetByUserID)
	noteRouter.GET("/by_id", noteHandler.GetByNoteID)

	return noteRouter
}
