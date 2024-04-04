package repository

import (
	"context"
	"fmt"
	"mth/internal/models"
	"mth/pkg/config"
	"mth/pkg/database"
	"testing"
)

func likeRoute(favourite Favourite) {
	like := models.Like{
		UserID:   1,
		EntityID: 1,
	}
	err := favourite.LikeRoute(context.TODO(), like)
	if err != nil {
		panic(fmt.Sprintf("error on like route, %v", err))
	}
}

func likePlace(favourite Favourite) {
	like := models.Like{
		UserID:   1,
		EntityID: 1,
	}
	err := favourite.LikePlace(context.TODO(), like)
	if err != nil {
		panic(fmt.Sprintf("error on like place, %v", err))
	}
}

func getLikedByUser(favourite Favourite) {
	placeIDs, routeIDs, err := favourite.GetLikedByUser(context.TODO(), 1)
	if err != nil {
		panic(fmt.Sprintf("error on getting, %v", err))
	}

	if placeIDs[0] != 1 || len(placeIDs) != 1 {
		panic(fmt.Sprintf("error on getting, wrong placeIDs get %v", err))
	}
	if routeIDs[0] != 1 || len(routeIDs) != 1 {
		panic(fmt.Sprintf("error on getting, wrong routeIDs get %v", err))
	}
}

func TestFavouriteRepo(t *testing.T) {
	config.InitConfig()
	db := database.GetDB()

	favouriteRepo := InitFavouriteRepo(db)

	likeRoute(favouriteRepo)
	likePlace(favouriteRepo)
	getLikedByUser(favouriteRepo)
}
