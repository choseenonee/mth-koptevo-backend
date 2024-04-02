package swagger

import "mth/internal/models"

type GetByAuthor struct {
	PlaceReviews []models.PlaceReview `json:"place_reviews"`
	RouteReviews []models.RouteReview `json:"route_reviews"`
}
