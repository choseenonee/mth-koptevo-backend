package handlers

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"mth/internal/models"
	"mth/internal/models/swagger"
	"mth/internal/service"
	tracing "mth/pkg/trace"
	"net/http"
	"strconv"
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

// GetAllWithFilter @Summary Get places by filter (or without)
// @Tags place
// @Accept  json
// @Produce  json
// @Param data body swagger.Filters true "Filters"
// @Success 200 {object} []models.Place "Successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /place/get_all_with_filter [put]
func (r PlaceHandler) GetAllWithFilter(c *gin.Context) {
	ctx, span := r.tracer.Start(c.Request.Context(), GetAllPlacesWithFilters)
	defer span.End()

	var filters swagger.Filters

	if err := c.ShouldBindJSON(&filters); err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.BindType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.AddEvent(tracing.CallToService)
	places, err := r.PlaceService.GetAllWithFilter(ctx, filters)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.ServiceError, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, places)
}

// GetByID @Summary Get place by id
// @Tags place
// @Accept  json
// @Produce  json
// @Param id query int true "Place id"
// @Success 200 {object} models.Place "Successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /place/by_id [get]
func (r PlaceHandler) GetByID(c *gin.Context) {
	ctx, span := r.tracer.Start(c.Request.Context(), GetPlaceById)
	defer span.End()

	idRaw := c.Query("id")
	id, err := strconv.Atoi(idRaw)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.Input, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.AddEvent(tracing.CallToService)
	place, err := r.PlaceService.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.ServiceError, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, place)
}
