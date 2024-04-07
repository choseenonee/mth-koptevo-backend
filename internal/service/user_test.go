package service

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"mth/internal/repository"
	"mth/pkg/config"
	"testing"
)

// TODO: проверить начало маршрута, проверить конец маршрута, проверить что место в неск маршрутах юзера сразу
func InitPlaces(tx *sqlx.Tx) {
	var err error

	createPlaceQuery := `INSERT INTO places (city_id, district_id, properties, name, variety) VALUES ($1, $2, $3, $4, $5);`

	_, err = tx.ExecContext(context.TODO(), createPlaceQuery, 1, 1, []byte("[1]"),
		"first", "restik")
	if err != nil {
		_ = tx.Rollback()
		panic(fmt.Errorf("error on creating place 1 %v", err))
	}

	_, err = tx.ExecContext(context.TODO(), createPlaceQuery, 1, 1, []byte("[1]"),
		"second", "restik")
	if err != nil {
		_ = tx.Rollback()
		panic(fmt.Errorf("error on creating place 2 %v", err))
	}

	_, err = tx.ExecContext(context.TODO(), createPlaceQuery, 1, 1, []byte("[1]"),
		"third", "restik")
	if err != nil {
		_ = tx.Rollback()
		panic(fmt.Errorf("error on creating place 3 %v", err))
	}

	_, err = tx.ExecContext(context.TODO(), createPlaceQuery, 1, 1, []byte("[1]"),
		"fourth", "restik")
	if err != nil {
		_ = tx.Rollback()
		panic(fmt.Errorf("error on creating place 4 %v", err))
	}
}

func InitUserLikes(tx *sqlx.Tx) {

}

func InitRoutes(tx *sqlx.Tx) {
	var err error
	createRouteQuery := `INSERT INTO routes (city_id, price, name, properties) VALUES ($1, $2, $3, $4) RETURNING id;`

	_, err = tx.ExecContext(context.TODO(), createRouteQuery, 1, 100, "first", []byte("[1]"))
	if err != nil {
		_ = tx.Rollback()
		panic(fmt.Errorf("error on creating route 1 %v", err))
	}

	firstRoute := []int{1, 2, 3, 4}

	var count = 0
	for _, placeID := range firstRoute {
		createRoutePlaceRelationQuery := `INSERT INTO routes_places (route_id, place_id, position) VALUES ($1, $2, $3);`

		_, err = tx.ExecContext(context.TODO(), createRoutePlaceRelationQuery, 1, placeID, count)
		if err != nil {
			panic(fmt.Errorf("1 route, err on creating relation to place %v, err: %v", placeID, err))
		}

		count++
	}

	_, err = tx.ExecContext(context.TODO(), createRouteQuery, 1, 222, "second", []byte("[1]"))
	if err != nil {
		_ = tx.Rollback()
		panic(fmt.Errorf("error on creating route 2 %v", err))
	}

	secondRoute := []int{2}

	count = 0
	for _, placeID := range secondRoute {
		createRoutePlaceRelationQuery := `INSERT INTO routes_places (route_id, place_id, position) VALUES ($1, $2, $3);`

		_, err = tx.ExecContext(context.TODO(), createRoutePlaceRelationQuery, 2, placeID, count)
		if err != nil {
			panic(fmt.Errorf("2 route, err on creating relation to place %v, err: %v", placeID, err))
		}

		count++
	}

	err = tx.Commit()
	if err != nil {
		panic(fmt.Errorf("error on commiting init data: %v", err))
	}
}

func testCaseStartRoute(service User, repo repository.User) {
	cipher, _ := vernamCipher("1 sdinasdiahsduia")
	_, err := service.CheckIn(context.TODO(), cipher, 1)
	if err != nil {
		panic(fmt.Errorf("error on checkining place 1, %v", err))
	}

	routeLogs, err := repo.GetRouteLogs(context.TODO(), 1)
	if err != nil {
		panic(fmt.Errorf("error on geting logs of user routes 1, %v", err))
	}

	if len(routeLogs) != 1 || routeLogs[0].RouteId != 1 {
		panic(fmt.Errorf("wanted: len 1, route_id 1, got: %v", routeLogs))
	}

	cipher, _ = vernamCipher("2 sdinasdiahsduia")
	_, err = service.CheckIn(context.TODO(), cipher, 1)
	if err != nil {
		panic(fmt.Errorf("error on checkining place 2, %v", err))
	}

	routeLogs, err = repo.GetRouteLogs(context.TODO(), 1)
	if err != nil {
		panic(fmt.Errorf("error on geting logs of user routes 1, %v", err))
	}

	if len(routeLogs) != 2 {
		panic(fmt.Errorf("wanted: len 2, route_id 1 and 2, got: %v", routeLogs))
	}
	var testFlag = false
	for _, routes := range routeLogs {
		if routes.RouteId == 2 {
			testFlag = true
		}
	}
	if !testFlag {
		panic(fmt.Errorf("wanted: len 2, route_id 1 and 2, got: %v", routeLogs))
	}

	cipher, _ = vernamCipher("3 sdinasdiahsduia")
	_, err = service.CheckIn(context.TODO(), cipher, 1)
	if err != nil {
		panic(fmt.Errorf("error on checkining place 3, %v", err))
	}

	routeLogsNew, err := repo.GetRouteLogs(context.TODO(), 1)
	if err != nil {
		panic(fmt.Errorf("error on geting logs of user 1 routes, %v", err))
	}

	fmt.Println("SHOULD BE IDENTIC")
	fmt.Println(routeLogs)
	fmt.Println(routeLogsNew)

	cipher, _ = vernamCipher("4 sdinasdiahsduia")
	_, err = service.CheckIn(context.TODO(), cipher, 1)
	if err != nil {
		panic(fmt.Errorf("error on checkining place 4, %v", err))
	}

	routeLogsNew, err = repo.GetRouteLogs(context.TODO(), 1)
	if err != nil {
		panic(fmt.Errorf("error on geting logs of user 1 routes, %v", err))
	}

	fmt.Println(routeLogsNew)
}

func TestUserService_CheckINLogs(t *testing.T) {
	config.InitConfig()

	fmt.Println(vernamCipher("1 1dinasdiahsduia"))
	//logger, _, _ := log.InitLogger()
	//db := database.GetDB()
	//
	//userRepo := repository.InitUserRepo(db)
	//favouriteRepo := repository.InitFavouriteRepo(db)
	//routeRepo := repository.InitRouteRepo(db)
	//placeRepo := repository.InitPlaceRepo(db)
	//tripRepo := repository.InitTripRepo(db)

	//tx, err := db.Beginx()
	//if err != nil {
	//	panic(fmt.Errorf("error on beginx, %v", err))
	//}
	//InitPlaces(tx)
	//InitRoutes(tx)

	//userService := InitUserService(userRepo, logger, favouriteRepo, routeRepo, placeRepo, tripRepo)
	//
	//testCaseStartRoute(userService, userRepo)
	//
	//_ = userService
}
