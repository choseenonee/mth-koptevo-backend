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

type NoteHandler struct {
	NoteService service.Note
	tracer      trace.Tracer
}

func InitNoteHandler(NoteService service.Note, tracer trace.Tracer) NoteHandler {
	return NoteHandler{
		NoteService: NoteService,
		tracer:      tracer,
	}
}

// Create @Summary Create note
// @Tags note
// @Accept  json
// @Produce  json
// @Param data body models.NoteCreate true "Note"
// @Success 200 {object} int "Successfully created note with id"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /note/create [post]
func (r NoteHandler) Create(c *gin.Context) {
	ctx, span := r.tracer.Start(c.Request.Context(), NoteCreate)
	defer span.End()

	var noteCreate models.NoteCreate

	if err := c.ShouldBindJSON(&noteCreate); err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.BindType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.AddEvent(tracing.CallToService)
	id, err := r.NoteService.Create(ctx, noteCreate)
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

// GetByIDs @Summary Get note by user_id and place_id
// @Tags note
// @Accept  json
// @Produce  json
// @Param user_id query int true "user_id"
// @Param place_id query int true "place_id"
// @Success 200 {object} models.Note "Successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /note/by_user_and_place_ids [get]
func (r NoteHandler) GetByIDs(c *gin.Context) {
	ctx, span := r.tracer.Start(c.Request.Context(), GetNoteByID)
	defer span.End()

	userIDRaw := c.Query("user_id")
	userID, err := strconv.Atoi(userIDRaw)
	placeIDRaw := c.Query("place_id")
	placeID, err := strconv.Atoi(placeIDRaw)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.Input, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.AddEvent(tracing.CallToService)
	note, err := r.NoteService.GetByID(ctx, userID, placeID)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.ServiceError, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, note)
}

// GetByUserID @Summary Get notes by user_id
// @Tags note
// @Accept  json
// @Produce  json
// @Param id query int true "User id"
// @Success 200 {object} []models.Note "Successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /note/by_user_id [get]
func (r NoteHandler) GetByUserID(c *gin.Context) {
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
	place, err := r.NoteService.GetByUser(ctx, id)
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
