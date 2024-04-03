package repository

//
//import (
//	"context"
//	"fmt"
//	"github.com/jmoiron/sqlx"
//	"mth/pkg/config"
//	"mth/pkg/database"
//	"testing"
//)
//
//func createCities(tx *sqlx.Tx) {
//	query := `INSERT INTO city (name) VALUES ($1);`
//	_, err := tx.Exec(query, "Moscow")
//	if err != nil {
//		_ = tx.Rollback()
//		panic(fmt.Sprintf("unable to create ct1, err: %v", err))
//	}
//
//	_, err = tx.Exec(query, "Koptevo")
//	if err != nil {
//		_ = tx.Rollback()
//		panic(fmt.Sprintf("unable to create ct2, err: %v", err))
//	}
//}
//
//func createTags(tx *sqlx.Tx) {
//	query := `INSERT INTO tags (name) VALUES ($1);`
//	_, err := tx.Exec(query, "Вкусно")
//	if err != nil {
//		_ = tx.Rollback()
//		panic(fmt.Sprintf("unable to create tag1, err: %v", err))
//	}
//
//	_, err = tx.Exec(query, "Невкусно")
//	if err != nil {
//		_ = tx.Rollback()
//		panic(fmt.Sprintf("unable to create tag2, err: %v", err))
//	}
//}
//
//func createDistricts(tx *sqlx.Tx) {
//	query := `INSERT INTO district (name) VALUES ($1);`
//	_, err := tx.Exec(query, "Сосенское")
//	if err != nil {
//		_ = tx.Rollback()
//		panic(fmt.Sprintf("unable to create ds1, err: %v", err))
//	}
//
//	_, err = tx.Exec(query, "Академический")
//	if err != nil {
//		_ = tx.Rollback()
//		panic(fmt.Sprintf("unable to create ds2, err: %v", err))
//	}
//}
//
//func createPlaceTagRelation(tx *sqlx.Tx, tagID, placeID int) {
//	query := `INSERT INTO places_tags (place_id, tag_id) VALUES ($1, $2);`
//	_, err := tx.Exec(query, placeID, tagID)
//	if err != nil {
//		_ = tx.Rollback()
//		panic(fmt.Sprintf("unable to create place tag relation, err: %v, tagID: %v, placeID: %v", err, tagID, placeID))
//	}
//}
//
//func createPlaces(tx *sqlx.Tx) {
//	query := `INSERT INTO places (city_id, district_id, properties, name) VALUES ($1, $2, $3, $4);`
//	_, err := tx.Exec(query, 1, 1, nil, "аб")
//	if err != nil {
//		_ = tx.Rollback()
//		panic(fmt.Sprintf("unable to create place 1, err: %v", err))
//	}
//
//	createPlaceTagRelation(tx, 1, 1)
//	createPlaceTagRelation(tx, 2, 1)
//
//	_, err = tx.Exec(query, 2, 2, nil, "вг")
//	if err != nil {
//		_ = tx.Rollback()
//		panic(fmt.Sprintf("unable to create place 2, err: %v", err))
//	}
//
//	createPlaceTagRelation(tx, 2, 2)
//}
//
//func initTestData(db *sqlx.DB) {
//	tx, err := db.Beginx()
//	if err != nil {
//		panic(fmt.Sprintf("error on tx begin: %v", err))
//	}
//	createCities(tx)
//	createDistricts(tx)
//	createTags(tx)
//	createPlaces(tx)
//
//	err = tx.Commit()
//	if err != nil {
//		panic(fmt.Sprintf("err on commiting tx: %v", err))
//	}
//}
//
//func testNoFilters(repo Place) {
//	places, err := repo.GetAllWithFilter(context.TODO(), 0, 0, []int{}, 0, "")
//	if err != nil {
//		panic(fmt.Sprintf("error on get places, err: %v", err))
//	}
//
//	if len(places) != 2 {
//		panic(fmt.Sprintf("ожидалось 2, testCaseDs1, places: %v", places))
//	}
//
//	fmt.Println("success on testNoFilters", places)
//}
//
//func testCaseDs1(repo Place) {
//	places, err := repo.GetAllWithFilter(context.TODO(), 1, 0, []int{}, 0, "")
//	if err != nil {
//		panic(fmt.Sprintf("error on get places, err: %v", err))
//	}
//
//	if len(places) != 1 {
//		panic(fmt.Sprintf("больше чем 1 вернулось мест хотя ожидалось 1, testCaseDs1, places: %v", places))
//	}
//
//	if !(places[0].Name == "аб") {
//		panic(fmt.Sprintf("не то место вернулось, places: %v", places))
//	}
//
//	fmt.Println("success on testCaseDs1", places)
//}
//
//func testCaseCt1(repo Place) {
//	places, err := repo.GetAllWithFilter(context.TODO(), 0, 1, []int{}, 0, "")
//	if err != nil {
//		panic(fmt.Sprintf("error on get places, err: %v", err))
//	}
//
//	if len(places) != 1 {
//		panic(fmt.Sprintf("больше чем 1 вернулось мест хотя ожидалось 1, testCaseCt2, places: %v", places))
//	}
//
//	if !(places[0].Name == "аб") {
//		panic(fmt.Sprintf("не то место вернулось, places: %v", places))
//	}
//
//	fmt.Println("success on testCaseCt2", places)
//}
//
//func testCaseTg1And2(repo Place) {
//	places, err := repo.GetAllWithFilter(context.TODO(), 0, 0, []int{2}, 0, "")
//	if err != nil {
//		panic(fmt.Sprintf("error on get places, err: %v", err))
//	}
//
//	if len(places) != 2 {
//		panic(fmt.Sprintf("больше чем 1 вернулось мест хотя ожидалось 2, testCaseTg1And2, places: %v", places))
//	}
//
//	fmt.Println("success on testCaseTg1And2", places)
//}
//
//func testCaseTg2(repo Place) {
//	places, err := repo.GetAllWithFilter(context.TODO(), 0, 0, []int{1, 2}, 0, "")
//	if err != nil {
//		panic(fmt.Sprintf("error on get places, err: %v", err))
//	}
//
//	if len(places) != 1 {
//		panic(fmt.Sprintf("больше чем 1 вернулось мест хотя ожидалось 1, testCaseTg2, places: %v", places))
//	}
//
//	if !(places[0].Name == "аб") {
//		panic(fmt.Sprintf("не то место вернулось, places: %v", places))
//	}
//
//	fmt.Println("success on testCaseTg2", places)
//}
//
//func testCaseNameLast(repo Place) {
//	places, err := repo.GetAllWithFilter(context.TODO(), 0, 0, []int{}, 0, "б")
//	if err != nil {
//		panic(fmt.Sprintf("error on get places, err: %v", err))
//	}
//
//	if len(places) != 1 {
//		panic(fmt.Sprintf("больше чем 1 вернулось мест хотя ожидалось 1, testCaseNameLast, places: %v", places))
//	}
//
//	if !(places[0].Name == "аб") {
//		panic(fmt.Sprintf("не то место вернулось, places: %v", places))
//	}
//
//	fmt.Println("success on testCaseNameLast", places)
//}
//
//func testCaseNameFirst(repo Place) {
//	places, err := repo.GetAllWithFilter(context.TODO(), 0, 0, []int{}, 0, "в")
//	if err != nil {
//		panic(fmt.Sprintf("error on get places, err: %v", err))
//	}
//
//	if len(places) != 1 {
//		panic(fmt.Sprintf("больше чем 1 вернулось мест хотя ожидалось 1, testCaseNameFirst, places: %v", places))
//	}
//
//	if !(places[0].Name == "вг") {
//		panic(fmt.Sprintf("не то место вернулось, places: %v", places))
//	}
//
//	fmt.Println("success on testCaseNameFirst", places)
//}
//
//func testAll(repo Place) {
//	places, err := repo.GetAllWithFilter(context.TODO(), 0, 0, []int{}, 0, "в")
//	if err != nil {
//		panic(fmt.Sprintf("error on get places, err: %v", err))
//	}
//
//	if len(places) != 1 {
//		panic(fmt.Sprintf("больше чем 1 вернулось мест хотя ожидалось 1, testCaseNameFirst, places: %v", places))
//	}
//
//	if !(places[0].Name == "вг") {
//		panic(fmt.Sprintf("не то место вернулось, places: %v", places))
//	}
//
//	fmt.Println("success on testCaseNameFirst", places)
//}
//
//func testCases(repo Place) {
//	testCaseDs1(repo)
//	testCaseCt1(repo)
//	testCaseTg1And2(repo)
//	testCaseTg2(repo)
//	testCaseNameLast(repo)
//	testCaseNameFirst(repo)
//	testAll(repo)
//	testNoFilters(repo)
//}
//
//func TestPlaceRepo_GetAllWithFilter(t *testing.T) {
//	config.InitConfig()
//	db := database.GetDB()
//
//	placeRepo := InitPlaceRepo(db)
//
//	//initTestData(db)
//
//	testCases(placeRepo)
//}
