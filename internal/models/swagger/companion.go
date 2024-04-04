package swagger

import "mth/internal/models"

type Companion struct {
	Places []models.CompanionsPlace `json:"places"`
	Routes []models.CompanionsRoute `json:"routes"`
}
