package repository

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"mth/pkg/config"
	"mth/pkg/database"
	"testing"
)

func createCity(tx *sqlx.Tx) {
	query := `INSERT INTO city (name) VALUES ($1);`
	_, err := tx.Exec(query, "Moscow")
	if err != nil {
		_ = tx.Rollback()
		panic("unable to create city")
	}

}

func TestPlaceRepo_GetAllWithFilter(t *testing.T) {
	config.InitConfig()
	db := database.GetDB()

	placeRepo := InitPlaceRepo(db)

	//places, err := placeRepo.GetAllWithFilter(context.TODO(), 0, 0, []int{}, 0)
	places, err := placeRepo.GetAllWithFilter(context.TODO(), 0, 0, []int{1}, 0, "")
	if err != nil {
		fmt.Println("err ", err)
	}
	fmt.Println(places)
}
