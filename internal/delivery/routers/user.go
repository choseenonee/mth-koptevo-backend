package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"
	"mth/internal/delivery/handlers"
	"mth/internal/repository"
	"mth/internal/service"
	"mth/pkg/log"
)

func RegisterUserRouter(r *gin.Engine, db *sqlx.DB, logger *log.Logs, tracer trace.Tracer) *gin.RouterGroup {
	userRouter := r.Group("/user")

	userRepo := repository.InitUserRepo(db)
	favouriteRepo := repository.InitFavouriteRepo(db)
	routeRepo := repository.InitRouteRepo(db)
	placeRepo := repository.InitPlaceRepo(db)
	tripRepo := repository.InitTripRepo(db)
	reviewRepo := repository.InitReviewRepo(db)

	userService := service.InitUserService(userRepo, logger, favouriteRepo, routeRepo, placeRepo, tripRepo, reviewRepo)
	userHandler := handlers.InitUserHandler(userService, tracer)

	userRouter.POST("/check_in", userHandler.CheckIn)
	userRouter.POST("/validate", userHandler.ValidateHash)
	userRouter.PUT("/login", userHandler.GetUser)
	userRouter.PUT("/register", userHandler.CreateUser)
	userRouter.GET("/checked_in", userHandler.GetCheckedPlaces)
	userRouter.GET("/properties", userHandler.GetMyProperties)
	userRouter.PUT("/update_properties", userHandler.UpdateProperties)
	userRouter.GET("/chrono", userHandler.GetChrono)
	userRouter.GET("/current_route", userHandler.GetCurrentRoute)
	userRouter.GET("/place_check_in_flag", userHandler.GetPlaceCheckInFlag)
	userRouter.GET("/route_check_in_flag", userHandler.GetRouteCheckInFlag)

	return userRouter
}
