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
	query := `INSERT INTO city (name) VALUES ('Moscow')`
	_, err := tx.Exec(query)
	if err != nil {
		_ = tx.Rollback()
		panic(fmt.Sprintf("unable to create city, err: %v", err))
	}
}

func createUser(tx *sqlx.Tx) {
	query := `INSERT INTO users (properties) VALUES (null);`
	_, err := tx.Exec(query)
	if err != nil {
		_ = tx.Rollback()
		panic(fmt.Sprintf("unable to create user 1, err: %v", err))
	}

	_, err = tx.Exec(query)
	if err != nil {
		_ = tx.Rollback()
		panic(fmt.Sprintf("unable to create user 2, err: %v", err))
	}
}

func createPlacesForNotes(tx *sqlx.Tx) {
	query := `INSERT INTO places (city_id, district_id, properties, name, variety) VALUES ($1, $2, null, $3, $4);`
	_, err := tx.Exec(query, 1, 0, "first", "first")
	if err != nil {
		_ = tx.Rollback()
		panic(fmt.Sprintf("unable to create place1, err: %v", err))
	}

	_, err = tx.Exec(query, 1, 0, "second", "second")
	if err != nil {
		_ = tx.Rollback()
		panic(fmt.Sprintf("unable to create place2, err: %v", err))
	}
}

func createCheckIn(tx *sqlx.Tx) {
	query := `INSERT INTO users_place_checkin (user_id, place_id, timestamp) VALUES ($1, $2, current_timestamp);`
	_, err := tx.Exec(query, 1, 1)
	if err != nil {
		_ = tx.Rollback()
		panic(fmt.Sprintf("unable to create checkin1, err: %v", err))
	}
}

func createNotes(tx *sqlx.Tx) {
	query := `INSERT INTO notes (user_id, place_id, properties) VALUES ($1, $2, null);`
	_, err := tx.Exec(query, 1, 1)
	if err != nil {
		_ = tx.Rollback()
		panic(fmt.Sprintf("unable to create note1, err: %v", err))
	}

	_, err = tx.Exec(query, 1, 2)
	if err != nil {
		_ = tx.Rollback()
		panic(fmt.Sprintf("unable to create note2, err: %v", err))
	}
}

func initTesNoteData(db *sqlx.DB) {
	tx, err := db.Beginx()
	if err != nil {
		panic(fmt.Sprintf("error on tx begin: %v", err))
	}
	//createCity(tx)
	//createUser(tx)
	//createPlacesForNotes(tx)
	//createCheckIn(tx)
	//createNotes(tx)

	err = tx.Commit()
	if err != nil {
		panic(fmt.Sprintf("err on commiting tx: %v", err))
	}
}

func testCase1CheckIn1Not(repo Note) {
	notes, err := repo.GetByUser(context.TODO(), 1)
	if err != nil {
		panic(fmt.Sprintf("error on get notes, err: %v", err))
	}

	count := 0
	for _, note := range notes {
		if note.IsCheckIn {
			count++
		}
	}

	if count == 1 {
		fmt.Println("test success!", notes)
	} else {
		fmt.Println("test wrong!!!!", count, notes)
	}
}

func testNoteCases(repo Note) {
	testCase1CheckIn1Not(repo)
}

func TestNoteRepo_GetByUser(t *testing.T) {
	config.InitConfig()
	db := database.GetDB()

	noteRepo := InitNoteRepo(db)

	initTesNoteData(db)

	testNoteCases(noteRepo)
}
