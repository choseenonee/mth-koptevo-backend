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

type FavouriteHandler struct {
	FavouriteService service.Favourite
	tracer           trace.Tracer
}

func InitFavouriteHandler(FavouriteService service.Favourite, tracer trace.Tracer) FavouriteHandler {
	return FavouriteHandler{
		FavouriteService: FavouriteService,
		tracer:           tracer,
	}
}

// LikePlace @Summary Like place
// @Tags favourite
// @Accept  json
// @Produce  json
// @Param like body models.Like true "Like"
// @Success 200 {object} string "Successfully!"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /favourite/like_place [post]
func (r FavouriteHandler) LikePlace(c *gin.Context) {
	ctx, span := r.tracer.Start(c.Request.Context(), LikePlace)
	defer span.End()

	var like models.Like

	if err := c.ShouldBindJSON(&like); err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.BindType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.AddEvent(tracing.CallToService)
	err := r.FavouriteService.LikePlace(ctx, like)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.ServiceError, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "successfully!")
}

// LikeRoute @Summary Like route
// @Tags favourite
// @Accept  json
// @Produce  json
// @Param like body models.Like true "Like"
// @Success 200 {object} string "Successfully!"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /favourite/like_route [post]
func (r FavouriteHandler) LikeRoute(c *gin.Context) {
	ctx, span := r.tracer.Start(c.Request.Context(), LikeRoute)
	defer span.End()

	var like models.Like

	if err := c.ShouldBindJSON(&like); err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.BindType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.AddEvent(tracing.CallToService)
	err := r.FavouriteService.LikeRoute(ctx, like)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.ServiceError, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "successfully!")
}

// GetLikedByUser @Summary Get liked
// @Tags favourite
// @Accept  json
// @Produce  json
// @Param id query int true "UserID"
// @Success 200 {object} string "Successfully!"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /favourite/by_user_id [get]
func (r FavouriteHandler) GetLikedByUser(c *gin.Context) {
	ctx, span := r.tracer.Start(c.Request.Context(), GetLiked)
	defer span.End()

	idRaw := c.Query("id")
	userID, err := strconv.Atoi(idRaw)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.Input, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.AddEvent(tracing.CallToService)
	places, routesRaw, err := r.FavouriteService.GetLikedByUser(ctx, userID)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.ServiceError, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, swagger.InitFavourite(places, routesRaw))
}

// DeleteOnPlace @Summary Delete liked on place
// @Tags favourite
// @Accept  json
// @Produce  json
// @Param delete body models.Like true "delete data"
// @Success 200 {object} string "Successfully!"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /favourite/like_place [delete]
func (r FavouriteHandler) DeleteOnPlace(c *gin.Context) {
	ctx, span := r.tracer.Start(c.Request.Context(), DeleteOnPlace)
	defer span.End()

	var like models.Like

	if err := c.ShouldBindJSON(&like); err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.BindType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.AddEvent(tracing.CallToService)
	err := r.FavouriteService.DeleteOnPlace(ctx, like)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.ServiceError, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "success")
}

// DeleteOnRoute @Summary Delete liked on route
// @Tags favourite
// @Accept  json
// @Produce  json
// @Param delete body models.Like true "delete data"
// @Success 200 {object} string "Successfully!"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /favourite/like_route [delete]
func (r FavouriteHandler) DeleteOnRoute(c *gin.Context) {
	ctx, span := r.tracer.Start(c.Request.Context(), DeleteOnRoute)
	defer span.End()

	var like models.Like

	if err := c.ShouldBindJSON(&like); err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.BindType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.AddEvent(tracing.CallToService)
	err := r.FavouriteService.DeleteOnRoute(ctx, like)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.ServiceError, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "success")
}
