package repository

import (
	"context"
	"fmt"
	"mth/pkg/config"
	"mth/pkg/database"
	"testing"
)

func TestPlaceRepo_GetAllWithFilter(t *testing.T) {
	config.InitConfig()
	db := database.GetDB()

	placeRepo := InitPlaceRepo(db)

	//places, err := placeRepo.GetAllWithFilter(context.TODO(), 0, 0, []int{}, 0)
	places, err := placeRepo.GetAllWithFilter(context.TODO(), 100, 1, []int{}, 0)
	if err != nil {
		fmt.Println("err ", err)
	}
	fmt.Println(places)
}
