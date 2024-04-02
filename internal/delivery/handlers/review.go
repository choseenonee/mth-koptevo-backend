package handlers

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"mth/internal/models"
	"mth/internal/models/swagger"
	"mth/internal/service"
	"mth/pkg/customerr"
	tracing "mth/pkg/trace"
	"net/http"
	"strconv"
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

// CreateOnRoute @Summary Create route review
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

// CreateOnPlace @Summary Create place review
// @Tags review
// @Accept  json
// @Produce  json
// @Param data body models.PlaceReviewCreate true "place review create"
// @Success 200 {object} int "Successfully created route review with id"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /review/create_on_place [post]
func (r ReviewHandler) CreateOnPlace(c *gin.Context) {
	ctx, span := r.tracer.Start(c.Request.Context(), CreatePlaceReview)
	defer span.End()

	var reviewCreate models.PlaceReviewCreate

	if err := c.ShouldBindJSON(&reviewCreate); err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.BindType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.AddEvent(tracing.CallToService)
	id, err := r.ReviewService.CreateOnPlace(ctx, reviewCreate)
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

// GetByAuthor @Summary Get reviews by author
// @Tags review
// @Accept  json
// @Produce  json
// @Param id query int true "author id"
// @Success 200 {object} swagger.GetByAuthor "Successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /review/author [get]
func (r ReviewHandler) GetByAuthor(c *gin.Context) {
	ctx, span := r.tracer.Start(c.Request.Context(), GetByAuthor)
	defer span.End()

	authorIDRaw := c.Query("id")
	authorID, err := strconv.Atoi(authorIDRaw)
	if err != nil {
		err := customerr.BadInput
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.Input, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.AddEvent(tracing.CallToService)
	placeReviews, routeReviews, err := r.ReviewService.GetByAuthor(ctx, authorID)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.ServiceError, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, swagger.GetByAuthor{
		PlaceReviews: placeReviews,
		RouteReviews: routeReviews,
	})
}

// GetByPlace @Summary Get reviews by place
// @Tags review
// @Accept  json
// @Produce  json
// @Param id query int true "place id"
// @Success 200 {object} []models.PlaceReview "Successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /review/place [get]
func (r ReviewHandler) GetByPlace(c *gin.Context) {
	ctx, span := r.tracer.Start(c.Request.Context(), GetByPlace)
	defer span.End()

	placeIDRaw := c.Query("id")
	placeID, err := strconv.Atoi(placeIDRaw)
	if err != nil {
		err := customerr.BadInput
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.Input, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.AddEvent(tracing.CallToService)
	placeReviews, err := r.ReviewService.GetByPlace(ctx, placeID)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.ServiceError, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, placeReviews)
}

// GetByRoute @Summary Get reviews by route
// @Tags review
// @Accept  json
// @Produce  json
// @Param id query int true "route id"
// @Success 200 {object} []models.RouteReview "Successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /review/route [get]
func (r ReviewHandler) GetByRoute(c *gin.Context) {
	ctx, span := r.tracer.Start(c.Request.Context(), GetByRoute)
	defer span.End()

	placeIDRaw := c.Query("id")
	placeID, err := strconv.Atoi(placeIDRaw)
	if err != nil {
		err := customerr.BadInput
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.Input, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.AddEvent(tracing.CallToService)
	placeReviews, err := r.ReviewService.GetByPlace(ctx, placeID)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.ServiceError, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, placeReviews)
}

// UpdateOnPlace @Summary Update review on place
// @Tags review
// @Accept  json
// @Produce  json
// @Param review body models.ReviewUpdate true "place review id"
// @Success 200 {object} string "Successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /review/update_on_place [put]
func (r ReviewHandler) UpdateOnPlace(c *gin.Context) {
	ctx, span := r.tracer.Start(c.Request.Context(), UpdateOnPlace)
	defer span.End()

	var reviewUpdate models.ReviewUpdate

	if err := c.ShouldBindJSON(&reviewUpdate); err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.BindType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.AddEvent(tracing.CallToService)
	err := r.ReviewService.UpdateOnPlace(ctx, reviewUpdate)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.ServiceError, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "Successfully!")
}

// UpdateOnRoute @Summary Update review on route
// @Tags review
// @Accept  json
// @Produce  json
// @Param review body models.ReviewUpdate true "route review id"
// @Success 200 {object} string "Successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /review/update_on_route [put]
func (r ReviewHandler) UpdateOnRoute(c *gin.Context) {
	ctx, span := r.tracer.Start(c.Request.Context(), UpdateOnRoute)
	defer span.End()

	var reviewUpdate models.ReviewUpdate

	if err := c.ShouldBindJSON(&reviewUpdate); err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.BindType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.AddEvent(tracing.CallToService)
	err := r.ReviewService.UpdateOnRoute(ctx, reviewUpdate)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.ServiceError, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "Successfully!")
}
