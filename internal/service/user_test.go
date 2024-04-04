package service

import (
	"fmt"
	"mth/internal/repository"
	"mth/pkg/config"
	"mth/pkg/database"
	"testing"
)

func TestUserService_(t *testing.T) {
	config.InitConfig()
	db := database.GetDB()

	userRepo := repository.InitUserRepo(db)

	encodedString, err := vernamCipher("1 saudasudasi1")
	fmt.Println(encodedString)

	encodedString, err = vernamCipher("2 saudasudasi1")
	fmt.Println(encodedString)

	_ = err

	userService := InitUserService(userRepo, nil)

	_ = userService
}
