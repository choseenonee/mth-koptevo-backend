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

type TripHandler struct {
	tripService service.Trip
	tracer      trace.Tracer
}

func InitTripHandler(tripService service.Trip, tracer trace.Tracer) TripHandler {
	return TripHandler{
		tripService: tripService,
		tracer:      tracer,
	}
}

// Create @Summary Create trip
// @Tags trip
// @Accept  json
// @Produce  json
// @Param data body models.TripCreate true "trip"
// @Success 200 {object} int "Successfully created trip with id"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /trip/create [post]
func (r TripHandler) Create(c *gin.Context) {
	ctx, span := r.tracer.Start(c.Request.Context(), "Create trip")
	defer span.End()

	var trip models.TripCreate

	if err := c.ShouldBindJSON(&trip); err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.BindType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.AddEvent(tracing.CallToService)
	id, err := r.tripService.Create(ctx, trip)
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

// GetByID @Summary Get trips
// @Tags trip
// @Accept  json
// @Produce  json
// @Param id query int true "trip id"
// @Success 200 {object} models.Trip "Successfully created trip with id"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /trip/by_id [get]
func (r TripHandler) GetByID(c *gin.Context) {
	ctx, span := r.tracer.Start(c.Request.Context(), "Get trip")
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
	trip, err := r.tripService.GetTripByID(ctx, id)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.ServiceError, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, trip)
}

// GetByUser @Summary Get trips
// @Tags trip
// @Accept  json
// @Produce  json
// @Param id query int true "user id"
// @Success 200 {object} []models.Trip "Successfully created trip with id"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /trip/by_user_id [get]
func (r TripHandler) GetByUser(c *gin.Context) {
	ctx, span := r.tracer.Start(c.Request.Context(), "Get user trips")
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
	trips, err := r.tripService.GetTripsByUser(ctx, id)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.ServiceError, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, trips)
}

// AddRoute @Summary Add route
// @Tags trip
// @Accept  json
// @Produce  json
// @Param data body swagger.TripAdd true "data"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /trip/add/route [put]
func (r TripHandler) AddRoute(c *gin.Context) {
	ctx, span := r.tracer.Start(c.Request.Context(), "Add route to trip")
	defer span.End()

	var trip swagger.TripAdd

	if err := c.ShouldBindJSON(&trip); err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.BindType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.AddEvent(tracing.CallToService)
	err := r.tripService.AddRoute(ctx, trip.TripID, trip.EntityID, trip.Day, trip.Position)
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

// AddPlace @Summary Add place
// @Tags trip
// @Accept  json
// @Produce  json
// @Param data body swagger.TripAdd true "data"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /trip/add/place [put]
func (r TripHandler) AddPlace(c *gin.Context) {
	ctx, span := r.tracer.Start(c.Request.Context(), "Add place to trip")
	defer span.End()

	var trip swagger.TripAdd

	if err := c.ShouldBindJSON(&trip); err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.BindType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.AddEvent(tracing.CallToService)
	err := r.tripService.AddPlace(ctx, trip.TripID, trip.EntityID, trip.Day, trip.Position)
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

// ChangeRouteDay @Summary Add place
// @Tags trip
// @Accept  json
// @Produce  json
// @Param data body swagger.TripChangeDay true "data"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /trip/change/route/day [put]
func (r TripHandler) ChangeRouteDay(c *gin.Context) {
	ctx, span := r.tracer.Start(c.Request.Context(), "Change route day in trip")
	defer span.End()

	var trip swagger.TripChangeDay

	if err := c.ShouldBindJSON(&trip); err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.BindType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.AddEvent(tracing.CallToService)
	err := r.tripService.ChangeRouteDay(ctx, trip.TripID, trip.EntityID, trip.Day)
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

// ChangePlaceDay @Summary Add place
// @Tags trip
// @Accept  json
// @Produce  json
// @Param data body swagger.TripChangeDay true "data"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /trip/change/place/day [put]
func (r TripHandler) ChangePlaceDay(c *gin.Context) {
	ctx, span := r.tracer.Start(c.Request.Context(), "Change place day in trip")
	defer span.End()

	var trip swagger.TripChangeDay

	if err := c.ShouldBindJSON(&trip); err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.BindType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.AddEvent(tracing.CallToService)
	err := r.tripService.ChangePlaceDay(ctx, trip.TripID, trip.EntityID, trip.Day)
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

// ChangeRoutePosition @Summary Add place
// @Tags trip
// @Accept  json
// @Produce  json
// @Param data body swagger.TripChangePosition true "data"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /trip/change/route/position [put]
func (r TripHandler) ChangeRoutePosition(c *gin.Context) {
	ctx, span := r.tracer.Start(c.Request.Context(), "Change route position in trip")
	defer span.End()

	var trip swagger.TripChangePosition

	if err := c.ShouldBindJSON(&trip); err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.BindType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.AddEvent(tracing.CallToService)
	err := r.tripService.ChangeRoutePosition(ctx, trip.TripID, trip.EntityID, trip.Position)
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

// ChangePlacePosition @Summary Add place
// @Tags trip
// @Accept  json
// @Produce  json
// @Param data body swagger.TripChangePosition true "data"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /trip/change/place/position [put]
func (r TripHandler) ChangePlacePosition(c *gin.Context) {
	ctx, span := r.tracer.Start(c.Request.Context(), "Change place position in trip")
	defer span.End()

	var trip swagger.TripChangePosition

	if err := c.ShouldBindJSON(&trip); err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.BindType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.AddEvent(tracing.CallToService)
	err := r.tripService.ChangePlacePosition(ctx, trip.TripID, trip.EntityID, trip.Position)
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

// DeleteRoute @Summary Add place
// @Tags trip
// @Accept  json
// @Produce  json
// @Param trip_id query int true "trip id"
// @Param route_id query int true "route id"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /trip/route [delete]
func (r TripHandler) DeleteRoute(c *gin.Context) {
	ctx, span := r.tracer.Start(c.Request.Context(), "Change place position in trip")
	defer span.End()

	tripIDRaw := c.Query("trip_id")
	tripID, err := strconv.Atoi(tripIDRaw)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.BindType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	routeIDRaw := c.Query("route_id")
	routeID, err := strconv.Atoi(routeIDRaw)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.BindType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.AddEvent(tracing.CallToService)
	err = r.tripService.DeleteRoute(ctx, tripID, routeID)
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

// DeletePlace @Summary Add place
// @Tags trip
// @Accept  json
// @Produce  json
// @Param trip_id query int true "trip id"
// @Param place_id query int true "place id"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /trip/place [delete]
func (r TripHandler) DeletePlace(c *gin.Context) {
	ctx, span := r.tracer.Start(c.Request.Context(), "Change place position in trip")
	defer span.End()

	tripIDRaw := c.Query("trip_id")
	tripID, err := strconv.Atoi(tripIDRaw)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.BindType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	placeIDRaw := c.Query("place_id")
	placeID, err := strconv.Atoi(placeIDRaw)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.BindType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.AddEvent(tracing.CallToService)
	err = r.tripService.DeletePlace(ctx, tripID, placeID)
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
