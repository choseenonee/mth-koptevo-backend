package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"mth/internal/service"
	tracing "mth/pkg/trace"
	"net/http"
	"strconv"
	"strings"
)

type UserHandler struct {
	userService service.User
	tracer      trace.Tracer
}

func InitUserHandler(userService service.User, tracer trace.Tracer) UserHandler {
	return UserHandler{
		userService: userService,
		tracer:      tracer,
	}
}

// CheckIn @Summary Checkin
// @Tags user
// @Accept  json
// @Produce  json
// @Param cipher query string true "Cipher"
// @Param user_id query int true "UserID"
// @Success 200 {object} string "Just valid hash"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /user/check_in [post]
func (u UserHandler) CheckIn(c *gin.Context) {
	ctx, span := u.tracer.Start(c.Request.Context(), "User check in")
	defer span.End()

	cipher := c.Query("cipher")
	userIDRaw := c.Query("user_id")
	if cipher == "" || userIDRaw == "" {
		err := errors.New("no userID or cipher provided")
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.Input, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := strconv.Atoi(userIDRaw)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.Input, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash, err := u.userService.CheckIn(ctx, cipher, userID)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.Input, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())

		if strings.Contains(err.Error(), "пользователь уже чекинился в этом месте") {
			c.JSON(http.StatusMethodNotAllowed, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, hash)
}

// ValidateHash @Summary Validate hash
// @Tags user
// @Accept  json
// @Produce  json
// @Param hash query string true "hash"
// @Success 200 {object} string "Just bool"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /user/validate [post]
func (u UserHandler) ValidateHash(c *gin.Context) {
	ctx, span := u.tracer.Start(c.Request.Context(), "User validate hash")
	defer span.End()

	hash := c.Query("hash")
	if hash == "" {
		err := errors.New("no userID or cipher provided")
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.Input, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if isHashValid := u.userService.ValidateHash(ctx, hash); isHashValid {
		c.Status(http.StatusOK)
	} else {
		c.Status(http.StatusTeapot)
	}
}
