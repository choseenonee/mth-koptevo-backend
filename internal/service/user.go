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
	reviewRepo    repository.Review
	logger        *log.Logs
	hashes        []string
}

func InitUserService(userRepo repository.User, logger *log.Logs, favouriteRepo repository.Favourite,
	routeRepo repository.Route, placeRepo repository.Place, tripRepo repository.Trip, reviewRepo repository.Review) User {
	return &userService{
		userRepo:      userRepo,
		favouriteRepo: favouriteRepo,
		routeRepo:     routeRepo,
		placeRepo:     placeRepo,
		tripRepo:      tripRepo,
		reviewRepo:    reviewRepo,
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

func containsPlaceIDWithPosition(placeID int, places []models.PlaceIDWithPosition) bool {
	for _, val := range places {
		if placeID == val.PlaceID {
			return true
		}
	}

	return false
}

func containsEntityIDWithDayAndPosition(placeID int, places []models.EntityWithDayAndPosition) bool {
	for _, val := range places {
		if placeID == val.EntityID {
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

		if containsPlaceIDWithPosition(placeID, routeRaw.PlaceIDsWithPosition) {
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

func containsWithPosition(places []models.PlaceIDWithPosition, placeID int) (bool, int) {
	for _, place := range places {
		if place.PlaceID == placeID {
			return true, place.Position
		}
	}

	return false, 0
}

var badPlaceVarieties = [7]string{"Редкое событие", "Площади", "Архитектура", "Памятники", "Набережные", "Улицы", "Природа"}

func isNonCheckinable(variety interface{}) bool {
	for _, badPlaceVariety := range badPlaceVarieties {
		if variety == badPlaceVariety {
			return true
		}
	}

	return false
}

// iterDownFromPosition flag нужен для того чтобы понимать встретили ли мы конец или
func iterDownFromPosition(places []map[string]interface{}, position int) []int {
	var nonCheckinablePlaceIDs []int

	position--
	for position > 0 {
		for _, place := range places {
			if place["position"] == position {
				if isNonCheckinable(place["variety"]) {
					nonCheckinablePlaceIDs = append(nonCheckinablePlaceIDs, place["id"].(int))
					break
				} else {
					return nonCheckinablePlaceIDs
				}
			}
		}
		position--
	}

	return nil
}

func (u *userService) getPlacesWithPosition(ctx context.Context, rawPlaces []models.PlaceIDWithPosition) ([]map[string]interface{}, error) {
	var placesWithPosition []map[string]interface{}
	for _, rawPlace := range rawPlaces {
		place, err := u.placeRepo.GetByID(ctx, rawPlace.PlaceID)
		if err != nil {
			return []map[string]interface{}{}, err
		}

		placesWithPosition = append(placesWithPosition, map[string]interface{}{
			"variety": place.Variety, "position": rawPlace.Position, "id": place.ID,
		})
	}

	return placesWithPosition, nil
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

	routeLogs, err := u.userRepo.GetRouteLogs(ctx, userID)
	if err != nil {
		u.logger.Error(err.Error())
		return "", err
	}

	for _, routeLog := range routeLogs {
		if routeLog.EndTime.Before(routeLog.StartTime) {
			route, err := u.routeRepo.GetByID(ctx, routeLog.RouteId)
			if err != nil {
				u.logger.Error(err.Error())
				return "", err
			}

			isInRoute, position := containsWithPosition(route.PlaceIDsWithPosition, placeID)
			if isInRoute {
				placesWithPosition, err := u.getPlacesWithPosition(ctx, route.PlaceIDsWithPosition)
				if err != nil {
					u.logger.Error(err.Error())
					return "", err
				}

				placesToCheckIn := iterDownFromPosition(placesWithPosition, position)
				for _, placeID := range placesToCheckIn {
					err := u.userRepo.CheckInPlace(ctx, userID, placeID)
					if err != nil {
						u.logger.Error(err.Error())
						return "", err
					}
				}
			}
		}
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

func (u *userService) timeBetween(timeStart, timeEnd, entityTime time.Time) bool {
	return entityTime.After(timeStart) && entityTime.Before(timeEnd)
}

func (u *userService) GetProperties(ctx context.Context, userID int) (string, time.Time, interface{}, error) {
	login, properties, err := u.userRepo.GetProperties(ctx, userID)
	if err != nil {
		u.logger.Error(err.Error())
		return "", time.Time{}, nil, err
	}

	trips, err := u.tripRepo.GetTripsByUser(ctx, userID)
	if err != nil {
		u.logger.Error(err.Error())
		return "", time.Time{}, nil, err
	}

	var currentStartTripDate time.Time
	for _, trip := range trips {
		if u.timeBetween(trip.DateStart, trip.DateEnd, time.Now()) {
			currentStartTripDate = trip.DateStart
			break
		}
	}

	return login, currentStartTripDate, properties, nil
}

func (u *userService) UpdateProperties(ctx context.Context, userID int, properties interface{}) error {
	err := u.userRepo.UpdateProperties(ctx, userID, properties)
	if err != nil {
		u.logger.Error(err.Error())
		return err
	}

	return nil
}

func timeBetween(timeStart time.Time, timeEnd time.Time, entityTime time.Time) bool {
	return entityTime.After(timeStart) && entityTime.Before(timeEnd)
}

func tripContainsEntity(entities []models.EntityWithDayAndPosition, entityID int) bool {
	for _, entity := range entities {
		if entity.EntityID == entityID {
			return true
		}
	}

	return false
}

func (u *userService) GetChrono(ctx context.Context, userID int) (models.Chrono, error) {
	var chrono models.Chrono

	trips, err := u.tripRepo.GetTripsByUser(ctx, userID)
	if err != nil {
		err = fmt.Errorf("error on getting user trips, %v", err)
		u.logger.Error(err.Error())
		return models.Chrono{}, err
	}

	for _, trip := range trips {
		for _, routeRaw := range trip.Routes {
			route, err := u.routeRepo.GetByID(ctx, routeRaw.EntityID)
			if err != nil {
				err = fmt.Errorf("error on getting route places, %v", err)
				u.logger.Error(err.Error())
				return models.Chrono{}, err
			}
			for _, place := range route.PlaceIDsWithPosition {
				placeRaw := models.EntityWithDayAndPosition{}
				placeRaw.EntityID = place.PlaceID
				trip.Places = append(trip.Places, placeRaw)
			}
		}
	}

	placeIDs, routeIDs, err := u.favouriteRepo.GetLikedByUser(ctx, userID)
	if err != nil {
		err = fmt.Errorf("error on getting user liked data, %v", err)
		u.logger.Error(err.Error())
		return models.Chrono{}, err
	}

	var placesChrono []models.ChronoEntity
	for _, placeID := range placeIDs {
		timeStamp, err := u.favouriteRepo.GetPlaceTimestamp(ctx, userID, placeID)
		if err != nil {
			err = fmt.Errorf("error on getting place timestamp, place: %v,  %v", placeID, err)
			u.logger.Error(err.Error())
			return models.Chrono{}, err
		}

		place := models.ChronoEntity{
			ID:        placeID,
			TimeStamp: timeStamp,
			TripID:    0,
		}

		for _, trip := range trips {
			if timeBetween(trip.DateStart, trip.DateEnd, timeStamp) && tripContainsEntity(trip.Places, placeID) {
				place.TripID = trip.ID
				break
			}
		}

		placesChrono = append(placesChrono, place)
	}

	chrono.LikedPlaces = placesChrono

	var routesChrono []models.ChronoEntity
	for _, routeID := range routeIDs {
		timeStamp, err := u.favouriteRepo.GetRouteTimestamp(ctx, userID, routeID)
		if err != nil {
			err = fmt.Errorf("error on getting place timestamp, route: %v,  %v", routeID, err)
			u.logger.Error(err.Error())
			return models.Chrono{}, err
		}

		route := models.ChronoEntity{
			ID:        routeID,
			TimeStamp: timeStamp,
			TripID:    0,
		}

		for _, trip := range trips {
			if timeBetween(trip.DateStart, trip.DateEnd, timeStamp) && tripContainsEntity(trip.Routes, routeID) {
				route.TripID = trip.ID
				break
			}
		}

		routesChrono = append(routesChrono, route)
	}

	chrono.LikedRoutes = routesChrono

	placeReviews, routeReviews, err := u.reviewRepo.GetByAuthor(ctx, userID)
	if err != nil {
		err = fmt.Errorf("error in getting reviews by author, %v", err)
		u.logger.Error(err.Error())
		return models.Chrono{}, err
	}

	var placeReviewsChrono []models.ChronoEntity
	for _, placeReview := range placeReviews {
		placeReviewChrono := models.ChronoEntity{
			ID:        placeReview.ID,
			TimeStamp: placeReview.TimeStamp,
			TripID:    0,
		}

		for _, trip := range trips {
			if timeBetween(trip.DateStart, trip.DateEnd, placeReview.TimeStamp) && tripContainsEntity(trip.Places, placeReview.PlaceID) {
				placeReviewChrono.TripID = trip.ID
				break
			}
		}

		placeReviewsChrono = append(placeReviewsChrono, placeReviewChrono)
	}

	chrono.PlaceReviews = placeReviewsChrono

	var routeReviewsChrono []models.ChronoEntity
	for _, routeReview := range routeReviews {
		routeReviewChrono := models.ChronoEntity{
			ID:        routeReview.ID,
			TimeStamp: routeReview.TimeStamp,
			TripID:    0,
		}

		for _, trip := range trips {
			if timeBetween(trip.DateStart, trip.DateEnd, routeReview.TimeStamp) && tripContainsEntity(trip.Routes, routeReview.RouteID) {
				routeReviewChrono.TripID = trip.ID
				break
			}
		}

		routeReviewsChrono = append(routeReviewsChrono, routeReviewChrono)
	}

	chrono.RouteReviews = routeReviewsChrono

	checkedInPlaceIDs, err := u.userRepo.GetCheckedInPlaceIDs(ctx, userID)
	if err != nil {
		err = fmt.Errorf("error in getting checked in places by userID, %v", err)
		u.logger.Error(err.Error())
		return models.Chrono{}, err
	}

	var checkedInPlacesChrono []models.ChronoEntity
	for _, checkedInPlaceID := range checkedInPlaceIDs {
		timeStamp, err := u.userRepo.GetCheckInTimeStamp(ctx, userID, checkedInPlaceID)
		if err != nil {
			err = fmt.Errorf("error in getting timeStamp for checked in place, %v", err)
			u.logger.Error(err.Error())
			return models.Chrono{}, err
		}

		checkedInPlace := models.ChronoEntity{
			ID:        checkedInPlaceID,
			TimeStamp: timeStamp,
			TripID:    0,
		}
		for _, trip := range trips {
			if containsEntityIDWithDayAndPosition(checkedInPlaceID, trip.Places) && timeBetween(trip.DateStart, trip.DateEnd, timeStamp) {
				checkedInPlace.TripID = trip.ID
				break
			}
		}

		checkedInPlacesChrono = append(checkedInPlacesChrono, checkedInPlace)
	}

	chrono.CheckIns = checkedInPlacesChrono

	routeLogs, err := u.userRepo.GetRouteLogs(ctx, userID)
	if err != nil {
		err = fmt.Errorf("error in getting route logs by userID, %v", err)
		u.logger.Error(err.Error())
		return models.Chrono{}, err
	}

	var routeLogsChrono []models.ChronoEntity

	for _, routeLog := range routeLogs {
		routeLogChrono := models.ChronoEntity{
			ID:        routeLog.RouteId,
			TimeStamp: routeLog.StartTime,
			TripID:    0,
		}

		for _, trip := range trips {
			if tripContainsEntity(trip.Routes, routeLog.RouteId) && timeBetween(trip.DateStart, trip.DateEnd, routeLog.StartTime) {
				routeLogChrono.TripID = trip.ID
				break
			}
		}

		routeLogsChrono = append(routeLogsChrono, routeLogChrono)
	}

	chrono.RouteLogs = routeLogsChrono

	return chrono, nil
}

func (u *userService) GetCurrentRoute(ctx context.Context, userID int) (models.RouteDisplay, error) {
	userRoutes, err := u.userRepo.GetRouteLogs(ctx, userID)
	if err != nil {
		u.logger.Error(err.Error())
		return models.RouteDisplay{}, err
	}

	for _, userRoute := range userRoutes {
		if userRoute.EndTime.Before(userRoute.StartTime) {
			route, err := u.routeRepo.GetByID(ctx, userRoute.RouteId)
			if err != nil {
				u.logger.Error(err.Error())
				return models.RouteDisplay{}, err
			}

			userCheckedInPlaces, err := u.userRepo.GetCheckedInPlaceIDs(ctx, userID)
			if err != nil {
				u.logger.Error(err.Error())
				return models.RouteDisplay{}, err
			}

			var nextPlace int
			var position = 1
			var breakFlag = false
			for position <= len(route.PlaceIDsWithPosition) {
				for _, routePlace := range route.PlaceIDsWithPosition {
					if routePlace.Position == position {
						if !containsInt(userCheckedInPlaces, routePlace.PlaceID) {
							nextPlace = routePlace.PlaceID
							breakFlag = true
							break
						}
					}
					if breakFlag {
						break
					}
				}
				if breakFlag {
					break
				}
				position++
			}

			res := models.RouteDisplay{
				ID:             route.ID,
				NextPlaceID:    nextPlace,
				CompletedPlace: position - 1,
			}

			return res, nil
		}
	}

	return models.RouteDisplay{}, nil
}
