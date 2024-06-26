package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"
	"mth/pkg/log"
)

func InitRouting(r *gin.Engine, db *sqlx.DB, logger *log.Logs, tracer trace.Tracer) {
	_ = RegisterTagRouter(r, db, logger, tracer)
	_ = RegisterReviewRouter(r, db, logger, tracer)
	_ = RegisterPlaceRouter(r, db, logger, tracer)
	_ = RegisterDistrictRouter(r, db, logger, tracer)
	_ = RegisterRouteRouter(r, db, logger, tracer)
	_ = RegisterNoteRouter(r, db, logger, tracer)
	_ = RegisterCompanionsRouter(r, db, logger, tracer)
	_ = RegisterFavouriteRouter(r, db, logger, tracer)
	_ = RegisterUserRouter(r, db, logger, tracer)
	_ = RegisterTripRouter(r, db, logger, tracer)
}
