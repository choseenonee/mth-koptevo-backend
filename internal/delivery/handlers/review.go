package handlers

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"mth/internal/models"
	"mth/internal/service"
	tracing "mth/pkg/trace"
	"net/http"
)

type ReviewHandler struct {
	ReviewService service.Review
	tracer        trace.Tracer
}

func InitReviewHandler(ReviewService service.Review, tracer trace.Tracer) ReviewHandler {
	return ReviewHandler{
		ReviewService: ReviewService,
		tracer:        tracer,
	}
}

// Create @Summary Create route review
// @Tags review
// @Accept  json
// @Produce  json
// @Param data body models.RouteReviewCreate true "Route review create"
// @Success 200 {object} int "Successfully created route review with id"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /review/create_on_route [post]
func (r ReviewHandler) CreateOnRoute(c *gin.Context) {
	ctx, span := r.tracer.Start(c.Request.Context(), CreateRouteReview)
	defer span.End()

	var reviewCreate models.RouteReviewCreate

	if err := c.ShouldBindJSON(&reviewCreate); err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.BindType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.AddEvent(tracing.CallToService)
	id, err := r.ReviewService.CreateOnRoute(ctx, reviewCreate)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.ServiceError, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, id)
}
