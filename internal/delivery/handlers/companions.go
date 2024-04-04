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

type CompanionsHandler struct {
	companionsService service.Companions
	tracer            trace.Tracer
}

func InitCompanionsHandler(companionsService service.Companions, tracer trace.Tracer) CompanionsHandler {
	return CompanionsHandler{
		companionsService: companionsService,
		tracer:            tracer,
	}
}

// CreateCompanionPlace creates a new companion place
// @Summary Create a new companion place
// @Description Adds a new companion place to the database
// @Tags companions
// @Accept json
// @Produce json
// @Param data body models.CompanionsPlaceCreate true "Companion Place Data"
// @Success 200 "Successfully created companion place with id"
// @Failure 400 {object} map[string]string "Invalid input data"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /companions/create_place_companion [post]
func (ch CompanionsHandler) CreateCompanionPlace(c *gin.Context) {
	ctx, span := ch.tracer.Start(c.Request.Context(), CompanionsCreate)
	defer span.End()

	var companionCreate models.CompanionsPlaceCreate

	if err := c.ShouldBindJSON(&companionCreate); err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.BindType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.AddEvent(tracing.CallToService)
	err := ch.companionsService.CreatePlaceCompanions(ctx, companionCreate)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.ServiceError, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// CreateCompanionRoute creates a new companion place
// @Summary Create a new companion place
// @Description Adds a new companion place to the database
// @Tags companions
// @Accept json
// @Produce json
// @Param data body models.CompanionsRouteCreate true "Companion Place Data"
// @Success 200 "Successfully created companion place with id"
// @Failure 400 {object} map[string]string "Invalid input data"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /companions/create_route_companion [post]
func (ch CompanionsHandler) CreateCompanionRoute(c *gin.Context) {
	ctx, span := ch.tracer.Start(c.Request.Context(), CompanionsCreate)
	defer span.End()

	var companionCreate models.CompanionsRouteCreate

	if err := c.ShouldBindJSON(&companionCreate); err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.BindType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.AddEvent(tracing.CallToService)
	err := ch.companionsService.CreateRouteCompanions(ctx, companionCreate)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.ServiceError, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
