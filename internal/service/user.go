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
	"time"
)

type userService struct {
	userRepo      repository.User
	favouriteRepo repository.Favourite
	routeRepo     repository.Route
	placeRepo     repository.Place
	tripRepo      repository.Trip
	logger        *log.Logs
	hashes        []string
}

func InitUserService(userRepo repository.User, logger *log.Logs, favouriteRepo repository.Favourite,
	routeRepo repository.Route, placeRepo repository.Place, tripRepo repository.Trip) User {
	return &userService{
		userRepo:      userRepo,
		favouriteRepo: favouriteRepo,
		routeRepo:     routeRepo,
		placeRepo:     placeRepo,
		tripRepo:      tripRepo,
		logger:        logger,
		hashes:        make([]string, 1),
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

func contains(placeID int, places []models.PlaceIDWithPosition) bool {
	for _, val := range places {
		if placeID == val.PlaceID {
			return true
		}
	}

	return false
}

// TODO: следить что intersection верный
func (u *userService) calculateCheckInsInRoute(places []models.PlaceIDWithPosition, placeIDs []int) int {
	intersection := 0

	for _, routePlace := range places {
		for _, checkedInPlace := range placeIDs {
			if checkedInPlace == routePlace.PlaceID {
				intersection++
			}
		}
	}

	return intersection
}

func containsInt(routeIDs []int, routeID int) bool {
	for _, val := range routeIDs {
		if routeID == val {
			return true
		}
	}

	return false
}

// TODO: проверить начало маршрута, проверить конец маршрута, проверить что место в неск маршрутах юзера сразу
func (u *userService) updateRouteLogStatus(ctx context.Context, userID, placeID int) error {
	_, routeIDs, err := u.favouriteRepo.GetLikedByUser(ctx, userID)
	if err != nil {
		u.logger.Error(err.Error())
		return err
	}

	// TODO: routeIDs = append(u.TripRepo.GetRoutesByUser)

	userTrips, err := u.tripRepo.GetTripsByUser(ctx, userID)
	if err != nil {
		u.logger.Error(err.Error())
		return err
	}

	for _, userTrip := range userTrips {
		for _, userTripRoute := range userTrip.Routes {
			if !containsInt(routeIDs, userTripRoute.EntityID) {
				routeIDs = append(routeIDs, userTripRoute.EntityID)
			}
		}
	}

	for _, routeID := range routeIDs {
		routeRaw, err := u.routeRepo.GetByID(ctx, routeID)
		if err != nil {
			u.logger.Error(err.Error())
			return err
		}

		if contains(placeID, routeRaw.PlaceIDsWithPosition) {
			routeLog := models.RouteLogWithOneTime{
				UserID:    userID,
				RouteID:   routeID,
				TimeStamp: time.Now(),
			}

			placeIDs, err := u.userRepo.GetCheckedInPlaceIDs(ctx, userID)
			if err != nil {
				u.logger.Error(err.Error())
				return err
			}

			intersection := u.calculateCheckInsInRoute(routeRaw.PlaceIDsWithPosition, placeIDs)

			switch intersection {
			case 1:
				err = u.userRepo.StartRoute(ctx, routeLog)
				if err != nil {
					u.logger.Error(err.Error())
					return err
				}
			case len(routeRaw.PlaceIDsWithPosition):
				err = u.userRepo.EndRoute(ctx, routeLog)
				if err != nil {
					u.logger.Error(err.Error())
					return err
				}
			}
		}
	}

	return nil
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

	err = u.updateRouteLogStatus(ctx, userID, placeID)
	if err != nil {
		u.logger.Error(err.Error())
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

func (u *userService) GetCheckedPlaces(ctx context.Context, userID int) ([]models.Place, error) {
	placeIDs, err := u.userRepo.GetCheckedInPlaceIDs(ctx, userID)
	if err != nil {
		u.logger.Error(err.Error())
		return []models.Place{}, err
	}

	var places []models.Place
	for _, placeID := range placeIDs {
		place, err := u.placeRepo.GetByID(ctx, placeID)
		if err != nil {
			u.logger.Error(err.Error())
			return []models.Place{}, err
		}

		places = append(places, place)
	}

	return places, nil
}

func (u *userService) GetProperties(ctx context.Context, userID int) (string, interface{}, error) {
	login, properties, err := u.userRepo.GetProperties(ctx, userID)
	if err != nil {
		u.logger.Error(err.Error())
		return "", nil, err
	}

	return login, properties, nil
}

func (u *userService) UpdateProperties(ctx context.Context, userID int, properties interface{}) error {
	err := u.userRepo.UpdateProperties(ctx, userID, properties)
	if err != nil {
		u.logger.Error(err.Error())
		return err
	}

	return nil
}
