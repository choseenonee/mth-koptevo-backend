package handlers

import (
	"errors"
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
	"strings"
	"time"
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
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
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

// GetUser @Summary Login user
// @Tags user
// @Accept  json
// @Produce  json
// @Param data body swagger.User true "login data"
// @Success 200 {object} int "user ID"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /user/login [put]
func (u UserHandler) GetUser(c *gin.Context) {
	ctx, span := u.tracer.Start(c.Request.Context(), "User login")
	defer span.End()

	var user swagger.User

	if err := c.ShouldBindJSON(&user); err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.BindType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := u.userService.GetUser(ctx, user.Login, user.Password)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.BindType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())

		var status int

		status = http.StatusInternalServerError

		if strings.Contains(err.Error(), "user password isn't correct") || strings.Contains(err.Error(), "no rows in result set") {
			status = http.StatusUnauthorized
		}

		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userID)
}

// CreateUser @Summary Create user
// @Tags user
// @Accept  json
// @Produce  json
// @Param data body models.UserCreate true "create user data"
// @Success 200 {object} int "user ID"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /user/register [put]
func (u UserHandler) CreateUser(c *gin.Context) {
	ctx, span := u.tracer.Start(c.Request.Context(), "User register")
	defer span.End()

	var user models.UserCreate

	if err := c.ShouldBindJSON(&user); err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.BindType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := u.userService.CreateUser(ctx, user)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.BindType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())

		if strings.Contains(err.Error(), "user password isn't correct") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userID)
}

// GetCheckedPlaces @Summary Получить места где юзер уже зачекинился
// @Tags user
// @Accept  json
// @Produce  json
// @Param id query string true "userID"
// @Success 200 {object} []models.Place "places"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /user/checked_in [get]
func (u UserHandler) GetCheckedPlaces(c *gin.Context) {
	ctx, span := u.tracer.Start(c.Request.Context(), "Get checked places")
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

	places, err := u.userService.GetCheckedPlaces(ctx, userID)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.BindType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, places)
}

// GetMyProperties @Summary Получить properties
// @Tags user
// @Accept  json
// @Produce  json
// @Param id query string true "userID"
// @Success 200 {object} swagger.UserMe "user properties json"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /user/properties [get]
func (u UserHandler) GetMyProperties(c *gin.Context) {
	ctx, span := u.tracer.Start(c.Request.Context(), "Get user properties")
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

	login, currentTripStartDate, properties, err := u.userService.GetProperties(ctx, userID)
	if err != nil {
		var status int
		if strings.Contains(err.Error(), "no rows in result set") {
			status = http.StatusUnauthorized
		} else {
			status = http.StatusInternalServerError
		}
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.BindType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())

		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	response := swagger.UserMe{
		Login:      login,
		Properties: properties,
	}

	if currentTripStartDate.After(time.Time{}) {
		response.CurrentTripStartDate = currentTripStartDate
	}

	c.JSON(http.StatusOK, response)
}

// UpdateProperties @Summary Получить properties
// @Tags user
// @Accept  json
// @Produce  json
// @Param data body swagger.UserUpdate true "user data"
// @Success 200 {object} string "user properties json"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /user/update_properties [put]
func (u UserHandler) UpdateProperties(c *gin.Context) {
	ctx, span := u.tracer.Start(c.Request.Context(), "Update user properties")
	defer span.End()

	var user swagger.UserUpdate

	if err := c.ShouldBindJSON(&user); err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.BindType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := u.userService.UpdateProperties(ctx, user.ID, user.Properties)
	if err != nil {
		var status int
		if strings.Contains(err.Error(), "user not found") {
			status = http.StatusUnauthorized
		} else {
			status = http.StatusInternalServerError
		}
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.BindType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())

		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// GetChrono @Summary Получить хронологию
// @Tags user
// @Accept  json
// @Produce  json
// @Param id query string true "userID"
// @Success 200 {object} models.Chrono "user properties json"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /user/chrono [get]
func (u UserHandler) GetChrono(c *gin.Context) {
	ctx, span := u.tracer.Start(c.Request.Context(), "Get user properties")
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

	chrono, err := u.userService.GetChrono(ctx, userID)
	if err != nil {
		var status int
		if strings.Contains(err.Error(), "no rows in result set") {
			status = http.StatusUnauthorized
		} else {
			status = http.StatusInternalServerError
		}
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.BindType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())

		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, chrono)
}

// GetCurrentRoute @Summary Получить хронологию
// @Tags user
// @Accept  json
// @Produce  json
// @Param id query string true "userID"
// @Success 200 {object} models.RouteDisplay "user properties json"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /user/current_route [get]
func (u UserHandler) GetCurrentRoute(c *gin.Context) {
	ctx, span := u.tracer.Start(c.Request.Context(), "Get user properties")
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

	currentRoute, err := u.userService.GetCurrentRoute(ctx, userID)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.BindType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, currentRoute)
}

// GetPlaceCheckInFlag @Summary Checkin
// @Tags user
// @Accept  json
// @Produce  json
// @Param place_id query string true "place_id"
// @Param user_id query int true "UserID"
// @Success 200 {object} string "Just valid hash"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /user/place_check_in_flag [get]
func (u UserHandler) GetPlaceCheckInFlag(c *gin.Context) {
	ctx, span := u.tracer.Start(c.Request.Context(), "User check in")
	defer span.End()

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

	userIDRaw := c.Query("user_id")
	userID, err := strconv.Atoi(userIDRaw)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.Input, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	flag, err := u.userService.GetPlaceCheckInFlag(ctx, userID, placeID)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.Input, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, flag)
}

// GetRouteCheckInFlag @Summary Checkin
// @Tags user
// @Accept  json
// @Produce  json
// @Param route_id query string true "route_id"
// @Param user_id query int true "UserID"
// @Success 200 {object} string "Just valid hash"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /user/route_check_in_flag [get]
func (u UserHandler) GetRouteCheckInFlag(c *gin.Context) {
	ctx, span := u.tracer.Start(c.Request.Context(), "User check in")
	defer span.End()

	routeIDRaw := c.Query("route_id")
	routeID, err := strconv.Atoi(routeIDRaw)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.Input, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDRaw := c.Query("user_id")
	userID, err := strconv.Atoi(userIDRaw)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.Input, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	flag, err := u.userService.GetRouteCheckInFlag(ctx, userID, routeID)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.Input, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, flag)
}
