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

// GetByUser
// @Summary get companion data by user id
// @Tags companions
// @Accept json
// @Produce json
// @Param id query int true "user id"
// @Success 200 {object} swagger.Companion "success"
// @Failure 400 {object} map[string]string "Invalid input data"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /companions/by_user [get]
func (ch CompanionsHandler) GetByUser(c *gin.Context) {
	ctx, span := ch.tracer.Start(c.Request.Context(), "Companions get by user")
	defer span.End()

	idRaw := c.Query("id")
	id, err := strconv.Atoi(idRaw)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.BindType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.AddEvent(tracing.CallToService)
	places, routes, err := ch.companionsService.GetByUser(ctx, id)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.ServiceError, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	companions := swagger.Companion{
		Places: places,
		Routes: routes,
	}

	c.JSON(http.StatusOK, companions)
}

// GetCompanionsPlace
// @Summary get places companions by filters
// @Tags companions
// @Accept json
// @Produce json
// @Param data body models.CompanionsFilters true "user id"
// @Success 200 {object} []models.CompanionsPlace "success"
// @Failure 400 {object} map[string]string "Invalid input data"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /companions/get_by_place [put]
func (ch CompanionsHandler) GetCompanionsPlace(c *gin.Context) {
	ctx, span := ch.tracer.Start(c.Request.Context(), "Companions get by place")
	defer span.End()

	var filters models.CompanionsFilters

	if err := c.ShouldBindJSON(&filters); err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.BindType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.AddEvent(tracing.CallToService)
	companions, err := ch.companionsService.GetCompanionsPlace(ctx, filters)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.ServiceError, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, companions)
}

// GetCompanionsRoute
// @Summary get places companions by filters
// @Tags companions
// @Accept json
// @Produce json
// @Param data body models.CompanionsFilters true "user id"
// @Success 200 {object} []models.CompanionsPlace "success"
// @Failure 400 {object} map[string]string "Invalid input data"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /companions/get_by_route [put]
func (ch CompanionsHandler) GetCompanionsRoute(c *gin.Context) {
	ctx, span := ch.tracer.Start(c.Request.Context(), "Companions get by route")
	defer span.End()

	var filters models.CompanionsFilters

	if err := c.ShouldBindJSON(&filters); err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.BindType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.AddEvent(tracing.CallToService)
	companions, err := ch.companionsService.GetCompanionsRoute(ctx, filters)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.ServiceError, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, companions)
}

// DeleteFromPlace
// @Summary get companion data by user id
// @Tags companions
// @Accept json
// @Produce json
// @Param id query int true "companion table id"
// @Success 200 "success"
// @Failure 400 {object} map[string]string "Invalid input data"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /companions/place [delete]
func (ch CompanionsHandler) DeleteFromPlace(c *gin.Context) {
	ctx, span := ch.tracer.Start(c.Request.Context(), "Companions delete by place")
	defer span.End()

	idRaw := c.Query("id")
	id, err := strconv.Atoi(idRaw)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.BindType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.AddEvent(tracing.CallToService)
	err = ch.companionsService.DeleteCompanionsPlace(ctx, id)
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

// DeleteFromRoute
// @Summary get companion data by user id
// @Tags companions
// @Accept json
// @Produce json
// @Param id query int true "companion table id"
// @Success 200 "success"
// @Failure 400 {object} map[string]string "Invalid input data"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /companions/route [delete]
func (ch CompanionsHandler) DeleteFromRoute(c *gin.Context) {
	ctx, span := ch.tracer.Start(c.Request.Context(), "Companions delete by route")
	defer span.End()

	idRaw := c.Query("id")
	id, err := strconv.Atoi(idRaw)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.BindType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.AddEvent(tracing.CallToService)
	err = ch.companionsService.DeleteCompanionsRoute(ctx, id)
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
