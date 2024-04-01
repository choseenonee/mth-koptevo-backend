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

type TagHandler struct {
	tagService service.Tag
	tracer     trace.Tracer
}

func InitTagHandler(tagService service.Tag, tracer trace.Tracer) TagHandler {
	return TagHandler{
		tagService: tagService,
		tracer:     tracer,
	}
}

// Create @Summary Create tag
// @Tags tag
// @Accept  json
// @Produce  json
// @Param data body models.TagCreate true "Tag create"
// @Success 200 {object} int "Successfully created tag with id"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /tag/create [post]
func (t TagHandler) Create(c *gin.Context) {
	ctx, span := t.tracer.Start(c.Request.Context(), CreateTag)
	defer span.End()

	var tagCreate models.TagCreate

	if err := c.ShouldBindJSON(&tagCreate); err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.BindType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.AddEvent(tracing.CallToService)
	id, err := t.tagService.Create(ctx, tagCreate)
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

// Create @Summary Get all tags
// @Tags tag
// @Accept  json
// @Produce  json
// @Success 200 {object} []models.Tag "Successfully created tag with id"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /tag/get_all [get]
func (t TagHandler) GetAll(c *gin.Context) {
	ctx, span := t.tracer.Start(c.Request.Context(), CreateTag)
	defer span.End()

	span.AddEvent(tracing.CallToService)
	tags, err := t.tagService.GetAll(ctx)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.ServiceError, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tags)
}
