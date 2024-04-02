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
	"strconv"
)

type RouteHandler struct {
	RouteService service.Route
	tracer       trace.Tracer
}

func InitRouteHandler(routeService service.Route, tracer trace.Tracer) RouteHandler {
	return RouteHandler{
		RouteService: routeService,
		tracer:       tracer,
	}
}

// Create @Summary Create routes with tags and places
// @Tags route
// @Accept  json
// @Produce  json
// @Param data body models.RouteCreate true "Route with tag ids and place ids create"
// @Success 200 {object} int "Successfully created route with id"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /route/create [post]
func (r RouteHandler) Create(c *gin.Context) {
	ctx, span := r.tracer.Start(c.Request.Context(), RouteCreate)
	defer span.End()

	var routeCreate models.RouteCreate

	if err := c.ShouldBindJSON(&routeCreate); err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.BindType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.AddEvent(tracing.CallToService)
	id, err := r.RouteService.Create(ctx, routeCreate)
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

// GetRouteByID @Summary Get route by id
// @Tags route
// @Accept  json
// @Produce  json
// @Param id query int true "id"
// @Success 200 {object} models.Route "Successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /route/by_id [get]
func (r RouteHandler) GetRouteByID(c *gin.Context) {
	ctx, span := r.tracer.Start(c.Request.Context(), GetRouteByID)
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
	route, err := r.RouteService.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.ServiceError, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, route)
}

// GetRouteByPage @Summary Get route by page
// @Tags route
// @Accept  json
// @Produce  json
// @Param page query int true "id"
// @Success 200 {object} []models.Route "Successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /route/by_page [get]
func (r RouteHandler) GetRouteByPage(c *gin.Context) {
	ctx, span := r.tracer.Start(c.Request.Context(), GetRoutesByPage)
	defer span.End()

	pageRaw := c.Query("page")

	page, err := strconv.Atoi(pageRaw)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.BindType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.AddEvent(tracing.CallToService)
	routes, err := r.RouteService.GetAll(ctx, page)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.ServiceError, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, routes)
}
