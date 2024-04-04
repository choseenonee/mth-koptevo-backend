package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/spf13/viper"
	"mth/internal/models"
	"mth/internal/repository"
	"mth/pkg/config"
	"mth/pkg/log"
	"strconv"
	"strings"
)

type userService struct {
	userRepo repository.User
	logger   *log.Logs
	hashes   []string
}

func InitUserService(userRepo repository.User, logger *log.Logs) User {
	return &userService{
		userRepo: userRepo,
		logger:   logger,
		hashes:   make([]string, 1),
	}
}

func vernamCipher(message string) (string, error) {
	key := viper.GetString(config.CipherKey)
	if len(message) > len(key) {
		return "", fmt.Errorf("сообщение не должно быть больше ключа msg: %v", message)
	}

	result := make([]byte, len(message))
	for i := 0; i < len(message); i++ {
		result[i] = message[i] ^ key[i]
	}

	return string(result), nil
}

func hashString(s string) string {
	hash := sha256.Sum256([]byte(s))
	return hex.EncodeToString(hash[:])
}

func writeHash(hash string, slice *[]string) {
	for idx, val := range *slice {
		if val == "" {
			(*slice)[idx] = hash
			return
		}
	}

	*slice = append(*slice, hash)
}

func validateHash(hash string, slice *[]string) bool {
	for idx, val := range *slice {
		if val == hash {
			(*slice)[idx] = ""
			return true
		}
	}

	return false
}

func (u *userService) CheckIn(ctx context.Context, cipher string, userID int) (string, error) {
	decodedString, err := vernamCipher(cipher)
	if err != nil {
		u.logger.Error(err.Error())
		return "", err
	}

	splittedStrings := strings.Split(decodedString, " ")
	if len(splittedStrings) != 2 {
		u.logger.Error(err.Error())
		return "", fmt.Errorf("расшифрованная строка не валидна")
	}

	placeID, err := strconv.Atoi(splittedStrings[0])
	if err != nil {
		u.logger.Error(err.Error())
		return "", fmt.Errorf("расшифрованная строка не содержит в себе валидный placeID")
	}

	err = u.userRepo.CheckInPlace(ctx, userID, placeID)
	if err != nil {
		u.logger.Error(err.Error())
		if strings.Contains(err.Error(), "unique") {
			return "", fmt.Errorf("пользователь уже чекинился в этом месте %v", err)
		}
		return "", err
	}

	hash := hashString(strconv.Itoa(placeID) + viper.GetString(config.CipherKey))

	writeHash(hash, &u.hashes)

	return hash, nil
}

func (u *userService) ValidateHash(ctx context.Context, hash string) bool {
	return validateHash(hash, &u.hashes)
}

func (u *userService) GetUser(ctx context.Context, login, password string) (int, error) {
	id, pwd, err := u.userRepo.GetUser(ctx, login)
	if err != nil {
		u.logger.Error(err.Error())
		return 0, err
	}

	if password == pwd {
		return id, nil
	} else {
		return 0, fmt.Errorf("user password isn't correct")
	}
}

func (u *userService) CreateUser(ctx context.Context, userCreate models.UserCreate) (int, error) {
	id, err := u.userRepo.CreateUser(ctx, userCreate)
	if err != nil {
		u.logger.Error(err.Error())
		return 0, err
	}

	return id, nil
}
