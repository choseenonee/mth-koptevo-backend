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

type PlaceHandler struct {
	PlaceService service.Place
	tracer       trace.Tracer
}

func InitPlaceHandler(PlaceService service.Place, tracer trace.Tracer) PlaceHandler {
	return PlaceHandler{
		PlaceService: PlaceService,
		tracer:       tracer,
	}
}

// Create @Summary Create place with tags
// @Tags place
// @Accept  json
// @Produce  json
// @Param data body models.PlaceCreate true "Place with tag ids create"
// @Success 200 {object} int "Successfully created place with id"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /place/create [post]
func (r PlaceHandler) Create(c *gin.Context) {
	ctx, span := r.tracer.Start(c.Request.Context(), PlaceCreate)
	defer span.End()

	var placeCreate models.PlaceCreate

	if err := c.ShouldBindJSON(&placeCreate); err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.BindType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.AddEvent(tracing.CallToService)
	id, err := r.PlaceService.Create(ctx, placeCreate)
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
