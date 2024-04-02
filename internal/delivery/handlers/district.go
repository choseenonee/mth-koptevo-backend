package handlers

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"mth/internal/service"
	tracing "mth/pkg/trace"
	"net/http"
	"strconv"
)

type DistrictHandler struct {
	DistrictService service.District
	tracer          trace.Tracer
}

func InitDistrictHandler(DistrictService service.District, tracer trace.Tracer) DistrictHandler {
	return DistrictHandler{
		DistrictService: DistrictService,
		tracer:          tracer,
	}
}

// GetByID @Summary Get districts by city id
// @Tags district
// @Accept  json
// @Produce  json
// @Param id query int true "City id"
// @Success 200 {object} []models.District "Successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /district/by_city_id [get]
func (r DistrictHandler) GetByID(c *gin.Context) {
	ctx, span := r.tracer.Start(c.Request.Context(), GetDistrictByCityID)
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
	districts, err := r.DistrictService.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.ServiceError, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, districts)
}
